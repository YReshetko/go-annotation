package generators

import (
	"github.com/YReshetko/go-annotation/annotations/constructor/annotations"
	annotation "github.com/YReshetko/go-annotation/pkg"
	"go/ast"
	"text/template"
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
		OptionalTypeName: g.optionalTypeName(),
		FunctionName:     g.optionalFunctionName(),
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

	data := must(execute(optionalTypeTpl, tplData))
	data = append(data, must(execute(optionalConstructorTpl, tplData))...)

	for fieldName, fieldValue := range argsToProcess.incoming {
		data = append(data, g.withFunction(fieldName, fieldValue, tplData)...)
	}

	return data, di.toSlice(), nil
}

func (g *OptionalGenerator) generateConstructor() ([]byte, distinctImports) {
	return nil, newDistinctImports()
}

func (g *OptionalGenerator) optionalTypeName() string {
	tpl := must(template.New(typeNameTpl).Parse(g.annotation.Name))
	typeNameData := map[string]string{"TypeName": g.node.Name.Name}

	return string(must(executeTpl(tpl, typeNameData)))
}

func (g *OptionalGenerator) optionalFunctionName() string {
	tpl := must(template.New(functionNameTpl).Parse(g.annotation.ConstructorName))
	typeNameData := map[string]string{"TypeName": g.node.Name.Name}

	return string(must(executeTpl(tpl, typeNameData)))
}

func (g *OptionalGenerator) withFunctionName(fieldName string) string {
	tpl := must(template.New(functionNameTpl).Parse(g.annotation.WithPattern))
	typeNameData := map[string]string{"FieldName": toPascalCase(fieldName)}

	return string(must(executeTpl(tpl, typeNameData)))
}

func (g *OptionalGenerator) withFunction(name string, value string, data OptionalValues) []byte {
	wv := WithValues{
		OptionalValues:   data,
		ArgumentType:     value,
		FieldName:        name,
		WithFunctionName: g.withFunctionName(name),
	}

	return must(execute(optionalWithTpl, wv))
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
