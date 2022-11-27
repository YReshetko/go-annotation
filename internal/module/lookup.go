package module

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-multierror"

	"github.com/YReshetko/go-annotation/internal/environment"
	"github.com/YReshetko/go-annotation/internal/logger"
)

var _ moduleLookup = (*selfLookup)(nil)
var _ moduleLookup = (*subLookup)(nil)
var _ moduleLookup = (*dependencyLookup)(nil)
var _ moduleLookup = (*stdLookup)(nil)
var _ moduleLookup = (*workspaceLookup)(nil)
var _ moduleLookup = (*vendorLookup)(nil)
var _ moduleLookup = (*chainLookup)(nil)

var lookup moduleLookup

func init() {
	cl := &chainLookup{}
	cl.add(selfLookup{})
	cl.add(subLookup{})
	//cl.add(workspaceLookup{})
	cl.add(dependencyLookup{})
	cl.add(stdLookup{})
	//cl.add(vendorLookup{})

	lookup = cl
}

type moduleLookup interface {
	find(m *module, importPath string) (*module, error)
}

type selfLookup struct{}

func (s selfLookup) find(m *module, importPath string) (*module, error) {
	if !m.hasModFile() {
		if m.hasImportPath(importPath) {
			return m, nil
		}
		return nil, nil
	}

	if m.isSubModule(importPath) {
		return nil, nil
	}

	if strings.Contains(importPath, moduleName(m.mod)) {
		return m, nil
	}
	return nil, nil
}

type subLookup struct{}

func (s subLookup) find(m *module, importPath string) (*module, error) {
	modPath := m.subModuleRoot(importPath)
	if len(modPath) == 0 {
		return nil, nil
	}

	// TODO Make it cacheable
	mod, err := moduleLoader.load(filepath.Dir(modPath))
	if err != nil {
		return nil, fmt.Errorf("unable to preload submodule by import %s: %w", importPath, err)
	}
	return &mod, nil
}

type dependencyLookup struct{}

func (d dependencyLookup) find(m *module, importPath string) (*module, error) {
	if !m.hasModFile() {
		return nil, nil
	}

	mn := m.findClosestModuleName(importPath)
	if len(mn) == 0 {
		return nil, nil
	}

	// Preload and return module
	path, ok := m.escapedPath(mn)
	if !ok {
		logger.Warnf("closest module was found: %s, but unable to build escaped path", mn)
		return nil, nil
	}
	modulePath := filepath.Join(environment.ModPath(), path)

	// TODO make it cacheable
	subModule, err := moduleLoader.load(modulePath)
	if err != nil {
		return nil, fmt.Errorf("unable to preload submodule %s, due to: %w", importPath, err)
	}

	return &subModule, nil
}

type stdLookup struct{}

func (s stdLookup) find(_ *module, importPath string) (*module, error) {
	nm, err := moduleLoader.load(stdLibKey)
	if err != nil {
		return nil, fmt.Errorf("unable to preload std lib")
	}
	for _, f := range nm.Files() {
		if strings.HasPrefix(f, importPath) {
			return &nm, nil
		}
	}
	return nil, nil
}

type workspaceLookup struct{}

func (w workspaceLookup) find(*module, string) (*module, error) {
	logger.Warn("workspace module lookup is not implemented, yet")
	return nil, nil
}

type vendorLookup struct{}

func (v vendorLookup) find(*module, string) (*module, error) {
	logger.Warn("vendor module lookup is not implemented, yet")
	return nil, nil
}

type chainLookup struct {
	lookup   moduleLookup
	fallback *chainLookup
}

func (c *chainLookup) find(m *module, importPath string) (*module, error) {
	if c == nil || c.lookup == nil {
		return nil, nil
	}

	mod, lookupErr := c.lookup.find(m, importPath)
	if mod != nil {
		return mod, lookupErr
	}

	cmod, chainErr := c.fallback.find(m, importPath)
	err := multierror.Append(lookupErr, chainErr)
	return cmod, err.ErrorOrNil()
}

func (c *chainLookup) add(lookup moduleLookup) {
	if c.lookup == nil {
		c.lookup = lookup
		return
	}

	if c.fallback == nil {
		c.fallback = &chainLookup{}
	}
	c.fallback.add(lookup)
}
