package annotation

import (
	"go/ast"
	"go/token"

	"github.com/YReshetko/go-annotation/internal/utils/astutils"
	. "github.com/YReshetko/go-annotation/internal/utils/stream"
)

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

func CastNode[T ast.Node](n Node) (T, bool) {
	v, ok := n.ASTNode().(T)
	return v, ok
}

func BytesToAST(data []byte) (ast.Node, *token.FileSet, error) {
	return astutils.BytesToAST(data)
}
