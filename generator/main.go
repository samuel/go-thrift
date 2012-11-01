package main

import (
	"flag"
	"fmt"
	"os"
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

	if flag.NArg() < 2 {
		fmt.Fprintf(os.Stderr, "Usage of %s: [options] inputfile outputfile\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	filename := flag.Arg(0)
	outfilename := flag.Arg(1)

	out := os.Stdout
	if outfilename != "-" {
		var err error
		out, err = os.OpenFile(outfilename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(2)
		}
	}

	p := &parser.Parser{}
	th, err := p.ParseFile(filename)
	if e, ok := err.(*parser.ErrSyntaxError); ok {
		fmt.Printf("%s\n", e.Left)
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(2)
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(2)
	}

	fp := strings.Split(filename, "/")
	name := strings.Split(fp[len(fp)-1], ".")[0]

	generator := &GoGenerator{th.MergeIncludes()}
	err = generator.Generate(name, out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(2)
	}

	if outfilename != "-" {
		out.Close()
	}
}
