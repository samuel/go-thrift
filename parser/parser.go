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
	"path/filepath"
	"strconv"
	"strings"

	"github.com/samuel/go-parser"
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
	ErrParserFail = errors.New("thrift.parser: parsing failed entirely")

	spec = parser.Spec{
		CommentStart:   "/*",
		CommentEnd:     "*/",
		CommentLine:    parser.Any(parser.String("#"), parser.String("//")),
		NestedComments: true,
		IdentStart: parser.Satisfy(
			func(c rune) bool {
				return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || c == '_'
			}),
		IdentLetter: parser.Satisfy(
			func(c rune) bool {
				return (c >= 'A' && c <= 'Z') ||
					(c >= 'a' && c <= 'z') ||
					(c >= '0' && c <= '9') ||
					c == '.' || c == '_'
			}),
		ReservedNames: []string{
			"namespace", "struct", "enum", "const", "service", "throws",
			"required", "optional", "exception", "list", "map", "set",
		},
	}
	simpleParser = buildParser()
)

func quotedString() parser.Parser {
	return func(st *parser.State) (parser.Output, bool, error) {
		next, err := st.Input.Next()
		if err != nil || next != '"' {
			return nil, false, err
		}

		st.Input.Pop(1)

		escaped := false
		runes := make([]rune, 1, 8)
		runes[0] = '"'
		for {
			next, err := st.Input.Next()
			if err != nil {
				return nil, false, err
			}
			st.Input.Pop(1)
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

		return string(runes), true, nil
	}
}

func integer() parser.Parser {
	return func(st *parser.State) (parser.Output, bool, error) {
		next, err := st.Input.Next()
		if err != nil || ((next < '0' || next > '9') && next != '-') {
			return nil, false, err
		}

		st.Input.Pop(1)

		runes := make([]rune, 1, 8)
		runes[0] = next
		for {
			next, err := st.Input.Next()
			if err == io.EOF || !(next >= '0' && next <= '9') {
				break
			} else if err != nil {
				return nil, false, err
			}
			st.Input.Pop(1)
			runes = append(runes, next)
		}

		// We're guaranteed to only have integers here so don't check the error
		i64, _ := strconv.ParseInt(string(runes), 10, 64)
		return i64, true, nil
	}
}

func float() parser.Parser {
	return func(st *parser.State) (parser.Output, bool, error) {
		next, err := st.Input.Next()
		if err != nil || ((next < '0' || next > '9') && next != '-') {
			return nil, false, err
		}

		st.Input.Pop(1)

		runes := make([]rune, 1, 8)
		runes[0] = next
		for {
			next, err := st.Input.Next()
			if err == io.EOF || !((next >= '0' && next <= '9') || next == '.') {
				break
			} else if err != nil {
				return nil, false, err
			}
			st.Input.Pop(1)
			runes = append(runes, next)
		}

		f64, err := strconv.ParseFloat(string(runes), 64)
		if err != nil {
			return nil, false, nil
		}
		return f64, true, nil
	}
}

type symbolValue struct {
	symbol string
	value  interface{}
}

func symbolDispatcher(table map[string]parser.Parser) parser.Parser {
	ws := parser.Whitespace()
	return func(st *parser.State) (parser.Output, bool, error) {
		next, err := st.Input.Next()
		if err != nil || !(next >= 'a' && next <= 'z') {
			return nil, false, err
		}
		st.Input.Pop(1)

		runes := make([]rune, 1, 8)
		runes[0] = next
		for {
			next, err := st.Input.Next()
			if err == io.EOF || next == ' ' {
				break
			} else if err != nil {
				return nil, false, err
			}
			st.Input.Pop(1)
			runes = append(runes, next)
		}

		sym := string(runes)
		par := table[sym]
		if par == nil {
			return nil, false, nil
		}
		_, ok, err := ws(st)
		if !ok || err != nil {
			return nil, false, err
		}
		out, ok, err := par(st)
		return symbolValue{sym, out}, ok, err
	}
}

func nilParser() parser.Parser {
	return func(st *parser.State) (parser.Output, bool, error) {
		return nil, true, nil
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

func buildParser() parser.Parser {
	constantValue := parser.Lexeme(parser.Any(quotedString(), integer(), float()))
	namespaceDef := parser.Collect(
		parser.Identifier(), parser.Identifier())
	includeDef := parser.Collect(
		parser.Lexeme(quotedString()))
	var typeDef func(st *parser.State) (parser.Output, bool, error)
	recurseTypeDef := func(st *parser.State) (parser.Output, bool, error) {
		return typeDef(st)
	}
	typeDef = parser.Any(
		parser.Identifier(),
		parser.Collect(parser.Symbol("list"),
			parser.Symbol("<"),
			recurseTypeDef,
			parser.Symbol(">")),
		parser.Collect(parser.Symbol("set"),
			parser.Symbol("<"),
			recurseTypeDef,
			parser.Symbol(">")),
		parser.Collect(parser.Symbol("map"),
			parser.Symbol("<"),
			recurseTypeDef,
			parser.Symbol(","),
			recurseTypeDef,
			parser.Symbol(">")),
	)
	typedefDef := parser.Collect(typeDef, parser.Identifier())
	constDef := parser.Collect(
		typeDef, parser.Identifier(), parser.Symbol("="), constantValue)
	enumItemDef := parser.Collect(
		parser.Identifier(),
		parser.Any(
			parser.All(parser.Symbol("="), parser.Lexeme(integer())),
			nilParser(),
		),
		parser.Any(parser.Symbol(","), parser.Symbol(";"), parser.Symbol("")),
	)
	enumDef := parser.Collect(
		parser.Identifier(),
		parser.Symbol("{"),
		parser.Many(enumItemDef),
		parser.Symbol("}"),
	)
	structFieldDef := parser.Collect(
		parser.Lexeme(integer()), parser.Symbol(":"),
		parser.Any(parser.Symbol("required"), parser.Symbol("optional"), parser.Symbol("")),
		typeDef, parser.Identifier(),
		// Default
		parser.Any(
			parser.All(parser.Symbol("="),
				parser.Lexeme(parser.Any(
					parser.Identifier(), quotedString(),
					parser.Try(float()), integer()))),
			nilParser(),
		),
		parser.Skip(parser.Any(parser.Symbol(","), parser.Symbol(";"), parser.Symbol(""))),
	)
	structDef := parser.Collect(
		parser.Identifier(),
		parser.Symbol("{"),
		parser.Many(structFieldDef),
		parser.Symbol("}"),
	)
	serviceMethodDef := parser.Collect(
		// // parser.Comments(),
		// parser.Whitespace(),
		typeDef, parser.Identifier(),
		parser.Symbol("("),
		parser.Many(structFieldDef),
		parser.Symbol(")"),
		// Exceptions
		parser.Any(
			parser.Collect(
				parser.Symbol("throws"),
				parser.Symbol("("),
				parser.Many(structFieldDef),
				parser.Symbol(")"),
			),
			nilParser(),
		),
		parser.Any(parser.Symbol(","), parser.Symbol(";"), parser.Symbol("")),
	)
	serviceDef := parser.Collect(
		parser.Identifier(),
		parser.Symbol("{"),
		parser.Many(serviceMethodDef),
		parser.Symbol("}"),
	)
	thriftSpec := parser.All(parser.Whitespace(), parser.Many(
		symbolDispatcher(map[string]parser.Parser{
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

func (p *Parser) outputToThrift(obj parser.Output) (*Thrift, error) {
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

	str := string(b)
	in := parser.NewStringInput(str)
	st := &parser.State{
		Input: in,
		Spec:  spec,
	}
	out, ok, err := simpleParser(st)

	if err != nil && err != io.EOF {
		return nil, err
	}
	if !ok {
		return nil, ErrParserFail
	}

	if err != io.EOF {
		_, err = st.Input.Next()
	}
	if err != io.EOF {
		pos := in.Position()
		return nil, &ErrSyntaxError{
			File:   pos.Name,
			Line:   pos.Line,
			Column: pos.Column,
			Offset: pos.Offset,
			Left:   str[pos.Offset:],
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
		filename, err = filepath.Abs(filename)
		if err != nil {
			return nil, err
		}
		filename = filepath.Clean(filename)
		r, err = os.Open(filename)
	}
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return p.Parse(r)
}
