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
	m module.Module // module for the explored node and file
	p string        // Absolute file path
	f *ast2.File    // ast.File related to path
	n ast2.Node     // explored ast.Node
	a []Annotation  // Node annotations if any known
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

func (n *node) ASTNode() ast2.Node {
	return n.n
}

func (n *node) AnnotatedNode(v ast2.Node) Node {
	return annotatedNode(n.m, n.p, n.f, v)
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

func (n *node) FindNodeByAlias(alias, nodeName string) (Node, error) {
	if len(alias) == 0 {
		return n.findNodeInCurrent(nodeName)
	}
	return n.findNodeByAlias(alias, nodeName)
}

func (n *node) findNodeInCurrent(nodeName string) (Node, error) {
	moduleDir := filepath.Dir(strings.TrimPrefix(n.p, n.m.Root()))
	astNode, astFile, filePath, err := lookup.FindTypeInDir(n.m, moduleDir, nodeName)
	if err != nil {
		return nil, fmt.Errorf("unable to find type %s in dir %s: %w", nodeName, moduleDir, err)
	}
	return annotatedNode(n.m, filePath, astFile, astNode), nil
}

func (n *node) findNodeByAlias(alias, nodeName string) (Node, error) {
	importPath, ok := n.FindImportByAlias(alias)
	if !ok {
		return nil, fmt.Errorf("import not found for alias %s", alias)
	}

	relatedModule, err := module.Find(n.m, importPath)
	if err != nil {
		return nil, fmt.Errorf("unable to find module for %s, %s: %w", importPath, nodeName, err)
	}

	astNode, astFile, filePath, err := lookup.FindTypeByImport(relatedModule, importPath, nodeName)
	if err != nil {
		return nil, fmt.Errorf("unable to find type %s by import: %w", nodeName, err)
	}

	return annotatedNode(relatedModule, filePath, astFile, astNode), nil
}

func annotatedNode(m module.Module, p string, f *ast2.File, n ast2.Node) Node {
	a, ok := annotationsByNode(n)
	if !ok {
		return newNode(m, p, f, n, nil)
	}
	return newNode(m, p, f, n, filledAnnotations(a))
}
