package generators

import (
	"fmt"
	"github.com/YReshetko/go-annotation/annotations/mapper/utils"
)

type Import struct {
	Alias  string
	Import string
}

// @Builder(constructor="newCache", build="set{{.FieldName}}", terminator="build", type="pointer")
type cache[K comparable, V any] struct {
	varPrefix string
	imports   map[Import]struct{}
	node      *utils.Node[K, V]
	index     int //@Exclude
}

func (c *cache[K, V]) addImport(i Import) {
	c.imports[i] = struct{}{}
}

func (c *cache[K, V]) nextVar() string {
	v := fmt.Sprintf("%s%d", c.varPrefix, c.index)
	c.index++
	return v
}

func (c *cache[K, V]) addIfClause(checks []K, line V) {
	c.node.Add(checks, line)
}

func (c *cache[K, V]) addCodeLine(line V) {
	c.node.Add(nil, line)
}

func (c *cache[K, V]) getNode() *utils.Node[K, V] {
	return c.node
}
