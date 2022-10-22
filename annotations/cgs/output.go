package cgs

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"text/template"
)

var tpl = `
func {{ .FunctionName }}({{ range .Arguments }} {{.}}, {{ end }}) {{ if .IsPointer }}*{{ end }}{{ .ReturnType }} {
	return {{ if .IsPointer }}&{{ end }}{{ .ReturnType }} {
		{{ range .Fields }} {{ .Name }}: {{ .Value }},
		{{ end }}}
}
`

type TemplateValues struct {
	FunctionName string
	Arguments    []string
	ReturnType   string
	IsPointer    bool
	Fields       []struct {
		Name  string
		Value string
	}
}

func generateConstructors(generates []toGenerate) []byte {
	pkgName := generates[0].packageName
	var out []byte
	for _, generate := range generates {
		out = append(out, generateConstructor(generate)...)
	}

	return append([]byte(`package `+pkgName+`
`), out...)
}

func generateConstructor(generate toGenerate) []byte {
	tv := TemplateValues{
		FunctionName: generate.annotation.Name,
		IsPointer:    generate.annotation.Type == "pointer",
		ReturnType:   generate.node.Name.Name,
	}

	templ, err := template.New("tpl").Parse(tpl)
	if err != nil {
		panic(err)
	}

	a := args(generate.node)

	for name, tpy := range a {
		fmt.Println(name, tpy)
		tv.Arguments = append(tv.Arguments, name+" "+tpy)
		tv.Fields = append(tv.Fields, struct {
			Name  string
			Value string
		}{Name: name, Value: name})
	}

	buff := &bytes.Buffer{}
	err = templ.Execute(buff, tv)
	if err != nil {
		panic(err)
	}

	return buff.Bytes()
}

func args(n *ast.TypeSpec) map[string]string {
	out := map[string]string{}

	strTpy, ok := n.Type.(*ast.StructType)
	if !ok {
		panic("not a *ast.StructType")
	}

	if strTpy.Fields == nil {
		fmt.Println("no TypeParams")
		return out
	}
	fields := strTpy.Fields.List
	if len(fields) == 0 {
		fmt.Println("no TypeParams.List")
		return out
	}

	for _, field := range fields {
		buff := bytes.NewBufferString("")
		err := printer.Fprint(buff, &token.FileSet{}, field.Type)
		if err != nil {
			panic(err)
		}
		for _, ident := range field.Names {
			out[ident.Name] = buff.String()
		}
	}
	return out
}
