package generators

import (
	"github.com/YReshetko/go-annotation/annotations/constructor/annotations"
	"github.com/YReshetko/go-annotation/annotations/constructor/templates"
	annotation "github.com/YReshetko/go-annotation/pkg"
	"go/ast"
	"sort"
)

type OptionalValues struct {
	OptionalTypeName     string
	FunctionName         string
	ReturnType           string
	IsParametrized       bool
	IsPointer            bool
	ParameterConstraints string
	Parameters           string
	Fields               []struct {
		Name  string
		Value string
	}
	PostConstructs []string
}

type WithValues struct {
	OptionalValues
	WithFunctionName string
	ArgumentType     string
	FieldName        string
}

type OptionalGenerator struct {
	node          *ast.TypeSpec
	annotation    annotations.Optional
	annotatedNode annotation.Node
}

func NewOptionalGenerator(node *ast.TypeSpec, annotation annotations.Optional, an annotation.Node) *OptionalGenerator {
	return &OptionalGenerator{
		node:          node,
		annotation:    annotation,
		annotatedNode: an,
	}
}

func (g *OptionalGenerator) Generate(pcvs []PostConstructValues) ([]byte, []Import, error) {
	di := newDistinctImports()
	tplData := OptionalValues{
		OptionalTypeName: g.annotation.BuildName(g.node.Name.String()),
		FunctionName:     g.annotation.BuildConstructorName(g.node.Name.String()),
		ReturnType:       g.node.Name.Name,
		IsPointer:        g.annotation.Type == "pointer",
		PostConstructs:   postConstructMethods(pcvs),
	}
	p, pdi := extractParameters(g.node, g.annotatedNode.FindImportByAlias)
	if p.isParametrised {
		tplData.IsParametrized = true
		tplData.ParameterConstraints = p.constraints
		tplData.Parameters = p.parameters
		di.merge(pdi)
	}

	argsToProcess, adi := extractArguments(g.node, g.annotatedNode.FindImportByAlias, g.annotatedNode)
	di.merge(adi)

	for name, value := range argsToProcess.toInit {
		tplData.Fields = append(tplData.Fields, struct {
			Name  string
			Value string
		}{Name: name, Value: value})
	}

	sort.Slice(tplData.Fields, func(i, j int) bool {
		return tplData.Fields[i].Name < tplData.Fields[j].Name
	})
	var incoming [][2]string
	for fn, fv := range argsToProcess.incoming {
		incoming = append(incoming, [2]string{fn, fv})
	}
	sort.Slice(incoming, func(i, j int) bool {
		return incoming[i][0] < incoming[j][0]
	})

	data := templates.Must(templates.Execute(templates.OptionalTypeTpl, tplData))
	data = append(data, templates.Must(templates.Execute(templates.OptionalConstructorTpl, tplData))...)
	for _, i := range incoming {
		data = append(data, g.withFunction(i[0], i[1], tplData)...)
	}

	return data, di.toSlice(), nil
}

func (g *OptionalGenerator) generateConstructor() ([]byte, distinctImports) {
	return nil, newDistinctImports()
}

func (g *OptionalGenerator) withFunction(name string, value string, data OptionalValues) []byte {
	wv := WithValues{
		OptionalValues:   data,
		ArgumentType:     value,
		FieldName:        name,
		WithFunctionName: g.annotation.BuildWithName(name),
	}

	return templates.Must(templates.Execute(templates.OptionalWithTpl, wv))
}

func toPascalCase(name string) string {
	if name[0] >= 'a' && name[0] <= 'z' {
		ch := name[0] - 'a' + 'A'
		return string(ch) + name[1:]
	}

	return name
}

func (g *OptionalGenerator) Name() string {
	return "OptionalGenerator"
}
