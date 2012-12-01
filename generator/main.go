// Copyright 2012 Samuel Stauffer. All rights reserved.
// Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"sort"
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

// Given a map with string keys, return a sorted list of keys.
// If m is not a map or doesn't have string keys then return nil.
func sortedKeys(m interface{}) []string {
	value := reflect.ValueOf(m)
	if value.Kind() != reflect.Map || value.Type().Key().Kind() != reflect.String {
		return nil
	}

	valueKeys := value.MapKeys()
	keys := make([]string, len(valueKeys))
	for i, k := range valueKeys {
		keys[i] = k.String()
	}
	sort.Strings(keys)
	return keys
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
