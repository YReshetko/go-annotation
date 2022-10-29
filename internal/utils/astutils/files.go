package astutils

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func BytesToAST(data []byte) (ast.Node, *token.FileSet, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", data, parser.ParseComments)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to parse go file AST: %w", err)
	}
	return f, fset, nil
}
