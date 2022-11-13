//go:build linux

package environment

import (
	"os"
)

func init() {
	goRoot = os.Getenv("GOROOT")
}
