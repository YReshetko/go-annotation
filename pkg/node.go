package annotation

import (
	"fmt"
	"github.com/YReshetko/go-annotation/internal/lookup"
	ast2 "go/ast"
	"path/filepath"
	"strings"

	"github.com/YReshetko/go-annotation/internal/module"
)

var _ Node = (*node)(nil)

type node struct {
	m       module.Module // module for the explored node and file
	p       string        // Absolute file path
	f       *ast2.File    // ast.File related to path
	n       ast2.Node     // explored ast.Node
	parents []ast2.Node   // explored parent nodes for above node
	a       []Annotation  // Node annotations if any known
}

func newNode(m module.Module, path string, f *ast2.File, n ast2.Node, parents []ast2.Node, a []Annotation) *node {
	return &node{
		m:       m,
		p:       path,
		f:       f,
		n:       n,
		parents: parents,
		a:       a,
	}
}

func (n *node) Annotations() []Annotation {
	return n.a
}

func (n *node) ASTNode() ast2.Node {
	return n.n
}

func (n *node) AnnotatedNode(v ast2.Node) Node {
	// Replacing parents from original node just because annotation processor already knows found parens for this (v ast.Node) node
	return annotatedNode(n.m, n.p, n.f, v, n.parents)
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
	return annotatedNode(n.m, n.p, n.f, nn, np), true
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

func (n *node) IsSamePackage(v Node) bool {
	return n.Root() == v.Root() && n.Dir() == v.Dir() && n.PackageName() == v.PackageName()
}

func (n *node) Imports() []*ast2.ImportSpec {
	return n.f.Imports
}

func (n *node) FindImportByAlias(alias string) (string, bool) {
	return lookup.FindImportByAlias(n.m, n.f, alias)
}

func (n *node) FindNodeByAlias(alias, nodeName string) (Node, string, error) {
	if len(alias) == 0 {
		n, err := n.findNodeInCurrent(nodeName)
		return n, "", err
	}
	return n.findNodeByAlias(alias, nodeName)
}

func (n *node) findNodeInCurrent(nodeName string) (Node, error) {
	moduleDir := filepath.Dir(strings.TrimPrefix(n.p, n.m.Root()))
	astNode, parents, astFile, filePath, err := lookup.FindTypeInDir(n.m, trimLeadingSlash(moduleDir), nodeName)
	if err != nil {
		return nil, fmt.Errorf(`unable to find type "%s"" in dir %s: %w`, nodeName, moduleDir, err)
	}
	return annotatedNode(n.m, filePath, astFile, astNode, parents), nil
}

func (n *node) findNodeByAlias(alias, nodeName string) (Node, string, error) {
	importPath, ok := n.FindImportByAlias(alias)
	if !ok {
		return nil, "", fmt.Errorf("import not found for alias %s", alias)
	}

	relatedModule, err := module.Find(n.m, importPath)
	if err != nil {
		return nil, "", fmt.Errorf(`unable to find module for "%s", %s: %w`, importPath, nodeName, err)
	}

	astNode, parents, astFile, filePath, err := lookup.FindTypeByImport(relatedModule, importPath, nodeName)
	if err != nil {
		return nil, "", fmt.Errorf(`unable to find type "%s" by import: %w`, nodeName, err)
	}

	return annotatedNode(relatedModule, filePath, astFile, astNode, parents), importPath, nil
}

func annotatedNode(m module.Module, p string, f *ast2.File, n ast2.Node, parents []ast2.Node) Node {
	a, ok := annotationsByNode(n)
	if !ok {
		return newNode(m, p, f, n, parents, nil)
	}
	return newNode(m, p, f, n, parents, filledAnnotations(a))
}

func trimLeadingSlash(s string) string {
	if s[0] == filepath.Separator {
		return s[1:]
	}
	return s
}
