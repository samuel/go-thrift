package main

// TODO:
// - Default arguments. Possibly don't bother...

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/samuel/go-thrift/parser"
)

var (
	f_go_binarystring = flag.Bool("go.binarystring", false, "Always use string for binary instead of []byte")
	f_go_packagename  = flag.String("go.packagename", "", "Override the package name")
	f_go_pointers     = flag.Bool("go.pointers", false, "Make all fields pointers")
)

type GoGenerator struct {
	Thrift *parser.Thrift
}

func (g *GoGenerator) FormatType(typ *parser.Type) string {
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
		valueType := g.FormatType(typ.ValueType)
		if valueType == "[]byte" {
			valueType = "string"
		}
		return "map[" + valueType + "]interface{}"
	case "list":
		return "[]" + g.FormatType(typ.ValueType)
	case "map":
		keyType := g.FormatType(typ.KeyType)
		if keyType == "[]byte" {
			// TODO: Log, warn, do something!
			// println("key type of []byte not supported for maps")
			keyType = "string"
		}
		return "map[" + keyType + "]" + g.FormatType(typ.ValueType)
	}

	if t := g.Thrift.Typedefs[typ.Name]; t != nil {
		return g.FormatType(t)
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

	panic("Unknown type " + typ.Name)
}

func (g *GoGenerator) FormatField(field *parser.Field) string {
	tags := ""
	if !field.Optional {
		tags = ",required"
	}
	return fmt.Sprintf(
		"%s %s `thrift:\"%d%s\" json:\"%s\"`",
		camelCase(field.Name), g.FormatType(field.Type), field.Id, tags, field.Name)
}

func (g *GoGenerator) FormatArguments(arguments []*parser.Field) string {
	args := make([]string, len(arguments))
	for i, arg := range arguments {
		args[i] = fmt.Sprintf("%s %s", camelCase(arg.Name), g.FormatType(arg.Type))
	}
	return strings.Join(args, ", ")
}

func (g *GoGenerator) FormatReturnType(typ *parser.Type) string {
	if typ == nil || typ.Name == "void" {
		return "error"
	}
	return fmt.Sprintf("(%s, error)", g.FormatType(typ))
}

func (g *GoGenerator) WriteConstant(out io.Writer, c *parser.Constant) error {
	if _, err := io.WriteString(out,
		fmt.Sprintf("\nconst %s = %+v\n",
			camelCase(c.Name), c.Value)); err != nil {
		return err
	}
	return nil
}

func (g *GoGenerator) WriteEnum(out io.Writer, enum *parser.Enum) error {
	enumName := camelCase(enum.Name)

	if _, err := io.WriteString(out, "\ntype "+enumName+" int32\n"); err != nil {
		return err
	}

	if _, err := io.WriteString(out, "\nvar (\n"); err != nil {
		return err
	}
	for name, val := range enum.Values {
		if _, err := io.WriteString(out,
			fmt.Sprintf(
				"\t%s%s = %s(%d)\n", enumName,
				camelCase(name), enumName, val.Value)); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(out, ")\n"); err != nil {
		return err
	}

	if _, err := io.WriteString(out, "\nfunc (e "+enumName+") String() string {\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "\tswitch e {\n"); err != nil {
		return err
	}
	for name := range enum.Values {
		name = camelCase(name)
		if _, err := io.WriteString(out,
			fmt.Sprintf(
				"\tcase %s%s:\n\t\treturn \"%s%s\"\n",
				enumName, name, enumName, name)); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(out, "\t}\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "\treturn fmt.Sprintf(\"Unknown value for "+enumName+": %d\", e)\n"); err != nil {
		return err
	}
	_, err := io.WriteString(out, "}\n")
	return err
}

func (g *GoGenerator) WriteStruct(out io.Writer, st *parser.Struct) error {
	structName := camelCase(st.Name)

	if _, err := io.WriteString(out, "\ntype "+structName+" struct {\n"); err != nil {
		return err
	}
	for _, field := range st.Fields {
		if _, err := io.WriteString(out, "\t"+g.FormatField(field)+"\n"); err != nil {
			return err
		}
	}
	_, err := io.WriteString(out, "}\n")
	return err
}

func (g *GoGenerator) WriteException(out io.Writer, ex *parser.Struct) error {
	if err := g.WriteStruct(out, ex); err != nil {
		return err
	}

	exName := camelCase(ex.Name)

	if _, err := io.WriteString(out, "\nfunc (e *"+exName+") Error() string {\n"); err != nil {
		return err
	}
	if len(ex.Fields) == 0 {
		if _, err := io.WriteString(out, "\treturn \""+exName+"{}\"\n"); err != nil {
			return err
		}
	} else {
		fieldNames := make([]string, len(ex.Fields))
		fieldVars := make([]string, len(ex.Fields))
		for i, field := range ex.Fields {
			fieldNames[i] = camelCase(field.Name) + ": %+v"
			fieldVars[i] = "e." + camelCase(field.Name)
		}
		if _, err := io.WriteString(out,
			fmt.Sprintf(
				"\treturn fmt.Sprintf(\"%s{%s}\", %s)\n",
				exName, strings.Join(fieldNames, ", "), strings.Join(fieldVars, ", "))); err != nil {
			return err
		}
	}
	_, err := io.WriteString(out, "}\n")
	return err
}

func (g *GoGenerator) WriteService(out io.Writer, svc *parser.Service) error {
	svcName := camelCase(svc.Name)

	if _, err := io.WriteString(out, "\ntype "+svcName+" interface {\n"); err != nil {
		return err
	}
	for _, method := range svc.Methods {
		if _, err := io.WriteString(out,
			fmt.Sprintf(
				"\t%s(%s) %s\n",
				camelCase(method.Name), g.FormatArguments(method.Arguments),
				g.FormatReturnType(method.ReturnType))); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(out, "}\n"); err != nil {
		return err
	}

	for _, method := range svc.Methods {
		if err := g.WriteStruct(out, &parser.Struct{svcName + camelCase(method.Name) + "Request", method.Arguments}); err != nil {
			return err
		}

		args := make([]*parser.Field, 0, len(method.Exceptions))
		if method.ReturnType != nil && method.ReturnType.Name != "void" {
			args = append(args, &parser.Field{0, "value", len(method.Exceptions) != 0, method.ReturnType, nil})
		}
		for _, ex := range method.Exceptions {
			args = append(args, ex)
		}

		res := &parser.Struct{svcName + camelCase(method.Name) + "Response", args}
		if err := g.WriteStruct(out, res); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(out, "\ntype "+svcName+"Client struct {\n\tClient RPCClient\n}\n"); err != nil {
		return err
	}

	for _, method := range svc.Methods {
		methodName := camelCase(method.Name)
		if _, err := io.WriteString(out,
			fmt.Sprintf("\nfunc (s *%sClient) %s(%s) %s {\n",
				svcName, methodName,
				g.FormatArguments(method.Arguments),
				g.FormatReturnType(method.ReturnType))); err != nil {
			return err
		}

		// Request
		if _, err := io.WriteString(out, fmt.Sprintf("\treq := &%s%sRequest{\n", svcName, methodName)); err != nil {
			return err
		}
		for _, arg := range method.Arguments {
			argName := camelCase(arg.Name)
			if _, err := io.WriteString(out,
				fmt.Sprintf("\t\t%s: %s,\n",
					argName, argName)); err != nil {
				return err
			}
		}
		if _, err := io.WriteString(out, "\t}\n"); err != nil {
			return err
		}

		// Response
		if _, err := io.WriteString(out, fmt.Sprintf("\tres := &%s%sResponse{}\n", svcName, methodName)); err != nil {
			return err
		}

		// Call
		if _, err := io.WriteString(out, "\terr := s.Client.Call(\""+method.Name+"\", req, res)\n"); err != nil {
			return err
		}

		// Exceptions
		if len(method.Exceptions) > 0 {
			if _, err := io.WriteString(out, "\tif err == nil {\n\t\tswitch {\n"); err != nil {
				return err
			}
			for _, ex := range method.Exceptions {
				exName := camelCase(ex.Name)
				if _, err := io.WriteString(out,
					fmt.Sprintf("\t\tcase res.%s != nil:\n\t\t\terr = res.%s\n",
						exName, exName)); err != nil {
					return err
				}
			}
			if _, err := io.WriteString(out, "\t\t}\n\t}\n"); err != nil {
				return err
			}
		}

		if method.ReturnType != nil && method.ReturnType.Name != "void" {
			if _, err := io.WriteString(out, "\treturn res.Value, err\n"); err != nil {
				return err
			}
		} else {
			if _, err := io.WriteString(out, "\treturn err\n"); err != nil {
				return err
			}
		}
		if _, err := io.WriteString(out, "}\n"); err != nil {
			return err
		}
	}

	return nil
}

func (g *GoGenerator) Generate(name string, out io.Writer) error {
	packageName := *f_go_packagename
	if packageName == "" {
		packageName = g.Thrift.Namespaces["go"]
		if packageName == "" {
			packageName = g.Thrift.Namespaces["perl"]
			if packageName == "" {
				packageName = g.Thrift.Namespaces["py"]
				if packageName == "" {
					packageName = name
				} else {
					parts := strings.Split(packageName, ".")
					packageName = parts[len(parts)-1]
				}
			}
		}
	}
	packageName = strings.ToLower(packageName)

	if _, err := io.WriteString(out, "// This file is automatically generated. Do not modify.\n"); err != nil {
		return err
	}
	if _, err := io.WriteString(out, "\npackage "+packageName+"\n"); err != nil {
		return err
	}

	// Imports
	imports := []string{"fmt"}
	if _, err := io.WriteString(out, "\nimport (\n"); err != nil {
		return err
	}
	for _, in := range imports {
		if _, err := io.WriteString(out, "\t\""+in+"\"\n"); err != nil {
			return err
		}
	}
	if _, err := io.WriteString(out, "\n)\n"); err != nil {
		return err
	}

	//

	for _, c := range g.Thrift.Constants {
		if err := g.WriteConstant(out, c); err != nil {
			return err
		}
	}

	for _, enum := range g.Thrift.Enums {
		if err := g.WriteEnum(out, enum); err != nil {
			return err
		}
	}

	for _, st := range g.Thrift.Structs {
		if err := g.WriteStruct(out, st); err != nil {
			return err
		}
	}

	for _, ex := range g.Thrift.Exceptions {
		if err := g.WriteException(out, ex); err != nil {
			return err
		}
	}

	if len(g.Thrift.Services) > 0 {
		if _, err := io.WriteString(out, "\ntype RPCClient interface {\n"+
			"\tCall(method string, request interface{}, response interface{}) error\n"+
			"}\n"); err != nil {
			return err
		}
	}

	for _, svc := range g.Thrift.Services {
		if err := g.WriteService(out, svc); err != nil {
			return err
		}
	}

	return nil
}
