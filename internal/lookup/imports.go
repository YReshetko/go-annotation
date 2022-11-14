package lookup

import (
	"go/ast"
	"path/filepath"
	"strings"

	ast2 "github.com/YReshetko/go-annotation/internal/ast"
	"github.com/YReshetko/go-annotation/internal/logger"
	"github.com/YReshetko/go-annotation/internal/module"
	"github.com/YReshetko/go-annotation/internal/utils"
	. "github.com/YReshetko/go-annotation/internal/utils/stream"
)

func FindImportByAlias(m module.Module, file *ast.File, alias string) (string, bool) {
	for _, spec := range file.Imports {
		if alias == getLocalPackageName(m, spec) {
			return utils.TrimQuotes(spec.Path.Value), true
		}
	}
	return "", false
}

func getLocalPackageName(m module.Module, spec *ast.ImportSpec) string {
	if spec.Name != nil {
		return spec.Name.String()
	}

	if spec.Path == nil {
		return ""
	}

	importPath := utils.TrimQuotes(spec.Path.Value)
	m, err := module.Find(m, importPath)
	if err != nil {
		logger.Warnf("unable to load module for import \"%s\", returning it as is: %w", importPath, err)
		return importPath
	}

	if m != nil {
		// TODO review the assamption that the first file is a good for check
		imp := OfSlice(m.Files()).
			//Map(utils.Root).
			Filter(hasPathSuffix(importPath)).
			Map(fileToPackageName(m.Root())).
			One()
		if strings.HasSuffix(imp, "_test") {
			imp = imp[:strings.Index(imp, "_test")]
		}
		return imp
	}

	return strings.ReplaceAll(utils.LastDir(importPath), "-", "_")
}

func fileToPackageName(root string) func(string) string {
	return func(file string) string {
		fast, err := ast2.LoadFileAST(filepath.Join(root, file))
		if err != nil {
			panic(err)
		}
		return fast.Name.Name
	}
}

func hasPathSuffix(path string) func(string) bool {
	return func(s string) bool {
		return strings.HasSuffix(path, utils.Root(s))
	}
}
