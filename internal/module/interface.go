package module

import (
	"errors"
	"fmt"
)

type Module interface {
	Root() string
	Files() []string
}

func Load(path string) (Module, error) {
	m, err := lookup(path)
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

	return mod.find(importPath)
}
