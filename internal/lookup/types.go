package lookup

import (
	"fmt"
	"github.com/YReshetko/go-annotation/internal/ast"
	"github.com/YReshetko/go-annotation/internal/module"
	. "github.com/YReshetko/go-annotation/internal/utils/stream"
	ast2 "go/ast"
)

// FindTypeByImport search for type/function/method (typeName) that declared in module.Module and located in importPath
func FindTypeByImport(m module.Module, importPath, typeName string) (ast2.Node, *ast2.File, string, error) {
	filePaths := OfSlice(module.FilesInPackage(m, importPath)).
		Map(toAbsolutePath(m)).
		Filter(IsNotEmpty[string]).
		ToSlice()

	if len(filePaths) == 0 {
		return nil, nil, "", fmt.Errorf("files for import %s not found", importPath)
	}

	return findType(filePaths, typeName)
}

// FindTypeInDir search for type/function/method (typeName) that declared in module.Module and located in dir
func FindTypeInDir(m module.Module, dir, typeName string) (ast2.Node, *ast2.File, string, error) {
	filePaths := OfSlice(module.FilesInDir(m, dir)).
		Map(toAbsolutePath(m)).
		Filter(IsNotEmpty[string]).
		ToSlice()

	if len(filePaths) == 0 {
		return nil, nil, "", fmt.Errorf("files for dir %s not found", dir)
	}

	return findType(filePaths, typeName)
}

func findType(filePaths []string, typeName string) (ast2.Node, *ast2.File, string, error) {
	for _, path := range filePaths {
		astFile, err := ast.LoadFileAST(path)
		if err != nil {
			return nil, nil, "", fmt.Errorf("unable to load ast.File for %s: %w", path, err)
		}
		node, ok := ast.FindTopNodeByName(astFile, typeName)
		if ok {
			return node, astFile, path, nil
		}
	}

	return nil, nil, "", fmt.Errorf("type %s not found", typeName)
}

func toAbsolutePath(m module.Module) func(string) string {
	return func(path string) string {
		s, _ := module.AbsolutePath(m, path)
		return s
	}
}
