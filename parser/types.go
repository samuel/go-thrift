// Copyright 2012-2015 Samuel Stauffer. All rights reserved.
// Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

package parser

import "fmt"

type Pos struct {
	Line int
	Col  int
}

type Type struct {
	Pos         Pos
	Name        string        `json:",omitempty"`
	KeyType     *Type         `json:",omitempty"` // If map
	ValueType   *Type         `json:",omitempty"` // If map, list, or set
	Annotations []*Annotation `json:",omitempty"`
}

type Typedef struct {
	*Type

	Pos         Pos
	Alias       string
	Annotations []*Annotation `json:",omitempty"`
}

type EnumValue struct {
	Pos         Pos
	Name        string
	Value       int
	Annotations []*Annotation `json:",omitempty"`
}

type Enum struct {
	Pos         Pos
	Name        string
	Values      map[string]*EnumValue
	Annotations []*Annotation `json:",omitempty"`
}

type Constant struct {
	Pos   Pos
	Name  string
	Type  *Type
	Value interface{}
}

type Field struct {
	Pos         Pos
	ID          int
	Name        string
	Optional    bool
	Type        *Type
	Default     interface{}   `json:",omitempty"`
	Annotations []*Annotation `json:",omitempty"`
}

type Struct struct {
	Pos         Pos
	Name        string
	Fields      []*Field
	Annotations []*Annotation `json:",omitempty"`
}

type Method struct {
	Pos         Pos
	Comment     string
	Name        string
	Oneway      bool
	ReturnType  *Type
	Arguments   []*Field
	Exceptions  []*Field      `json:",omitempty"`
	Annotations []*Annotation `json:",omitempty"`
}

type Service struct {
	Pos         Pos
	Name        string
	Extends     string `json:",omitempty"`
	Methods     map[string]*Method
	Annotations []*Annotation `json:",omitempty"`
}

type Thrift struct {
	Filename   string
	Includes   map[string]string    `json:",omitempty"` // name -> unique identifier (absolute path generally)
	Imports    map[string]*Thrift   `json:",omitempty"` // name -> imported file
	Typedefs   map[string]*Typedef  `json:",omitempty"`
	Namespaces map[string]string    `json:",omitempty"`
	Constants  map[string]*Constant `json:",omitempty"`
	Enums      map[string]*Enum     `json:",omitempty"`
	Structs    map[string]*Struct   `json:",omitempty"`
	Exceptions map[string]*Struct   `json:",omitempty"`
	Unions     map[string]*Struct   `json:",omitempty"`
	Services   map[string]*Service  `json:",omitempty"`
}

type Identifier string

type KeyValue struct {
	Key   interface{}
	Value interface{}
}

type Annotation struct {
	Pos   Pos
	Name  string
	Value string
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
