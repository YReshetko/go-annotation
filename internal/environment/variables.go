package environment

import "os"

const (
	modSubPath = "/pkg/mod"
	version    = "0.0.x"
)

func GoPath() string {
	return os.Getenv("GOPATH")
}

func ModPath() string {
	return GoPath() + modSubPath
}

func ToolVersion() string {
	return version
}
