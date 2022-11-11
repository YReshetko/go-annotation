package environment

import (
	"os"
	"path/filepath"
)

const (
	modSubPath    = "/pkg/mod"
	goHomeSubPath = "/src"
	version       = "0.0.3-alpha"
)

func GoPath() string {
	return os.Getenv("GOPATH")
}

func GoHome() string {
	return filepath.Join(GoPath(), goHomeSubPath)
}

func ModPath() string {
	return filepath.Join(GoPath(), modSubPath)
}

func ToolVersion() string {
	return version
}
