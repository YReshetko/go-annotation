package annotation

import (
	ast2 "go/ast"
	"path/filepath"

	"github.com/YReshetko/go-annotation/internal/module"
)

var _ Meta = (*meta)(nil)

type meta struct {
	module   module.Module // module for the explored node and file
	filePath string        // Absolute file path
	fileAST  *ast2.File    // ast.File related to path
}

func (n *meta) Root() string {
	return n.module.Root()
}

func (n *meta) Dir() string {
	return filepath.Dir(n.filePath)
}
func (n *meta) FileName() string {
	return filepath.Base(n.filePath)
}

func (n *meta) PackageName() string {
	if n.fileAST.Name != nil {
		return n.fileAST.Name.Name
	}
	return ""
}
