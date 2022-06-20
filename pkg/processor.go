package pkg

import (
	"go/ast"

	"github.com/YReshetko/go-annotation/internal/nodes"
)

type Path string
type Data []byte

type Annotation any

type NodeType string

const (
	Interface NodeType = "interface"
	Structure NodeType = "structure"
	Field     NodeType = "field"
	Function  NodeType = "function"
	Variable  NodeType = "variable"
)

type Node interface {
	Name() string
	Dir() string
	FileName() string
	GoNode() ast.Node
	FileSpec() *ast.File
	NodeType() NodeType
	InnerNodes() []Node
}

type AnnotationProcessor interface {
	Process(annotation Annotation, node Node) error
	Output() map[Path]Data
}

var _ Node = (*internalNode)(nil)

type internalNode struct {
	n     nodes.Node
	inner []Node
}

func newInternalNode(n nodes.Node) internalNode {
	intNode := make([]Node, len(n.Inner))
	for i, node := range n.Inner {
		intNode[i] = newInternalNode(node)
	}
	return internalNode{
		n:     n,
		inner: intNode,
	}
}

func (i internalNode) Name() string {
	return i.n.Metadata.Name
}
func (i internalNode) Dir() string {
	return i.n.Metadata.Dir
}
func (i internalNode) FileName() string {
	return i.n.Metadata.FileName
}

func (i internalNode) GoNode() ast.Node {
	return i.n.GoNode
}

func (i internalNode) FileSpec() *ast.File {
	return i.n.Metadata.FileSpec
}

func (i internalNode) NodeType() NodeType {
	return map[nodes.NodeType]NodeType{
		nodes.Field:     Field,
		nodes.Structure: Structure,
		nodes.Interface: Interface,
		nodes.Function:  Function,
		nodes.Variable:  Variable,
	}[i.n.Metadata.Type]
}

func (i internalNode) InnerNodes() []Node {
	return i.inner
}
