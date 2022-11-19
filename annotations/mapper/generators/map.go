package generators

import (
	"fmt"
	"github.com/YReshetko/go-annotation/annotations/mapper/generators/nodes"
)

func isEqualMaps(f1, f2 nodes.Type) bool {
	mt1, ok1 := f1.(*nodes.MapType)
	mt2, ok2 := f2.(*nodes.MapType)
	if !(ok1 && ok2) {
		return false
	}

	return mt1.Equal(mt2)
}

func assignMap(toName, fromName string, toField, fromField *nodes.MapType, fromPrefix []string, c *cache) error {
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
