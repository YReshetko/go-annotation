package generators

import (
	"fmt"
	"go/ast"
	"sort"

	"github.com/YReshetko/go-annotation/annotations/constructor/annotations"
	"github.com/YReshetko/go-annotation/annotations/constructor/templates"
	annotation "github.com/YReshetko/go-annotation/pkg"
)

type BuilderValues struct {
	BuilderTypeName      string
	ConstructorName      string
	BuildMethodName      string
	ReturnType           string
	IsParametrized       bool
	IsPointer            bool
	ParameterConstraints string
	Parameters           string
	Arguments            []struct {
		FakeName string
		Name     string
		Type     string
	}
	Fields []struct {
		FakeName string
		Name     string
		Value    string
	}
	PostConstructs []PostConstructValues
	ReturnsError   bool
}

type BuildValues struct {
	BuilderValues
	BuilderMethodName string
	ArgumentType      string
	FakeName          string
}

type BuilderGenerator struct {
	node          *ast.TypeSpec
	annotation    annotations.Builder
	annotatedNode annotation.Node
}

func NewBuilderGenerator(node *ast.TypeSpec, annotation annotations.Builder, an annotation.Node) *BuilderGenerator {
	return &BuilderGenerator{
		node:          node,
		annotation:    annotation,
		annotatedNode: an,
	}
}

func (g *BuilderGenerator) Generate(pcvs []PostConstructValues) ([]byte, []Import, error) {
	di := newDistinctImports()
	tplData := BuilderValues{
		BuilderTypeName: g.annotation.BuildStructureName(g.node.Name.String()),
		ConstructorName: g.annotation.BuildConstructorName(g.node.Name.String()),
		BuildMethodName: g.annotation.BuilderName,
		ReturnType:      g.node.Name.Name,
		IsPointer:       g.annotation.Type == "pointer",
		PostConstructs:  sortPostConstructs(pcvs),
		ReturnsError:    hasAnyErrorReturned(pcvs),
	}

	lookup := g.annotatedNode.Lookup().FindImportByAlias
	p, pdi := extractParameters(g.node, lookup)
	if p.isParametrised {
		tplData.IsParametrized = true
		tplData.ParameterConstraints = p.constraints
		tplData.Parameters = p.parameters
		di.merge(pdi)
	}

	argsToProcess, adi := extractArguments(g.node, lookup, g.annotatedNode)
	di.merge(adi)

	index := 0
	for fName, fType := range argsToProcess.incoming {
		index++
		fakeName := fmt.Sprintf("_%s_", fName)
		tplData.Arguments = append(tplData.Arguments, struct {
			FakeName string
			Name     string
			Type     string
		}{FakeName: fakeName, Name: fName, Type: fType})
		if fValue, ok := argsToProcess.toInit[fName]; ok {
			tplData.Fields = append(tplData.Fields, struct {
				FakeName string
				Name     string
				Value    string
			}{FakeName: fakeName, Name: fName, Value: fValue})
		}
	}

	sort.Slice(tplData.Arguments, func(i, j int) bool {
		return tplData.Arguments[i].Name < tplData.Arguments[j].Name
	})
	sort.Slice(tplData.Fields, func(i, j int) bool {
		return tplData.Fields[i].Name < tplData.Fields[j].Name
	})

	data := templates.Must(templates.Execute(templates.BuilderTypeTpl, tplData))
	data = append(data, templates.Must(templates.Execute(templates.BuilderConstructorTpl, tplData))...)
	for _, argument := range tplData.Arguments {
		data = append(data, g.buildMethod(argument.FakeName, argument.Name, argument.Type, tplData)...)
	}
	data = append(data, templates.Must(templates.Execute(templates.BuilderBuildMethodTpl, tplData))...)

	return data, di.toSlice(), nil
}

func (g *BuilderGenerator) buildMethod(fakeName, name, value string, data BuilderValues) []byte {
	wv := BuildValues{
		BuilderValues:     data,
		ArgumentType:      value,
		FakeName:          fakeName,
		BuilderMethodName: g.annotation.BuildBuildName(name),
	}

	return templates.Must(templates.Execute(templates.BuilderMethodTpl, wv))
}

func (g *BuilderGenerator) Name() string {
	return "BuilderGenerator"
}
