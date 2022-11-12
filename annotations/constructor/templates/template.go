package templates

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
	returnValue := {{ if .IsPointer }}&{{ end }}{{ .ReturnType }}{{ if .IsParametrized }}[{{ .Parameters }}]{{ end }} {
		{{ range .Fields }} {{ .Name }}: {{ .Value }},
	{{ end }}}
	{{ range .PostConstructs }}returnValue.{{ . }}()
	{{ end }}
	return returnValue
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
	{{ range .PostConstructs }}rt.{{ . }}()
	{{ end }}
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

const builderTypeTemplate = `
type {{ .BuilderTypeName }}{{ if .IsParametrized }}[{{ .ParameterConstraints }}]{{ end }} struct {
	{{ range .Arguments }} {{ .FakeName }} {{ .Type }} 
	{{ end }}
}
`

const builderConstructorTemplate = `
func {{ .ConstructorName }}{{ if .IsParametrized }}[{{ .ParameterConstraints }}]{{ end }}() *{{ .BuilderTypeName }}{{ if .IsParametrized }}[{{ .Parameters }}]{{ end }}  {
	return &{{ .BuilderTypeName }}{{ if .IsParametrized }}[{{ .Parameters }}]{{ end }}{}
}
`

const builderMethodTemplate = `
func (b *{{ .BuilderTypeName }}{{ if .IsParametrized }}[{{ .Parameters }}]{{ end }}){{ .BuilderMethodName }}(v {{ .ArgumentType }}) *{{ .BuilderTypeName }}{{ if .IsParametrized }}[{{ .Parameters }}]{{ end }} {
	b.{{ .FakeName }} = v
	return b
}
`

const builderBuildMethodTemplate = `
func (b *{{ .BuilderTypeName }}{{ if .IsParametrized }}[{{ .Parameters }}]{{ end }}) {{ .BuildMethodName }}() {{ if .IsPointer }}*{{ end }}{{ .ReturnType }}{{ if .IsParametrized }}[{{ .Parameters }}]{{ end }}{
	out := {{ .ReturnType }}{{ if .IsParametrized }}[{{ .Parameters }}]{{ end }}{}
	{{ range .Fields }}if b.{{ .FakeName }} == nil {
		b.{{ .FakeName }} = {{ .Value }}
	}
	{{ end }}
	{{ range .Arguments }}out.{{ .Name }} = b.{{ .FakeName }}
	{{ end }}
	{{ range .PostConstructs }}out.{{ . }}()
	{{ end }}
	return {{ if .IsPointer }}&{{ end }}out
}
`

const (
	FileTpl                = "fileTemplate"
	ConstructorTpl         = "constructorTemplate"
	OptionalTypeTpl        = "optionalTypeTemplate"
	OptionalConstructorTpl = "optionalConstructorTemplate"
	OptionalWithTpl        = "optionalWithTemplate"
	BuilderTypeTpl         = "builderTypeTemplate"
	BuilderConstructorTpl  = "builderConstructorTemplate"
	BuilderMethodTpl       = "builderMethodTemplate"
	BuilderBuildMethodTpl  = "builderBuildMethodTemplate"
	TemporaryTpl           = "temporaryTemplate"
)

var dataTemplate *template.Template

func init() {
	dataTemplate = Must(template.New(FileTpl).Parse(fileTemplate))
	dataTemplate = Must(dataTemplate.New(ConstructorTpl).Parse(constructorTemplate))
	dataTemplate = Must(dataTemplate.New(OptionalTypeTpl).Parse(optionalTypeTemplate))
	dataTemplate = Must(dataTemplate.New(OptionalConstructorTpl).Parse(optionalConstructorTemplate))
	dataTemplate = Must(dataTemplate.New(OptionalWithTpl).Parse(optionalWithTemplate))
	dataTemplate = Must(dataTemplate.New(BuilderTypeTpl).Parse(builderTypeTemplate))
	dataTemplate = Must(dataTemplate.New(BuilderConstructorTpl).Parse(builderConstructorTemplate))
	dataTemplate = Must(dataTemplate.New(BuilderMethodTpl).Parse(builderMethodTemplate))
	dataTemplate = Must(dataTemplate.New(BuilderBuildMethodTpl).Parse(builderBuildMethodTemplate))

}

func Must[T any](t T, e error) T {
	if e != nil {
		panic(e)
	}
	return t
}

func ExecuteTempTemplate(tpl string, data any) string {
	//typeNameData := map[string]string{"TypeName": g.node.Name.Name}
	return string(Must(ExecuteTemplate(Must(template.New(TemporaryTpl).Parse(tpl)), data)))
}

func Execute(templateName string, data any) ([]byte, error) {
	tpl := dataTemplate.Lookup(templateName)
	if tpl == nil {
		return nil, fmt.Errorf("template %s not found", templateName)
	}

	return ExecuteTemplate(tpl, data)
}

func ExecuteTemplate(tpl *template.Template, data any) ([]byte, error) {
	b := bytes.NewBufferString("")
	err := tpl.Execute(b, data)
	if err != nil {
		return nil, fmt.Errorf("unable to process template %s: %w", tpl.Name(), err)
	}
	return b.Bytes(), nil
}
