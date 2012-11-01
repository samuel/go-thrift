package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/samuel/go-thrift"
	"github.com/samuel/go-thrift/parser"
)

var (
	goTemplate = `{{define "field"}}{{.Name|camelCase}} {{.Type|mapType}} ` + "`" + `thrift:"{{.Id}}{{if .Optional}}{{else}},required{{end}}" json:"{{.Name}}"` + "`" + `{{end}}` +
		`{{define "argumentList"}}{{range $i, $a := .}}{{if $i|first|not}}, {{end}}{{$a.Name|lowerCamelCase}} {{$a.Type|mapType}}{{end}}{{end}}` +
		`{{range $name, $enum := .Enums}}
type {{$name|camelCase}} int32

var ({{range $vname, $vid := .Values}}
	{{$name|camelCase}}{{$vname|camelCase}} = {{$name|camelCase}}({{$vid.Value}}){{end}}
)

func (e {{$name|camelCase}}) String() string {
	switch e {{"{"}}{{range $vname, $vid := .Values}}
		case {{$name|camelCase}}{{$vname|camelCase}}: return "{{$name|camelCase}}{{$vname|camelCase}}"{{end}}
	}
	return fmt.Sprintf("Unknown value for {{$name|camelCase}}: %d", e)
}
{{end}}{{range $name, $field := .Structs}}
type {{$name|camelCase}} struct {{"{"}}{{range .Fields}}
	{{template "field" .}}{{end}}
}
{{end}}{{range $name, $field := .Exceptions}}
type {{$name|camelCase}} struct {{"{"}}{{range .Fields}}
	{{template "field" .}}{{end}}
}

func (e *{{$name|camelCase}}) Error() string {
	return fmt.Sprintf("{{$name|camelCase}}{{"{"}}{{range $i, $f := .Fields}}{{if $i|first|not}}, {{end}}{{$f.Name|camelCase}}: %%+v{{end}}{{"}"}}"{{if .Fields}}, {{end}}{{range $i, $f := .Fields}}{{if $i|first|not}}, {{end}}e.{{$f.Name|camelCase}}{{end}})
}
{{end}}{{range $name, $svc := .Services}}
type {{$name|camelCase}} interface {{"{"}}{{range .Methods}}
	{{.Name|camelCase}}({{template "argumentList" .Fields}}) {{.ReturnType|returnType}}{{end}}
}{{range .Methods}}

type {{$name|camelCase}}{{.Name|camelCase}}Request struct {{"{"}}{{range .Fields}}
	{{template "field" .}}{{end}}
}

type {{$name|camelCase}}{{.Name|camelCase}}Response struct {{"{"}}{{if .ReturnType}}
	Value {{.ReturnType|mapType}} ` + "`" + `thrift:"0" json:"value"` + "`" + `{{end}}{{range .Exceptions}}
	{{.Name|camelCase}} {{.Type|mapType}} ` + "`" + `thrift:"{{.Id}}" json:"{{.Name}}"` + "`" + `{{end}}
}{{end}}
{{end}}{{if .Services}}
type RPCClient interface {
	Call(method string, request interface{}, response interface{}) error
}
{{end}}{{range $svc := .Services}}
type {{$svc.Name|camelCase}}Client struct {
	Client RPCClient
}{{range $svc.Methods}}

func (s *{{$svc.Name|camelCase}}Client) {{.Name|camelCase}}({{template "argumentList" .Fields}}) {{.ReturnType|returnType}} {
	req := &{{$svc.Name|camelCase}}{{.Name|camelCase}}Request {{"{"}}{{range .Fields}}
		{{.Name|camelCase}}: {{.Name|lowerCamelCase}},{{end}}
	}
	res := &{{$svc.Name|camelCase}}{{.Name|camelCase}}Response {}
	err := s.Client.Call("{{.Name}}", req, res){{if .Exceptions}}
	if err == nil {
		switch {{"{"}}{{range .Exceptions}}
		case res.{{.Name|camelCase}} != nil:
			err = res.{{.Name|camelCase}}{{end}}
		}
	}
{{end}}
	{{if .ReturnType}}return res.Value, err{{else}}return err{{end}}
}{{end}}
{{end}}`
)

func camelCase(st string) string {
	if strings.ToUpper(st) == st {
		st = strings.ToLower(st)
	}
	return thrift.CamelCase(st)
}

func lowerCamelCase(st string) string {
	// // Assume st is not unicode
	// if strings.ToUpper(st) == st {
	// 	return strings.ToLower(st)
	// }
	// st = thrift.CamelCase(st)
	// return strings.ToLower(st[:1]) + st[1:]
	return camelCase(st)
}

func mapType(def *parser.Thrift, typ *parser.Type) string {
	switch typ.Name {
	case "byte", "bool", "string":
		return typ.Name
	case "binary":
		return "[]byte"
	case "i16":
		return "int16"
	case "i32":
		return "int32"
	case "i64":
		return "int64"
	case "double":
		return "float64"
	case "list":
		return "[]" + mapType(def, typ.ValueType)
	case "map":
		keyType := mapType(def, typ.KeyType)
		if keyType == "[]byte" {
			// TODO: Log, warn, do something besides println!
			println("key type of []byte not supported for maps")
			keyType = "string"
		}
		return "map[" + keyType + "]" + mapType(def, typ.ValueType)
	}
	if e := def.Enums[typ.Name]; e != nil {
		return typ.Name
	}
	// TODO: References to types in includes
	return "*" + typ.Name
}

func returnType(def *parser.Thrift, typ *parser.Type) string {
	if typ == nil || typ.Name == "void" {
		return "error"
	}
	return fmt.Sprintf("(%s, error)", mapType(def, typ))
}

func main() {
	filename := os.Args[1]

	p := &parser.Parser{}
	th, err := p.ParseFile(filename)
	if e, ok := err.(*parser.ErrSyntaxError); ok {
		fmt.Printf("%s\n", e.Left)
		panic(err)
	} else if err != nil {
		panic(err)
	}

	out := &bytes.Buffer{}

	// Package docs and package name

	packageName := th.Namespaces["go"]
	if packageName == "" {
		packageName = th.Namespaces["perl"]
		if packageName == "" {
			packageName = th.Namespaces["py"]
			if packageName == "" {
				packageName = strings.Split(filename, ".")[0]
			} else {
				parts := strings.Split(packageName, ".")
				packageName = parts[len(parts)-1]
			}
		}
	}
	packageName = strings.ToLower(packageName)

	out.WriteString("// This file is automatically generated. Do not modify.\n\n")
	out.WriteString("package " + packageName + "\n")

	// Imports

	imports := []string{"fmt"}
	out.WriteString("\nimport (\n")
	for _, in := range imports {
		out.WriteString("\t\"" + in + "\"\n")
	}
	out.WriteString(")\n")

	funcMap := map[string]interface{}{
		"camelCase":      camelCase,
		"lowerCamelCase": lowerCamelCase,
		"mapType":        func(typ *parser.Type) string { return mapType(th, typ) },
		"returnType":     func(typ *parser.Type) string { return returnType(th, typ) },
		"first":          func(i int) bool { return i == 1 },
	}

	tmpl := template.New("go")
	tmpl.Funcs(funcMap)
	_, err = tmpl.Parse(goTemplate)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(out, th)
	if err != nil {
		panic(err)
	}

	fmt.Printf(out.String())
}
