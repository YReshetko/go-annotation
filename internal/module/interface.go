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
	m, err := loadModule(path)
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

	out, err := mod.find(importPath)
	if err != nil {
		return nil, fmt.Errorf("unable to find module: %w", err)
	}

	if out == (*module)(nil) {
		return nil, nil
	}

	return out, nil
}
