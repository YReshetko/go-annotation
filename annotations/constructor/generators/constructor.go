package generators

import (
	"bytes"
	"fmt"
	"github.com/YReshetko/go-annotation/annotations/constructor/annotations"
	"go/ast"
	"go/printer"
	"go/token"
	"strings"
	"text/template"

	annotation "github.com/YReshetko/go-annotation/pkg"
)

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

type arguments struct {
	incoming map[string]string
	toInit   map[string]string
}

type ConstructorGenerator struct {
	node          *ast.TypeSpec
	annotation    annotations.Constructor
	annotatedNode annotation.Node
}

func NewConstructorGenerator(node *ast.TypeSpec, annotation annotations.Constructor, an annotation.Node) *ConstructorGenerator {
	return &ConstructorGenerator{
		node:          node,
		annotation:    annotation,
		annotatedNode: an,
	}
}

func (g *ConstructorGenerator) Generate() ([]byte, []Import, error) {
	data, imports := g.generateConstructor()
	return data, imports.toSlice(), nil
}

func (g *ConstructorGenerator) generateConstructor() ([]byte, distinctImports) {
	tpl := must(template.New(functionNameTpl).Parse(g.annotation.Name))
	data := map[string]string{"TypeName": g.node.Name.Name}
	di := newDistinctImports()

	tv := ConstructorValues{
		FunctionName: string(must(executeTpl(tpl, data))),
		IsPointer:    g.annotation.Type == "pointer",
		ReturnType:   g.node.Name.Name,
	}

	a, adi := args(g.node, g.annotatedNode.FindImportByAlias, g.annotatedNode)
	for name, tpy := range a.incoming {
		if _, ok := a.toInit[name]; ok {
			continue
		}
		tv.Arguments = append(tv.Arguments, name+" "+tpy)
		tv.Fields = append(tv.Fields, struct {
			Name  string
			Value string
		}{Name: name, Value: name})
	}
	for name, tpy := range a.toInit {
		tv.Fields = append(tv.Fields, struct {
			Name  string
			Value string
		}{Name: name, Value: tpy})
	}

	di.merge(adi)

	c, p, pdi, ok := params(g.node, g.annotatedNode.FindImportByAlias)
	if ok {
		tv.IsParametrized = true
		tv.ParameterConstraints = c
		tv.Parameters = p
		di.merge(pdi)
	}

	return must(execute(constructorTpl, tv)), di
}

func args(n *ast.TypeSpec, fn func(string) (string, bool), node annotation.Node) (arguments, distinctImports) {
	out := arguments{
		incoming: map[string]string{},
		toInit:   map[string]string{},
	}
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
		toCheck := node.AnnotatedNode(field)
		if len(annotation.FindAnnotations[annotations.Exclude](toCheck.Annotations())) > 0 {
			continue
		}
		buff := bytes.NewBufferString("")
		err := printer.Fprint(buff, &token.FileSet{}, field.Type)
		if err != nil {
			panic(err)
		}
		initAnnotations := annotation.FindAnnotations[annotations.Init](toCheck.Annotations())
		if len(initAnnotations) == 0 {
			for _, ident := range field.Names {
				out.incoming[ident.Name] = buff.String()
			}
		} else {
			result := ""
			origVal := buff.String()
			switch field.Type.(type) {
			case *ast.ArrayType:
				result = sliceInitialisation(initAnnotations[0], origVal)
			case *ast.MapType:
				result = mapInitialisation(initAnnotations[0], origVal)
			case *ast.ChanType:
				result = chanInitialisation(initAnnotations[0], origVal)
			default:
				continue
			}
			for _, ident := range field.Names {
				out.toInit[ident.Name] = result
				out.incoming[ident.Name] = origVal
			}
		}

		imps.merge(getImports(field.Type, fn))

	}
	return out, imps
}

func sliceInitialisation(a annotations.Init, t string) string {
	switch {
	case a.Cap > -1 && a.Len > -1:
		return fmt.Sprintf("make(%s, %d, %d)", t, a.Len, a.Cap)
	case a.Len > -1:
		return fmt.Sprintf("make(%s, %d)", t, a.Len)
	case a.Cap > -1:
		return fmt.Sprintf("make(%s, 0, %d)", t, a.Cap)
	default:
		return fmt.Sprintf("%s{}", t)
	}
}

func mapInitialisation(a annotations.Init, t string) string {
	switch {
	case a.Cap > -1:
		return fmt.Sprintf("make(%s, %d)", t, a.Cap)
	default:
		return fmt.Sprintf("make(%s)", t)
	}
}

func chanInitialisation(a annotations.Init, t string) string {
	switch {
	case a.Cap > -1:
		return fmt.Sprintf("make(%s, %d)", t, a.Cap)
	default:
		return fmt.Sprintf("make(%s)", t)
	}
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
