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

	m, err := Load(environment.GoLibs())
	if err != nil {
		fmt.Println("Unable to load std libs: ", environment.GoLibs())
	}
	stdModule = m.(*module)
	return stdModule
}
