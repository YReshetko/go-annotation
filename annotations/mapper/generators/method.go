package generators

import (
	"bytes"
	"fmt"
	"github.com/YReshetko/go-annotation/annotations/mapper/templates"
	"github.com/YReshetko/go-annotation/annotations/mapper/utils"
	"go/ast"
	"strings"

	annotation "github.com/YReshetko/go-annotation/pkg"
)

// @Builder(constructor="newMethodGeneratorBuilder", build="set{{ .FieldName }}", terminator="build", type="pointer")
type methodGenerator struct {
	node        annotation.Node
	name        string
	overloading *overloading
	input       []*ast.Field
	output      []*ast.Field

	inputGenerators []*fieldGenerator // @Init
	outputGenerator []*fieldGenerator // @Init
}

// @PostConstruct
func (mg *methodGenerator) buildOutput() {
	for _, field := range mg.output {
		if len(field.Names) == 0 {
			mg.outputGenerator = append(mg.outputGenerator, buildRootFieldGenerator(field, "", mg.node))
		} else {
			for _, name := range field.Names {
				mg.outputGenerator = append(mg.outputGenerator, buildRootFieldGenerator(field, name.String(), mg.node))
			}
		}
	}
}

// @PostConstruct
func (mg *methodGenerator) buildInput() {
	for i, field := range mg.input {
		if len(field.Names) == 0 {
			mg.inputGenerators = append(mg.inputGenerators, buildRootFieldGenerator(field, fmt.Sprintf("in%d", i), mg.node))
		} else {
			for _, name := range field.Names {
				mg.inputGenerators = append(mg.inputGenerators, buildRootFieldGenerator(field, name.String(), mg.node))
			}
		}
	}
}

func (mg *methodGenerator) generate(receiverName string, imports map[Import]struct{}) ([]byte, error) {
	var args []string
	c := newCache().
		setNode(utils.NewNode[string, string]()).
		setVarPrefix("_var_").
		setImports(imports).
		build()
	for _, generator := range mg.inputGenerators {
		a, err := generator.argumentNameAndType(c)
		if err != nil {
			return nil, fmt.Errorf("unable to build argument %s for %s.%s", generator.name, receiverName, mg.name)
		}
		args = append(args, string(a))
	}

	variables := make([]string, len(mg.outputGenerator))
	returnTypes := make([]string, len(mg.outputGenerator))
	for i, generator := range mg.outputGenerator {
		generator.appendImport(c)
		if len(generator.name) == 0 {
			variables[i] = c.nextVar()
		} else {
			variables[i] = generator.name
		}
		returnTypes[i] = generator.buildReturnType()
		err := generator.generate(variables[i], mg.inputGenerators, c, mg.overloading)
		if err != nil {
			return nil, fmt.Errorf("unable to generate maping for %s: %w", mg.name, err)
		}
	}

	// TODO process node lines: Optimize, Execute
	n := c.getNode()
	n.Optimize()
	buffer := bytes.NewBufferString("")
	t := &tabs{}
	n.Execute(preProcessing(buffer, t), postProcessing(buffer, t))

	returnTypeLine := strings.Join(returnTypes, ", ")
	if len(returnTypes) > 1 {
		returnTypeLine = "(" + returnTypeLine + ")"
	}
	return templates.Execute(templates.MapperMethodTemplate, map[string]interface{}{
		"TypeName":        receiverName,
		"MethodName":      mg.name,
		"Arguments":       strings.Join(args, ", "),
		"ReturnTypes":     returnTypeLine,
		"Block":           buffer.String(),
		"ReturnVariables": strings.Join(variables, ", "),
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
