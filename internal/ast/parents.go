package ast

import (
	"go/ast"
)

type parentCache []ast.Node

func newParentCache() *parentCache {
	return new(parentCache)
}

func (pc *parentCache) push(v ast.Node) {
	if pc == nil {
		return
	}
	n := append(*pc, v)
	*pc = n
}

func (pc *parentCache) pop() ast.Node {
	if pc == nil {
		return nil
	}
	ds := *pc
	if len(ds) == 0 {
		return nil
	}
	o := ds[len(ds)-1]
	if len(ds) == 1 {
		var n []ast.Node
		*pc = n
		return o
	}
	n := ds[:len(ds)-1]
	*pc = n
	return o
}

func (pc *parentCache) peek() ast.Node {
	if pc == nil {
		return nil
	}
	ds := *pc
	if len(ds) == 0 {
		return nil
	}
	return ds[len(ds)-1]
}

func (pc *parentCache) add(n ast.Node) {
	start := n.Pos()
	for !pc.empty() && pc.peek().End() <= start {
		pc.pop()
	}
	pc.push(n)
}

func (pc *parentCache) copy() []ast.Node {
	if pc == nil || len(*pc) == 0 {
		return nil
	}
	ds := *pc
	c := make([]ast.Node, len(ds), len(ds))
	copy(c, ds)
	return c
}

func (pc *parentCache) empty() bool {
	if pc == nil || len(*pc) == 0 {
		return true
	}
	return false
}
