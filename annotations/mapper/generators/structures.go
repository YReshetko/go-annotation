package generators

import (
	"fmt"
	"github.com/YReshetko/go-annotation/annotations/mapper/generators/nodes"
)

func isBothStructures(f1, f2 nodes.Type) bool {
	_, ok1 := f1.(*nodes.StructType)
	_, ok2 := f2.(*nodes.StructType)
	return ok1 && ok2
}

func mapStructures(toName, fromName string, toField, fromField *nodes.StructType, fromPrefix []string, c *cache) error {
	if !toField.Equal(fromField) {
		// TODO resolve it by type mappers
		return nil
	}

	if toField.IsPointer() == fromField.IsPointer() {
		c.addIfClause(fromPrefix, fmt.Sprintf("%s = %s", toName, fromName))
		return nil
	}
	if toField.IsPointer() {
		c.addIfClause(fromPrefix, fmt.Sprintf("%s = &%s", toName, fromName))
		return nil
	}

	c.addIfClause(append(fromPrefix, fromName), fmt.Sprintf("%s = *%s", toName, fromName))
	return nil
}
