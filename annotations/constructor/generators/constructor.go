package generators

import (
	"go/ast"
	"sort"
	"strings"

	"github.com/YReshetko/go-annotation/annotations/constructor/annotations"
	"github.com/YReshetko/go-annotation/annotations/constructor/templates"
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
	PostConstructs []PostConstructValues
	ReturnsError   bool
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
	di := newDistinctImports()
	tv := ConstructorValues{
		FunctionName:   g.annotation.BuildName(g.node.Name.String()),
		IsPointer:      g.annotation.Type == "pointer",
		ReturnType:     g.node.Name.Name,
		PostConstructs: sortPostConstructs(pcvs),
		ReturnsError:   hasAnyErrorReturned(pcvs),
	}

	lookup := g.annotatedNode.Lookup().FindImportByAlias
	a, adi := extractArguments(g.node, lookup, g.annotatedNode)
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

	sort.Slice(tv.Arguments, func(i, j int) bool {
		is := tv.Arguments[i]
		js := tv.Arguments[j]
		in := is[:strings.Index(is, " ")]
		jn := js[:strings.Index(js, " ")]
		return in < jn
	})

	sort.Slice(tv.Fields, func(i, j int) bool {
		return tv.Fields[i].Name < tv.Fields[j].Name
	})

	di.merge(adi)

	p, pdi := extractParameters(g.node, lookup)
	if p.isParametrised {
		tv.IsParametrized = true
		tv.ParameterConstraints = p.constraints
		tv.Parameters = p.parameters
		di.merge(pdi)
	}

	return templates.Must(templates.Execute(templates.ConstructorTpl, tv)), di
}

func (g *ConstructorGenerator) Name() string {
	return "ConstructorGenerator"
}
