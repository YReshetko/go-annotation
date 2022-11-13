package ast

import "go/ast"

func LoadFileAST(path string) (*ast.File, error) {
	return load(path)
}

func FindTopNodeByName(f *ast.File, name string) (ast.Node, []ast.Node, bool) {
	return findHighLevelNode(f, name)
}

func Comment(n ast.Node) (string, bool) {
	return extractComment(n)
}

func Walk(f *ast.File, visitor func(node ast.Node, parents []ast.Node) bool) {
	walk(f, visitor)
}
