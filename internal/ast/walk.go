package ast

import (
	"go/ast"
)

func walk(file *ast.File, visitor func(node ast.Node, parents []ast.Node) bool) {
	pc := newParentCache()
	ast.Inspect(file, func(node ast.Node) bool {
		if node == nil {
			return false
		}
		// Ignoring comment groups as it should be a part of another nodes
		if _, ok := node.(*ast.CommentGroup); ok {
			return false
		}

		switch n := node.(type) {
		case *ast.GenDecl:
			// TODO inspect the hacky replacement of the doc
			if len(n.Specs) == 1 && n.Doc != nil && len(n.Doc.Text()) > 0 {
				switch rs := n.Specs[0].(type) {
				case *ast.ValueSpec:
					rs.Doc = mergeCommentGroups(rs.Doc, n.Doc)
				case *ast.TypeSpec:
					rs.Doc = mergeCommentGroups(rs.Doc, n.Doc)
				}
				return true
			}
		}
		parents := pc.copy()
		pc.add(node)
		return visitor(node, parents)
	})
}
