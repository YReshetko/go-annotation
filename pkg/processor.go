package pkg

import (
	"go/ast"

	"github.com/YReshetko/go-annotation/internal/annotation/tag"
	"github.com/YReshetko/go-annotation/internal/module"
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
	Method    NodeType = "method"
	Variable  NodeType = "variable"
	Type      NodeType = "type"
)

type Node interface {
	Name() string
	Dir() string
	FileName() string
	GoNode() ast.Node
	FileSpec() *ast.File
	NodeType() NodeType
	InnerNodes() []Node
	Annotations() []Annotation
	ModuleName() string
}

type Selector struct {
	PackageImport string
	TypeName      string
}

type AnnotationProcessor interface {
	Process(annotation Annotation, node Node) error
	Output() map[Path]Data
	SetLookup(lookup Lookup)
}

type Lookup func(Node, Selector) Node

var _ Node = (*node)(nil)

type node struct {
	n           module.Node
	inner       []node
	annotations []Annotation
}

func newNode(n module.Node) node {
	intNode := make([]node, len(n.Inner))
	for i, node := range n.Inner {
		intNode[i] = newNode(node)
	}
	var annptations []Annotation
	for _, annotation := range n.Annotations {
		a, ok := annotations[annotation.Name()]
		if !ok {
			continue
		}
		annptations = append(annptations, tag.Parse(a, annotation))
	}
	return node{
		n:           n,
		inner:       intNode,
		annotations: annptations,
	}
}

func (i node) Name() string {
	return i.n.Metadata.Name
}
func (i node) Dir() string {
	return i.n.Metadata.Dir
}
func (i node) FileName() string {
	return i.n.Metadata.FileName
}

func (i node) GoNode() ast.Node {
	return i.n.GoNode
}

func (i node) FileSpec() *ast.File {
	return i.n.Metadata.FileSpec
}

func (i node) NodeType() NodeType {
	return map[module.NodeType]NodeType{
		module.Field:     Field,
		module.Structure: Structure,
		module.Interface: Interface,
		module.Function:  Function,
		module.Method:    Method,
		module.Variable:  Variable,
		module.Type:      Type,
	}[i.n.Metadata.Type]
}

func (i node) InnerNodes() []Node {
	out := make([]Node, len(i.inner))
	for i, n := range i.inner {
		out[i] = n
	}
	return out
}

func (i node) Annotations() []Annotation {
	return i.annotations
}

func (i node) ModuleName() string {
	return i.n.ModuleName
}
