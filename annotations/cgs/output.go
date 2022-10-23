package cgs

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"strings"
	"text/template"
)

type FileValues struct {
	PackageName string
	HasImports  bool
	Imports     []Import
	Data        string
}

type ConstructorValues struct {
	FunctionName   string
	Arguments      []string
	ReturnType     string
	IsPointer      bool
	IsParametrized bool
	Fields         []struct {
		Name  string
		Value string
	}
	ParameterConstraints string
	Parameters           string
}

func generateConstructors(generates []toGenerate) []byte {
	pkgName := generates[0].packageName
	di := newDistinctImports()
	var out []byte
	for _, generate := range generates {
		data, imps := generateConstructor(generate)
		out = append(out, data...)
		di.merge(imps)
	}

	fv := FileValues{
		PackageName: pkgName,
		HasImports:  !di.isEmpty(),
		Imports:     di.toSlice(),
		Data:        string(out),
	}

	return must(execute(file, fv))
}

func generateConstructor(generate toGenerate) ([]byte, distinctImports) {
	tpl := must(template.New(functionName).Parse(generate.annotation.Name))
	data := map[string]string{"TypeName": generate.node.Name.Name}
	di := newDistinctImports()

	tv := ConstructorValues{
		FunctionName: string(must(executeTpl(tpl, data))),
		IsPointer:    generate.annotation.Type == pointerType,
		ReturnType:   generate.node.Name.Name,
	}

	a, adi := args(generate.node, generate.an.FindImportByAlias)

	for name, tpy := range a {
		//fmt.Println(name, tpy)
		tv.Arguments = append(tv.Arguments, name+" "+tpy)
		tv.Fields = append(tv.Fields, struct {
			Name  string
			Value string
		}{Name: name, Value: name})
	}
	di.merge(adi)

	c, p, pdi, ok := params(generate.node, generate.an.FindImportByAlias)
	if ok {
		tv.IsParametrized = true
		tv.ParameterConstraints = c
		tv.Parameters = p
		di.merge(pdi)
	}

	return must(execute(constructor, tv)), di
}

func args(n *ast.TypeSpec, fn func(string) (string, bool)) (map[string]string, distinctImports) {
	out := map[string]string{}
	imps := newDistinctImports()

	strTpy, ok := n.Type.(*ast.StructType)
	if !ok {
		panic("not a *ast.StructType")
	}

	if strTpy.Fields == nil {
		fmt.Println("no TypeParams")
		return out, imps
	}
	fields := strTpy.Fields.List
	if len(fields) == 0 {
		fmt.Println("no TypeParams.List")
		return out, imps
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
		imps.merge(getImports(field.Type, fn))

	}
	return out, imps
}

func params(n *ast.TypeSpec, fn func(string) (string, bool)) (string, string, distinctImports, bool) {
	if n.TypeParams == nil || len(n.TypeParams.List) == 0 {
		return "", "", nil, false
	}
	var c []string
	var p []string
	imps := newDistinctImports()
	for _, field := range n.TypeParams.List {
		buff := bytes.NewBufferString("")
		err := printer.Fprint(buff, &token.FileSet{}, field.Type)
		if err != nil {
			panic(err)
		}
		for _, name := range field.Names {
			if name == nil || len(name.Name) == 0 {
				continue
			}
			c = append(c, name.Name+" "+buff.String())
			p = append(p, name.Name)
		}
		imps.merge(getImports(field.Type, fn))
	}

	return strings.Join(c, ","), strings.Join(p, ","), imps, true
}
