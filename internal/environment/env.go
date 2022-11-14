package environment

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/YReshetko/go-annotation/internal/logger"
)

type env struct {
	GoBin       string `json:"GOBIN"`
	GoMod       string `json:"GOMOD"`
	GoModCache  string `json:"GOMODCACHE"`
	GoRoot      string `json:"GOROOT"`
	GoPath      string `json:"GOPATH"`
	GoVersion   string `json:"GOVERSION"`
	GoWorkspace string `json:"GOWORK"` // TODO implement module lookup by workspace
	ProjectRoot string
}

var debug = map[string]logger.Level{
	"-v": logger.DebugLvl,
	"-f": logger.FatalLvl,
	"-i": logger.InfoLvl,
	"-w": logger.WarnLvl,
	"-e": logger.ErrorLvl,
}

func init() {
	enironment = &env{}
	initArgs()
	initEnv()
}
func initArgs() {
	if len(os.Args) < 2 {
		panic("Requires project path to process annotations")
	}
	args := make([]string, len(os.Args))
	for i, v := range os.Args {
		args[i] = v
	}
	enironment.ProjectRoot = args[1]

	for i := 2; i < len(args); i++ {
		if lvl, ok := debug[args[i]]; ok {
			logger.LogLevel(lvl)
			break
		}
	}
	logger.Debug("app is configured")
}

func initEnv() {
	cmd := exec.Command("go", "env", "-json")
	cmd.Dir = enironment.ProjectRoot
	data, err := cmd.Output()
	if err != nil {
		logger.Errorf("unable to get environment variables: %w", err)
	} else {
		err = json.Unmarshal(data, enironment)
		if err != nil {
			logger.Errorf("unable to unmarshal environment variables: %w", err)
		}
	}

	if len(enironment.GoRoot) > 0 {
		logger.Debugf("env is preloaded %v", enironment)
		return
	}

	r := os.Getenv("GOROOT")
	if len(r) != 0 {
		enironment.GoRoot = r
	}

	p := os.Getenv("GOPATH")
	if len(p) != 0 {
		enironment.GoPath = p
		enironment.GoModCache = filepath.Join(p, modSubPath)
	}
}
