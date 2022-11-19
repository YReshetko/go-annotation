package generators

import (
	"fmt"
	"strings"

	"github.com/YReshetko/go-annotation/annotations/mapper/generators/nodes"
	"github.com/YReshetko/go-annotation/annotations/mapper/templates"
)

func override(m *mapping, longFieldName string, toType nodes.Type, in []*nodes.Field, c *cache, imp importCache) error {
	switch m.mappingType {
	case source:
		mapperLine := m.source
		err := overrideSource(longFieldName, mapperLine, in, toType, c)
		if err != nil {
			return fmt.Errorf("unable to generate overriden source %s: %w", mapperLine, err)
		}
	case function:
		mapperLine := m.funcOrThisLine()
		err := overrideFunction(longFieldName, mapperLine, in, c)
		if err != nil {
			return fmt.Errorf("unable to generate overloaded function %s: %w", longFieldName, err)
		}
	case constant:
		mapperLine := m.constant
		prType, ok := toType.(*nodes.PrimitiveType)
		if !ok {
			return fmt.Errorf("unable to set constant for non-primitive %s = %s", longFieldName, mapperLine)
		}
		err := mapConstant(longFieldName, prType.Name(), mapperLine, prType.IsPointer(), c, imp)
		if err != nil {
			return fmt.Errorf("unable to build constant mapping %s: %w", longFieldName, err)
		}
	case slice:
		fromField := m.source
		funcName := m.funcOrThisLine()
		_, ok := toType.(*nodes.ArrayType)
		if !ok {
			return fmt.Errorf("unable to prepare mapping for non-slice field %s = %s(%s)", longFieldName, funcName, fromField)
		}
		err := overrideArray(longFieldName, fromField, funcName, toType, in, c)
		if err != nil {
			return fmt.Errorf("unable to build slice mapping %s: %w", longFieldName, err)
		}
	case dictionary:
		fromField := m.source
		funcName := m.funcOrThisLine()
		_, ok := toType.(*nodes.MapType)
		if !ok {
			return fmt.Errorf("unable to prepare mapping for non-map field %s = %s(%s)", longFieldName, funcName, fromField)
		}
		err := overrideMap(longFieldName, fromField, funcName, toType, in, c)
		if err != nil {
			return fmt.Errorf("unable to build map mapping %s: %w", longFieldName, err)
		}
	}
	return nil
}

func overrideSource(toName, fromLine string, in []*nodes.Field, to nodes.Type, c *cache) error {
	buff, isFromPointer, err := extractNilCheck(fromLine, in)
	if err != nil {
		return fmt.Errorf("field path %s: %w", fromLine, err)
	}

	assignLine := toName + " = " + fromLine
	if to.IsPointer() != isFromPointer {
		if to.IsPointer() {
			assignLine = toName + " = &" + fromLine
		} else {
			assignLine = toName + " = *" + fromLine
		}
	}

	if len(buff) > 0 {
		c.addIfClause(buff, assignLine)
		return nil
	}
	c.addCodeLine(assignLine)
	return nil
}

func overrideFunction(toName, mappingLine string, in []*nodes.Field, c *cache) error {
	cbIndex := strings.Index(mappingLine, ")")
	obIndex := strings.Index(mappingLine, "(")
	if cbIndex == -1 || obIndex == -1 || cbIndex < obIndex {
		return fmt.Errorf("invalid function call %s", mappingLine)
	}

	args := strings.Split(mappingLine[strings.Index(mappingLine, "(")+1:strings.Index(mappingLine, ")")], ",")
	var nilableLines []string
	for _, arg := range args {
		buff, _, err := extractNilCheck(strings.TrimSpace(arg), in)
		if err != nil {
			return fmt.Errorf("field path %s: %w", arg, err)
		}
		if len(buff) > 0 {
			nilableLines = append(nilableLines, buff...)
		}
	}

	assignLine := toName + " = " + mappingLine
	if len(nilableLines) > 0 {
		var ls []string
		distinct := map[string]struct{}{}
		for _, line := range nilableLines {
			if _, ok := distinct[line]; !ok {
				ls = append(ls, line)
				distinct[line] = struct{}{}
			}
		}
		c.addIfClause(ls, assignLine)
		return nil
	}
	c.addCodeLine(assignLine)
	return nil
}

func overrideArray(toName, fromPath, funcName string, toField nodes.Type, in []*nodes.Field, c *cache) error {
	buff, isSourcePointer, err := extractNilCheck(fromPath, in)
	if err != nil {
		return fmt.Errorf("unable to build mapping for %s: %w", toName, err)
	}

	data, err := templates.Execute(templates.SliceMappingTemplate, map[string]interface{}{
		"VariableInterimSlice": c.nextVar(),
		"IsSourcePointer":      isSourcePointer,
		"VariableName":         c.nextVar(),
		"ReceiverType":         toField.DeclaredType(),
		"SourceName":           fromPath,
		"VariableIndex":        c.nextVar(),
		"VariableValue":        c.nextVar(),
		"FunctionName":         funcName,
		"ReceiverName":         toName,
		"IsPointerReceiver":    toField.IsPointer(),
	})
	if err != nil {
		return fmt.Errorf("unable to prepare slice mapping template %s: %w", toName, err)
	}

	if len(buff) > 0 {
		c.addIfClause(buff, string(data))
		return nil
	}

	c.addCodeLine(string(data))
	return nil
}

func overrideMap(toName, fromPath, funcName string, toField nodes.Type, in []*nodes.Field, c *cache) error {
	buff, isSourcePointer, err := extractNilCheck(fromPath, in)
	if err != nil {
		return fmt.Errorf("unable to build mapping for %s: %w", toName, err)
	}

	data, err := templates.Execute(templates.MapMappingTemplate, map[string]interface{}{
		"VariableInterimMap": c.nextVar(),
		"IsSourcePointer":    isSourcePointer,
		"SourceName":         fromPath,
		"VariableName":       c.nextVar(),
		"ReceiverType":       toField.DeclaredType(),
		"VariableKey":        c.nextVar(),
		"VariableValue":      c.nextVar(),
		"VariableNewKey":     c.nextVar(),
		"VariableNewValue":   c.nextVar(),
		"FunctionName":       funcName,
		"ReceiverName":       toName,
		"IsPointerReceiver":  toField.IsPointer(),
	})
	if err != nil {
		return fmt.Errorf("unable to prepare slice mapping template %s: %w", toName, err)
	}

	if len(buff) > 0 {
		c.addIfClause(buff, string(data))
		return nil
	}

	c.addCodeLine(string(data))
	return nil
}
