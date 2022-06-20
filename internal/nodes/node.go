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

type Metadata struct {
	Name     string
	Dir      string
	FileName string
	Type     NodeType
	FileSpec *ast.File
}

type Node struct {
	Annotations []annotation.Annotation
	Metadata    Metadata
	GoNode      ast.Node
	Inner       []Node
}
