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
	name     string
	goNode   ast.Node
	dir      string
	fileName string
	fileSpec *ast.File
	inner    []Node
	nodeType NodeType
}

func newInternalNode(n nodes.Node) internalNode {
	intNode := make([]Node, len(n.Inner))
	for i, node := range n.Inner {
		intNode[i] = newInternalNode(node)
	}
	return internalNode{
		name:     n.Name,
		goNode:   n.GoNode,
		dir:      n.Dir,
		fileName: n.FileName,
		fileSpec: n.FileSpec,
		inner:    intNode,
		nodeType: map[nodes.NodeType]NodeType{
			nodes.Field:     Field,
			nodes.Structure: Structure,
			nodes.Interface: Interface,
			nodes.Function:  Function,
			nodes.Variable:  Variable,
		}[n.Type],
	}
}

func (i internalNode) Name() string {
	return i.name
}
func (i internalNode) Dir() string {
	return i.dir
}
func (i internalNode) FileName() string {
	return i.fileName
}

func (i internalNode) GoNode() ast.Node {
	return i.goNode
}

func (i internalNode) FileSpec() *ast.File {
	return i.fileSpec
}

func (i internalNode) NodeType() NodeType {
	return i.nodeType
}

func (i internalNode) InnerNodes() []Node {
	return i.inner
}
