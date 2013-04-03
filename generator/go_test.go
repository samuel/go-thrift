// Copyright 2013 Samuel Stauffer. All rights reserved.
// Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"github.com/samuel/go-thrift/parser"
	"io"
	"regexp"
	"testing"
)

const TEST_SIMPLE_THRIFT = `struct UserProfile {
  1: i32 uid,
  2: string name,
  3: string blurb
}`

func GenerateThrift(name string, in io.Reader) (out string, err error) {
	var (
		p  *parser.Parser
		th *parser.Thrift
		g  *GoGenerator
		b  *bytes.Buffer
	)
	if th, err = p.Parse(in); err != nil {
		return
	}
	g = &GoGenerator{th.MergeIncludes()}
	b = new(bytes.Buffer)
	if err = g.Generate(name, b); err != nil {
		return
	}
	out = b.String()
	return
}

func Includes(pattern string, in string) bool {
	matched, err := regexp.MatchString(pattern, in)
	return matched == true && err == nil
}

func TestPackageNameWithDash(t *testing.T) {
	var (
		in  *bytes.Buffer
		out string
		err error
	)
	in = bytes.NewBufferString(TEST_SIMPLE_THRIFT)
	if out, err = GenerateThrift("foo-bar", in); err != nil {
		t.Fatalf("Could not generate Thrift: %v", err)
	}
	t.Logf("Generated Thrift:\n%v\n", out)
	if Includes("package [A-Za-z_0-9]*-[A-Za-z_0-9]*", out) {
		t.Errorf("Package name must not contain dashes")
	}
	if !Includes("package foo_bar", out) {
		t.Errorf("Package name must convert dashes to underscores")
	}
}
