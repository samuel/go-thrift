// Copyright 2012-2015 Samuel Stauffer. All rights reserved.
// Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/samuel/go-thrift/parser"
	"github.com/samuel/go-thrift/thrift"
)

func camelCase(st string) string {
	if strings.ToUpper(st) == st {
		st = strings.ToLower(st)
	}
	return thrift.CamelCase(st)
}

func lowerCamelCase(st string) string {
	if len(st) <= 1 {
		return strings.ToLower(st)
	}
	st = thrift.CamelCase(st)
	return strings.ToLower(st[:1]) + st[1:]
}

// Converts a string to a valid Golang identifier, as defined in
// http://golang.org/ref/spec#identifier
// by converting invalid characters to the value of replace.
// If the first character is a Unicode digit, then replace is
// prepended to the string.
func validIdentifier(st string, replace string) string {
	var (
		invalidRune  = regexp.MustCompile("[^\\pL\\pN_]")
		invalidStart = regexp.MustCompile("^\\pN")
		out          string
	)
	out = invalidRune.ReplaceAllString(st, "_")
	if invalidStart.MatchString(out) {
		out = fmt.Sprintf("%v%v", replace, out)
	}
	return out
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
		fmt.Fprintf(os.Stderr, "Usage of %s: [options] inputfile outputpath\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	filename := flag.Arg(0)
	outpath := flag.Arg(1)

	p := &parser.Parser{}
	parsedThrift, _, err := p.ParseFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(2)
	}

	generator := &GoGenerator{
		ThriftFiles: parsedThrift,
		Format:      true,
		SignedBytes: *flagGoSignedBytes,
	}
	err = generator.Generate(outpath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(2)
	}
}
