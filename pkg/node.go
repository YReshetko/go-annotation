package pkg

import (
	"github.com/YReshetko/go-annotation/internal/lookup"
	ast2 "go/ast"
	"path/filepath"

	"github.com/YReshetko/go-annotation/internal/module"
)

var _ Node = (*node)(nil)

type node struct {
	m module.Module
	p string
	f *ast2.File
	n ast2.Node
	a []Annotation
}

func newNode(m module.Module, path string, f *ast2.File, n ast2.Node, a []Annotation) *node {
	return &node{
		m: m,
		p: path,
		f: f,
		n: n,
		a: a,
	}
}

func (n *node) Annotations() []Annotation {
	return n.a
}

func (n *node) Node() ast2.Node {
	return n.n
}

func (n *node) AnnotatedNode(v ast2.Node) Node {
	a, ok := annotationsByNode(v)
	if !ok {
		return newNode(n.m, n.p, n.f, n.n, nil)
	}
	return newNode(n.m, n.p, n.f, n.n, filledAnnotations(a))
}

func (n *node) Root() string {
	return n.m.Root()
}

func (n *node) Dir() string {
	return filepath.Dir(n.p)
}
func (n *node) FileName() string {
	return filepath.Base(n.p)
}

func (n *node) PackageName() string {
	if n.f.Name != nil {
		return n.f.Name.Name
	}
	return ""
}

func (n *node) Imports() []*ast2.ImportSpec {
	return n.f.Imports
}

func (n *node) FindImportByAlias(alias string) (string, bool) {
	return lookup.FindImportByAlias(n.m, n.f, alias)
}
