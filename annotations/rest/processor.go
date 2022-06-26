package rest

import (
	"fmt"
	"go/ast"
	"net/http"

	annotation "github.com/YReshetko/go-annotation/pkg"
)

type Processor struct {
	mapping map[handlerMetadata]handlerMapping
}

func init() {
	annotation.Register(Rest{}, &Processor{
		mapping: make(map[handlerMetadata]handlerMapping),
	})
}

func (p *Processor) Process(an annotation.Annotation, node annotation.Node) error {
	a := annotation.CastAnnotation[Rest](an)
	switch node.NodeType() {
	case annotation.Structure:
		return p.processStructure(a, node)
	case annotation.Method:
		return p.processMethod(a, node)
	default:
		panic("Rest annotation can be used for structure or function only")
	}
}

type handlerMetadata struct {
	pkg        string
	structName string
	dir        string
	fileName   string
}

type handlerMapping struct {
	path    string
	mapping map[string]string
}

func (p *Processor) processStructure(rest Rest, node annotation.Node) error {
	fmt.Println("Structure processing: ", rest, node)
	key := handlerMetadata{
		pkg:        node.FileSpec().Name.Name,
		structName: node.Name(),
		dir:        node.Dir(),
		fileName:   node.FileName(),
	}

	_, ok := p.mapping[key]
	if !ok {
		p.mapping[key] = handlerMapping{
			path:    rest.Path,
			mapping: make(map[string]string),
		}
	}

	return nil
}

func (p *Processor) processMethod(rest Rest, node annotation.Node) error {
	fmt.Println("Function processing: ", rest, node)
	if !p.validateHTTPMethod(rest.Method) {
		return fmt.Errorf("invalid HTTP method: %s", rest.Method)
	}

	fnNode, ok := node.GoNode().(*ast.FuncDecl)
	if !ok {
		return fmt.Errorf("expected ast.FuncDecl node, but got %T", node.GoNode())
	}

	recvName := annotation.MethodReceiver(fnNode)
	if recvName == "" {
		return fmt.Errorf("expected method receiver, but got empty for %s", node.Name())
	}

	key := handlerMetadata{
		pkg:        node.FileSpec().Name.Name,
		structName: recvName,
		dir:        node.Dir(),
		fileName:   node.FileName(),
	}

	v, ok := p.mapping[key]
	if !ok {
		return fmt.Errorf("no mapping for %s", key)
	}

	v.mapping[rest.Method] = fnNode.Name.Name

	return nil
}

func (p *Processor) validateHTTPMethod(m string) bool {
	_, ok := map[string]struct{}{
		http.MethodGet:     {},
		http.MethodHead:    {},
		http.MethodPost:    {},
		http.MethodPut:     {},
		http.MethodPatch:   {},
		http.MethodDelete:  {},
		http.MethodConnect: {},
		http.MethodOptions: {},
		http.MethodTrace:   {},
	}[m]
	return ok
}

func (p *Processor) Output() map[annotation.Path]annotation.Data {
	fmt.Println(p.mapping)
	o, err := newOutput()
	if err != nil {
		panic(err)
	}

	for k, v := range p.mapping {
		err = o.append(k, v)
		if err != nil {
			panic(err)
		}
	}

	return o.get()
}

func (p *Processor) SetLookup(lookup annotation.Lookup) {}
