package environment

import (
	"os"
	"path/filepath"
)

const (
	modSubPath = "/pkg/mod"
	version    = "0.0.2-alpha"
)

func GoPath() string {
	return os.Getenv("GOPATH")
}

func ModPath() string {
	return filepath.Join(GoPath(), modSubPath)
}

func ToolVersion() string {
	return version
}
