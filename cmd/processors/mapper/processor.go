package mapper

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/YReshetko/go-annotation/pkg"
)

type Processor struct{}

func init() {
	p := &Processor{}
	pkg.Register(Mapper{}, p)
	pkg.Register(Mapping{}, p)
}

func (p *Processor) Process(annotation pkg.Annotation, node pkg.Node) error {
	switch v := annotation.(type) {
	case Mapper:
		return p.processMapper(v, node)
	case Mapping:
		return p.processMapping(v, node)
	}
	return nil
}

func (p *Processor) processMapper(m Mapper, node pkg.Node) error {
	if node.NodeType() != pkg.Interface {
		return fmt.Errorf("expected node type Interface for @Mapper, but got: %s", node.NodeType())
	}

	mapperName := strings.ReplaceAll(m.Name, "*", node.Name())
	fmt.Println("New mapper name:", mapperName)

	fmt.Println("Mapper: ", m.Name, node.NodeType(), len(node.InnerNodes()))
	for _, n := range node.InnerNodes() {
		fmt.Printf("%T\n", n.GoNode())
		f, ok := n.GoNode().(*ast.Field)
		if !ok {
			return fmt.Errorf("expected *ast.Field, but got: %T", n.GoNode())
		}
		ft, ok := f.Type.(*ast.FuncType)
		if !ok {
			return fmt.Errorf("expected *ast.Field, but got: %T", f.Type)
		}

		fmt.Println("Function name", f.Names[0].Name, "names:", f.Names)

		for _, v := range ft.Params.List {
			fmt.Println("Params", v.Names)
		}
		/*	for _, v := range ft.TypeParams.List {
			fmt.Println("Type", v.Names)
		}*/
		for _, v := range ft.Results.List {
			fmt.Println("Results", v.Names)
		}

		fmt.Printf("%T\n", f.Type)

		for _, a2 := range n.Annotations() {
			fmt.Println("internal Mapping:", pkg.CastAnnotation[Mapping](a2), n.NodeType(), len(n.InnerNodes()))
		}

	}
	return nil
}

func (p *Processor) processMapping(m Mapping, node pkg.Node) error {
	fmt.Println("Mapping: ", m.Source, m.Target, node.NodeType(), len(node.InnerNodes()))
	return nil
}

func (p *Processor) Output() map[pkg.Path]pkg.Data {
	return nil
}
