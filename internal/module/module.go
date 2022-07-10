package module

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/YReshetko/go-annotation/internal/debug"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

type Module struct {
	moduleRoot string
	mod        *modfile.File
	// package - []files;
	//for example: if package is github.com/YReshetko/go-annotation/annotations/rest and module is github.com/YReshetko/go-annotation
	//then mapping is /annotations/rest - ["model.go", "data.go"]
	files map[string][]string
}

func (m Module) annotatedNodes() []Node {
	out := []Node{}
	for pkg, fileNames := range m.files {
		for _, fileName := range fileNames {
			path := pkg + "/" + fileName
			fset := token.NewFileSet()
			fileSpec, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				debug.Critical("unable parse %s: %w", path, err)
			}

			ast.Inspect(fileSpec, func(node ast.Node) bool {
				processedNodes, proceed := processNode(node)
				for _, n := range processedNodes {
					if !n.hasAnnotations() {
						continue
					}
					n.Metadata.FileSpec = fileSpec
					n.Metadata.Dir = pkg
					n.Metadata.FileName = fileName
					out = append(out, n)
				}
				return proceed
			})
		}
	}
	return out
}

type importedSubModule struct {
	moduleName string
	version    string
	path       string
}

func (i *importedSubModule) escapedSubModuleDir() string {
	p, _ := module.EscapePath(i.moduleName)
	v, _ := module.EscapeVersion(i.version)
	return p + "@" + v
}

func (m Module) subModuleAndPath(importedPackage string) importedSubModule {
	if strings.HasPrefix(importedPackage, m.moduleName()) {
		if importedPackage == m.moduleName() {
			return importedSubModule{
				moduleName: m.moduleName(),
			}
		}
		path := strings.TrimPrefix(importedPackage, m.moduleName()+"/")
		return importedSubModule{
			moduleName: m.moduleName(),
			path:       path,
		}
	}

	for _, require := range m.mod.Require {
		// TODO select max path because the submodule can have it's own sub modules
		if strings.HasPrefix(importedPackage, require.Mod.Path) {
			if importedPackage == require.Mod.Path {
				return importedSubModule{
					moduleName: require.Mod.Path,
					version:    require.Mod.Version,
				}
			}
			path := strings.TrimPrefix(importedPackage, require.Mod.Path+"/")
			return importedSubModule{
				moduleName: require.Mod.Path,
				version:    require.Mod.Version,
				path:       path,
			}
		}
	}
	debug.Critical("submodule not found for %s", importedPackage)
	return importedSubModule{}
}

func (m Module) moduleName() string {
	if m.mod == nil ||
		m.mod.Module == nil ||
		m.mod.Module.Syntax == nil ||
		len(m.mod.Module.Syntax.Token) < 2 {
		return ""
	}
	return m.mod.Module.Mod.Path
}

func (m Module) findNode(nodePath, typeName string) Node {
	var out Node
	var found bool
	filepath.Walk(m.moduleRoot, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() || !isGoFile(info) || !strings.HasSuffix(path, nodePath+"/"+info.Name()) || found {
			return nil
		}
		fset := token.NewFileSet()
		fileSpec, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			debug.Critical("unable parse %s: %w", path, err)
		}

		ast.Inspect(fileSpec, func(node ast.Node) bool {
			if found {
				return false
			}
			nodes, _ := processNode(node)
			for _, n := range nodes {
				if n.Metadata.Name == typeName {
					out = n
					out.Metadata.FileSpec = fileSpec
					out.Metadata.Dir = nodePath
					out.Metadata.FileName = info.Name()
					found = true
					return false
				}
			}
			return true
		})
		return nil
	})
	return out
}

func modSpec(path string) *modfile.File {
	if len(path) == 0 {
		debug.Warn("go.mod file not found")
		return nil
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		debug.Error("unable to load go.mod file: %w", err)
		return nil
	}
	f, err := modfile.Parse(path, data, func(_, version string) (string, error) {
		return version, nil
	})
	if err != nil {
		debug.Error("unable to parse go.mod file: %w", err)
		return nil
	}
	return f
}

func isGoFile(info fs.FileInfo) bool {
	return filepath.Ext(info.Name()) == ".go"
}
