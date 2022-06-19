package output

import (
	"go/format"
	"os"

	"github.com/YReshetko/go-annotation/internal/environment"
)

func Save(_ environment.Arguments, data map[string][]byte) error {
	for path, content := range data {
		if isGoFile(path) {
			formattedContent, err := format.Source(content)
			if err != nil {
				return err
			}
			content = formattedContent
		}

		if err := os.WriteFile(path, content, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func isGoFile(path string) bool {
	return path[len(path)-3:] == ".go"
}
