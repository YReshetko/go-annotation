package module

import (
	"fmt"
	"github.com/YReshetko/go-annotation/internal/environment"
)

var stdModule *module

func loadStdModule() *module {
	if stdModule != nil {
		return stdModule
	}

	m, err := Load(environment.GoStdLibs())
	if err != nil {
		fmt.Println("Unable to load std libs: ", environment.GoStdLibs())
	}
	stdModule = m.(*module)
	return stdModule
}
