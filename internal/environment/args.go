package environment

import "os"

type Arguments struct {
	Args []string

	GenToolPath string
	ProjectPath string

	GoArch    string
	GoOS      string
	GoFile    string
	GoLine    string
	GoPackage string
	GoPath    string
	GoBin     string
}

func LoadArguments() Arguments {
	if len(os.Args) < 2 {
		panic("Requires project path to process annotations")
	}
	args := make([]string, len(os.Args))
	for i, v := range os.Args {
		args[i] = v
	}
	return Arguments{
		Args:        args,
		GenToolPath: args[0],
		ProjectPath: args[1],
		GoArch:      os.Getenv("GOARCH"),
		GoOS:        os.Getenv("GOOS"),
		GoFile:      os.Getenv("GOFILE"),
		GoLine:      os.Getenv("GOLINE"),
		GoPackage:   os.Getenv("GOPACKAGE"),
		GoPath:      os.Getenv("GOPATH"),
		GoBin:       os.Getenv("GOBIN"),
	}
}
