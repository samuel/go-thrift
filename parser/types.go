// Copyright 2012 Samuel Stauffer. All rights reserved.
// Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

package parser

type Type struct {
	Name      string
	KeyType   *Type // If map
	ValueType *Type // If map or list
	IncludeName string
}

type EnumValue struct {
	Name  string
	Value int
}

type Enum struct {
	Name   string
	Values map[string]*EnumValue
	IncludeName string
}

type Constant struct {
	Name  string
	Type  *Type
	Value interface{}
}

type Field struct {
	Id       int
	Name     string
	Optional bool
	Type     *Type
	Default  interface{}
}

type Struct struct {
	Name   string
	Fields []*Field
	IncludeName string
}

type Method struct {
	Comment    string
	Name       string
	Oneway     bool
	ReturnType *Type
	Arguments  []*Field
	Exceptions []*Field
}

type Service struct {
	Name    string
	Methods map[string]*Method
}

type Thrift struct {
	Includes   map[string]*Thrift
	Typedefs   map[string]*Type
	Namespaces map[string]string
	Constants  map[string]*Constant
	Enums      map[string]*Enum
	Structs    map[string]*Struct
	Exceptions map[string]*Struct
	Services   map[string]*Service
}

// Generate a combined Thrift struct with includes merged into the namespace
func (t *Thrift) MergeIncludes() *Thrift {
	if len(t.Includes) == 0 {
		return t
	}

	newT := &Thrift{
		Namespaces: make(map[string]string),
		Typedefs:   make(map[string]*Type),
		Constants:  make(map[string]*Constant),
		Enums:      make(map[string]*Enum),
		Structs:    make(map[string]*Struct),
		Exceptions: make(map[string]*Struct),
		Services:   make(map[string]*Service),
		Includes:   make(map[string]*Thrift),
	}

	for k, v := range t.Namespaces {
		newT.Namespaces[k] = v
	}

	for name, inc := range t.Includes {
		inc = inc.MergeIncludes()
		for n, t := range inc.Typedefs {
			newT.Typedefs[name+"."+n] = t
		}
		for _, c := range inc.Constants {
			newT.Constants[name+"."+c.Name] = c
		}
		for _, e := range inc.Enums {
			newT.Enums[name+"."+e.Name] = e
		}
		for _, s := range inc.Structs {
			newT.Structs[name+"."+s.Name] = s
		}
		for _, e := range inc.Exceptions {
			newT.Exceptions[name+"."+e.Name] = e
		}
		for _, s := range inc.Services {
			newT.Services[name+"."+s.Name] = s
		}
	}

	for n, t := range t.Typedefs {
		newT.Typedefs[n] = t
	}
	for _, c := range t.Constants {
		newT.Constants[c.Name] = c
	}
	for _, e := range t.Enums {
		newT.Enums[e.Name] = e
	}
	for _, s := range t.Structs {
		newT.Structs[s.Name] = s
	}
	for _, e := range t.Exceptions {
		newT.Exceptions[e.Name] = e
	}
	for _, s := range t.Services {
		newT.Services[s.Name] = s
	}

	return newT
}
