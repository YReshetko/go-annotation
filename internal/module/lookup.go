package module

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

func loadModule(path string) (module, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return module{}, fmt.Errorf("unable to define absolut path %s: %w", path, err)
	}
	var modFile string
	var goFiles []string
	err = filepath.Walk(path, func(filePath string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		switch filepath.Ext(info.Name()) {
		case ".go":
			p, err := filepath.Rel(path, filePath)
			if err != nil {
				return err
			}
			goFiles = append(goFiles, p)
		case ".mod":
			if len(modFile) == 0 {
				modFile = filePath
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	modeFile, err := loadModSpec(modFile)
	if err != nil {
		return newModule(path, nil, goFiles), nil
	}

	return newModule(path, modeFile, goFiles), nil
}
