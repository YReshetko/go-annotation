package todo

import (
	"bytes"
	"encoding/json"
	"fmt"
	annotation "github.com/YReshetko/go-annotation/pkg"
	"go/ast"
	"go/printer"
	"go/token"
	"path/filepath"
)

func init() {
	annotation.Register[TODO](&Processor{infos: make([]Info, 0, 0)})
}

var _ annotation.AnnotationProcessor = (*Processor)(nil)

type Processor struct {
	root  string
	infos []Info
}

type Info struct {
	Message  string `json:"message"`
	CodeLine string `json:"code_line"`
	TypeDecl string `json:"type_decl"`
}

func (p *Processor) Process(node annotation.Node) error {
	if len(p.root) == 0 {
		p.root = node.Meta().Root()
	}
	ans := annotation.FindAnnotations[TODO](node.Annotations())
	fmt.Println(ans)

	ts, ok := annotation.ParentType[*ast.TypeSpec](node)
	if !ok {
		return nil
	}
	gt, _ := annotation.CastNode[*ast.TypeSpec](ts)
	msg := ""
	for _, todo := range ans {
		msg += todo.Message + "\n"
	}

	f, ok := annotation.CastNode[*ast.Field](node)
	if !ok {
		return nil
	}

	fieldLines := ""
	for _, n := range f.Names {
		typeLine := bytes.NewBufferString("")
		err := printer.Fprint(typeLine, token.NewFileSet(), f.Type)
		if err != nil {
			return err
		}
		fieldLines += n.String() + " " + typeLine.String() + "\n"
	}
	typeDecl := gt.Name.String()

	p.infos = append(p.infos, Info{
		Message:  msg,
		CodeLine: fieldLines,
		TypeDecl: typeDecl,
	})
	return nil
}

func (p *Processor) Output() map[string][]byte {
	data, _ := json.Marshal(p.infos)
	return map[string][]byte{
		filepath.Join(p.root, "todo.list.json"): data,
	}
}

func (p *Processor) Version() string {
	return "0.0.1"
}

func (p *Processor) Name() string {
	return "TODO"
}
