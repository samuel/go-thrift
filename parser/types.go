// Copyright 2012-2015 Samuel Stauffer. All rights reserved.
// Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

package parser

import "fmt"

type Type struct {
	Name        string        `json:"name,omitempty"`
	KeyType     *Type         `json:"key_type,omitempty"`   // If map
	ValueType   *Type         `json:"value_type,omitempty"` // If map, list, or set
	Annotations []*Annotation `json:"annotations,omitempty"`
}

type Typedef struct {
	*Type

	Alias       string        `json:"alias"`
	Annotations []*Annotation `json:"annotations,omitempty"`
}

type EnumValue struct {
	Name        string        `json:"name"`
	Value       int           `json:"value"`
	Annotations []*Annotation `json:"annotations,omitempty"`
}

type Enum struct {
	Name        string                `json:"name"`
	Values      map[string]*EnumValue `json:"values"`
	Annotations []*Annotation         `json:"annotations,omitempty"`
}

type Constant struct {
	Name  string `json:"name"`
	Type  *Type  `json:"type"`
	Value interface{}
}

type Field struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Optional    bool          `json:"optional"`
	Type        *Type         `json:"type"`
	Default     interface{}   `json:"default,omitempty"`
	Annotations []*Annotation `json:"annotations,omitempty"`
}

type Struct struct {
	Name        string        `json:"name"`
	Fields      []*Field      `json:"fields"`
	Annotations []*Annotation `json:"annotations,omitempty"`
}

type Method struct {
	Comment     string        `json:"comment"`
	Name        string        `json:"name"`
	Oneway      bool          `json:"one_way"`
	ReturnType  *Type         `json:"return_type"`
	Arguments   []*Field      `json:"arguments"`
	Exceptions  []*Field      `json:"exceptions,omitempty"`
	Annotations []*Annotation `json:"annotations,omitempty"`
}

type Service struct {
	Name        string             `json:"name"`
	Extends     string             `json:"extends,omitempty"`
	Methods     map[string]*Method `json:"methods"`
	Annotations []*Annotation      `json:"annotations,omitempty"`
}

type Thrift struct {
	Includes   map[string]string    `json:"includes,omitempty"` // name -> unique identifier (absolute path generally)
	Typedefs   map[string]*Typedef  `json:"typedefs,omitempty"`
	Namespaces map[string]string    `json:"namespaces,omitempty"`
	Constants  map[string]*Constant `json:"constants,omitempty"`
	Enums      map[string]*Enum     `json:"enums,omitempty"`
	Structs    map[string]*Struct   `json:"structs,omitempty"`
	Exceptions map[string]*Struct   `json:"exceptions,omitempty"`
	Unions     map[string]*Struct   `json:"unions,omitempty"`
	Services   map[string]*Service  `json:"services,omitempty"`
}

type Identifier string

type KeyValue struct {
	Key   interface{} `json:"key"`
	Value interface{} `json:"value"`
}

type Annotation struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (t *Type) String() string {
	switch t.Name {
	case "map":
		return fmt.Sprintf("map<%s,%s>", t.KeyType.String(), t.ValueType.String())
	case "list":
		return fmt.Sprintf("list<%s>", t.ValueType.String())
	case "set":
		return fmt.Sprintf("set<%s>", t.ValueType.String())
	}
	return t.Name
}
