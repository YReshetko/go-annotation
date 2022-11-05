package constructor

import (
	"fmt"

	"github.com/YReshetko/go-annotation/annotations/constructor/generators"
)

type key struct {
	dir string
	pkg string
}

type cache struct {
	data map[key]map[string]*typeCache
}

func newCache() *cache {
	return &cache{data: map[key]map[string]*typeCache{}}
}

func (c *cache) addGenerator(dir, pkg, typeName string, g generator) {
	cll := c.initOnGetCollector(dir, pkg, typeName)
	cll.addGenerator(g)
}

func (c *cache) addPostConstruct(dir, pkg, typeName string, p generators.PostConstructValues) {
	cll := c.initOnGetCollector(dir, pkg, typeName)
	cll.addPostConstruct(p)
}

func (c *cache) generate() (map[key][]generators.Data, error) {
	out := make(map[key][]generators.Data)
	for k, m := range c.data {
		var data []generators.Data
		for tn, tc := range m {
			gd, err := tc.generate()
			if err != nil {
				return nil, fmt.Errorf("unable to generate data for type %s", tn)
			}
			data = append(data, gd...)
		}

		out[k] = data
	}
	return out, nil
}

func (c *cache) initOnGetCollector(dir, pkg, typeName string) *typeCache {
	k := key{
		dir: dir,
		pkg: pkg,
	}

	v, ok := c.data[k]
	if !ok {
		v = map[string]*typeCache{}
		c.data[k] = v
	}

	cll, ok := v[typeName]
	if !ok {
		cll = newTypeCache()
		v[typeName] = cll
	}
	return cll
}

type typeCache struct {
	generators     []generator
	postConstructs []generators.PostConstructValues
}

func newTypeCache() *typeCache {
	return &typeCache{
		generators:     []generator{},
		postConstructs: []generators.PostConstructValues{},
	}
}

func (c *typeCache) addGenerator(g generator) {
	c.generators = append(c.generators, g)
}

func (c *typeCache) addPostConstruct(p generators.PostConstructValues) {
	c.postConstructs = append(c.postConstructs, p)
}

func (c *typeCache) generate() ([]generators.Data, error) {
	data := make([]generators.Data, len(c.generators))
	for i, g := range c.generators {
		d, im, err := g.Generate(c.postConstructs)
		if err != nil {
			return nil, fmt.Errorf("unable to generate code sample for %s", g.Name())
		}
		data[i] = generators.Data{
			Data:    d,
			Imports: im,
		}
	}
	return data, nil
}
