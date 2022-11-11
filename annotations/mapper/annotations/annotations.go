package annotations

import (
	"fmt"
	"github.com/YReshetko/go-annotation/annotations/mapper/templates"
)

type Mapper struct {
	Name string `annotation:"name=name,default={{.TypeName}}Impl"`
}

type Mapping struct {
	Target   string `annotation:"name=target"`
	Source   string `annotation:"name=source"`
	This     string `annotation:"name=this"`
	Func     string `annotation:"name=func"`
	Constant string `annotation:"name=const"`
}

type IgnoreDefaultMapping struct{}

func (m Mapper) BuildName(typeName string) (string, error) {
	name, err := templates.ExecuteTemplate(m.Name, map[string]string{"TypeName": typeName})
	if err != nil {
		return "", fmt.Errorf("unable to prepare mapper name for %s: %w", typeName, err)
	}
	return string(name), nil
}
