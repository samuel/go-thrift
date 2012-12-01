// Copyright 2012 Samuel Stauffer. All rights reserved.
// Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

package parser

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/samuel/go-parse"
)

type Filesystem interface {
	Open(filename string) (io.ReadCloser, error)
}

type Parser struct {
	Filesystem Filesystem // For handling includes. Can be set to nil to fall back to os package.
}

type ErrSyntaxError struct {
	File   string
	Line   int
	Column int
	Offset int
	Left   string
}

func (e *ErrSyntaxError) Error() string {
	return fmt.Sprintf("Syntax Error %s:%d column %d offset %d",
		e.File, e.Line, e.Column, e.Offset)
}

var (
	ErrParserFail = errors.New("Parsing failed entirely")

	spec = parsec.Spec{
		CommentStart:   "/*",
		CommentEnd:     "*/",
		CommentLine:    parsec.Any(parsec.String("#"), parsec.String("//")),
		NestedComments: true,
		IdentStart: parsec.Satisfy(
			func(c rune) bool {
				return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || c == '_'
			}),
		IdentLetter: parsec.Satisfy(
			func(c rune) bool {
				return (c >= 'A' && c <= 'Z') ||
					(c >= 'a' && c <= 'z') ||
					(c >= '0' && c <= '9') ||
					c == '.' || c == '_'
			}),
		ReservedNames: []parsec.Output{
			"namespace", "struct", "enum", "const", "service", "throws",
			"required", "optional", "exception", "list", "map", "set",
		},
	}
	simpleParser = buildParser()
)

func quotedString() parsec.Parser {
	return func(in parsec.Vessel) (parsec.Output, bool) {
		next, ok := in.Next()
		if !ok || next != '"' {
			return nil, false
		}

		in.Pop(1)

		escaped := false
		runes := make([]rune, 1, 8)
		runes[0] = '"'
		for {
			next, ok := in.Next()
			if !ok {
				return nil, false
			}
			in.Pop(1)
			if escaped {
				switch next {
				case 'n':
					next = '\n'
				case 'r':
					next = '\r'
				case 't':
					next = '\t'
				}
				runes = append(runes, next)
				escaped = false
			} else {
				if next == '\\' {
					escaped = true
				} else {
					runes = append(runes, next)
				}

				if next == '"' {
					break
				}
			}
		}

		return string(runes), true
	}
}

func integer() parsec.Parser {
	return func(in parsec.Vessel) (parsec.Output, bool) {
		next, ok := in.Next()
		if !ok || ((next < '0' || next > '9') && next != '-') {
			return nil, false
		}

		in.Pop(1)

		runes := make([]rune, 1, 8)
		runes[0] = next
		for {
			next, ok := in.Next()
			if !ok || !(next >= '0' && next <= '9') {
				break
			}
			in.Pop(1)
			runes = append(runes, next)
		}

		// We're guaranteed to only have integers here so don't check the error
		i64, _ := strconv.ParseInt(string(runes), 10, 64)
		return i64, true
	}
}

func float() parsec.Parser {
	return func(in parsec.Vessel) (parsec.Output, bool) {
		next, ok := in.Next()
		if !ok || ((next < '0' || next > '9') && next != '-') {
			return nil, false
		}

		in.Pop(1)

		runes := make([]rune, 1, 8)
		runes[0] = next
		for {
			next, ok := in.Next()
			if !ok || !((next >= '0' && next <= '9') || next == '.') {
				break
			}
			in.Pop(1)
			runes = append(runes, next)
		}

		f64, err := strconv.ParseFloat(string(runes), 64)
		if err != nil {
			return nil, false
		}
		return f64, true
	}
}

type symbolValue struct {
	symbol string
	value  interface{}
}

func symbolDispatcher(table map[string]parsec.Parser) parsec.Parser {
	ws := parsec.Whitespace()
	return func(in parsec.Vessel) (parsec.Output, bool) {
		next, ok := in.Next()
		if !ok || !(next >= 'a' && next <= 'z') {
			return nil, false
		}
		in.Pop(1)

		runes := make([]rune, 1, 8)
		runes[0] = next
		for {
			next, ok := in.Next()
			if !ok || next == ' ' {
				break
			}
			in.Pop(1)
			runes = append(runes, next)
		}

		sym := string(runes)
		par := table[sym]
		if par == nil {
			return nil, false
		}
		_, ok = ws(in)
		if !ok {
			return nil, false
		}
		out, ok := par(in)
		return symbolValue{sym, out}, ok
	}
}

func nilParser() parsec.Parser {
	return func(in parsec.Vessel) (parsec.Output, bool) {
		return nil, true
	}
}

func parseType(t interface{}) *Type {
	typ := &Type{}
	switch t2 := t.(type) {
	case string:
		if t2 == "void" {
			return nil
		}
		typ.Name = t2
	case []interface{}:
		typ.Name = t2[0].(string)
		if typ.Name == "map" {
			typ.KeyType = parseType(t2[2])
			typ.ValueType = parseType(t2[4])
		} else if typ.Name == "list" || typ.Name == "set" {
			typ.ValueType = parseType(t2[2])
		} else {
			panic("Basic type should never not be map or list: " + typ.Name)
		}
	default:
		panic("Type should never be anything but string or []interface{}")
	}
	return typ
}

func parseFields(fi []interface{}) []*Field {
	fields := make([]*Field, len(fi))
	for i, f := range fi {
		parts := f.([]interface{})
		field := &Field{}
		field.Id = int(parts[0].(int64))
		field.Optional = strings.ToLower(parts[2].(string)) == "optional"
		field.Type = parseType(parts[3])
		field.Name = parts[4].(string)
		field.Default = parts[5]
		fields[i] = field
	}
	return fields
}

func buildParser() parsec.Parser {
	constantValue := parsec.Lexeme(parsec.Any(quotedString(), integer(), float()))
	namespaceDef := parsec.Collect(
		parsec.Identifier(), parsec.Identifier())
	includeDef := parsec.Collect(
		parsec.Lexeme(quotedString()))
	var typeDef func(in parsec.Vessel) (parsec.Output, bool)
	recurseTypeDef := func(in parsec.Vessel) (parsec.Output, bool) {
		return typeDef(in)
	}
	typeDef = parsec.Any(
		parsec.Identifier(),
		parsec.Collect(parsec.Symbol("list"),
			parsec.Symbol("<"),
			recurseTypeDef,
			parsec.Symbol(">")),
		parsec.Collect(parsec.Symbol("set"),
			parsec.Symbol("<"),
			recurseTypeDef,
			parsec.Symbol(">")),
		parsec.Collect(parsec.Symbol("map"),
			parsec.Symbol("<"),
			recurseTypeDef,
			parsec.Symbol(","),
			recurseTypeDef,
			parsec.Symbol(">")),
	)
	typedefDef := parsec.Collect(typeDef, parsec.Identifier())
	constDef := parsec.Collect(
		typeDef, parsec.Identifier(), parsec.Symbol("="), constantValue)
	enumItemDef := parsec.Collect(
		parsec.Identifier(),
		parsec.Any(
			parsec.All(parsec.Symbol("="), parsec.Lexeme(integer())),
			nilParser(),
		),
		parsec.Any(parsec.Symbol(","), parsec.Symbol("")),
	)
	enumDef := parsec.Collect(
		parsec.Identifier(),
		parsec.Symbol("{"),
		parsec.Many(enumItemDef),
		parsec.Symbol("}"),
	)
	structFieldDef := parsec.Collect(
		parsec.Lexeme(integer()), parsec.Symbol(":"),
		parsec.Any(parsec.Symbol("required"), parsec.Symbol("optional"), parsec.Symbol("")),
		typeDef, parsec.Identifier(),
		// Default
		parsec.Any(
			parsec.All(parsec.Symbol("="),
				parsec.Lexeme(parsec.Any(
					parsec.Identifier(), quotedString(),
					parsec.Try(float()), integer()))),
			nilParser(),
		),
		parsec.Skip(parsec.Many(parsec.Symbol(","))),
	)
	structDef := parsec.Collect(
		parsec.Identifier(),
		parsec.Symbol("{"),
		parsec.Many(structFieldDef),
		parsec.Symbol("}"),
	)
	serviceMethodDef := parsec.Collect(
		typeDef, parsec.Identifier(),
		parsec.Symbol("("),
		parsec.Many(structFieldDef),
		parsec.Symbol(")"),
		// Exceptions
		parsec.Any(
			parsec.Collect(
				parsec.Symbol("throws"),
				parsec.Symbol("("),
				parsec.Many(structFieldDef),
				parsec.Symbol(")"),
			),
			nilParser(),
		),
		parsec.Any(parsec.Symbol(","), parsec.Symbol("")),
	)
	serviceDef := parsec.Collect(
		parsec.Identifier(),
		parsec.Symbol("{"),
		parsec.Many(serviceMethodDef),
		parsec.Symbol("}"),
	)
	thriftSpec := parsec.All(parsec.Whitespace(), parsec.Many(
		symbolDispatcher(map[string]parsec.Parser{
			"namespace": namespaceDef,
			"typedef":   typedefDef,
			"const":     constDef,
			"include":   includeDef,
			"enum":      enumDef,
			"exception": structDef,
			"struct":    structDef,
			"service":   serviceDef,
		}),
	))
	return thriftSpec
}

func (p *Parser) outputToThrift(obj parsec.Output) (*Thrift, error) {
	thrift := &Thrift{
		Namespaces: make(map[string]string),
		Typedefs:   make(map[string]*Type),
		Constants:  make(map[string]*Constant),
		Enums:      make(map[string]*Enum),
		Structs:    make(map[string]*Struct),
		Exceptions: make(map[string]*Struct),
		Services:   make(map[string]*Service),
		Includes:   make(map[string]*Thrift),
	}

	for _, symI := range obj.([]interface{}) {
		sym := symI.(symbolValue)
		val := sym.value.([]interface{})
		switch sym.symbol {
		case "namespace":
			thrift.Namespaces[strings.ToLower(val[0].(string))] = val[1].(string)
		case "typedef":
			thrift.Typedefs[val[1].(string)] = parseType(val[0])
		case "const":
			thrift.Constants[val[1].(string)] = &Constant{val[1].(string), &Type{Name: val[0].(string)}, val[3]}
		case "enum":
			en := &Enum{
				Name:   val[0].(string),
				Values: make(map[string]*EnumValue),
			}
			next := 0
			for _, e := range val[2].([]interface{}) {
				parts := e.([]interface{})
				name := parts[0].(string)
				val := -1
				if parts[1] != nil {
					val = int(parts[1].(int64))
				} else {
					val = next
				}
				if val >= next {
					next = val + 1
				}
				en.Values[name] = &EnumValue{name, val}
			}
			thrift.Enums[en.Name] = en
		case "struct":
			thrift.Structs[val[0].(string)] = &Struct{
				Name:   val[0].(string),
				Fields: parseFields(val[2].([]interface{})),
			}
		case "exception":
			thrift.Exceptions[val[0].(string)] = &Struct{
				Name:   val[0].(string),
				Fields: parseFields(val[2].([]interface{})),
			}
		case "service":
			s := &Service{
				Name:    val[0].(string),
				Methods: make(map[string]*Method),
			}
			for _, m := range val[2].([]interface{}) {
				parts := m.([]interface{})
				var exc []*Field = nil
				if parts[5] != nil {
					exc = parseFields((parts[5].([]interface{}))[2].([]interface{}))
				} else {
					exc = make([]*Field, 0)
				}
				for _, f := range exc {
					f.Optional = true
				}
				method := &Method{
					Name:       parts[1].(string),
					ReturnType: parseType(parts[0]),
					Arguments:  parseFields(parts[3].([]interface{})),
					Exceptions: exc,
				}
				s.Methods[method.Name] = method
			}
			thrift.Services[s.Name] = s
		case "include":
			filename := val[0].(string)
			filename = filename[1 : len(filename)-1]
			tr, err := p.ParseFile(filename)
			if err != nil {
				return nil, err
			}
			thrift.Includes[strings.Split(filename, ".")[0]] = tr
		default:
			panic("Should never have an unhandled symbol: " + sym.symbol)
		}
	}
	return thrift, nil
}

func (p *Parser) Parse(r io.Reader) (*Thrift, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	in := &parsec.StringVessel{}
	in.SetSpec(spec)
	in.SetInput(string(b))
	out, ok := simpleParser(in)

	if !ok {
		return nil, ErrParserFail
	}

	_, unfinished := in.Next()
	if unfinished {
		pos := in.GetPosition()
		return nil, &ErrSyntaxError{
			File:   pos.Name,
			Line:   pos.Line,
			Column: pos.Column,
			Offset: pos.Offset,
			Left:   in.GetInput().(string),
		}
	}

	return p.outputToThrift(out)
}

func (p *Parser) ParseFile(filename string) (*Thrift, error) {
	var r io.ReadCloser
	var err error
	if p.Filesystem != nil {
		r, err = p.Filesystem.Open(filename)
	} else {
		r, err = os.Open(filename)
	}
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return p.Parse(r)
}
