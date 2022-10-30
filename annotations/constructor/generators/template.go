package generators

import (
	"bytes"
	"fmt"
	"text/template"
)

const fileTemplate = `
package {{ .PackageName }}

{{ if .HasImports }} import (
		{{ range .Imports }} {{ .Alias }} "{{ .Package }}"
		{{ end }}
){{ end }}

{{ .Data }}
`

const constructorTemplate = `
func {{ .FunctionName }}{{ if .IsParametrized }}[{{ .ParameterConstraints }}]{{ end }}({{ range .Arguments }} {{.}}, {{ end }}) {{ if .IsPointer }}*{{ end }}{{ .ReturnType }}{{ if .IsParametrized }}[{{ .Parameters }}]{{ end }} {
	return {{ if .IsPointer }}&{{ end }}{{ .ReturnType }}{{ if .IsParametrized }}[{{ .Parameters }}]{{ end }} {
		{{ range .Fields }} {{ .Name }}: {{ .Value }},
		{{ end }}}
}
`

const optionalTypeTemplate = `
type {{ .OptionalTypeName }}{{ if .IsParametrized }}[{{ .ParameterConstraints }}]{{ end }} func(*{{ .ReturnType }}{{ if .IsParametrized }}[{{ .Parameters }}]{{ end }})
`

const optionalConstructorTemplate = `
func {{ .FunctionName }}{{ if .IsParametrized }}[{{ .ParameterConstraints }}]{{ end }}(opts ...{{ .OptionalTypeName }}{{ if .IsParametrized }}[{{ .Parameters }}]{{ end }}){{ if .IsPointer }}*{{ end }}{{ .ReturnType }}{{ if .IsParametrized }}[{{ .Parameters }}]{{ end }} {
	rt := &{{ .ReturnType }}{{ if .IsParametrized }}[{{ .Parameters }}]{{ end }}{
		{{ range .Fields }} {{ .Name }}: {{ .Value }},
		{{ end }}}
	for _, o := range opts{
		o(rt)
	}
	return {{ if not .IsPointer }}*{{ end }}rt
}
`

const optionalWithTemplate = `
func {{ .WithFunctionName }}{{ if .IsParametrized }}[{{ .ParameterConstraints }}]{{ end }}(v {{ .ArgumentType }}) {{ .OptionalTypeName }}{{ if .IsParametrized }}[{{ .Parameters }}]{{ end }} {
	return func(rt *{{ .ReturnType }}{{ if .IsParametrized }}[{{ .Parameters }}]{{ end }}) {
		rt.{{ .FieldName }} = v
	}
}
`

const (
	fileTpl                = "fileTemplate"
	constructorTpl         = "constructorTemplate"
	optionalTypeTpl        = "optionalTypeTemplate"
	optionalConstructorTpl = "optionalConstructorTemplate"
	optionalWithTpl        = "optionalWithTemplate"
	functionNameTpl        = "functionName"
	typeNameTpl            = "typeName"
)

var dataTemplate *template.Template

func init() {
	dataTemplate = must(template.New(fileTpl).Parse(fileTemplate))
	dataTemplate = must(dataTemplate.New(constructorTpl).Parse(constructorTemplate))
	dataTemplate = must(dataTemplate.New(optionalTypeTpl).Parse(optionalTypeTemplate))
	dataTemplate = must(dataTemplate.New(optionalConstructorTpl).Parse(optionalConstructorTemplate))
	dataTemplate = must(dataTemplate.New(optionalWithTpl).Parse(optionalWithTemplate))

}

func must[T any](t T, e error) T {
	if e != nil {
		panic(e)
	}
	return t
}

func execute(templateName string, data any) ([]byte, error) {
	tpl := dataTemplate.Lookup(templateName)
	if tpl == nil {
		return nil, fmt.Errorf("template %s not found", templateName)
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
