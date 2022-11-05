package generators

import (
	"github.com/YReshetko/go-annotation/annotations/constructor/annotations"
	"go/ast"
	"text/template"

	annotation "github.com/YReshetko/go-annotation/pkg"
)

type ConstructorValues struct {
	FunctionName         string
	Arguments            []string
	ReturnType           string
	IsPointer            bool
	IsParametrized       bool
	ParameterConstraints string
	Parameters           string
	Fields               []struct {
		Name  string
		Value string
	}
	PostConstructs []string
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

func (g *ConstructorGenerator) Generate(pcvs []PostConstructValues) ([]byte, []Import, error) {
	data, imports := g.generateConstructor(pcvs)
	return data, imports.toSlice(), nil
}

func (g *ConstructorGenerator) generateConstructor(pcvs []PostConstructValues) ([]byte, distinctImports) {
	tpl := must(template.New(functionNameTpl).Parse(g.annotation.Name))
	data := map[string]string{"TypeName": g.node.Name.Name}
	di := newDistinctImports()

	tv := ConstructorValues{
		FunctionName:   string(must(executeTpl(tpl, data))),
		IsPointer:      g.annotation.Type == "pointer",
		ReturnType:     g.node.Name.Name,
		PostConstructs: postConstructMethods(pcvs),
	}

	a, adi := extractArguments(g.node, g.annotatedNode.FindImportByAlias, g.annotatedNode)
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

	p, pdi := extractParameters(g.node, g.annotatedNode.FindImportByAlias)
	if p.isParametrised {
		tv.IsParametrized = true
		tv.ParameterConstraints = p.constraints
		tv.Parameters = p.parameters
		di.merge(pdi)
	}

	return must(execute(constructorTpl, tv)), di
}

func (g *ConstructorGenerator) Name() string {
	return "ConstructorGenerator"
}
