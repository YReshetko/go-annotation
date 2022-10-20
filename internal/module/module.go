package module

import (
	"errors"
	"fmt"
	"github.com/YReshetko/go-annotation/internal/environment"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
	module2 "golang.org/x/mod/module"
)

type module struct {
	root       string
	files      []string
	mod        *modfile.File
	subModules map[string]*module
}

func newModule(root string, mod *modfile.File, goFiles []string) module {
	return module{
		root:       root,
		files:      goFiles,
		mod:        mod,
		subModules: map[string]*module{},
	}
}

func (m *module) Files() []string {
	return m.files
}

func (m *module) Root() string {
	return m.root
}

func (m *module) find(importPath string) (*module, error) {
	if m.mod == nil {
		for _, file := range m.files {
			if strings.HasSuffix(importPath, file) {
				return m, nil
			}
		}
		return nil, nil
	}

	mn := m.findClosestModuleName(importPath)
	if len(mn) == 0 {
		return nil, errors.New("module not found for: " + importPath)
	}

	// Return self module
	if mn == moduleName(m.mod) {
		return m, nil
	}

	// Return already loaded module
	if o, ok := m.subModules[mn]; ok {
		return o, nil
	}

	// Preload and return module
	path, ok := m.escapedPath(mn)
	if !ok {
		return nil, errors.New("unable to build submodule path for: " + importPath)
	}
	modulePath := filepath.Join(environment.ModPath(), path)
	subModule, err := lookup(modulePath)
	if err != nil {
		return nil, fmt.Errorf("unable to preload submodule %s, due to: %w", importPath, err)
	}

	m.subModules[mn] = &subModule

	return &subModule, nil

}

func (m *module) findClosestModuleName(importPath string) string {
	var out []string

	fn := func(pmn string) {
		if len(pmn) > 0 && strings.HasPrefix(importPath, pmn) {
			out = append(out, pmn)
		}
	}

	fn(moduleName(m.mod))

	for mn, _ := range m.subModules {
		fn(mn)
	}

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
				fn(mn)
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
		return "", false
	}
	line := findModuleLine(m.mod.Require, moduleName)
	if line == nil {
		return "", false
	}

	if len(line.Token) < 2 {
		return "", false
	}
	p, err := module2.EscapePath(line.Token[0])
	if err != nil {
		return "", false
	}
	v, err := module2.EscapeVersion(line.Token[1])
	if err != nil {
		return "", false
	}
	return p + "@" + v, true
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
