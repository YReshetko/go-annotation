package environment

import (
	"path/filepath"
)

const (
	modSubPath = "/pkg/mod"
	srcSubPath = "/src"
	version    = "0.1.0"
)

var enironment *env

func GoPath() string {
	return enironment.GoPath
}

func GoHome() string {
	return filepath.Join(enironment.GoPath, srcSubPath)
}

func ModPath() string {
	return enironment.GoModCache
}

func ToolVersion() string {
	return version
}

func GoVersion() string {
	return enironment.GoVersion
}

func GoStdLibs() string {
	return filepath.Join(enironment.GoRoot, srcSubPath)
}

func ProjectRoot() string {
	return enironment.ProjectRoot
}
