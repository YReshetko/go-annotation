package module

import (
	"go/ast"

	"github.com/YReshetko/go-annotation/internal/annotation"
	"github.com/YReshetko/go-annotation/internal/debug"
)

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

	Selector   selector
	ModuleName string
}

func (n Node) hasAnnotations() bool {
	if len(n.Annotations) > 0 {
		return true
	}
	for _, node := range n.Inner {
		if node.hasAnnotations() {
			return true
		}
	}
	return false
}

func processNode(node ast.Node) ([]Node, bool) {
	switch v := node.(type) {
	case *ast.FuncDecl:
		a, _ := annotation.Parse(v.Doc.Text())
		nodeType := Function
		if v.Recv != nil {
			nodeType = Method
		}
		return []Node{
			{
				Metadata: Metadata{
					Name: v.Name.Name,
					Type: nodeType,
				},
				Annotations: a,
				GoNode:      v,
			},
		}, false

	case *ast.GenDecl:
		nodes := []Node{}
		for _, spec := range v.Specs {
			node := processSpec(v, spec)
			if node != nil {
				nodes = append(nodes, *node)
			}
		}
		return nodes, false

	default:
		debug.Debug("%T, %+v\n", v, v)
	}
	return nil, true
}

func processSpec(n *ast.GenDecl, spec ast.Spec) *Node {
	a, _ := annotation.Parse(n.Doc.Text())
	switch v := spec.(type) {
	case *ast.TypeSpec:
		switch t := v.Type.(type) {
		case *ast.StructType:
			ad, ok := annotation.Parse(v.Doc.Text())
			if ok {
				a = append(a, ad...)
			}
			fields := processFields(t.Fields, Field)
			return &Node{
				Metadata: Metadata{
					Name: v.Name.Name,
					Type: Structure,
				},
				Annotations: a,
				GoNode:      v,
				Inner:       fields,
			}

		case *ast.InterfaceType:
			ad, ok := annotation.Parse(v.Doc.Text())
			if ok {
				a = append(a, ad...)
			}
			fields := processFields(t.Methods, Method)
			return &Node{
				Metadata: Metadata{
					Name: v.Name.Name,
					Type: Interface,
				},
				Annotations: a,
				GoNode:      v,
				Inner:       fields,
			}
		case *ast.ArrayType:
			ad, ok := annotation.Parse(v.Doc.Text())
			if ok {
				a = append(a, ad...)
			}
			return &Node{
				Metadata: Metadata{
					Name: v.Name.Name,
					Type: Type,
				},
				Annotations: a,
				GoNode:      v,
			}
		}
	case *ast.ValueSpec:
		doc := n.Doc.Text()
		if len(doc) == 0 {
			doc = v.Doc.Text()
		}
		ad, ok := annotation.Parse(doc)
		if ok {
			a = append(a, ad...)
		}
		return &Node{
			Metadata: Metadata{
				Name: v.Names[0].Name,
				Type: Variable,
			},
			Annotations: a,
			GoNode:      v,
		}
	}
	return nil
}

func processFields(params *ast.FieldList, nodeType NodeType) []Node {
	if params == nil {
		return nil
	}

	if len(params.List) == 0 {
		return nil
	}

	nodes := []Node{}
	for _, field := range params.List {
		a, _ := annotation.Parse(field.Doc.Text())
		nodes = append(nodes, Node{
			Metadata: Metadata{
				Name: field.Names[0].Name,
				Type: nodeType,
			},
			Annotations: a,
			GoNode:      field,
		})
	}
	return nodes
}
