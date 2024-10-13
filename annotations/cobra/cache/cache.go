package cache

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/YReshetko/go-annotation/annotations/cobra/annotations"
	"github.com/YReshetko/go-annotation/annotations/cobra/templates"
)

type key struct {
	pkg      string
	typeName string
}

type Cache struct {
	m map[key]item
}

func NewCache() *Cache {
	return &Cache{m: map[key]item{}}
}

type item struct {
	Handlers []templates.Handler
	Flags    []templates.Flag
	Command  annotations.Cobra
}

func (c *Cache) AddHandler(pkg, typeName string, handler templates.Handler) {
	k := key{
		pkg:      pkg,
		typeName: typeName,
	}
	i := c.m[k]
	i.Handlers = append(i.Handlers, handler)
	c.m[k] = i
}

func (c *Cache) AddFlag(pkg, typeName string, flag templates.Flag) {
	k := key{
		pkg:      pkg,
		typeName: typeName,
	}
	i := c.m[k]
	i.Flags = append(i.Flags, flag)
	c.m[k] = i
}

func (c *Cache) AddCommandAnnotation(pkg, typeName string, cmd annotations.Cobra) {
	k := key{
		pkg:      pkg,
		typeName: typeName,
	}
	i := c.m[k]
	i.Command = cmd
	c.m[k] = i
}

// Print TODO Remove as it's created for debug
func (c *Cache) Print() {
	for k, v := range c.m {
		fmt.Printf("%v: %v\n", k, v)
	}
}

type buildTagName = string

type sortedItem struct {
	Key       key
	Value     item
	Usage     []string
	BuildTags []string
}

type itemNode struct {
	Value sortedItem
	Nodes []*itemNode
}

func (n *itemNode) add(usages []string, item sortedItem) bool {
	if len(usages) == 0 {
		return false
	}
	if len(usages) == 1 {
		item.Value.Command.Usage = usages[0]
		n.Nodes = append(n.Nodes, &itemNode{
			Value: item,
		})
		return true
	}
	for _, node := range n.Nodes {

		if node.Value.Value.Command.Usage == usages[0] {
			node.add(usages[1:], item)
			return true
		}
	}
	return false
}

func (c *Cache) GetInitCommands() (map[buildTagName]templates.InitCommands, error) {
	sortedItems := c.sortedItems()
	builds, err := buildTrees(sortedItems)
	if err != nil {
		return nil, err
	}
	out := map[buildTagName]templates.InitCommands{}
	for name, node := range builds {
		tmp := initCommands(node)
		tmp.BuildTag = name
		out[name] = tmp
	}

	//printJson(builds)
	return out, nil
}

type importCache struct {
	imports map[string]templates.Import
	index   int
	aliases map[string]struct{}
}

func newImportCache() *importCache {
	return &importCache{
		imports: map[string]templates.Import{},
		index:   0,
		aliases: map[string]struct{}{},
	}
}

func (c *importCache) add(pkg string) string {
	imp, ok := c.imports[pkg]
	if ok {
		return imp.Alias
	}
	pathItems := strings.Split(pkg, "/")
	alias := pathItems[len(pathItems)-1]
	if _, ok := c.aliases[alias]; !ok {
		alias = fmt.Sprintf("_imp%d", c.index)
		c.index++
	}
	c.aliases[alias] = struct{}{}
	c.imports[pkg] = templates.Import{
		Alias:   alias,
		Package: pkg,
	}
	return alias
}

func (c *importCache) getImports() []templates.Import {
	out := make([]templates.Import, 0, len(c.imports))
	for _, t := range c.imports {
		out = append(out, t)
	}
	return out
}

func namesGenerator() func() string {
	index := 0
	return func() string {
		index++
		return fmt.Sprintf("_cmd%d", index)
	}
}

func initCommands(root *itemNode) templates.InitCommands {
	if len(root.Nodes) > 1 {
		root.Value = sortedItem{
			Key:   root.Nodes[0].Value.Key,
			Value: item{Command: annotations.Cobra{Usage: "cli"}},
			Usage: []string{"cli"},
		}
	} else if len(root.Nodes) == 1 {
		root = root.Nodes[0]
	} else {
		return templates.InitCommands{}
	}

	cache := newImportCache()
	names := namesGenerator()
	commands := collectCommands(cache, names, "root", "", root)

	return templates.InitCommands{
		Imports:  cache.getImports(),
		Commands: commands,
	}
}

func collectCommands(impCache *importCache, names func() string, varName, parentVarName string, root *itemNode) []templates.Command {
	cmd := templates.Command{
		IsRoot:        varName == "root",
		VarName:       varName,
		ParentVarName: parentVarName,
		Use:           root.Value.Value.Command.Usage,
		Example:       root.Value.Value.Command.Example,
		Short:         root.Value.Value.Command.Short,
		Long:          root.Value.Value.Command.Long,
		SilenceUsage:  root.Value.Value.Command.SilenceUsage,
		SilenceErrors: root.Value.Value.Command.SilenceError,
		Flags:         root.Value.Value.Flags,
		Handlers:      root.Value.Value.Handlers,
	}
	for i, handler := range cmd.Handlers {
		handler.ExecutorPackageAlias = impCache.add(root.Value.Key.pkg)
		handler.ExecutorTypeName = root.Value.Key.typeName
		cmd.Handlers[i] = handler
	}
	out := []templates.Command{cmd}
	for _, node := range root.Nodes {
		newVarName := names()
		out = append(
			out,
			collectCommands(impCache, names, newVarName, varName, node)...,
		)
	}
	return out
}

func buildTrees(sortedItems []sortedItem) (map[buildTagName]*itemNode, error) {
	builds := map[buildTagName]*itemNode{}
	for _, s := range sortedItems {
		for _, tag := range s.BuildTags {
			root, ok := builds[tag]
			if !ok {
				root = &itemNode{}
			}
			ok = root.add(s.Usage, s)
			if !ok {
				return nil, fmt.Errorf("unable to find parent command for %s", strings.Join(s.Usage, " "))
			}
			builds[tag] = root
		}
	}
	return builds, nil
}

func (c *Cache) sortedItems() []sortedItem {
	sortedItems := make([]sortedItem, 0, len(c.m))
	for k, v := range c.m {
		usages := strings.Split(v.Command.Usage, " ")
		usages = slices.DeleteFunc(usages, func(s string) bool {
			return len(strings.TrimSpace(s)) == 0
		})
		buildTags := strings.Split(v.Command.Build, ",")
		for i, tag := range buildTags {
			buildTags[i] = strings.TrimSpace(tag)
		}
		sortedItems = append(sortedItems, sortedItem{
			Key:       k,
			Value:     v,
			Usage:     usages,
			BuildTags: buildTags,
		})
	}
	slices.SortFunc(sortedItems, func(a, b sortedItem) int {
		return len(a.Usage) - len(b.Usage)
	})
	return sortedItems
}

func printJson(v any) {
	d, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(d))
}
