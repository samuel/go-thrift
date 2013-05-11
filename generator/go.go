// Copyright 2012 Samuel Stauffer. All rights reserved.
// Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

package main

// TODO:
// - Default arguments. Possibly don't bother...

import (
	"flag"
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/samuel/go-thrift/parser"
)

var (
	f_go_binarystring = flag.Bool("go.binarystring", false, "Always use string for binary instead of []byte")
	f_go_json_enumnum = flag.Bool("go.json.enumnum", false, "For JSON marshal enums by number instead of name")
	f_go_packagename  = flag.String("go.packagename", "", "Override the package name")
	f_go_pointers     = flag.Bool("go.pointers", false, "Make all fields pointers")
)

var (
	goNamespaceOrder = []string{"go", "perl", "py"}
)

type ErrUnknownType string

func (e ErrUnknownType) Error() string {
	return fmt.Sprintf("{ErrUnknownType %s}", string(e))
}

type GoGenerator struct {
	Thrift *parser.Thrift
}

func (g *GoGenerator) error(err error) {
	panic(err)
}

func (g *GoGenerator) write(w io.Writer, f string, a ...interface{}) error {
	if _, err := io.WriteString(w, fmt.Sprintf(f, a...)); err != nil {
		g.error(err)
		return err
	}
	return nil
}

func (g *GoGenerator) formatType(typ *parser.Type) string {
	ptr := ""
	if *f_go_pointers {
		ptr = "*"
	}
	switch typ.Name {
	case "byte", "bool", "string":
		return ptr + typ.Name
	case "binary":
		if *f_go_binarystring {
			return ptr + "string"
		}
		return "[]byte"
	case "i16":
		return ptr + "int16"
	case "i32":
		return ptr + "int32"
	case "i64":
		return ptr + "int64"
	case "double":
		return ptr + "float64"
	case "set":
		valueType := g.formatType(typ.ValueType)
		if valueType == "[]byte" {
			valueType = "string"
		}
		return "map[" + valueType + "]interface{}"
	case "list":
		return "[]" + g.formatType(typ.ValueType)
	case "map":
		keyType := g.formatType(typ.KeyType)
		if keyType == "[]byte" {
			// TODO: Log, warn, do something!
			// println("key type of []byte not supported for maps")
			keyType = "string"
		}
		return "map[" + keyType + "]" + g.formatType(typ.ValueType)
	}

	if t := g.Thrift.Typedefs[typ.Name]; t != nil {
		return g.formatType(t)
	}
	if e := g.Thrift.Enums[typ.Name]; e != nil {
		return ptr + e.Name
	}
	if s := g.Thrift.Structs[typ.Name]; s != nil {
		return "*" + s.Name
	}
	if e := g.Thrift.Exceptions[typ.Name]; e != nil {
		return "*" + e.Name
	}

	g.error(ErrUnknownType(typ.Name))
	return ""
}

func (g *GoGenerator) formatField(field *parser.Field) string {
	tags := ""
	if !field.Optional {
		tags = ",required"
	}
	return fmt.Sprintf(
		"%s %s `thrift:\"%d%s\" json:\"%s\"`",
		camelCase(field.Name), g.formatType(field.Type), field.Id, tags, field.Name)
}

func (g *GoGenerator) formatArguments(arguments []*parser.Field) string {
	args := make([]string, len(arguments))
	for i, arg := range arguments {
		args[i] = fmt.Sprintf("%s %s", camelCase(arg.Name), g.formatType(arg.Type))
	}
	return strings.Join(args, ", ")
}

func (g *GoGenerator) formatReturnType(typ *parser.Type) string {
	if typ == nil || typ.Name == "void" {
		return "error"
	}
	return fmt.Sprintf("(%s, error)", g.formatType(typ))
}

func (g *GoGenerator) writeConstant(out io.Writer, c *parser.Constant) error {
	return g.write(out, "\nconst %s = %+v\n", camelCase(c.Name), c.Value)
}

func (g *GoGenerator) writeEnum(out io.Writer, enum *parser.Enum) error {
	enumName := camelCase(enum.Name)

	g.write(out, "\ntype %s int32\n\nvar(\n", enumName)

	valueNames := sortedKeys(enum.Values)

	for _, name := range valueNames {
		val := enum.Values[name]
		g.write(out, "\t%s%s = %s(%d)\n", enumName, camelCase(name), enumName, val.Value)
	}

	// EnumByName
	g.write(out, "\t%sByName = map[string]%s{\n", enumName, enumName)
	for _, name := range valueNames {
		realName := enum.Name + "." + name
		fullName := enumName + camelCase(name)
		g.write(out, "\t\t\"%s\": %s,\n", realName, fullName)
	}
	g.write(out, "\t}\n")

	// EnumByValue
	g.write(out, "\t%sByValue = map[%s]string{\n", enumName, enumName)
	for _, name := range valueNames {
		realName := enum.Name + "." + name
		fullName := enumName + camelCase(name)
		g.write(out, "\t\t%s: \"%s\",\n", fullName, realName)
	}
	g.write(out, "\t}\n")

	// end var
	g.write(out, ")\n")

	g.write(out, `
func (e %s) String() string {
	name := %sByValue[e]
	if name == "" {
		name = fmt.Sprintf("Unknown enum value %s(%%d)", e)
	}
	return name
}
`, enumName, enumName, enumName)

	if !*f_go_json_enumnum {
		g.write(out, `
func (e %s) MarshalJSON() ([]byte, error) {
	name := %sByValue[e]
	if name == "" {
		name = strconv.Itoa(int(e))
	}
	return []byte("\""+name+"\""), nil
}
`, enumName, enumName)
	}

	g.write(out, `
func (e *%s) UnmarshalJSON(b []byte) error {
	st := string(b)
	if st[0] == '"' {
		*e = %s(%sByName[st[1:len(st)-1]])
		return nil
	}
	i, err := strconv.Atoi(st)
	*e = %s(i)
	return err
}
`, enumName, enumName, enumName, enumName)

	return nil
}

func (g *GoGenerator) writeStruct(out io.Writer, st *parser.Struct) error {
	structName := camelCase(st.Name)

	g.write(out, "\ntype %s struct {\n", structName)
	for _, field := range st.Fields {
		g.write(out, "\t%s\n", g.formatField(field))
	}
	return g.write(out, "}\n")
}

func (g *GoGenerator) writeException(out io.Writer, ex *parser.Struct) error {
	if err := g.writeStruct(out, ex); err != nil {
		return err
	}

	exName := camelCase(ex.Name)

	g.write(out, "\nfunc (e *%s) Error() string {\n", exName)
	if len(ex.Fields) == 0 {
		g.write(out, "\treturn \"%s{}\"\n", exName)
	} else {
		fieldNames := make([]string, len(ex.Fields))
		fieldVars := make([]string, len(ex.Fields))
		for i, field := range ex.Fields {
			fieldNames[i] = camelCase(field.Name) + ": %+v"
			fieldVars[i] = "e." + camelCase(field.Name)
		}
		g.write(out, "\treturn fmt.Sprintf(\"%s{%s}\", %s)\n",
			exName, strings.Join(fieldNames, ", "), strings.Join(fieldVars, ", "))
	}
	return g.write(out, "}\n")
}

func (g *GoGenerator) writeService(out io.Writer, svc *parser.Service) error {
	svcName := camelCase(svc.Name)

	// Service interface

	g.write(out, "\ntype %s interface {\n", svcName)
	methodNames := sortedKeys(svc.Methods)
	for _, k := range methodNames {
		method := svc.Methods[k]
		g.write(out,
			"\t%s(%s) %s\n",
			camelCase(method.Name), g.formatArguments(method.Arguments),
			g.formatReturnType(method.ReturnType))
	}
	g.write(out, "}\n")

	// Server

	g.write(out, "\ntype %sServer struct {\n\tImplementation %s\n}\n", svcName, svcName)

	// Server method wrappers

	for _, k := range methodNames {
		method := svc.Methods[k]
		mName := camelCase(method.Name)
		resArg := ""
		if !method.Oneway {
			resArg = fmt.Sprintf(", res *%s%sResponse", svcName, mName)
		}
		g.write(out, "\nfunc (s *%sServer) %s(req *%s%sRequest%s) error {\n", svcName, mName, svcName, mName, resArg)
		args := make([]string, 0)
		for _, arg := range method.Arguments {
			aName := camelCase(arg.Name)
			args = append(args, "req."+aName)
		}
		isVoid := method.ReturnType == nil || method.ReturnType.Name == "void"
		val := ""
		if !isVoid {
			val = "val, "
		}
		g.write(out, "\t%serr := s.Implementation.%s(%s)\n", val, mName, strings.Join(args, ", "))
		if len(method.Exceptions) > 0 {
			g.write(out, "\tswitch e := err.(type) {\n")
			for _, ex := range method.Exceptions {
				g.write(out, "\tcase %s:\n\t\tres.%s = e\n\t\terr = nil\n", g.formatType(ex.Type), camelCase(ex.Name))
			}
			g.write(out, "\t}\n")
		}
		if !isVoid {
			g.write(out, "\tres.Value = val\n")
		}
		g.write(out, "\treturn err\n}\n")
	}

	for _, k := range methodNames {
		// Request struct
		method := svc.Methods[k]
		reqStructName := svcName + camelCase(method.Name) + "Request"
		if err := g.writeStruct(out, &parser.Struct{reqStructName, method.Arguments}); err != nil {
			return err
		}

		if method.Oneway {
			g.write(out, "\nfunc (r *%s) Oneway() bool {\n\treturn true\n}\n", reqStructName)
		} else {
			// Response struct
			args := make([]*parser.Field, 0, len(method.Exceptions))
			if method.ReturnType != nil && method.ReturnType.Name != "void" {
				args = append(args, &parser.Field{0, "value", len(method.Exceptions) != 0, method.ReturnType, nil})
			}
			for _, ex := range method.Exceptions {
				args = append(args, ex)
			}
			res := &parser.Struct{svcName + camelCase(method.Name) + "Response", args}
			if err := g.writeStruct(out, res); err != nil {
				return err
			}
		}
	}

	g.write(out, "\ntype %sClient struct {\n\tClient RPCClient\n}\n", svcName)

	for _, k := range methodNames {
		method := svc.Methods[k]
		methodName := camelCase(method.Name)
		returnType := "error"
		if !method.Oneway {
			returnType = g.formatReturnType(method.ReturnType)
		}
		g.write(out, "\nfunc (s *%sClient) %s(%s) %s {\n",
			svcName, methodName,
			g.formatArguments(method.Arguments),
			returnType)

		// Request
		g.write(out, "\treq := &%s%sRequest{\n", svcName, methodName)
		for _, arg := range method.Arguments {
			argName := camelCase(arg.Name)
			g.write(out, "\t\t%s: %s,\n", argName, argName)
		}
		g.write(out, "\t}\n")

		// Response
		if method.Oneway {
			// g.write(out, "\tvar res *%s%sResponse = nil\n", svcName, methodName)
			g.write(out, "\tvar res interface{} = nil\n")
		} else {
			g.write(out, "\tres := &%s%sResponse{}\n", svcName, methodName)
		}

		// Call
		g.write(out, "\terr := s.Client.Call(\"%s\", req, res)\n", method.Name)

		// Exceptions
		if len(method.Exceptions) > 0 {
			g.write(out, "\tif err == nil {\n\t\tswitch {\n")
			for _, ex := range method.Exceptions {
				exName := camelCase(ex.Name)
				g.write(out, "\t\tcase res.%s != nil:\n\t\t\terr = res.%s\n", exName, exName)
			}
			g.write(out, "\t\t}\n\t}\n")
		}

		if method.ReturnType != nil && method.ReturnType.Name != "void" {
			g.write(out, "\treturn res.Value, err\n")
		} else {
			g.write(out, "\treturn err\n")
		}
		g.write(out, "}\n")
	}

	return nil
}

func (g *GoGenerator) Generate(name string, out io.Writer) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
		}
	}()

	packageName := *f_go_packagename
	if packageName == "" {
		for _, k := range goNamespaceOrder {
			packageName = g.Thrift.Namespaces[k]
			if packageName != "" {
				parts := strings.Split(packageName, ".")
				packageName = parts[len(parts)-1]
				break
			}
		}
		if packageName == "" {
			packageName = name
		}
	}
	packageName = validIdentifier(strings.ToLower(packageName), "_")
	g.write(out, "// This file is automatically generated. Do not modify.\n")
	g.write(out, "\npackage %s\n", packageName)

	// Imports
	imports := []string{}
	if len(g.Thrift.Enums) > 0 {
		imports = append(imports, "strconv")
	}
	if len(g.Thrift.Enums) > 0 || len(g.Thrift.Exceptions) > 0 {
		imports = append(imports, "fmt")
	}
	if len(imports) > 0 {
		g.write(out, "\nimport (\n")
		for _, in := range imports {
			g.write(out, "\t\"%s\"\n", in)
		}
		g.write(out, ")\n")
	}

	//

	for _, k := range sortedKeys(g.Thrift.Constants) {
		c := g.Thrift.Constants[k]
		if err := g.writeConstant(out, c); err != nil {
			return err
		}
	}

	for _, k := range sortedKeys(g.Thrift.Enums) {
		enum := g.Thrift.Enums[k]
		if err := g.writeEnum(out, enum); err != nil {
			return err
		}
	}

	for _, k := range sortedKeys(g.Thrift.Structs) {
		st := g.Thrift.Structs[k]
		if err := g.writeStruct(out, st); err != nil {
			return err
		}
	}

	for _, k := range sortedKeys(g.Thrift.Exceptions) {
		ex := g.Thrift.Exceptions[k]
		if err := g.writeException(out, ex); err != nil {
			return err
		}
	}

	if len(g.Thrift.Services) > 0 {
		g.write(out, "\ntype RPCClient interface {\n"+
			"\tCall(method string, request interface{}, response interface{}) error\n"+
			"}\n")
	}

	for _, k := range sortedKeys(g.Thrift.Services) {
		svc := g.Thrift.Services[k]
		if err := g.writeService(out, svc); err != nil {
			return err
		}
	}

	return nil
}
