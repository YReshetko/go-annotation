package mapper

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type BuildType int

const (
	Array BuildType = iota
	Map
	Instance
	Value
)

type ResultBuilder struct {
	pkg         string
	selector    string
	typeName    string
	fieldName   string
	isPointer   bool
	builderType BuildType
	children    []ResultBuilder
}

func (rb ResultBuilder) print(prefix string) {
	if rb.fieldName != "" {
		if prefix == "" {
			prefix += rb.fieldName
		} else {
			prefix += "." + rb.fieldName
		}
	}

	switch rb.builderType {
	case Array:
		prefix += "[]"
		if rb.isPointer {
			prefix += "*"
		}
		fmt.Println("FIELD TO MAP", prefix, "- create array")
	case Map:
		prefix += "[map]"

		if rb.isPointer {
			prefix += "*"
		}
		fmt.Println("FIELD TO MAP", prefix, "- create map")
	case Instance:
		if rb.isPointer {
			prefix += "*"
		}
		fmt.Println("FIELD TO MAP", prefix, " - create new instance or map value")
	case Value:
		if rb.isPointer {
			prefix += "*"
		}
		fmt.Println("FIELD TO MAP", prefix, " - map value")
		prefix += " - map value"
	}
	for _, child := range rb.children {
		child.print(prefix)
	}
}

func (rb ResultBuilder) find(reference string) (ResultBuilder, bool) {
	selector := reference
	restRef := ""
	isFinal := true
	if strings.Contains(reference, ".") {
		isFinal = false
		selector, restRef = reference[:strings.Index(reference, ".")], reference[strings.Index(reference, ".")+1:]
	}
	if rb.fieldName != selector {
		return ResultBuilder{}, false
	}

	if isFinal {
		return rb, true
	}

	for _, child := range rb.children {
		out, ok := child.find(restRef)
		if ok {
			return out, true
		}
	}
	return ResultBuilder{}, false
}

type BuildContext struct {
	// package - used selector
	imports map[string]string

	lines []string
}

func (b *BuildContext) addImport(pkg string) string {
	selector, ok := b.imports[pkg]
	if !ok {
		// TODO selector deduplication
		// TODO replace special symbols
		selector = pkg[strings.LastIndex(pkg, "/")+1:]
		b.imports[pkg] = selector
	}
	return selector
}

func (b *BuildContext) addLine(l string) {
	b.lines = append(b.lines, l)
}

func (rb ResultBuilder) mapping(prefix string, c *BuildContext, source map[string][]ResultBuilder, m []Mapping) {
	if rb.fieldName != "" {
		if prefix == "" {
			prefix += rb.fieldName
		} else {
			prefix += "." + rb.fieldName
		}
	}

	switch rb.builderType {
	case Array:
		if len(rb.children) != 1 {
			return
		}
		selector := c.addImport(rb.children[0].pkg)
		if len(selector) > 0 {
			selector = selector + "."
		}
		eq := "="
		if !strings.Contains(prefix, ".") {
			eq = ":="
		}
		modifier := ""
		if rb.children[0].isPointer {
			modifier = "&"
		}
		l := prefix + eq + "[]" + modifier + selector + rb.typeName + "{}"
		c.addLine(l)

		sourcePath := findSourceFromMapping(prefix, m)
		var sourceBuilder ResultBuilder
		var ok bool
		var sourceParamName string
		for k, builders := range source {
			for _, builder := range builders {
				sourceBuilder, ok = builder.find(k + "." + sourcePath)
				if ok {
					sourceParamName = k
					break
				}
			}
		}

		if ok {
			c.addLine("for _, v := range " + sourceParamName + "." + sourcePath + " {")
			defer c.addLine("}")
			c.addLine(prefix + "= append(" + prefix + "," + ")")
			// TODO make param check (int - *int, int - int32, float32 - string...)
			tmp := fromToBaseTypeMapping(c, sourceBuilder, rb)
			templ, err := template.New("exp").Parse(tmp)
			if err != nil {
				return
			}

			buff := bytes.NewBufferString("")
			err = templ.Execute(buff, map[string]string{"L": prefix, "R": sourceParamName + "." + sourcePath})
			if err != nil {
				return
			}

			c.addLine(buff.String())
		}
		for _, child := range rb.children {
			child.mapping(prefix, c, source, m)
		}

	case Map:
	case Instance:
		selector := c.addImport(rb.pkg)
		if len(selector) > 0 {
			selector = selector + "."
		}
		eq := "="
		if !strings.Contains(prefix, ".") {
			eq = ":="
		}
		modifier := ""
		if rb.isPointer {
			modifier = "&"
		}
		l := prefix + eq + modifier + selector + rb.typeName + "{}"
		c.addLine(l)
		for _, child := range rb.children {
			child.mapping(prefix, c, source, m)
		}
	case Value:
		sourcePath := findSourceFromMapping(prefix, m)
		var sourceBuilder ResultBuilder
		var ok bool
		var sourceParamName string
		for k, builders := range source {
			for _, builder := range builders {
				sourceBuilder, ok = builder.find(k + "." + sourcePath)
				if ok {
					sourceParamName = k
					break
				}
			}
		}

		if ok {
			// TODO make param check (int - *int, int - int32, float32 - string...)
			tmp := fromToBaseTypeMapping(c, sourceBuilder, rb)
			templ, err := template.New("exp").Parse(tmp)
			if err != nil {
				return
			}

			buff := bytes.NewBufferString("")
			err = templ.Execute(buff, map[string]string{"L": prefix, "R": sourceParamName + "." + sourcePath})
			if err != nil {
				return
			}

			c.addLine(buff.String())
		}
		for _, child := range rb.children {
			child.mapping(prefix, c, source, m)
		}
	}
}

func fromToBaseTypeMapping(c *BuildContext, from, to ResultBuilder) string {
	// The same type mapping
	if from.typeName == to.typeName {
		if from.isPointer == to.isPointer {
			return "{{ .L }} = {{ .R }}"
		}
		if from.isPointer {
			return "if {{ .R }} != nil { {{ .L }} = *{{ .R }} }"
		}
		if to.isPointer {
			return "{{ .L }} = &{{ .R }}"
		}
	}
	// Numeric type mapping
	nType := map[string]struct{}{
		"int":     {},
		"uint":    {},
		"uintptr": {},
		"byte":    {},
		"rune":    {},
		"uint8":   {},
		"uint16":  {},
		"uint32":  {},
		"uint64":  {},
		"int8":    {},
		"int16":   {},
		"int32":   {},
		"int64":   {},
		"float32": {},
		"float64": {},
	}
	_, ft := nType[from.typeName]
	_, tt := nType[to.typeName]
	if ft && tt {
		if !from.isPointer && !to.isPointer {
			return "{{ .L }} = " + to.typeName + "({{ .R }})"
		}
		if from.isPointer && !to.isPointer {
			return "if {{ .R }} != nil { {{ .L }} = " + to.typeName + "(*{{ .R }}) }"
		}
		if from.isPointer && to.isPointer {
			return "if {{ .R }} != nil { tmp := " + to.typeName + "(*{{ .R }}); {{ .L }} = &tmp }"
		}
		if !from.isPointer && to.isPointer {
			return "tmp := " + to.typeName + "({{ .R }}); {{ .L }} = &tmp"
		}
	}

	// From type is string
	if to.typeName == "string" && ft {
		_, ok := map[string]struct{}{
			"float32": {},
			"float64": {},
		}[from.typeName]
		formatter := "d"
		if ok {
			formatter = "f"
		}
		fmtImp := c.addImport("fmt")

		if !from.isPointer && !to.isPointer {
			return `{{ .L }} = ` + fmtImp + `.Sprintf("%` + formatter + `", {{ .R }})`
		}
		if from.isPointer && !to.isPointer {
			return `if {{ .R }} != nil { {{ .L }} = ` + fmtImp + `.Sprintf("%` + formatter + `", *{{ .R }}) }`
		}
		if from.isPointer && to.isPointer {
			return `if {{ .R }} != nil { tmp := ` + fmtImp + `.Sprintf("%` + formatter + `", *{{ .R }}); {{ .L }} = &tmp }`
		}
		if !from.isPointer && to.isPointer {
			return `tmp := ` + fmtImp + `.Sprintf("%` + formatter + `", {{ .R }}); {{ .L }} = &tmp`
		}
	}

	return ""
}

func findSourceFromMapping(prefix string, m []Mapping) string {
	rest := ""
	if strings.Contains(prefix, ".") {
		rest = prefix[strings.Index(prefix, ".")+1:]
	}
	for _, mapping := range m {
		if mapping.Target == rest {
			return mapping.Source
		}
	}
	return rest
}
