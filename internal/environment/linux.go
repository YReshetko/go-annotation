//go:build linux

package environment

import (
	"fmt"
	"github.com/YReshetko/go-annotation/internal/utils"
	"os"
	"os/exec"
	"strings"
)

func init() {
	r := os.Getenv("GOROOT")
	if len(r) != 0 {
		goRoot = r
		return
	}
	data, err := exec.Command("go", "env").Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, kv := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(kv, "GOROOT") {
			r = strings.Split(kv, "=")[1]
			goRoot = utils.TrimQuotes(r)
		}
	}

}
