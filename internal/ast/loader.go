package ast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func load(path string) (*ast.File, error) {
	fset := token.NewFileSet()
	fileSpec, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("unable to parse go file AST: %w", err)
	}
	return fileSpec, nil
}
