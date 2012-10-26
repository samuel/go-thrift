package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/samuel/go-parse"
)

type Constant struct {
	Type  string
	Value interface{}
}

type Field struct {
	Name     string
	Optional bool
	Type     interface{} // TODO: Use something more convenient than interface{}
	Default  interface{}
}

type Struct struct {
	Fields map[int]Field
}

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
		if !ok || !(next >= '0' && next <= '9') {
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

		// We're guarantted to only have integers here so don't check the error
		i64, _ := strconv.ParseInt(string(runes), 10, 64)
		return i64, true
	}
}

func float() parsec.Parser {
	return func(in parsec.Vessel) (parsec.Output, bool) {
		next, ok := in.Next()
		if !ok || !(next >= '0' && next <= '9') {
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

func main() {
	in := new(parsec.StringVessel)
	in.SetSpec(parsec.Spec{
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
			"required", "optional", "exception", "list", "map",
		},
	})

	f, err := os.Open("cassandra.thrift")
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	in.SetInput(string(b))

	constantValue := parsec.Lexeme(parsec.Any(quotedString(), integer(), float()))
	namespaceDef := parsec.Collect(
		// parsec.Symbol("namespace"),
		parsec.Identifier(), parsec.Identifier())
	includeDef := parsec.Collect(
		// parsec.Symbol("include"),
		parsec.Lexeme(quotedString()))
	// TODO: Should be recursive
	var typeDef func(in parsec.Vessel) (parsec.Output, bool)
	recurseTypeDef := func(in parsec.Vessel) (parsec.Output, bool) {
		return typeDef(in)
	}
	typeDef = parsec.Any(
		parsec.Identifier(),
		parsec.Try(parsec.Collect(parsec.Symbol("list"),
			parsec.Symbol("<"),
			recurseTypeDef,
			parsec.Symbol(">"))),
		parsec.Collect(parsec.Symbol("map"),
			parsec.Symbol("<"),
			recurseTypeDef,
			parsec.Symbol(","),
			recurseTypeDef,
			parsec.Symbol(">")),
	)
	constDef := parsec.Collect(
		// parsec.Symbol("const"),
		typeDef, parsec.Identifier(), parsec.Symbol("="), constantValue)
	enumItemDef := parsec.Collect(
		parsec.Identifier(),
		parsec.Any(
			parsec.All(parsec.Symbol("="), parsec.Lexeme(integer())),
			parsec.Symbol(""),
		))
	enumDef := parsec.Collect(
		// parsec.Symbol("enum"),
		parsec.Identifier(),
		parsec.Symbol("{"),
		parsec.SepBy(parsec.Symbol(","), enumItemDef),
		parsec.Symbol("}"),
	)
	structFieldDef := parsec.Collect(
		parsec.Lexeme(integer()), parsec.Symbol(":"),
		parsec.Any(parsec.Symbol("required"), parsec.Symbol("optional"), parsec.Symbol("")),
		typeDef, parsec.Identifier(),
		parsec.Any(
			parsec.All(parsec.Symbol("="),
				parsec.Lexeme(parsec.Any(
					parsec.Identifier(), quotedString(),
					parsec.Try(float()), integer()))),
			parsec.Symbol(""),
		),
		parsec.Skip(parsec.Many(parsec.Symbol(","))),
	)
	exceptionDef := parsec.Collect(
		// parsec.Symbol("exception"),
		parsec.Identifier(),
		parsec.Symbol("{"),
		parsec.Many(structFieldDef),
		parsec.Symbol("}"),
	)
	structDef := parsec.Collect(
		// parsec.Symbol("struct"),
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
			parsec.Symbol(""),
		),
		parsec.Any(parsec.Symbol(","), parsec.Symbol("")),
	)
	serviceDef := parsec.Collect(
		// parsec.Symbol("service"),
		parsec.Identifier(),
		parsec.Symbol("{"),
		parsec.Many(serviceMethodDef),
		parsec.Symbol("}"),
	)
	// thriftSpec := parsec.All(parsec.Whitespace(), parsec.Many(parsec.Any(
	// 	namespaceDef,
	// 	constDef,
	// 	includeDef,
	// 	parsec.Try(enumDef),
	// 	exceptionDef,
	// 	parsec.Try(structDef),
	// 	serviceDef,
	// )))
	thriftSpec := parsec.All(parsec.Whitespace(), parsec.Many(
		symbolDispatcher(map[string]parsec.Parser{
			"namespace": namespaceDef,
			"const":     constDef,
			"include":   includeDef,
			"enum":      enumDef,
			"exception": exceptionDef,
			"struct":    structDef,
			"service":   serviceDef,
		}),
	))

	out, ok := thriftSpec(in)

	// if _, unfinished := in.Next(); unfinished {
	// 	fmt.Printf("Incomplete parse: %+v\n", out)
	// 	fmt.Println("Parse error.")
	// 	fmt.Printf("Position: %+v\n", in.GetPosition())
	// 	fmt.Printf("State: %+v\n", in.GetState())
	// 	fmt.Printf("Rest: `%s`\n", in.GetInput())
	// 	return
	// }

	fmt.Printf("Parsed: %#v\n", ok)
	// fmt.Printf("Tree: %+v\n", out)
	// fmt.Printf("Rest: %#v\n", in.GetInput())

	namespaces := make(map[string]string)
	constants := make(map[string]Constant)
	structs := make(map[string]Struct)

	thrift := out.([]interface{})
	for _, symI := range thrift {
		sym := symI.(symbolValue)
		val := sym.value.([]interface{})
		switch sym.symbol {
		case "namespace":
			namespaces[val[0].(string)] = val[1].(string)
		case "const":
			constants[val[1].(string)] = Constant{val[0].(string), val[3]}
		case "struct":
			st := Struct{
				Fields: make(map[int]Field),
			}
			// fmt.Printf("%+v\n", val)
			for _, f := range val[2].([]interface{}) {
				fmt.Printf("%+v\n", f)
			}
			structs[val[0].(string)] = st
			// default:
			// 	fmt.Printf("UNKNOWN symbol %s: %+v\n", sym.symbol, val)
		}
	}
	fmt.Printf("Namespaces: %+v\n", namespaces)
	fmt.Printf("Constants: %+v\n", constants)
}
