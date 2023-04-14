package generators

import (
	"github.com/YReshetko/go-annotation/annotations/validator/model"
	"github.com/YReshetko/go-annotation/annotations/validator/templates"
	"strings"
)

// Output @Constructor
type Output struct {
	v model.Validatable
}

func (o Output) Generate() string {
	data := model.TemplateData{
		ValidatorName:  o.v.ValidatorName, //TODO get it from annotation
		TargetTypeName: o.v.Name,
		Snippets:       fieldsSamples("value", o.v.Fields),
	}

	b, _ := templates.Execute(templates.ValidatorTypeTpl, data)
	b1, _ := templates.Execute(templates.ValidatorMethodTpl, data)
	return string(append(b, b1...))
}

func fieldsSamples(prefix string, fields []model.Field) []string {
	var data []string
	for _, field := range fields {
		fName := fieldName(prefix, field.Name)
		//fmt.Println(field.String())
		if field.Tag.IsIgnore() {
			continue
		}

		//TODO extends types
		switch {
		case isString(field.FType):
			b, _ := templates.Execute(templates.StringZeroSnippetTpl, map[string]string{"FieldName": fName})
			data = append(data, string(b))
		case isNumeric(field.FType):
			b, _ := templates.Execute(templates.NumericZeroSnippetTpl, map[string]string{"FieldName": fName})
			data = append(data, string(b))
		case isPointer(field.FType):
			b, _ := templates.Execute(templates.PointerZeroSnippetTpl, map[string]string{"FieldName": fName})
			data = append(data, string(b))
			if field.Validatable != nil {
				strs := fieldsSamples(fName, field.Fields)
				data = append(data, strs...)
			}
		case isStatic(field.FType):
			if field.Validatable != nil {
				strs := fieldsSamples(fName, field.Fields)
				data = append(data, strs...)
			}
			/*		default:
					fmt.Println(field.String())*/
		}
	}
	return data
}

func fieldName(prefix, name string) string {
	if prefix == "" {
		return name
	}
	if name == "" {
		return prefix
	}
	return strings.Join([]string{prefix, name}, ".")
}

func isStatic(t string) bool {
	return t == "static"
}

func isString(t string) bool {
	return t == "string"
}

func isPointer(t string) bool {
	return t == "pointer" || t == "any"
}

func isNumeric(t string) bool {
	numeric := map[string]struct{}{
		"uint":    {},
		"uint8":   {},
		"uint16":  {},
		"uint32":  {},
		"uint64":  {},
		"byte":    {},
		"int":     {},
		"int8":    {},
		"int16":   {},
		"int32":   {},
		"int64":   {},
		"float32": {},
		"float64": {},
		"uintptr": {},
		"rune":    {},
	}
	_, ok := numeric[t]
	return ok
}
