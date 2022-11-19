package cache

import "fmt"

type Import string
type AliasItem string

// ImportCache - @Constructor(type="pointer")
type ImportCache struct {
	aliasTpl    string
	index       int                  // @Exclude
	imports     map[Import]AliasItem // @Init
	baseImports map[Import]struct{}  //@Init
}

func (i *ImportCache) StoreImport(imp string) string {
	if len(imp) == 0 {
		return ""
	}
	v, ok := i.imports[Import(imp)]
	if ok {
		return string(v)
	}
	i.index++
	v = AliasItem(fmt.Sprintf(i.aliasTpl, i.index))
	i.imports[Import(imp)] = v
	return string(v)
}

func (i *ImportCache) BuildImports() [][2]string {
	var imps [][2]string
	for imp, alias := range i.imports {
		imps = append(imps, [2]string{string(alias), string(imp)})
	}
	for imp, _ := range i.baseImports {
		imps = append(imps, [2]string{"", string(imp)})
	}
	return imps
}

func (i *ImportCache) BuildReplaceMap() map[string]string {
	imps := make(map[string]string)
	for _, alias := range i.imports {
		imps[string(alias)] = string(alias)
	}
	return imps
}

func (i *ImportCache) AddImport(imp string) {
	i.baseImports[Import(imp)] = struct{}{}
}
