// Copyright 2012-2015 Samuel Stauffer. All rights reserved.
// Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

package parser

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"
	"testing"
)

func TestServiceParsing(t *testing.T) {
	parser := &Parser{}
	thrift, err := parser.Parse(bytes.NewBuffer([]byte(`
		include "other.thrift"

		namespace go somepkg
		namespace python some.module123
		namespace python.py-twisted another

		const map<string,string> M1 = {"hello": "world", "goodnight": "moon"}
		const string S1 = "foo\"\tbar"
		const string S2 = 'foo\'\tbar'
		const list<i64> L = [1, 2, 3];

		union myUnion {
			1: double dbl = 1.1;
			2: string str = "2";
			3: i32 int32 = 3;
		}

		service ServiceNAME extends SomeBase {
			# authenticate method
			// comment2
			/* some other
			   comments */
			string login(1:string password) throws (1:AuthenticationException authex),
			oneway void explode();
			blah something()
		}

		struct SomeStruct {
			1: double dbl = 1.2,
			2: optional string abc
		}`)))
	if err != nil {
		t.Fatalf("Service parsing failed with error %s", err.Error())
	}

	if thrift.Includes["other"] != "other.thrift" {
		t.Errorf("Include not parsed: %+v", thrift.Includes)
	}

	if c := thrift.Constants["M1"]; c == nil {
		t.Errorf("M1 constant missing")
	} else if c.Name != "M1" {
		t.Errorf("M1 name not M1, got '%s'", c.Name)
	} else if v, e := c.Type.String(), "map<string,string>"; v != e {
		t.Errorf("Expected type '%s' for M1, got '%s'", e, v)
	} else if _, ok := c.Value.([]KeyValue); !ok {
		t.Errorf("Expected []KeyValue value for M1, got %T", c.Value)
	}

	if c := thrift.Constants["S1"]; c == nil {
		t.Errorf("S1 constant missing")
	} else if v, e := c.Value.(string), "foo\"\tbar"; e != v {
		t.Errorf("Excepted %s for constnat S1, got %s", strconv.Quote(e), strconv.Quote(v))
	}
	if c := thrift.Constants["S2"]; c == nil {
		t.Errorf("S2 constant missing")
	} else if v, e := c.Value.(string), "foo'\tbar"; e != v {
		t.Errorf("Excepted %s for constnat S2, got %s", strconv.Quote(e), strconv.Quote(v))
	}

	expConst := &Constant{
		Name: "L",
		Type: &Type{
			Name:      "list",
			ValueType: &Type{Name: "i64"},
		},
		Value: []interface{}{int64(1), int64(2), int64(3)},
	}
	if c := thrift.Constants["L"]; c == nil {
		t.Errorf("L constant missing")
	} else if !reflect.DeepEqual(c, expConst) {
		t.Errorf("Expected for L:\n%s\ngot\n%s", pprint(expConst), pprint(c))
	}

	expectedStruct := &Struct{
		Name: "SomeStruct",
		Fields: []*Field{
			{
				ID:      1,
				Name:    "dbl",
				Default: 1.2,
				Type: &Type{
					Name: "double",
				},
			},
			{
				ID:       2,
				Name:     "abc",
				Optional: true,
				Type: &Type{
					Name: "string",
				},
			},
		},
	}
	if s := thrift.Structs["SomeStruct"]; s == nil {
		t.Errorf("SomeStruct missing")
	} else if !reflect.DeepEqual(s, expectedStruct) {
		t.Errorf("Expected\n%s\ngot\n%s", pprint(expectedStruct), pprint(s))
	}

	expectedUnion := &Struct{
		Name: "myUnion",
		Fields: []*Field{
			{
				ID:      1,
				Name:    "dbl",
				Default: 1.1,
				Type: &Type{
					Name: "double",
				},
			},
			{
				ID:      2,
				Name:    "str",
				Default: "2",
				Type: &Type{
					Name: "string",
				},
			},
			{
				ID:      3,
				Name:    "int32",
				Default: int64(3),
				Type: &Type{
					Name: "i32",
				},
			},
		},
	}
	if u := thrift.Unions["myUnion"]; u == nil {
		t.Errorf("myUnion missing")
	} else if !reflect.DeepEqual(u, expectedUnion) {
		t.Errorf("Expected\n%s\ngot\n%s", pprint(expectedUnion), pprint(u))
	}

	if len(thrift.Services) != 1 {
		t.Fatalf("Parsing service returned %d services rather than 1 as expected", len(thrift.Services))
	}
	svc := thrift.Services["ServiceNAME"]
	if svc == nil || svc.Name != "ServiceNAME" {
		t.Fatalf("Parsing service expected to find 'ServiceNAME' rather than '%+v'", thrift.Services)
	} else if svc.Extends != "SomeBase" {
		t.Errorf("Expected extends 'SomeBase' got '%s'", svc.Extends)
	}

	expected := map[string]*Service{
		"ServiceNAME": &Service{
			Name:    "ServiceNAME",
			Extends: "SomeBase",
			Methods: map[string]*Method{
				"login": &Method{
					Name: "login",
					ReturnType: &Type{
						Name: "string",
					},
					Arguments: []*Field{
						&Field{
							ID:       1,
							Name:     "password",
							Optional: false,
							Type: &Type{
								Name: "string",
							},
						},
					},
					Exceptions: []*Field{
						&Field{
							ID:       1,
							Name:     "authex",
							Optional: true,
							Type: &Type{
								Name: "AuthenticationException",
							},
						},
					},
				},
				"explode": &Method{
					Name:       "explode",
					ReturnType: nil,
					Oneway:     true,
					Arguments:  []*Field{},
				},
			},
		},
	}
	for n, m := range expected["ServiceNAME"].Methods {
		if !reflect.DeepEqual(svc.Methods[n], m) {
			t.Fatalf("Parsing service returned method\n%s\ninstead of\n%s", pprint(svc.Methods[n]), pprint(m))
		}
	}
}

func TestParseConstant(t *testing.T) {
	parser := &Parser{}
	thrift, err := parser.Parse(bytes.NewBuffer([]byte(`
		const string C1 = "test"
		const string C2 = C1
		`)))
	if err != nil {
		t.Fatalf("Service parsing failed with error %s", err.Error())
	}

	expected := map[string]*Constant{
		"C1": &Constant{
			Name:  "C1",
			Type:  &Type{Name: "string"},
			Value: "test",
		},
		"C2": &Constant{
			Name:  "C2",
			Type:  &Type{Name: "string"},
			Value: Identifier("C1"),
		},
	}
	if got := thrift.Constants; !reflect.DeepEqual(expected, got) {
		t.Errorf("Unexpected constant parsing got\n%s\ninstead of\n%s", pprint(expected), pprint(got))
	}
}

// func TestParseFile(t *testing.T) {
// 	th, err := ParseFile("../testfiles/full.thrift")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	b, err := json.MarshalIndent(th, "", "    ")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	_ = b
// }

func pprint(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(b)
}
