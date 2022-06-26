package mapper

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/YReshetko/go-annotation/pkg"
)

type Processor struct {
	lookup pkg.Lookup
}

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
		fs, err := pkg.FunctionSignature(n.GoNode())
		if err != nil {
			return err
		}

		for _, param := range fs.Params {
			if param.Selector != "" {
				packageImport := findImport(node.FileSpec(), param.Selector)
				externalNode := p.lookup(node, pkg.Selector{
					PackageImport: packageImport,
					TypeName:      param.TypeName,
				})

				fmt.Printf("External node: %+v\n", externalNode)
			}
		}

		fmt.Printf("Function signature: %+v\n", fs)

		/*		for _, a2 := range n.Annotations() {
				fmt.Println("internal Mapping:", pkg.CastAnnotation[Mapping](a2), n.NodeType(), len(n.InnerNodes()))
			}*/

	}
	return nil
}

func findImport(f *ast.File, alias string) string {
	var found bool
	var out string
	ast.Inspect(f, func(node ast.Node) bool {
		if found {
			return false
		}
		imp, ok := node.(*ast.ImportSpec)
		if !ok {
			return true
		}
		impPath := unquote(imp.Path.Value)
		if imp.Name != nil && imp.Name.Name == alias {
			found = true
			out = impPath
			return false
		}
		if strings.HasSuffix(impPath, "/"+alias) || strings.HasSuffix(strings.ReplaceAll(impPath, "-", "_"), "/"+alias) {
			found = true
			out = impPath
			return false
		}
		return true
	})

	return out
}

func unquote(s string) string {
	out := strings.TrimSpace(s)
	if out[0] == '"' && out[len(out)-1] == '"' {
		return out[1 : len(out)-1]
	}
	return out
}

func (p *Processor) processMapping(m Mapping, node pkg.Node) error {
	// fmt.Println("Mapping: ", m.Source, m.Target, node.NodeType(), len(node.InnerNodes()))
	return nil
}

func (p *Processor) Output() map[pkg.Path]pkg.Data {
	return nil
}
func (p *Processor) SetLookup(lookup pkg.Lookup) {
	p.lookup = lookup
}
