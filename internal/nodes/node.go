package nodes

import (
	"go/ast"

	"github.com/YReshetko/go-annotation/internal/annotation"
)

type NodeType string

const (
	Interface NodeType = "interface"
	Structure NodeType = "structure"
	Field     NodeType = "field"
	Function  NodeType = "function"
	Variable  NodeType = "variable"
)

type Node struct {
	Annotations []annotation.Annotation
	Name        string
	GoNode      ast.Node
	FileSpec    *ast.File
	Inner       []Node
	Type        NodeType
}
