// Copyright 2012-2015 Samuel Stauffer. All rights reserved.
// Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

package parser

import "fmt"

type Type struct {
	Name        string
	KeyType     *Type // If map
	ValueType   *Type // If map, list, or set
	Annotations []*Annotation
}

type Typedef struct {
	*Type

	Alias       string
	Annotations []*Annotation
}

type EnumValue struct {
	Name        string
	Value       int
	Annotations []*Annotation
}

type Enum struct {
	Name        string
	Values      map[string]*EnumValue
	Annotations []*Annotation
}

type Constant struct {
	Name  string
	Type  *Type
	Value interface{}
}

type Field struct {
	ID          int
	Name        string
	Optional    bool
	Type        *Type
	Default     interface{}
	Annotations []*Annotation
}

type Struct struct {
	Name        string
	Fields      []*Field
	Annotations []*Annotation
}

type Method struct {
	Comment     string
	Name        string
	Oneway      bool
	ReturnType  *Type
	Arguments   []*Field
	Exceptions  []*Field
	Annotations []*Annotation
}

type Service struct {
	Name        string
	Extends     string
	Methods     map[string]*Method
	Annotations []*Annotation
}

type Thrift struct {
	Includes   map[string]string // name -> unique identifier (absolute path generally)
	Typedefs   map[string]*Typedef
	Namespaces map[string]string
	Constants  map[string]*Constant
	Enums      map[string]*Enum
	Structs    map[string]*Struct
	Exceptions map[string]*Struct
	Unions     map[string]*Struct
	Services   map[string]*Service
}

type Identifier string

type KeyValue struct {
	Key, Value interface{}
}

type Annotation struct {
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
