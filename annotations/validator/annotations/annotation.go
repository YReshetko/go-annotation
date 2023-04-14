package annotations

import (
	"github.com/YReshetko/go-annotation/annotations/constructor/templates"
	"strings"
)

const (
	All    = "all"
	Tagged = "tagged"
)

type Validator struct {
	Required string `annotation:"name=type,default=all,oneOf=all;tagged"`
	Name     string `annotation:"name=name,default={{.TypeName}}Validator"`
}

func (v Validator) BuildName(typeName string) string {
	return generateName(v.Name, "TypeName", typeName, true)
}

func generateName(tpl, replacerName, value string, isExported bool) string {
	tpl, value = correctExportCases(tpl, value, isExported)
	return templates.ExecuteTempTemplate(tpl, map[string]string{replacerName: value})
}

func correctExportCases(tpl, value string, isExported bool) (string, string) {
	tpl = strings.TrimSpace(tpl)
	value = strings.TrimSpace(value)
	if isExported {
		return firstToUpper(tpl), firstToUpper(value)
	}
	if tpl[0] == '{' {
		return tpl, firstToLower(value)
	}
	return firstToLower(tpl), firstToUpper(value)
}

const letterDiff = 'a' - 'A'

func firstToUpper(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] >= 'a' && s[0] <= 'z' {
		return string(s[0]-letterDiff) + s[1:]
	}
	return s
}

func firstToLower(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] >= 'A' && s[0] <= 'A' {
		return string(s[0]+letterDiff) + s[1:]
	}
	return s
}
