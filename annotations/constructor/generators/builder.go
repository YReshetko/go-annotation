package generators

import (
	"fmt"
	"go/ast"
	"text/template"

	"github.com/YReshetko/go-annotation/annotations/constructor/annotations"
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
	PostConstructs []string
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
		BuilderTypeName: g.builderTypeName(),
		ConstructorName: g.builderConstructorName(),
		BuildMethodName: g.annotation.BuilderName,
		ReturnType:      g.node.Name.Name,
		IsPointer:       g.annotation.Type == "pointer",
		PostConstructs:  postConstructMethods(pcvs),
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

	index := 0
	for fName, fType := range argsToProcess.incoming {
		index++
		fakeName := fmt.Sprintf("_field_%d_", index)
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

	data := must(execute(builderTypeTpl, tplData))
	data = append(data, must(execute(builderConstructorTpl, tplData))...)
	for _, argument := range tplData.Arguments {
		data = append(data, g.buildMethod(argument.FakeName, argument.Name, argument.Type, tplData)...)
	}
	data = append(data, must(execute(builderBuildMethodTpl, tplData))...)

	return data, di.toSlice(), nil
}

func (g *BuilderGenerator) buildMethod(fakeName, name, value string, data BuilderValues) []byte {
	wv := BuildValues{
		BuilderValues:     data,
		ArgumentType:      value,
		FakeName:          fakeName,
		BuilderMethodName: g.builderMethodName(name),
	}

	return must(execute(builderMethodTpl, wv))
}

func (g *BuilderGenerator) builderTypeName() string {
	tpl := must(template.New(typeNameTpl).Parse(g.annotation.StructureName))
	typeNameData := map[string]string{"TypeName": g.node.Name.Name}

	return string(must(executeTpl(tpl, typeNameData)))
}

func (g *BuilderGenerator) builderConstructorName() string {
	tpl := must(template.New(typeNameTpl).Parse(g.annotation.ConstructorName))
	typeNameData := map[string]string{"TypeName": g.node.Name.Name}

	return string(must(executeTpl(tpl, typeNameData)))
}

func (g *BuilderGenerator) builderMethodName(fieldName string) string {
	tpl := must(template.New(functionNameTpl).Parse(g.annotation.BuildPattern))
	methodNameData := map[string]string{"FieldName": toPascalCase(fieldName)}

	return string(must(executeTpl(tpl, methodNameData)))
}

func (g *BuilderGenerator) Name() string {
	return "BuilderGenerator"
}
