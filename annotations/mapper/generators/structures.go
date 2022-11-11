package generators

import (
	"fmt"
)

func isBothStructures(f1, f2 *fieldGenerator) bool {
	return f1.structGen != nil && f2.structGen != nil
}

func mapStructures(toName, fromName string, toField, fromField *fieldGenerator, fromPrefix []string, c *cache) error {
	toNode, _, err := toField.node.FindNodeByAlias(toField.alias, toField.structGen.name)
	if err != nil {
		return fmt.Errorf("unable to preload node for %s: %w", toName, err)
	}
	fromNode, _, err := fromField.node.FindNodeByAlias(fromField.alias, fromField.structGen.name)
	if err != nil {
		return fmt.Errorf("unable to preload node for %s: %w", fromName, err)
	}

	if !toNode.IsSamePackage(fromNode) {
		return nil
	}

	if toField.isPointer == fromField.isPointer {
		c.addIfClause(fromPrefix, fmt.Sprintf("%s = %s", toName, fromName))
		return nil
	}
	if toField.isPointer {
		c.addIfClause(fromPrefix, fmt.Sprintf("%s = &%s", toName, fromName))
		return nil
	}

	c.addIfClause(append(fromPrefix, fromName), fmt.Sprintf("%s = *%s", toName, fromName))
	return nil
}
