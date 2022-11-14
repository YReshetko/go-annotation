package annotations

import (
	"github.com/YReshetko/go-annotation/annotations/constructor/templates"
	"strings"
)

const (
	structType  = "struct"
	pointerType = "pointer"
)

type Constructor struct {
	Name     string `annotation:"name=name,default=New{{.TypeName}}"`
	Type     string `annotation:"name=type,default=struct,oneOf=struct;pointer"`
	Exported bool   `annotation:"name=exported,default=true"`
}

type Optional struct {
	Name            string `annotation:"name=name,default={{.TypeName}}Option"`
	ConstructorName string `annotation:"name=constructor,default=New{{.TypeName}}"`
	WithPattern     string `annotation:"name=with,default=With{{.FieldName}}"`
	Type            string `annotation:"name=type,default=struct,oneOf=struct;pointer"`
	Exported        bool   `annotation:"name=exported,default=true"`
}

type Builder struct {
	StructureName   string `annotation:"name=name,default={{.TypeName}}Builder"`
	ConstructorName string `annotation:"name=constructor,default=New{{.TypeName}}Builder"`
	BuildPattern    string `annotation:"name=build,default={{.FieldName}}"`
	BuilderName     string `annotation:"name=terminator,default=Build"`
	Type            string `annotation:"name=type,default=struct,oneOf=struct;pointer"`
	Exported        bool   `annotation:"name=exported,default=true"`
}

type PostConstruct struct {
	Priority int `annotation:"name=priority,default=1"`
}

// Init is used for fields initialisation such as slice, map, chan
// If Init.Len and Init.Cap then the values are set by default (chan is non-buffered)
type Init struct {
	Len int `annotation:"name=len,default=-1"`
	Cap int `annotation:"name=cap,default=-1"`
}

type Exclude struct{} // Excludes structure field from constructors

func (c Constructor) BuildName(typeName string) string {
	return generateName(c.Name, "TypeName", typeName, c.Exported)
}

func (o Optional) BuildName(typeName string) string {
	return generateName(o.Name, "TypeName", typeName, o.Exported)
}

func (o Optional) BuildConstructorName(typeName string) string {
	return generateName(o.ConstructorName, "TypeName", typeName, o.Exported)
}

func (o Optional) BuildWithName(fieldName string) string {
	return generateName(o.WithPattern, "FieldName", fieldName, o.Exported)
}

func (b Builder) BuildStructureName(typeName string) string {
	return generateName(b.StructureName, "TypeName", typeName, b.Exported)
}

func (b Builder) BuildConstructorName(typeName string) string {
	return generateName(b.ConstructorName, "TypeName", typeName, b.Exported)
}

func (b Builder) BuildBuildName(fieldName string) string {
	return generateName(b.BuildPattern, "FieldName", fieldName, b.Exported)
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
