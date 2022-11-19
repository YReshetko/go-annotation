package annotation

import (
	"github.com/YReshetko/go-annotation/internal/module"
	ast2 "go/ast"
)

var _ Node = (*node)(nil)

type node struct {
	nodeAST     ast2.Node    // explored ast.Node
	parents     []ast2.Node  // explored parent nodes for above node
	annotations []Annotation // Node annotations if any known

	meta   *meta
	lookup *nodeLookup
}

func newNode(loadedModule module.Module, path string, fileAST *ast2.File, nodeAST ast2.Node, parents []ast2.Node, annotations []Annotation) *node {
	return &node{
		nodeAST:     nodeAST,
		parents:     parents,
		annotations: annotations,
		lookup: &nodeLookup{
			module:   loadedModule,
			filePath: path,
			fileAST:  fileAST,
		},
		meta: &meta{
			module:   loadedModule,
			filePath: path,
			fileAST:  fileAST,
		},
	}
}

func (n *node) Annotations() []Annotation {
	return n.annotations
}

func (n *node) ASTNode() ast2.Node {
	return n.nodeAST
}

func (n *node) AnnotatedNode(v ast2.Node) Node {
	// Replacing parents from original node just because annotation processor already knows found parens for this (v ast.Node) node
	return annotatedNode(n.meta.module, n.meta.filePath, n.meta.fileAST, v, n.parents)
}

func (n *node) ParentNode() (Node, bool) {
	if len(n.parents) == 0 {
		return nil, false
	}
	nn := n.parents[len(n.parents)-1]
	var np []ast2.Node
	if len(n.parents) > 1 {
		np = n.parents[:len(n.parents)-2]
	}
	return annotatedNode(n.meta.module, n.meta.filePath, n.meta.fileAST, nn, np), true
}

func (n *node) IsSamePackage(v Node) bool {
	nv, ok := v.(*node)
	if !ok {
		return false
	}
	return n.meta.Root() == nv.meta.Root() && n.meta.Dir() == nv.meta.Dir() && n.meta.PackageName() == nv.meta.PackageName()
}

func (n *node) Imports() []*ast2.ImportSpec {
	return n.meta.fileAST.Imports
}

func (n *node) Lookup() Lookup {
	return n.lookup
}

func (n *node) Meta() Meta {
	return n.meta
}

func annotatedNode(m module.Module, p string, f *ast2.File, n ast2.Node, parents []ast2.Node) Node {
	a, ok := annotationsByNode(n)
	if !ok {
		return newNode(m, p, f, n, parents, nil)
	}
	return newNode(m, p, f, n, parents, filledAnnotations(a))
}
