package cache

import (
	"encoding/json"
	"errors"
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
	handlers []templates.Handler
	flags    []templates.Flag
	command  annotations.Cobra
}

func (c *Cache) AddHandler(pkg, typeName string, handler templates.Handler) {
	k := key{
		pkg:      pkg,
		typeName: typeName,
	}
	i := c.m[k]
	i.handlers = append(i.handlers, handler)
	c.m[k] = i
}

func (c *Cache) AddFlag(pkg, typeName string, flag templates.Flag) {
	k := key{
		pkg:      pkg,
		typeName: typeName,
	}
	i := c.m[k]
	i.flags = append(i.flags, flag)
	c.m[k] = i
}

func (c *Cache) AddCommandAnnotation(pkg, typeName string, cmd annotations.Cobra) {
	k := key{
		pkg:      pkg,
		typeName: typeName,
	}
	i := c.m[k]
	i.command = cmd
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
	value sortedItem
	nodes []itemNode
}

func (c *Cache) GetInitCommands() (map[buildTagName]templates.InitCommands, error) {
	sortedItems := c.sortedItems()
	err := validateItems(sortedItems)
	if err != nil {
		return nil, err
	}

	printJson(sortedItems)
	return map[buildTagName]templates.InitCommands{}, nil
}

func validateItems(items []sortedItem) error {
	if len(items) == 0 {
		return nil
	}
	if len(items[0].Usage) != 1 {
		return errors.New("root command must have single word in usage section of @Cobra annotation")
	}
	if len(items) == 1 {
		return nil
	}
	
	return nil
}

func (c *Cache) sortedItems() []sortedItem {
	sortedItems := make([]sortedItem, 0, len(c.m))
	for k, v := range c.m {
		usages := strings.Split(v.command.Usage, " ")
		usages = slices.DeleteFunc(usages, func(s string) bool {
			return len(strings.TrimSpace(s)) == 0
		})
		buildTags := strings.Split(v.command.Build, ",")
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
