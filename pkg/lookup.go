package annotation

import (
	"fmt"
	ast2 "go/ast"
	"path/filepath"
	"strings"

	"github.com/YReshetko/go-annotation/internal/lookup"
	"github.com/YReshetko/go-annotation/internal/module"
)

var _ Lookup = (*nodeLookup)(nil)

type nodeLookup struct {
	module   module.Module // module for the explored node and file
	filePath string        // Absolute file path
	fileAST  *ast2.File    // ast.File related to path
}

func (n *nodeLookup) FindImportByAlias(alias string) (string, bool) {
	return lookup.FindImportByAlias(n.module, n.fileAST, alias)
}

func (n *nodeLookup) FindNodeByAlias(alias, nodeName string) (Node, string, error) {
	if len(alias) == 0 {
		n, err := n.findNodeInCurrent(nodeName)
		return n, "", err
	}
	return n.findNodeByAlias(alias, nodeName)
}

func (n *nodeLookup) findNodeInCurrent(nodeName string) (Node, error) {
	moduleDir := filepath.Dir(strings.TrimPrefix(n.filePath, n.module.Root()))
	astNode, parents, astFile, filePath, err := lookup.FindTypeInDir(n.module, trimLeadingSlash(moduleDir), nodeName)
	if err != nil {
		return nil, fmt.Errorf(`unable to find type "%s"" in dir %s: %w`, nodeName, moduleDir, err)
	}
	return annotatedNode(n.module, filePath, astFile, astNode, parents), nil
}

func (n *nodeLookup) findNodeByAlias(alias, nodeName string) (Node, string, error) {
	importPath, ok := n.FindImportByAlias(alias)
	if !ok {
		return nil, "", fmt.Errorf("import not found for alias %s", alias)
	}

	relatedModule, err := module.Find(n.module, importPath)
	if err != nil {
		return nil, "", fmt.Errorf(`unable to find module for "%s", %s: %w`, importPath, nodeName, err)
	}

	astNode, parents, astFile, filePath, err := lookup.FindTypeByImport(relatedModule, importPath, nodeName)
	if err != nil {
		return nil, "", fmt.Errorf(`unable to find type "%s" by import: %w`, nodeName, err)
	}

	return annotatedNode(relatedModule, filePath, astFile, astNode, parents), importPath, nil
}

func trimLeadingSlash(s string) string {
	if s[0] == filepath.Separator {
		return s[1:]
	}
	return s
}
