package templates

import (
	"bytes"
	"fmt"
	"text/template"
)

const validatorType = `type {{ .ValidatorName }} struct{}
`

const validatorMethod = `func (v {{ .ValidatorName }}) IsValid(value {{ .TargetTypeName }}) bool {
	{{ range .Snippets }}{{ . }}
	{{ end }}
	return true
}
`

const stringZeroSnippet = `if {{ .FieldName }} == "" {
	return false
}`

const numericZeroSnippet = `if {{ .FieldName }} == 0 {
	return false
}`

const pointerZeroSnippet = `if {{ .FieldName }} == nil {
	return false
}`

var dataTemplate *template.Template

const (
	ValidatorTypeTpl      = "validatorType"
	ValidatorMethodTpl    = "validatorMethod"
	StringZeroSnippetTpl  = "stringZeroSnippet"
	NumericZeroSnippetTpl = "numericZeroSnippet"
	PointerZeroSnippetTpl = "pointerZeroSnippet"
)

func init() {
	dataTemplate = Must(template.New(ValidatorTypeTpl).Parse(validatorType))
	dataTemplate = Must(dataTemplate.New(ValidatorMethodTpl).Parse(validatorMethod))
	dataTemplate = Must(dataTemplate.New(StringZeroSnippetTpl).Parse(stringZeroSnippet))
	dataTemplate = Must(dataTemplate.New(NumericZeroSnippetTpl).Parse(numericZeroSnippet))
	dataTemplate = Must(dataTemplate.New(PointerZeroSnippetTpl).Parse(pointerZeroSnippet))
}

func Must[T any](t T, e error) T {
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
