package environment

import (
	"os"
)

var (
	projectPath = ""
)

func init() {
	if len(os.Args) < 2 {
		panic("Requires project path to process annotations")
	}
	args := make([]string, len(os.Args))
	for i, v := range os.Args {
		args[i] = v
	}
	/*	for _, env := range os.Environ() {
		fmt.Println(env)
	}*/

	projectPath = args[1]
}

func ProjectRoot() string {
	return projectPath
}
