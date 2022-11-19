package templates

import (
	"bytes"
	"fmt"
	"text/template"
)

const fileTemplate = `
package {{ .PackageName }}

{{ if .HasImports }} import (
		{{ range .Imports }} {{ .Alias }} "{{ .Import }}"
		{{ end }}
){{ end }}

{{ .Data }}
`

const mapperStructureTemplate = `
var _ {{ .InterfaceName }} = (*{{ .TypeName }})(nil)
type {{ .TypeName }} struct{}

{{ range .Methods }}{{.}}
{{ end }}
`

const mapperMethodTemplate = `
func (_this_ {{ .TypeName }}) {{ .MethodName }}({{ .Arguments }}) {{ .ReturnTypes }} {
	{{ .Block }}
	return {{ .ReturnVariables }}
}
`

const newVariableTemplate = `{{ .VariableName }} := {{ .Declaration }}
`

const notNilSourceTemplate = `if {{ .SourceValue }} != nil {
	{{ .Line }}
}`

const notNilBlockTemplate = `if {{ .ComparationValue }} {
	{{ .Line }}
}`

const primitiveConverterFuncTemplate = `{{ .ReceiverName }} = func(v {{ .SourceType }}) {{ if .IsPointerReceiver }}*{{ end }}{{ .ReceiverType }}{
	res := {{ .MappingLine }}
	return {{ if .IsPointerReceiver }}&{{ end }}res
}({{ if .IsPointerSource }}*{{ end }}{{ .SourceName }})`

const argumentNameAndTypeTemplate = `{{ .Name }} {{ if .IsPointer }}*{{ end }}{{ .Type }}`

const sliceMappingTemplate = `
	{{ .VariableInterimSlice }} := {{ if .IsSourcePointer }}*{{end}}{{ .SourceName }}
	{{ .VariableName }} := make({{ .ReceiverType }}, len({{ .VariableInterimSlice }}), len({{ .VariableInterimSlice }}))
	for {{ .VariableIndex }}, {{ .VariableValue }} := range {{ .VariableInterimSlice }} {
		{{ .VariableName }}[{{ .VariableIndex }}] = {{ .FunctionName }}({{ .VariableValue }})
	}
	{{ .ReceiverName }} = {{ if .IsPointerReceiver }}&{{ end }}{{ .VariableName }}
`

const mapMappingTemplate = `
	{{ .VariableInterimMap }} := {{ if .IsSourcePointer }}*{{end}}{{ .SourceName }}
	{{ .VariableName }} := make({{ .ReceiverType }}, len({{ .VariableInterimMap }}))
	for {{ .VariableKey }}, {{ .VariableValue }} := range {{ .VariableInterimMap }} {
		{{ .VariableNewKey }}, {{ .VariableNewValue }} := {{ .FunctionName }}({{ .VariableKey }}, {{ .VariableValue }})
		{{ .VariableName }}[{{ .VariableNewKey }}] = {{ .VariableNewValue }}
	}
	{{ .ReceiverName }} = {{ if .IsPointerReceiver }}&{{ end }}{{ .VariableName }}
`

const (
	FileTemplate                   = "fileTemplate"
	MapperStructureTemplate        = "mapperStructureTemplate"
	MapperMethodTemplate           = "mapperMethodTemplate"
	NewVariableTemplate            = "newVariableTemplate"
	NotNilSourceTemplate           = "notNilSourceTemplate"
	NotNilBlockTemplate            = "notNilBlockTemplate"
	PrimitiveConverterFuncTemplate = "primitiveConverterFuncTemplate"
	ArgumentNameAndTypeTemplate    = "argumentNameAndTypeTemplate"
	SliceMappingTemplate           = "sliceMappingTemplate"
	MapMappingTemplate             = "mapMappingTemplate"
)

var dataTemplate *template.Template

func init() {
	dataTemplate = must(template.New(FileTemplate).Parse(fileTemplate))
	dataTemplate = must(dataTemplate.New(MapperStructureTemplate).Parse(mapperStructureTemplate))
	dataTemplate = must(dataTemplate.New(MapperMethodTemplate).Parse(mapperMethodTemplate))
	dataTemplate = must(dataTemplate.New(NewVariableTemplate).Parse(newVariableTemplate))
	dataTemplate = must(dataTemplate.New(NotNilSourceTemplate).Parse(notNilSourceTemplate))
	dataTemplate = must(dataTemplate.New(NotNilBlockTemplate).Parse(notNilBlockTemplate))
	dataTemplate = must(dataTemplate.New(PrimitiveConverterFuncTemplate).Parse(primitiveConverterFuncTemplate))
	dataTemplate = must(dataTemplate.New(ArgumentNameAndTypeTemplate).Parse(argumentNameAndTypeTemplate))
	dataTemplate = must(dataTemplate.New(SliceMappingTemplate).Parse(sliceMappingTemplate))
	dataTemplate = must(dataTemplate.New(MapMappingTemplate).Parse(mapMappingTemplate))
}

func must[T any](t T, e error) T {
	if e != nil {
		panic(e)
	}
	return t
}

func Execute(templateName string, data any) ([]byte, error) {
	tpl := dataTemplate.Lookup(templateName)
	if tpl == nil {
		return nil, fmt.Errorf("template %s not found", templateName)
	}

	return executeTpl(tpl, data)
}

func ExecuteTemplate(tplStr string, data any) ([]byte, error) {
	tpl, err := template.New("temp").Parse(tplStr)
	if err != nil {
		return nil, fmt.Errorf("unable to parse template %s", tplStr)
	}
	return executeTpl(tpl, data)
}

func executeTpl(tpl *template.Template, data any) ([]byte, error) {
	b := bytes.NewBufferString("")
	err := tpl.Execute(b, data)
	if err != nil {
		return nil, fmt.Errorf("unable to process template %s: %w", tpl.Name(), err)
	}
	return b.Bytes(), nil
}
