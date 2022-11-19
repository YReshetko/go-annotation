package generators

import (
	"fmt"
	"github.com/YReshetko/go-annotation/annotations/mapper/utils"
)

type importCache interface {
	AddImport(string)
}

type Import struct {
	Alias  string
	Import string
}

// @Builder(constructor="newCache", build="set{{.FieldName}}", terminator="build", type="pointer", exported="false")
type cache struct {
	varPrefix string
	node      *utils.Node[string, string]
	index     int //@Exclude
}

func (c *cache) nextVar() string {
	v := fmt.Sprintf("%s%d", c.varPrefix, c.index)
	c.index++
	return v
}

func (c *cache) addIfClause(checks []string, line string) {
	c.node.Add(checks, line)
}

func (c *cache) addCodeLine(line string) {
	c.node.Add(nil, line)
}

func (c *cache) getNode() *utils.Node[string, string] {
	return c.node
}
