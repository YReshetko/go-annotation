package cgs

import (
	"go/ast"
)

type Import struct {
	Alias   string
	Package string
}

type distinctImports map[Import]struct{}

func newDistinctImports() distinctImports {
	return map[Import]struct{}{}
}

func (d distinctImports) append(i Import) {
	d[i] = struct{}{}
}

func (d distinctImports) merge(new distinctImports) {
	for k, _ := range new {
		d[k] = struct{}{}
	}
}

func (d distinctImports) isEmpty() bool {
	return len(d) == 0
}

func (d distinctImports) toSlice() []Import {
	if d.isEmpty() {
		return nil
	}
	i := make([]Import, len(d))
	index := 0
	for k, _ := range d {
		i[index] = k
		index++
	}
	return i
}

func getImports(e ast.Expr, fn func(string) (string, bool)) distinctImports {
	out := newDistinctImports()
	ast.Inspect(e, func(node ast.Node) bool {

		switch n := node.(type) {
		case *ast.SelectorExpr:
			switch i := n.X.(type) {
			case *ast.Ident:
				alias := i.String()
				pkg, ok := fn(alias)
				if !ok {
					return true
				}
				out.append(Import{
					Alias:   alias,
					Package: pkg,
				})
			}
		}
		return true
	})
	return out
}