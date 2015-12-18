// Copyright 2012-2015 Samuel Stauffer. All rights reserved.
// Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

package parser

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestServiceParsing(t *testing.T) {
	thrift, err := parse(`
		include "other.thrift"

		namespace go somepkg
		namespace python some.module123
		namespace python.py-twisted another

		const map<string,string> M1 = {"hello": "world", "goodnight": "moon"}
		const string S1 = "foo\"\tbar"
		const string S2 = 'foo\'\tbar'
		const list<i64> L = [1, 2, 3];

		union myUnion
		{
			1: double dbl = 1.1;
			2: string str = "2";
			3: i32 int32 = 3;
			4: i64 int64
				= 5;
		}

		enum Operation
		{
			ADD = 1,
			SUBTRACT = 2
		}

		enum NoNewLineBeforeBrace {
			ADD = 1,
			SUBTRACT = 2
		}

		service ServiceNAME extends SomeBase
		{
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
		}

		struct NewLineBeforeBrace
		{
			1: double dbl = 1.2,
			2: optional string abc
		}`)

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
				ID:       1,
				Name:     "dbl",
				Default:  1.1,
				Optional: true,
				Type: &Type{
					Name: "double",
				},
			},
			{
				ID:       2,
				Name:     "str",
				Default:  "2",
				Optional: true,
				Type: &Type{
					Name: "string",
				},
			},
			{
				ID:       3,
				Name:     "int32",
				Default:  int64(3),
				Optional: true,
				Type: &Type{
					Name: "i32",
				},
			},
			{
				ID:       4,
				Name:     "int64",
				Default:  int64(5),
				Optional: true,
				Type: &Type{
					Name: "i64",
				},
			},
		},
	}
	if u := thrift.Unions["myUnion"]; u == nil {
		t.Errorf("myUnion missing")
	} else if !reflect.DeepEqual(u, expectedUnion) {
		t.Errorf("Expected\n%s\ngot\n%s", pprint(expectedUnion), pprint(u))
	}

	expectedEnum := &Enum{
		Name: "Operation",
		Values: map[string]*EnumValue{
			"ADD": &EnumValue{
				Name:  "ADD",
				Value: 1,
			},
			"SUBTRACT": &EnumValue{
				Name:  "SUBTRACT",
				Value: 2,
			},
		},
	}
	if e := thrift.Enums["Operation"]; e == nil {
		t.Errorf("enum Operation missing")
	} else if !reflect.DeepEqual(e, expectedEnum) {
		t.Errorf("Expected\n%s\ngot\n%s", pprint(expectedEnum), pprint(e))
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

func TestParseTypeAnnotations(t *testing.T) {
	thrift, err := parse(`
typedef i64 (
	ann1 = "a1",
	ann2  =  "a2",
	js.type = 'Long'
) long (tAnn1="tv1")

typedef list<string> (a1 = "v1") listT (a2="v2")
typedef map<string,i64> (a1 = "v1") mapT (a2="v2")
typedef set<string> (a1 = "v1") setT (a2="v2")
`)
	if err != nil {
		t.Fatalf("Parse annotations failed: %v", err)
	}

	expected := map[string]*Typedef{
		"long": &Typedef{
			Alias: "long",
			Type: &Type{
				Name: "i64",
				Annotations: []*Annotation{
					{"ann1", "a1"},
					{"ann2", "a2"},
					{"js.type", "Long"},
				},
			},
			Annotations: []*Annotation{{"tAnn1", "tv1"}},
		},
		"listT": &Typedef{
			Alias: "listT",
			Type: &Type{
				Name:        "list",
				ValueType:   &Type{Name: "string"},
				Annotations: []*Annotation{{"a1", "v1"}},
			},
			Annotations: []*Annotation{{"a2", "v2"}},
		},
		"mapT": &Typedef{
			Alias: "mapT",
			Type: &Type{
				Name:        "map",
				KeyType:     &Type{Name: "string"},
				ValueType:   &Type{Name: "i64"},
				Annotations: []*Annotation{{"a1", "v1"}},
			},
			Annotations: []*Annotation{{"a2", "v2"}},
		},
		"setT": &Typedef{
			Alias: "setT",
			Type: &Type{
				Name:        "set",
				ValueType:   &Type{Name: "string"},
				Annotations: []*Annotation{{"a1", "v1"}},
			},
			Annotations: []*Annotation{{"a2", "v2"}},
		},
	}
	if got := thrift.Typedefs; !reflect.DeepEqual(expected, got) {
		t.Errorf("Unexpected annotation parsing got\n%s\n instead of\n%v", pprint(got), pprint(expected))
	}
}

func TestParseEnumAnnotations(t *testing.T) {
	thrift, err := parse(`
		enum E {
			ONE (a1="v1"),
			TWO = 2 (a2 = "v2"),
			THREE (a3 = "v3")
		} (a4 = "v4")
	`)
	if err != nil {
		t.Fatalf("Parse enum annotations failed: %v", err)
	}

	expected := map[string]*Enum{
		"E": &Enum{
			Name: "E",
			Values: map[string]*EnumValue{
				"ONE": &EnumValue{
					Name:        "ONE",
					Value:       0,
					Annotations: []*Annotation{{"a1", "v1"}},
				},
				"TWO": &EnumValue{
					Name:        "TWO",
					Value:       2,
					Annotations: []*Annotation{{"a2", "v2"}},
				},
				"THREE": &EnumValue{
					Name:        "THREE",
					Value:       3,
					Annotations: []*Annotation{{"a3", "v3"}},
				},
			},
			Annotations: []*Annotation{{"a4", "v4"}},
		},
	}
	if got := thrift.Enums; !reflect.DeepEqual(expected, got) {
		t.Errorf("Unexpected annotation parsing got\n%s\n instead of\n%v", pprint(got), pprint(expected))
	}
}

func TestParseFieldAnnotations(t *testing.T) {
	thrift, err := parse(`
		struct S {
			1: optional i32 f1 (a1 = "v1")
		}
	`)
	if err != nil {
		t.Fatalf("Parse struct like annotations failed: %v", err)
	}

	expected := map[string]*Struct{
		"S": &Struct{
			Name: "S",
			Fields: []*Field{
				&Field{
					ID:          1,
					Name:        "f1",
					Optional:    true,
					Type:        &Type{Name: "i32"},
					Annotations: []*Annotation{{"a1", "v1"}},
				},
			},
		},
	}

	if got := thrift.Structs; !reflect.DeepEqual(expected, got) {
		t.Errorf("Unexpected annotation parsing got\n%s\n instead of\n%v", pprint(got), pprint(expected))
	}
}

func TestParseStructLikeAnnotations(t *testing.T) {
	thrift, err := parse(`
		struct S {
			1: optional i32 f1
			2: optional string f2
		} (a1 = "v1")
		union U {
			1: optional i32 f1
			2: optional string f2
		} (a2 = "v2")
		exception E {
			1: optional i32 f1
			2: optional string f2
		} (a3 = "v3")
	`)
	if err != nil {
		t.Fatalf("Parse struct like annotations failed: %v", err)
	}

	expected, _ := parse("")
	fields := []*Field{
		&Field{
			ID:       1,
			Name:     "f1",
			Optional: true,
			Type:     &Type{Name: "i32"},
		},
		&Field{
			ID:       2,
			Name:     "f2",
			Optional: true,
			Type:     &Type{Name: "string"},
		},
	}
	expected.Structs = map[string]*Struct{
		"S": &Struct{
			Name:        "S",
			Fields:      fields,
			Annotations: []*Annotation{{"a1", "v1"}},
		},
	}
	expected.Unions = map[string]*Struct{
		"U": &Struct{
			Name:        "U",
			Fields:      fields,
			Annotations: []*Annotation{{"a2", "v2"}},
		},
	}
	expected.Exceptions = map[string]*Struct{
		"E": &Struct{
			Name:        "E",
			Fields:      fields,
			Annotations: []*Annotation{{"a3", "v3"}},
		},
	}
	if !reflect.DeepEqual(expected, thrift) {
		t.Errorf("Unexpected annotation parsing got\n%s\n instead of\n%v", pprint(thrift), pprint(expected))
	}
}

func TestParseServiceAnnotations(t *testing.T) {
	thrift, err := parse(`
		service S {
			void foo(1: i32 f1) (a1="v1")
		} (a2 = "v2")
	`)
	if err != nil {
		t.Fatalf("Parse service annotations failed: %v", err)
	}

	expected := map[string]*Service{
		"S": &Service{
			Name: "S",
			Methods: map[string]*Method{
				"foo": &Method{
					Name: "foo",
					Arguments: []*Field{
						&Field{
							ID:   1,
							Name: "f1",
							Type: &Type{Name: "i32"},
						},
					},
					Annotations: []*Annotation{{"a1", "v1"}},
				},
			},
			Annotations: []*Annotation{{"a2", "v2"}},
		},
	}
	if got := thrift.Services; !reflect.DeepEqual(expected, got) {
		t.Errorf("Unexpected annotation parsing got\n%s\n instead of\n%v", pprint(got), pprint(expected))
	}
}

func TestParseConstant(t *testing.T) {
	thrift, err := parse(`
		const string C1 = "test"
		const string C2 = C1
		`)
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
		t.Errorf("Unexpected constant parsing got\n%s\ninstead of\n%s", pprint(got), pprint(expected))
	}
}

func TestParseFiles(t *testing.T) {
	files := []string{
		"cassandra.thrift",
		"Hbase.thrift",
		"include_test.thrift",
	}

	for _, f := range files {
		_, err := ParseFile(filepath.Join("../testfiles", f))
		if err != nil {
			t.Errorf("Failed to parse file %q: %v", f, err)
		}
	}
}

func pprint(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func parse(contents string) (*Thrift, error) {
	parser := &Parser{}
	thrift, err := parser.Parse(strings.NewReader(contents))
	return thrift, err
}
