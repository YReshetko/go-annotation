package generators

import (
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
	Fields               []struct {
		Name  string
		Value string
	}
}

type BuildValues struct {
	BuilderValues
	BuilderMethodName string
	ArgumentType      string
	FieldName         string
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

func (g *BuilderGenerator) Generate() ([]byte, []Import, error) {
	di := newDistinctImports()
	tplData := BuilderValues{
		BuilderTypeName: g.builderTypeName(),
		ConstructorName: g.builderConstructorName(),
		BuildMethodName: g.annotation.BuilderName,
		ReturnType:      g.node.Name.Name,
		IsPointer:       g.annotation.Type == "pointer",
	}

	c, p, pdi, ok := params(g.node, g.annotatedNode.FindImportByAlias)
	if ok {
		tplData.IsParametrized = true
		tplData.ParameterConstraints = c
		tplData.Parameters = p
		di.merge(pdi)
	}

	argsToProcess, adi := args(g.node, g.annotatedNode.FindImportByAlias, g.annotatedNode)
	di.merge(adi)

	for name, value := range argsToProcess.toInit {
		tplData.Fields = append(tplData.Fields, struct {
			Name  string
			Value string
		}{Name: name, Value: value})
	}

	data := must(execute(builderTypeTpl, tplData))
	data = append(data, must(execute(builderConstructorTpl, tplData))...)
	for fieldName, fieldValue := range argsToProcess.incoming {
		data = append(data, g.buildMethod(fieldName, fieldValue, tplData)...)
	}
	data = append(data, must(execute(builderBuildMethodTpl, tplData))...)

	return data, di.toSlice(), nil
}

func (g *BuilderGenerator) buildMethod(name string, value string, data BuilderValues) []byte {
	wv := BuildValues{
		BuilderValues:     data,
		ArgumentType:      value,
		FieldName:         name,
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
