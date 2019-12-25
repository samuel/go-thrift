// Copyright 2012-2015 Samuel Stauffer. All rights reserved.
// Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

package main

// TODO:
// - Default arguments. Possibly don't bother...

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/samuel/go-thrift/parser"
)

var (
	flagGoBinarystring    = flag.Bool("go.binarystring", false, "Always use string for binary instead of []byte")
	flagGoJSONEnumnum     = flag.Bool("go.json.enumnum", false, "For JSON marshal enums by number instead of name")
	flagGoPointers        = flag.Bool("go.pointers", false, "Make all fields pointers")
	flagGoImportPrefix    = flag.String("go.importprefix", "", "Prefix for thrift-generated go package imports")
	flagGoGenerateMethods = flag.Bool("go.generate", false, "Add testing/quick compatible Generate methods to enum types")
	flagGoSignedBytes     = flag.Bool("go.signedbytes", false, "Interpret Thrift byte as Go signed int8 type")
	flagGoRPCContext      = flag.Bool("go.rpccontext", false, "Add context.Context objects to rpc wrappers")
)

var (
	goNamespaceOrder = []string{"go", "perl", "py", "cpp", "rb", "java", "cpp2"}
)

type ErrUnknownType string

func (e ErrUnknownType) Error() string {
	return fmt.Sprintf("Unknown type %s", string(e))
}

type ErrMissingInclude string

func (e ErrMissingInclude) Error() string {
	return fmt.Sprintf("Missing include %s", string(e))
}

type GoPackage struct {
	Path string
	Name string
}

type GoGenerator struct {
	thrift *parser.Thrift
	pkg    string

	ThriftFiles map[string]*parser.Thrift
	Packages    map[string]GoPackage
	Format      bool
	Pointers    bool
	SignedBytes bool

	// package names imported
	packageNames map[string]bool
}

var goKeywords = map[string]bool{
	"break":       true,
	"default":     true,
	"func":        true,
	"interface":   true,
	"select":      true,
	"case":        true,
	"defer":       true,
	"go":          true,
	"map":         true,
	"struct":      true,
	"chan":        true,
	"else":        true,
	"goto":        true,
	"package":     true,
	"switch":      true,
	"const":       true,
	"fallthrough": true,
	"if":          true,
	"range":       true,
	"type":        true,
	"continue":    true,
	"for":         true,
	"import":      true,
	"return":      true,
	"var":         true,

	// request arguments are hardcoded to 'req' and the response to 'res'
	"req": true,
	"res": true,
	// ctx is passed as the first argument, SetContext methods are generated iff flagGoRPCContext is set
	"ctx":        true,
	"SetContext": true,
}

var basicTypes = map[string]bool{
	"byte":   true,
	"bool":   true,
	"string": true,
	"i16":    true,
	"i32":    true,
	"i64":    true,
	"double": true,
}

func init() {
	if *flagGoBinarystring {
		basicTypes["binary"] = true
	}
}

func validGoIdent(id string) string {
	if goKeywords[id] {
		return "_" + id
	}
	return id
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

type typeOption int

const (
	toOptional typeOption = 1 << iota
	toNoPointer
)

func (to typeOption) has(opt typeOption) bool {
	return (to & opt) != 0
}

func (g *GoGenerator) formatType(pkg string, thrift *parser.Thrift, typ *parser.Type, opt typeOption) string {
	// Follow includes
	if strings.Contains(typ.Name, ".") {
		// <include>.<type>
		parts := strings.SplitN(typ.Name, ".", 2)
		// Get Thrift struct for the given include
		thriftFilename := thrift.Includes[parts[0]]
		if thriftFilename == "" {
			g.error(ErrMissingInclude(parts[0]))
		}
		thrift = g.ThriftFiles[thriftFilename]
		if thrift == nil {
			g.error(ErrMissingInclude(thriftFilename))
		}
		pkg = g.Packages[thriftFilename].Name
		typ = &parser.Type{
			Name:      parts[1],
			KeyType:   typ.KeyType,
			ValueType: typ.ValueType,
		}
	}

	ptr := ""
	if !opt.has(toNoPointer) && (g.Pointers || opt.has(toOptional)) {
		ptr = "*"
	}
	switch typ.Name {
	case "binary":
		if *flagGoBinarystring {
			return ptr + "string"
		}
		return "[]byte"
	case "bool", "string":
		return ptr + typ.Name
	case "byte":
		if g.SignedBytes {
			return ptr + "int8"
		}
		return ptr + typ.Name
	case "i16":
		return ptr + "int16"
	case "i32":
		return ptr + "int32"
	case "i64":
		return ptr + "int64"
	case "double":
		return ptr + "float64"
	case "set":
		valueType := g.formatKeyType(pkg, thrift, typ.ValueType)
		return "map[" + valueType + "]struct{}"
	case "list":
		return "[]" + g.formatType(pkg, thrift, typ.ValueType, toNoPointer)
	case "map":
		keyType := g.formatKeyType(pkg, thrift, typ.KeyType)
		return "map[" + keyType + "]" + g.formatType(pkg, thrift, typ.ValueType, toNoPointer)
	}

	if t := thrift.Typedefs[typ.Name]; t != nil {
		name := camelCase(typ.Name)
		if pkg != g.pkg {
			name = pkg + "." + name
		}
		return ptr + name
	}
	if e := thrift.Enums[typ.Name]; e != nil {
		name := e.Name
		if pkg != g.pkg {
			name = pkg + "." + name
		}
		return ptr + name
	}
	if s := thrift.Structs[typ.Name]; s != nil {
		name := camelCase(s.Name)
		if pkg != g.pkg {
			name = pkg + "." + name
		}
		return ptr + name
	}
	if e := thrift.Exceptions[typ.Name]; e != nil {
		name := e.Name
		if pkg != g.pkg {
			name = pkg + "." + name
		}
		return ptr + name
	}
	if u := thrift.Unions[typ.Name]; u != nil {
		name := u.Name
		if pkg != g.pkg {
			name = pkg + "." + name
		}
		return ptr + name
	}

	g.error(ErrUnknownType(typ.Name))
	return ""
}

func (g *GoGenerator) formatKeyType(pkg string, thrift *parser.Thrift, typ *parser.Type) string {
	keyType := g.formatType(pkg, thrift, typ, toNoPointer)

	// We can't use the []byte type as a map key.  Use string instead.
	if t := thrift.Typedefs[typ.Name]; t != nil && t.Name == "binary" {
		keyType = "string"
	} else if keyType == "[]byte" {
		keyType = "string"
	}

	return keyType
}

// Follow typedefs to the actual type
func (g *GoGenerator) resolveType(typ *parser.Type) string {
	if t := g.thrift.Typedefs[typ.Name]; t != nil {
		return g.resolveType(t.Type)
	}
	return typ.Name
}

func isContainer(typ *parser.Type) bool {
	return typ.Name == "map" || typ.Name == "list" || typ.Name == "set" || (!*flagGoBinarystring && typ.Name == "binary")
}

func (g *GoGenerator) formatField(field *parser.Field) string {
	tags := ""
	jsonTags := ""
	if !field.Optional {
		tags = ",required"
	} else {
		jsonTags = ",omitempty"
	}
	var opt typeOption
	if field.Optional {
		opt |= toOptional
	}
	return fmt.Sprintf(
		"%s %s `thrift:\"%d%s\" json:\"%s%s\"`",
		camelCase(field.Name), g.formatType(g.pkg, g.thrift, field.Type, opt), field.ID, tags, field.Name, jsonTags)
}

var getterTmpl = template.Must(template.New("").Parse(`func ({{ .Receiver }} *{{ .Typ }}) Get{{ .Field }}() (val {{ .FieldTyp }}) {
	if {{ .Receiver }} != nil {{ if .Ptr }}&& {{ .Receiver }}.{{ .Field }} != nil{{ end }} {
		return {{ if .Ptr }}*{{ end }}{{ .Receiver }}.{{ .Field }}
	}

{{ if .Zero }}
	val = {{ .Zero }}
{{ end }}
	return
}
`))

func (g *GoGenerator) formatFieldGetter(receiver, typ string, field *parser.Field) string {
	buf := new(bytes.Buffer)

	zero := ""
	if field.Default != nil {
		var err error
		zero, err = g.formatValue(field.Default, field.Type)
		if err != nil {
			g.error(err)
		}
	}

	isPtr := (g.Pointers || field.Optional) && !isContainer(field.Type)
	if err := getterTmpl.Execute(buf, struct {
		Receiver, Typ   string
		Field, FieldTyp string
		Zero            string
		Ptr             bool
	}{
		Receiver: receiver,
		Typ:      typ,
		Field:    camelCase(field.Name),
		FieldTyp: g.formatType(g.pkg, g.thrift, field.Type, toNoPointer),
		Zero:     zero,
		Ptr:      isPtr,
	}); err != nil {
		g.error(err)
	}

	return buf.String()
}

func (g *GoGenerator) formatArguments(arguments []*parser.Field) string {
	args := make([]string, len(arguments))
	for i, arg := range arguments {
		var opt typeOption
		if arg.Optional {
			opt |= toOptional
		}
		args[i] = fmt.Sprintf("%s %s", validGoIdent(lowerCamelCase(arg.Name)), g.formatType(g.pkg, g.thrift, arg.Type, opt))
	}
	return strings.Join(args, ", ")
}

func (g *GoGenerator) formatReturnType(typ *parser.Type, named bool) string {
	if typ == nil || typ.Name == "void" {
		if named {
			return "(err error)"
		}
		return "error"
	}
	if named {
		return fmt.Sprintf("(ret %s, err error)", g.formatType(g.pkg, g.thrift, typ, 0))
	}
	return fmt.Sprintf("(%s, error)", g.formatType(g.pkg, g.thrift, typ, 0))
}

func (g *GoGenerator) formatValue(v interface{}, t *parser.Type) (string, error) {
	switch v2 := v.(type) {
	case string:
		return strconv.Quote(v2), nil
	case int:
		return strconv.Itoa(v2), nil
	case bool:
		if v2 {
			return "true", nil
		}
		return "false", nil
	case int64:
		if t.Name == "bool" {
			if v2 == 0 {
				return "false", nil
			}
			return "true", nil
		}
		return strconv.FormatInt(v2, 10), nil
	case float64:
		return strconv.FormatFloat(v2, 'f', -1, 64), nil
	case []interface{}:
		buf := &bytes.Buffer{}
		buf.WriteString(g.formatType(g.pkg, g.thrift, t, 0))
		buf.WriteString("{\n")
		for _, v := range v2 {
			buf.WriteString("\t\t")
			s, err := g.formatValue(v, t.ValueType)
			if err != nil {
				return "", err
			}
			buf.WriteString(s)
			if t.Name == "set" {
				buf.WriteString(": struct{}{}")
			}
			buf.WriteString(",\n")
		}
		buf.WriteString("\t}")
		return buf.String(), nil
	case []parser.KeyValue:
		buf := &bytes.Buffer{}
		buf.WriteString(g.formatType(g.pkg, g.thrift, t, toNoPointer))
		buf.WriteString("{\n")
		for _, kv := range v2 {
			buf.WriteString("\t\t")

			s, err := g.formatValue(kv.Key, t.KeyType)
			if err != nil {
				return "", err
			}
			buf.WriteString(s)
			buf.WriteString(": ")
			s, err = g.formatValue(kv.Value, t.ValueType)
			if err != nil {
				return "", err
			}

			// struct values are pointers
			if t.ValueType == nil && *flagGoPointers {
				s += ".Ptr()"
			}

			buf.WriteString(s)
			buf.WriteString(",\n")
		}
		buf.WriteString("\t}")
		return buf.String(), nil
	case parser.Identifier:
		ident := string(v2)
		idx := strings.LastIndex(ident, ".")
		if idx == -1 {
			return camelCase(ident), nil
		}

		scope := ident[:idx]
		if g.packageNames[scope] {
			scope += "."
		}

		return scope + camelCase(ident[idx+1:]), nil
	}
	return "", fmt.Errorf("unsupported value type %T", v)
}

func (g *GoGenerator) writeEnum(out io.Writer, enum *parser.Enum) error {
	enumName := camelCase(enum.Name)

	g.write(out, "\ntype %s int32\n", enumName)

	valueNames := sortedKeys(enum.Values)
	g.write(out, "\nconst (\n")
	for _, name := range valueNames {
		val := enum.Values[name]
		g.write(out, "\t%s%s %s = %d\n", enumName, camelCase(name), enumName, val.Value)
	}
	g.write(out, ")\n")

	// begin var
	g.write(out, "\nvar (\n")

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

	if !*flagGoJSONEnumnum {
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

	g.write(out, `
func (e %s) Ptr() *%s {
	return &e
}
`, enumName, enumName)

	if *flagGoGenerateMethods {
		valueStrings := make([]string, 0, len(enum.Values))
		for _, val := range enum.Values {
			valueStrings = append(valueStrings, strconv.FormatInt(int64(val.Value), 10))
		}
		sort.Strings(valueStrings)
		valueStringsName := strings.ToLower(enumName) + "Values"

		g.write(out, `
var %s = []int32{%s}

func (e *%s) Generate(rand *rand.Rand, size int) reflect.Value {
	v := %s(%s[rand.Intn(%d)])
	return reflect.ValueOf(&v)
}
`, valueStringsName, strings.Join(valueStrings, ", "), enumName, enumName, valueStringsName, len(valueNames))
	}

	return nil
}

func (g *GoGenerator) writeStruct(out io.Writer, st *parser.Struct, includeContext bool) error {
	structName := camelCase(st.Name)

	g.write(out, "\ntype %s struct {\n", structName)
	if includeContext {
		g.write(out, "\tctx context.Context\n")
	}

	for _, field := range st.Fields {
		g.write(out, "\t%s\n", g.formatField(field))
	}
	g.write(out, "}\n")

	receiver := strings.ToLower(structName[:1])

	// generate Get methods for all fields
	for _, field := range st.Fields {
		g.write(out, "%s\n", g.formatFieldGetter(receiver, structName, field))
	}

	if includeContext {
		g.write(out, `
func (%s *%s) SetContext(ctx context.Context) {
	%s.ctx = ctx
}
`, receiver, structName, receiver)
	}

	return g.write(out, "\n")
}

func (g *GoGenerator) writeException(out io.Writer, ex *parser.Struct) error {
	if err := g.writeStruct(out, ex, false); err != nil {
		return err
	}

	exName := camelCase(ex.Name)

	g.write(out, "\nfunc (e %s) Error() string {\n", exName)
	if len(ex.Fields) == 0 {
		g.write(out, "\treturn \"%s{}\"\n", exName)
	} else {
		fieldNames := make([]string, len(ex.Fields))
		fieldVars := make([]string, len(ex.Fields))
		for i, field := range ex.Fields {
			fieldNames[i] = camelCase(field.Name) + ": %+v"
			fieldVars[i] = "e.Get" + camelCase(field.Name) + "()"
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
	if svc.Extends != "" {
		g.write(out, "\t%s\n", svc.Extends)
	}
	methodNames := sortedKeys(svc.Methods)
	for _, k := range methodNames {
		method := svc.Methods[k]
		args := g.formatArguments(method.Arguments)
		if *flagGoRPCContext {
			args = "ctx context.Context, " + args
		}

		g.write(out,
			"\t%s(%s) %s\n",
			camelCase(method.Name), args,
			g.formatReturnType(method.ReturnType, false))
	}
	g.write(out, "}\n")

	// Server

	if svc.Extends == "" {
		g.write(out, "\ntype %sServer struct {\n\tImplementation %s\n}\n", svcName, svcName)
	} else {
		g.write(out, "\ntype %sServer struct {\n\t%sServer\n\tImplementation %s\n}\n", svcName, svc.Extends, svcName)
	}

	// Server method wrappers

	for _, k := range methodNames {
		method := svc.Methods[k]
		mName := camelCase(method.Name)

		requestStructName := "InternalRPC" + svcName + camelCase(method.Name) + "Request"
		responseStructName := "InternalRPC" + svcName + camelCase(method.Name) + "Response"

		resArg := ""
		if !method.Oneway {
			resArg = fmt.Sprintf(", res *%s", responseStructName)
		}
		g.write(out, "\nfunc (s *%sServer) %s(req *%s%s) error {\n", svcName, mName, requestStructName, resArg)
		var args []string

		if *flagGoRPCContext {
			args = append(args, "req.ctx")
		}

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
				g.write(out, "\tcase %s:\n\t\tres.%s = e\n\t\terr = nil\n", g.formatType(g.pkg, g.thrift, ex.Type, toOptional), camelCase(ex.Name))
			}
			g.write(out, "\t}\n")
		}
		if !isVoid {
			if !g.Pointers && !isContainer(method.ReturnType) {
				g.write(out, "\tres.Value = &val\n")
			} else {
				g.write(out, "\tres.Value = val\n")
			}
		}
		g.write(out, "\treturn err\n}\n")
	}

	for _, k := range methodNames {
		// Request struct
		method := svc.Methods[k]

		requestStructName := "InternalRPC" + svcName + camelCase(method.Name) + "Request"
		responseStructName := "InternalRPC" + svcName + camelCase(method.Name) + "Response"

		if err := g.writeStruct(out, &parser.Struct{Name: requestStructName, Fields: method.Arguments}, *flagGoRPCContext); err != nil {
			return err
		}

		if method.Oneway {
			g.write(out, "\nfunc (r *%s) Oneway() bool {\n\treturn true\n}\n", requestStructName)
		} else {
			// Response struct
			args := make([]*parser.Field, 0, len(method.Exceptions))
			if method.ReturnType != nil && method.ReturnType.Name != "void" {
				args = append(args, &parser.Field{ID: 0, Name: "value", Optional: true /*len(method.Exceptions) != 0*/, Type: method.ReturnType, Default: nil})
			}
			for _, ex := range method.Exceptions {
				args = append(args, ex)
			}
			res := &parser.Struct{Name: responseStructName, Fields: args}
			if err := g.writeStruct(out, res, false); err != nil {
				return err
			}
		}
	}

	if svc.Extends == "" {
		g.write(out, "\ntype %sClient struct {\n\tClient RPCClient\n}\n", svcName)
	} else {
		g.write(out, "\ntype %sClient struct {\n\t%sClient\n}\n", svcName, svc.Extends)
	}

	for _, k := range methodNames {
		method := svc.Methods[k]

		requestStructName := "InternalRPC" + svcName + camelCase(method.Name) + "Request"
		responseStructName := "InternalRPC" + svcName + camelCase(method.Name) + "Response"

		methodName := camelCase(method.Name)
		returnType := "(err error)"
		if !method.Oneway {
			returnType = g.formatReturnType(method.ReturnType, true)
		}
		g.write(out, "\nfunc (s *%sClient) %s(%s) %s {\n",
			svcName, methodName,
			g.formatArguments(method.Arguments),
			returnType)

		// Request
		g.write(out, "\treq := &%s{\n", requestStructName)
		for _, arg := range method.Arguments {
			g.write(out, "\t\t%s: %s,\n", camelCase(arg.Name), validGoIdent(lowerCamelCase(arg.Name)))
		}
		g.write(out, "\t}\n")

		// Response
		if method.Oneway {
			// g.write(out, "\tvar res *%s%sResponse = nil\n", svcName, methodName)
			g.write(out, "\tvar res interface{} = nil\n")
		} else {
			g.write(out, "\tres := &%s{}\n", responseStructName)
		}

		// Call
		g.write(out, "\terr = s.Client.Call(\"%s\", req, res)\n", method.Name)

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
			if !g.Pointers && !isContainer(method.ReturnType) {
				g.write(out, "\tif err == nil && res.Value != nil {\n\t ret = *res.Value\n}\n")
			} else {
				g.write(out, "\tif err == nil {\n\tret = res.Value\n}\n")
			}
		}
		g.write(out, "\treturn\n")
		g.write(out, "}\n")
	}

	return nil
}

var validMapKeys = map[string]bool{
	"string": true,
	"i32":    true,
	"i64":    true,
	"bool":   true,
	"double": true,
}

func (g *GoGenerator) isValidGoType(typ *parser.Type) bool {
	if typ.KeyType == nil {
		return true
	}

	if _, ok := g.thrift.Enums[g.resolveType(typ.KeyType)]; ok {
		return true
	}

	return validMapKeys[g.resolveType(typ.KeyType)]
}

func (g *GoGenerator) generateSingle(out io.Writer, thriftPath string, thrift *parser.Thrift) {
	packageName := g.Packages[thriftPath].Name
	g.thrift = thrift
	g.pkg = packageName

	g.write(out, "// This file is automatically generated. Do not modify.\n")
	g.write(out, "\npackage %s\n", packageName)

	// Imports
	imports := []string{"fmt"}
	if len(thrift.Enums) > 0 {
		imports = append(imports, "strconv")

		if *flagGoGenerateMethods {
			imports = append(imports, "math/rand", "reflect")
		}
	}

	if len(thrift.Services) > 0 && *flagGoRPCContext {
		imports = append(imports, "golang.org/x/net/context")
	}

	if len(thrift.Includes) > 0 {
		for _, path := range thrift.Includes {
			pkg := g.Packages[path].Name
			if pkg != packageName {
				if *flagGoImportPrefix != "" {
					pkg = filepath.Join(*flagGoImportPrefix, pkg)
				}

				imports = append(imports, pkg)
			}
		}
	}
	if len(imports) > 0 {
		g.write(out, "\nimport (\n")
		for _, in := range imports {
			g.write(out, "\t\"%s\"\n", in)
		}
		g.write(out, ")\n")
	}

	g.write(out, "\nvar _ = fmt.Sprintf\n")

	if len(thrift.Typedefs) > 0 {
		g.write(out, "\n")
		for _, k := range sortedKeys(thrift.Typedefs) {
			t := thrift.Typedefs[k]
			g.write(out, "type %s %s\n", camelCase(k), g.formatType(g.pkg, g.thrift, t.Type, toNoPointer))
			g.write(out, fmt.Sprintf("func (e %s) Ptr() *%s { return &e }\n", camelCase(k), camelCase(k)))
		}
	}

	if len(thrift.Constants) > 0 {
		for _, k := range sortedKeys(thrift.Constants) {
			c := thrift.Constants[k]

			if !g.isValidGoType(c.Type) {
				log.Printf("Skipping generation for constant %s - type is not a valid go type (%s)\n", c.Name, g.resolveType(c.Type.KeyType))
				continue
			}

			v, err := g.formatValue(c.Value, c.Type)
			if err != nil {
				g.error(err)
			}

			if c.Type.Name == "list" || c.Type.Name == "map" || c.Type.Name == "set" {
				g.write(out, "var ")
			} else {
				g.write(out, "const ")
			}
			g.write(out, "\t%s = %+v\n", camelCase(c.Name), v)
		}
	}

	for _, k := range sortedKeys(thrift.Enums) {
		enum := thrift.Enums[k]
		if err := g.writeEnum(out, enum); err != nil {
			g.error(err)
		}
	}

	for _, k := range sortedKeys(thrift.Structs) {
		st := thrift.Structs[k]
		if err := g.writeStruct(out, st, false); err != nil {
			g.error(err)
		}
	}

	for _, k := range sortedKeys(thrift.Exceptions) {
		ex := thrift.Exceptions[k]
		if err := g.writeException(out, ex); err != nil {
			g.error(err)
		}
	}

	for _, k := range sortedKeys(thrift.Unions) {
		un := thrift.Unions[k]
		if err := g.writeStruct(out, un, false); err != nil {
			g.error(err)
		}
	}

	for _, k := range sortedKeys(thrift.Services) {
		svc := thrift.Services[k]
		if err := g.writeService(out, svc); err != nil {
			g.error(err)
		}
	}
}

func (g *GoGenerator) Generate(outPath string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
		}
	}()

	// Generate package namespace mapping if necessary
	if g.Packages == nil {
		g.Packages = map[string]GoPackage{}
		g.packageNames = map[string]bool{}
	}
	for path, th := range g.ThriftFiles {
		if pkg, ok := g.Packages[path]; !ok || pkg.Name == "" {
			pkg := GoPackage{}
			for _, k := range goNamespaceOrder {
				pkg.Name = th.Namespaces[k]
				if pkg.Name != "" {
					parts := strings.Split(pkg.Name, ".")
					if len(parts) > 1 {
						pkg.Path = strings.Join(parts[:len(parts)-1], "/")
						pkg.Name = parts[len(parts)-1]
					}
					break
				}
			}
			if pkg.Name == "" {
				pkg.Name = filepath.Base(path)
			}
			pkg.Name = validIdentifier(strings.ToLower(pkg.Name), "_")
			g.Packages[path] = pkg
			g.packageNames[pkg.Name] = true
		}
	}

	rpcPackages := map[string]string{}

	for path, th := range g.ThriftFiles {
		pkg := g.Packages[path]
		filename := strings.ToLower(filepath.Base(path))
		for i := len(filename) - 1; i >= 0; i-- {
			if filename[i] == '.' {
				filename = filename[:i]
			}
		}
		if strings.HasSuffix(filename, "_test") {
			filename = filename[:len(filename)-len("_test")]
		}

		filename += ".go"
		pkgpath := filepath.Join(outPath, pkg.Path, pkg.Name)
		outfile := filepath.Join(pkgpath, filename)

		if err := os.MkdirAll(pkgpath, 0755); err != nil {
			g.error(err)
		}

		out := &bytes.Buffer{}
		g.generateSingle(out, path, th)

		if len(th.Services) > 0 {
			rpcPackages[pkgpath] = pkg.Name
		}

		outBytes := out.Bytes()
		if g.Format {
			outBytes, err = format.Source(outBytes)
			if err != nil {
				g.error(err)
			}
		}

		fi, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(2)
		}
		if _, err := fi.Write(outBytes); err != nil {
			fi.Close()
			g.error(err)
		}
		fi.Close()
	}

	for path, name := range rpcPackages {
		outfile := filepath.Join(path, "rpc_stub.go")

		fi, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(2)
		}
		_, err = fi.WriteString(fmt.Sprintf("package %s\n\ntype RPCClient interface {\n"+
			"\tCall(method string, request interface{}, response interface{}) error\n"+
			"}\n", name))
		fi.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(2)
		}
	}

	return nil
}
