package model

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
)

type Validatable struct {
	ValidatorName string
	Name          string
	Fields        []Field
}

type Field struct {
	*Validatable
	Name     string
	FType    string
	Tag      *Tag // TODO prepare a tag type
	TypeNode ast.Node
}

func (f Field) String() string {
	buf := bytes.NewBufferString("")
	_ = printer.Fprint(buf, token.NewFileSet(), f.TypeNode)
	return fmt.Sprintf("{name: %s, type: %s, tag: %s, node: %s}", f.Name, f.FType, f.Tag, buf.String())
}

type TemplateData struct {
	ValidatorName  string
	TargetTypeName string
	Snippets       []string
}

type Tag struct {
	Ignore   bool
	Range    Range
	FuncName string
}

type Range struct {
	Left, Right int
}

func (t *Tag) IsIgnore() bool {
	if t == nil {
		return false
	}
	return t.Ignore
}

func (t *Tag) String() string {
	if t == nil {
		return ""
	}
	return fmt.Sprintf("Tag{ignore:%v, range:%v, func:%s}", t.Ignore, t.Range, t.FuncName)
}
