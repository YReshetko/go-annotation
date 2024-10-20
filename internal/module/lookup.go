package module

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-multierror"
	"golang.org/x/tools/go/packages"
)

var _ moduleLookup = (*stdLookup)(nil)
var _ moduleLookup = (*chainLookup)(nil)
var _ moduleLookup = (*packagesLookup)(nil)

var lookup moduleLookup

func init() {
	cl := &chainLookup{}
	cl.add(&packagesLookup{})
	cl.add(&stdLookup{})
	lookup = cl
}

type moduleLookup interface {
	find(m *module, importPath string) (*module, error)
}

type stdLookup struct{}

func (s stdLookup) find(_ *module, importPath string) (*module, error) {
	nm, err := moduleLoader.load(stdLibKey)
	if err != nil {
		return nil, fmt.Errorf("unable to preload std lib")
	}
	importPath = filepath.Clean(importPath)

	for _, f := range nm.Files() {
		if strings.HasPrefix(f, importPath) {
			return &nm, nil
		}
	}
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
	if cmod != nil {
		return cmod, nil
	}
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

type packagesLookup struct{}

func (p packagesLookup) find(_ *module, importPath string) (*module, error) {
	pkgs, err := packages.Load(&packages.Config{
		Mode:    packages.NeedModule,
		Context: context.Background(),
	}, importPath)

	if err != nil {
		return nil, fmt.Errorf("unable to load package for %s import: %w", importPath, err)
	}

	if len(pkgs) == 0 {
		return nil, fmt.Errorf("module for import '%s' not found", importPath)
	}
	if len(pkgs) > 1 {
		return nil, fmt.Errorf("found more than one module for import '%s': %v", importPath, pkgs)
	}

	if pkgs[0].Module == nil {
		return nil, fmt.Errorf("trying to load non-module pakage: %s", importPath)
	}

	subModule, err := moduleLoader.load(pkgs[0].Module.Dir)
	return &subModule, nil
}
