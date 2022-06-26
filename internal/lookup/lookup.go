package lookup

import (
	"io/fs"
	"path/filepath"
)

func ModuleStructure(projectPath string) (string, []string) {
	var modFile string
	var goFiles []string
	err := filepath.Walk(projectPath, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if isGoFile(info) {
			goFiles = append(goFiles, path)
			return nil
		}

		if isModFile(info) && len(modFile) == 0 {
			modFile = path
			return nil
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return modFile, goFiles
}

func isGoFile(info fs.FileInfo) bool {
	return !info.IsDir() && len(info.Name()) > 3 && info.Name()[len(info.Name())-3:] == ".go"
}

func isModFile(info fs.FileInfo) bool {
	return !info.IsDir() && info.Name() == "go.mod"
}
