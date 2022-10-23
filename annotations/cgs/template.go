package cgs

import (
	"bytes"
	"fmt"
	"text/template"
)

var fileTemplate = `
package {{ .PackageName }}

{{ if .HasImports }} import (
		{{ range .Imports }} {{ .Alias }} "{{ .Package }}"
		{{ end }}
){{ end }}

{{ .Data }}
`

var constructorTemplate = `
func {{ .FunctionName }}{{ if .IsParametrized }}[{{ .ParameterConstraints }}]{{ end }}({{ range .Arguments }} {{.}}, {{ end }}) {{ if .IsPointer }}*{{ end }}{{ .ReturnType }}{{ if .IsParametrized }}[{{ .Parameters }}]{{ end }} {
	return {{ if .IsPointer }}&{{ end }}{{ .ReturnType }}{{ if .IsParametrized }}[{{ .Parameters }}]{{ end }} {
		{{ range .Fields }} {{ .Name }}: {{ .Value }},
		{{ end }}}
}
`

var getterTemplate = ``
var setterTemplate = ``

const (
	file         = "fileTemplate"
	constructor  = "constructorTemplate"
	getter       = "getterTemplate"
	setter       = "setterTemplate"
	functionName = "functionName"
)

var dataTemplate *template.Template

func init() {
	dataTemplate = must(template.New(file).Parse(fileTemplate))
	dataTemplate = must(dataTemplate.New(constructor).Parse(constructorTemplate))
	dataTemplate = must(dataTemplate.New(getter).Parse(getterTemplate))
	dataTemplate = must(dataTemplate.New(setter).Parse(setterTemplate))

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
