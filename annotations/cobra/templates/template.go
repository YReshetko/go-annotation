package templates

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

const (
	InitCommandsTpl = "initCommands"
	MainFileTpl     = "mainFile"
)

var dataTemplate *template.Template

func init() {
	dataTemplate = must(template.New(InitCommandsTpl).Parse(initCommands))
	dataTemplate = must(dataTemplate.New(MainFileTpl).Parse(mainFile))
}

func must[T any](t T, e error) T {
	if e != nil {
		panic(e)
	}
	return t
}

func Execute(templateName string, data any) ([]byte, error) {
	tpl := dataTemplate.Lookup(templateName)
	if tpl == nil {
		return nil, fmt.Errorf("template %s not found", templateName)
	}

	return ExecuteTemplate(tpl, data)
}

func ExecuteTemplate(tpl *template.Template, data any) ([]byte, error) {
	b := bytes.NewBufferString("")
	err := tpl.Execute(b, data)
	if err != nil {
		return nil, fmt.Errorf("unable to process template %s: %w", tpl.Name(), err)
	}
	scanner := bufio.NewScanner(b)
	out := bytes.NewBufferString("")
	for ok := scanner.Scan(); ok; ok = scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		out.WriteString(line)
		out.WriteString("\n")
	}
	return out.Bytes(), nil
}
