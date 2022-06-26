package module

import (
	"fmt"
	"strings"

	"github.com/YReshetko/go-annotation/internal/debug"
	"github.com/YReshetko/go-annotation/internal/environment"
	"github.com/pkg/errors"
)

type Loader struct {
	args         environment.Arguments
	moduleLoader func(string) (string, []string)

	tree *tree
}

func NewLoader(args environment.Arguments, loader func(string) (string, []string)) (*Loader, selector, error) {
	l := &Loader{
		args:         args,
		moduleLoader: loader,
		tree:         newTree(),
	}
	rootModule := l.load(args.ProjectPath)
	s := selector{pkg: pkg(rootModule.moduleName())}
	err := l.tree.put(selector{pkg: pkg(rootModule.moduleName())}, &rootModule)
	return l, s, errors.Wrap(err, "unable to prepare root module")
}

func (l *Loader) AllAnnotatedNodes(s selector) ([]Node, error) {
	m, ok := l.tree.get(s)
	if !ok {
		return nil, fmt.Errorf("unable to find module:\n %s\n", s.String())
	}
	nodes := m.annotatedNodes()
	for i, _ := range nodes {
		nodes[i].Selector = s
	}

	return nodes, nil
}

func (l *Loader) FindNode(s selector, pkg, typeName string) (Node, error) {
	debug.Debug("selector: %v", s)
	debug.Debug("import: %s", pkg)
	debug.Debug("type: %s", typeName)
	m, ok := l.tree.get(s)
	if !ok {
		return Node{}, fmt.Errorf("unable to find module:\n %s\n", s.String())
	}

	moduleAndPath := m.subModuleAndPath(pkg)
	if pkg == m.moduleName() {
		out := m.findNode(moduleAndPath.path, typeName)
		out.Selector = s
		return out, nil
	}

	newSelector := s.put(pkg)
	nm, ok := l.tree.get(newSelector)
	if ok {
		out := nm.findNode(moduleAndPath.path, typeName)
		out.Selector = newSelector
		return out, nil
	}

	modulePath := l.args.GoPath + "/pkg/mod/" + moduleAndPath.escapedSubModuleDir()
	loadedModule := l.load(modulePath)
	err := l.tree.put(newSelector, &loadedModule)
	if err != nil {
		return Node{}, errors.Wrap(err, "unable to put loaded module to tree")
	}
	out := loadedModule.findNode(moduleAndPath.path, typeName)
	out.Selector = newSelector

	return out, nil
}

// load module metadata by path
func (l *Loader) load(moduleRoot string) Module {
	mod := Module{
		moduleRoot: moduleRoot,
		files:      make(map[string][]string),
	}
	modFile, goFiles := l.moduleLoader(moduleRoot)
	debug.Debug("Loaded mod file:%s", modFile)
	for i, file := range goFiles {
		ind := strings.LastIndex(file, "/")
		pkg := ""
		fileName := file
		if ind != -1 {
			pkg, fileName = file[:ind], file[ind+1:]
		}
		debug.Debug("append go file to module [%d]: map[%s]:%s", i, pkg, fileName)
		mod.files[pkg] = append(mod.files[pkg], fileName)
	}

	mod.mod = modSpec(modFile)
	return mod
}
