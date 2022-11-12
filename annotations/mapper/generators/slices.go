package generators

import (
	"fmt"
	"github.com/YReshetko/go-annotation/annotations/mapper/templates"
	"strings"
)

func mapSlices(toName, fromPath, funcName string, toField *fieldGenerator, in []*fieldGenerator, c *cache) error {
	names := strings.Split(fromPath, ".")
	buff, isSourcePointer, err := findPointersInPath("", names, in, []string{})
	if err != nil {
		return fmt.Errorf("unable to build mapping for %s: %w", toName, err)
	}

	data, err := templates.Execute(templates.SliceMappingTemplate, map[string]interface{}{
		"VariableInterimSlice": c.nextVar(),
		"IsSourcePointer":      isSourcePointer,
		"VariableName":         c.nextVar(),
		"ReceiverType":         toField.buildArgType(),
		"SourceName":           fromPath,
		"VariableIndex":        c.nextVar(),
		"VariableValue":        c.nextVar(),
		"FunctionName":         funcName,
		"ReceiverName":         toName,
		"IsPointerReceiver":    toField.isPointer,
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
