package module

import (
	"errors"
	"fmt"
	"github.com/YReshetko/go-annotation/internal/logger"
	"path/filepath"
	"strings"

	. "github.com/YReshetko/go-annotation/internal/utils/stream"
)

type Module interface {
	Root() string
	Files() []string
}

func Load(path string) (Module, error) {
	m, err := moduleLoader.load(path)
	if err != nil {
		return nil, fmt.Errorf("unable to load module: %w", err)
	}
	return &m, nil
}

func Find(m Module, importPath string) (Module, error) {
	mod, ok := m.(*module)
	if !ok {
		return nil, errors.New("can not cast module to required internal type")
	}

	out, err := lookup.find(mod, importPath)

	if err != nil {
		return nil, fmt.Errorf("unable to find module: %w", err)
	}

	if out == (*module)(nil) {
		return nil, nil
	}

	return out, nil
}

// FilesInPackage finds all files in importPath for Module
// For example:
// Module.Root - /home/<some-home>/go/src/github.com/YReshetko/go-annotation
// Module.Files - [internal/module/lookup.go...]
// importPath - github.com/YReshetko/go-annotation/internal/module
// Then the function returns all files in internal/module dir with package prefix:
// 		"github.com/YReshetko/go-annotation/internal/module/load.go"
//		"github.com/YReshetko/go-annotation/internal/module/interface.go"
//		"github.com/YReshetko/go-annotation/internal/module/interface_test.go"
//		"github.com/YReshetko/go-annotation/internal/module/lookup.go"
//		"github.com/YReshetko/go-annotation/internal/module/module.go"
func FilesInPackage(m Module, importPath string) []string {
	if m == nil {
		return nil
	}

	nativeModule, ok := m.(*module)
	if !ok {
		logger.Warnf("unable to cast module %s to native module", m.Root())
		return nil
	}

	if nativeModule.isFromModCache() {
		files := OfSlice(m.Files()).
			Filter(isPotentialImport(importPath)).
			Map(joinPath(m.Root())).
			ToSlice()
		//logger.Debugf("found files in cached module: %v", files)
		return files
	}
	return OfSlice(m.Files()).
		Map(joinPath(m.Root())).
		Filter(contains(importPath)).
		Map(trimImportPath(importPath)).
		Filter(hasNoSubPath()).
		Map(joinPath(importPath)).
		ToSlice()
}

// FilesInDir finds all files in a particular directory for Module
// For example:
// Module.Root - /home/<some-home>/go/src/github.com/YReshetko/go-annotation
// Module.Files - [internal/lookup/imports.go...]
// dir - internal/lookup
// Then the function returns all files in internal/lookup dir with no prefixes:
//		"internal/lookup/imports.go"
//		"internal/lookup/imports_test.go"
//		"internal/lookup/types.go"
//		"internal/lookup/types_test.go"
func FilesInDir(m Module, dir string) []string {
	return OfSlice(m.Files()).
		Filter(hasPrefix(dir)).
		Map(trimPrefix(dir)).
		Filter(hasNoSubPath()).
		Map(joinPath(dir)).
		ToSlice()
}

// AbsolutePath checks if filePath exists in module and provides absolute path to the file
// For example:
// Module.Root - /home/<some-home>/go/src/github.com/YReshetko/go-annotation
// Module.Files - [internal/module/lookup.go...]
// filePath - github.com/YReshetko/go-annotation/internal/module/module.go
// return - /home/<some-home>/go/src/github.com/YReshetko/go-annotation/internal/module/module.go
func AbsolutePath(m Module, filePath string) (string, bool) {
	path := OfSlice(m.Files()).
		Filter(hasPostfix(filePath)).
		Map(joinPath(m.Root())).
		One()

	return path, len(path) != 0
}

func trimPrefix(prefix string) func(string) string {
	return func(path string) string {
		s := strings.TrimPrefix(path, prefix)
		if s[0] == filepath.Separator {
			return s[1:]
		}
		return s
	}
}

func hasPrefix(prefix string) func(string) bool {
	return func(path string) bool {
		return strings.HasPrefix(path, prefix)
	}
}

func hasPostfix(filePath string) func(string) bool {
	return func(subPath string) bool {
		return strings.HasSuffix(filePath, subPath)
	}
}

func joinPath(importPath string) func(string) string {
	return func(fileName string) string {
		return filepath.Join(importPath, fileName)
	}
}

func hasNoSubPath() func(string) bool {
	return func(s string) bool {
		return !strings.Contains(s, "/")
	}
}

func trimImportPath(importPath string) func(string) string {
	return func(path string) string {
		ind := strings.Index(path, importPath)
		return path[ind+len(importPath)+1:]
	}
}

func contains(importPath string) func(string) bool {
	return func(file string) bool {
		return strings.Contains(file, importPath)
	}
}

func isPotentialImport(importPath string) func(string) bool {
	return func(file string) bool {
		dir := filepath.Dir(file)
		if dir == "." {
			return true
		}
		return strings.HasSuffix(importPath, dir)
	}
}
