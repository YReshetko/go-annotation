package module

import (
	"github.com/YReshetko/go-annotation/internal/environment"
	"github.com/YReshetko/go-annotation/internal/logger"
	"github.com/YReshetko/go-annotation/internal/utils/imports"
	"strings"

	"golang.org/x/mod/modfile"
	module2 "golang.org/x/mod/module"

	. "github.com/YReshetko/go-annotation/internal/utils/stream"
)

type module struct {
	root  string
	files []string
	mod   *modfile.File
	// file path (absolute path) - module name
	subModFiles map[string]string
}

func newModule(root string, mod *modfile.File, goFiles []string, subModFiles map[string]string) module {
	return module{
		root:        root,
		files:       goFiles,
		mod:         mod,
		subModFiles: subModFiles,
	}
}

func (m *module) Files() []string {
	return m.files
}

func (m *module) Root() string {
	return m.root
}

func (m *module) LocalPackageOf(path string) string {
	return imports.Of(path).Right(imports.Of(environment.GoHome())).String()
}

func (m *module) isFromModCache() bool {
	return strings.Contains(m.root, environment.ModPath())
}

func (m *module) hasModFile() bool {
	return m.mod != nil
}

func (m *module) hasSubModules() bool {
	return len(m.subModFiles) > 0
}

func (m *module) isSubModule(importPath string) bool {
	for _, modName := range m.subModFiles {
		if strings.Contains(importPath, modName) {
			return true
		}
	}
	return false
}

func (m *module) subModuleRoot(importPath string) string {
	var modPath string
	var modName string
	for fp, mn := range m.subModFiles {
		if strings.Contains(importPath, mn) && len(mn) > len(modName) {
			modPath = fp
			modName = mn
		}
	}
	return modPath
}

// Usecase:
// m.root: /home/yury/go/src/github.com/YReshetko/go-annotation/examples/constructor
// m.files: /internal/common/common.go, ...
// importPath: github.com/YReshetko/go-annotation/examples/constructor/internal/common
// result - true as m.root + m.files[i] contains importPath
func (m *module) hasImportPath(importPath string) bool {
	path := OfSlice(m.files).
		Map(joinPath(m.root)).
		Filter(contains(importPath)).
		One()
	return len(path) > 0
}

func (m *module) findClosestModuleName(importPath string) string {
	var out []string
	if m.mod != nil && len(m.mod.Require) > 0 {
		for _, sm := range m.mod.Require {
			if sm.Indirect || sm.Syntax == nil {
				continue
			}
			for _, t := range sm.Syntax.Token {
				mn, _, ok := module2.SplitPathVersion(t)
				if !ok {
					continue
				}
				if len(mn) > 0 && strings.HasPrefix(importPath, mn) {
					out = append(out, mn)
				}
			}

		}
	}
	if len(out) == 0 {
		return ""
	}

	ind := 0
	l := len(out[ind])

	for i := 1; i < len(out); i++ {
		if len(out[i]) > l {
			l = len(out[i])
			ind = i
		}
	}

	return out[ind]
}

func (m *module) escapedPath(moduleName string) (string, bool) {
	if m.mod == nil || len(m.mod.Require) == 0 {
		logger.Debug("escapedPath: mod file nil or no required")
		return "", false
	}
	line := findModuleLine(m.mod.Require, moduleName)
	if line == nil {
		logger.Debug("escapedPath: no module line")
		return "", false
	}

	if len(line.Token) < 2 {
		logger.Debug("escapedPath: invalid line token")
		return "", false
	}
	p, err := escapedPath(line)
	if err != nil {
		logger.Errorf("escapedPath: exc path err on %s: %s", line.Token[0], err.Error())
		return "", false
	}
	v, err := escapedVersion(line)
	if err != nil {
		logger.Errorf("escapedPath: exc version err: %s", err.Error())
		return "", false
	}
	return p + "@" + v, true
}

func escapedPath(l *modfile.Line) (string, error) {
	if l.Token[0] == "require" {
		return module2.EscapePath(l.Token[1])
	}
	return module2.EscapePath(l.Token[0])
}

func escapedVersion(l *modfile.Line) (string, error) {
	if l.Token[0] == "require" {
		return module2.EscapeVersion(l.Token[2])
	}
	return module2.EscapeVersion(l.Token[1])
}

func findModuleLine(r []*modfile.Require, moduleName string) *modfile.Line {
	for _, sm := range r {
		if sm.Indirect || sm.Syntax == nil {
			continue
		}
		for _, t := range sm.Syntax.Token {
			mn, _, ok := module2.SplitPathVersion(t)
			if !ok || mn != moduleName {
				continue
			}
			return sm.Syntax
		}
	}
	return nil
}

func moduleName(mod *modfile.File) string {
	if mod == nil ||
		mod.Module == nil ||
		mod.Module.Syntax == nil ||
		len(mod.Module.Syntax.Token) < 2 {
		return ""
	}
	return mod.Module.Mod.Path
}
