package output

import (
	"github.com/YReshetko/go-annotation/internal/environment"
	"os"
)

func Save(_ environment.Arguments, data map[string][]byte) error {
	for path, content := range data {
		if err := os.WriteFile(path, content, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}
