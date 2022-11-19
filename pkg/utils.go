package annotation

import (
	"go/ast"
	"go/token"

	"github.com/YReshetko/go-annotation/internal/utils/astutils"
	. "github.com/YReshetko/go-annotation/internal/utils/stream"
)

// FindAnnotations returns all annotation with particular type
func FindAnnotations[T any](a []Annotation) []T {
	return Map(OfSlice(a).Filter(ofType[T]), toType[T]).ToSlice()
}

func ofType[T any](a Annotation) bool {
	_, ok := a.(T)
	return ok
}

func toType[T any](a Annotation) T {
	return a.(T)
}

// CastNode quick cast annotation.Node to ast.Node with type check
func CastNode[T ast.Node](n Node) (T, bool) {
	v, ok := n.ASTNode().(T)
	return v, ok
}

// ParentType quick retrieve parent ast.Node by type
func ParentType[T ast.Node](n Node) (Node, bool) {
	current := n
	p, ok := current.ParentNode()
	for ; ok; p, ok = current.ParentNode() {
		_, cok := CastNode[T](p)
		if cok {
			return p, true
		}
		current = p
	}
	return nil, false
}

// BytesToAST conversion some bytes to ast.Node
func BytesToAST(data []byte) (ast.Node, *token.FileSet, error) {
	return astutils.BytesToAST(data)
}
