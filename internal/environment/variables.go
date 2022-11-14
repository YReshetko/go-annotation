package environment

import (
	"path/filepath"
)

const (
	modSubPath    = "/pkg/mod"
	goHomeSubPath = "/src"
	version       = "0.0.4-alpha"
)

var enironment *env

func GoPath() string {
	return enironment.GoPath
}

func GoHome() string {
	return filepath.Join(enironment.GoPath, goHomeSubPath)
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
	return filepath.Join(enironment.GoRoot, goHomeSubPath)
}

func ProjectRoot() string {
	return enironment.ProjectRoot
}
