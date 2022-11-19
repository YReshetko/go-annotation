package generators

import (
	"fmt"
	"github.com/YReshetko/go-annotation/annotations/mapper/generators/nodes"
)

func isEqualSlices(f1, f2 nodes.Type) bool {
	at1, ok1 := f1.(*nodes.ArrayType)
	at2, ok2 := f2.(*nodes.ArrayType)
	if !(ok1 && ok2) {
		return false
	}

	return at1.Equal(at2)
}

func assignSlice(toName, fromName string, toField, fromField *nodes.ArrayType, fromPrefix []string, c *cache) error {
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
