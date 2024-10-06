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

	//logger.Debugf("module is found for %s", importPath)

	if m != nil {
		filesInPackage := module.FilesInPackage(m, importPath)
		// TODO review the assamption that the first file is a good for check
		imp := OfSlice(m.Files()).
			Filter(containsFile(filesInPackage)).
			Map(fileToPackageName(m.Root())).
			One()
		//logger.Debugf("module is found for %s, checking %s", importPath, imp)
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

func containsFile(files []string) func(string) bool {
	return func(file string) bool {
		for _, f := range files {
			if strings.HasSuffix(f, file) && isCompletePathSuffixMatch(f, file) {
				return true
			}
		}
		return false
	}
}

// The check is required, because if we just check suffix for similar package ends:
// log/example_test.go
// log/slog/example_test.go
// In this case both have equel suffix log/example_test.go, but this is completely different packages
func isCompletePathSuffixMatch(str, suffix string) bool {
	prefix := strings.TrimSuffix(str, suffix)
	if len(prefix) == 0 {
		return true
	}
	return prefix[len(prefix)-1] == '/' || prefix[len(prefix)-1] == '\\'
}
