package generators

import (
	"bytes"
	"fmt"
	cache2 "github.com/YReshetko/go-annotation/annotations/mapper/generators/cache"
	"github.com/YReshetko/go-annotation/annotations/mapper/generators/nodes"
	"github.com/YReshetko/go-annotation/annotations/mapper/templates"
	"github.com/YReshetko/go-annotation/annotations/mapper/utils"
	"go/ast"
	"strings"

	annotation "github.com/YReshetko/go-annotation/pkg"
)

// @Builder(constructor="newMethodGeneratorBuilder", build="set{{ .FieldName }}", terminator="build", type="pointer", exported="false")
type methodGenerator struct {
	impCache    *cache2.ImportCache
	node        annotation.Node
	name        string
	overloading *overloading
	input       []*ast.Field
	output      []*ast.Field

	inputFields  []*nodes.Field // @Init
	outputFields []*nodes.Field // @Init
}

// @PostConstruct
func (mg *methodGenerator) buildOutput() {
	for i, field := range mg.output {
		builder := nodes.NewFieldBuilder().ImpCache(mg.impCache).Node(mg.node).AstExpr(field.Type)
		if len(field.Names) == 0 {
			mg.outputFields = append(mg.outputFields, builder.Name(fmt.Sprintf("out%d", i)).Build())
			//mg.outputGenerator = append(mg.outputGenerator, buildRootFieldGenerator(field, "", mg.node))
		} else {
			for _, name := range field.Names {
				mg.outputFields = append(mg.outputFields, builder.Name(name.String()).Build())
				//mg.outputGenerator = append(mg.outputGenerator, buildRootFieldGenerator(field, name.String(), mg.node))
			}
		}
	}
}

// @PostConstruct
func (mg *methodGenerator) buildInput() {
	for i, field := range mg.input {
		builder := nodes.NewFieldBuilder().ImpCache(mg.impCache).Node(mg.node).AstExpr(field.Type)
		if len(field.Names) == 0 {
			mg.inputFields = append(mg.inputFields, builder.Name(fmt.Sprintf("in%d", i)).Build())
			//mg.inputGenerators = append(mg.inputGenerators, buildRootFieldGenerator(field, fmt.Sprintf("in%d", i), mg.node))
		} else {
			for _, name := range field.Names {
				mg.inputFields = append(mg.inputFields, builder.Name(name.String()).Build())
				//mg.inputGenerators = append(mg.inputGenerators, buildRootFieldGenerator(field, name.String(), mg.node))
			}
		}
	}
}

func (mg *methodGenerator) generate(receiverName string) ([]byte, error) {
	var inArgs []string
	for _, field := range mg.inputFields {
		inArgs = append(inArgs, field.Name()+" "+field.DeclaredType())
	}
	var outArgs []string
	var retVars []string
	var data []byte
	for _, field := range mg.outputFields {
		outArgs = append(outArgs, field.DeclaredType())
		retVars = append(retVars, field.Name())
		varDeclaration, err := templates.Execute(templates.NewVariableTemplate, map[string]interface{}{
			"VariableName": field.Name(),
			"Declaration":  field.Declaration(),
		})
		if err != nil {
			return nil, fmt.Errorf("unable to build out field declaration: %w", err)
		}
		data = append(data, varDeclaration...)

	}
	outLine := strings.Join(outArgs, ", ")
	if len(outArgs) > 1 {
		outLine = "(" + outLine + ")"
	}

	c := newCache().
		setNode(utils.NewNode[string, string]()).
		setVarPrefix("_var_").
		build()

	for _, field := range mg.outputFields {
		// TODO implement method body generator
		err := generate(field, mg.inputFields, c, mg.overloading, mg.impCache)
		if err != nil {
			return nil, fmt.Errorf("unable to generate maping for %s: %w", mg.name, err)
		}
	}

	n := c.getNode()
	n.Optimize()
	buffer := bytes.NewBufferString("")
	t := &tabs{}
	n.Execute(preProcessing(buffer, t), postProcessing(buffer, t))
	data = append(data, buffer.Bytes()...)

	return templates.Execute(templates.MapperMethodTemplate, map[string]interface{}{
		"TypeName":        receiverName,
		"MethodName":      mg.name,
		"Arguments":       strings.Join(inArgs, ", "),
		"ReturnTypes":     outLine,
		"Block":           string(data), //buffer.String(),
		"ReturnVariables": strings.Join(retVars, ", "),
	})
}

type tabs struct {
	t []string
}

func preProcessing(buffer *bytes.Buffer, t *tabs) func([]string, []string) {
	return func(k []string, value []string) {
		var keys []string
		for _, s := range k {
			keys = append(keys, s+" != nil")
		}
		ifLine := strings.Join(keys, " && ")
		if len(k) != 0 {
			buffer.Write([]byte(strings.Join(t.t, "") + "if " + ifLine + " {\n"))
			t.t = append(t.t, "\t")
		}
		if len(value) != 0 {
			buffer.Write([]byte(strings.Join(t.t, "") + strings.Join(value, "\n"+strings.Join(t.t, "")) + "\n"))
		}
	}
}

func postProcessing(buffer *bytes.Buffer, t *tabs) func([]string, []string) {
	return func(k []string, value []string) {
		if len(t.t) > 0 {
			t.t = t.t[1:]
		}
		if len(k) > 0 {
			buffer.Write([]byte(strings.Join(t.t, "") + "}\n"))
		}
	}
}

func fieldFirstName(f *ast.Field) string {
	if f == nil || len(f.Names) == 0 {
		return ""
	}
	return f.Names[0].String()
}
