package module

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"

	"github.com/YReshetko/go-annotation/internal/environment"
	"github.com/YReshetko/go-annotation/internal/logger"
)

const stdLibKey = "__std__"

type cacheableLoader[T any] struct {
	cache  map[string]T
	loader func(string) (T, error)
}

var moduleLoader cacheableLoader[module]
var modFileLoader cacheableLoader[*modfile.File]

func init() {
	moduleCache := map[string]module{}
	modFileCache := map[string]*modfile.File{}

	moduleLoader = cacheableLoader[module]{
		cache:  moduleCache,
		loader: loadModule,
	}

	modFileLoader = cacheableLoader[*modfile.File]{
		cache:  modFileCache,
		loader: loadModSpec,
	}
}

func (l cacheableLoader[T]) load(path string) (T, error) {
	t, ok := l.cache[path]
	if ok {
		logger.Debugf("hit cache %s", path)
		return t, nil
	}

	t, err := l.loader(path)
	if err != nil {
		return t, fmt.Errorf("loading failed: %w", err)
	}
	l.cache[path] = t
	return t, nil
}

func loadModule(path string) (module, error) {
	if path == stdLibKey {
		return moduleLoader.load(environment.GoStdLibs())
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return module{}, fmt.Errorf("unable to define absolut path %s: %w", path, err)
	}
	var modFilePath string
	var goFiles []string
	var subModFiles map[string]string
	err = filepath.Walk(path, func(filePath string, info fs.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}
		switch filepath.Ext(info.Name()) {
		case ".go":
			p, err := filepath.Rel(path, filePath)
			if err != nil {
				return err
			}
			goFiles = append(goFiles, p)
		case ".mod":
			if isInRoot(path, filePath) {
				modFilePath = filePath
				return nil
			}
			subModFile, err := modFileLoader.load(filePath)
			if err != nil {
				logger.Errorf("unable to load submodule %s: %s", filePath, err)
				return nil
			}
			subModFiles[filePath] = moduleName(subModFile)
		}
		return nil
	})

	if err != nil {
		return module{}, fmt.Errorf("unable to load module: %w", err)
	}

	if len(modFilePath) == 0 {
		logger.Warnf("loading module %s with no mod file", path)
		return newModule(path, nil, goFiles, subModFiles), nil
	}

	modeFile, err := modFileLoader.load(modFilePath)
	if err != nil {
		logger.Errorf("unable to preload base mod file %s: %s", modFilePath, err)
		return newModule(path, nil, goFiles, subModFiles), nil
	}

	return newModule(path, modeFile, goFiles, subModFiles), nil
}

func isInRoot(root, file string) bool {
	rest := strings.TrimPrefix(file, root)
	return strings.Contains(rest, string(filepath.Separator))
}

func loadModSpec(path string) (*modfile.File, error) {
	if len(path) == 0 {
		return nil, errors.New("go.mod file not found")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to load go.mod file: %w", err)
	}
	f, err := modfile.Parse(path, data, func(_, version string) (string, error) {
		return version, nil
	})
	if err != nil {
		return nil, fmt.Errorf("unable to parse go.mod file: %w", err)
	}
	return f, nil
}
