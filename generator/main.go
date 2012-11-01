package main

import (
	"bytes"
	"flag"
	"fmt"
	"strings"

	"github.com/samuel/go-thrift"
	"github.com/samuel/go-thrift/parser"
)

func camelCase(st string) string {
	if strings.ToUpper(st) == st {
		st = strings.ToLower(st)
	}
	return thrift.CamelCase(st)
}

func lowerCamelCase(st string) string {
	// // Assume st is not unicode
	// if strings.ToUpper(st) == st {
	// 	return strings.ToLower(st)
	// }
	// st = thrift.CamelCase(st)
	// return strings.ToLower(st[:1]) + st[1:]
	return camelCase(st)
}

func main() {
	flag.Parse()

	filename := flag.Arg(0)

	p := &parser.Parser{}
	th, err := p.ParseFile(filename)
	if e, ok := err.(*parser.ErrSyntaxError); ok {
		fmt.Printf("%s\n", e.Left)
		panic(err)
	} else if err != nil {
		panic(err)
	}

	out := &bytes.Buffer{}

	fp := strings.Split(filename, "/")
	name := strings.Split(fp[len(fp)-1], ".")[0]

	generator := &GoGenerator{th.MergeIncludes()}
	err = generator.Generate(name, out)
	if err != nil {
		panic(err)
	}

	fmt.Println(out.String())
}
