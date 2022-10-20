package module

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/mod/modfile"
)

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
