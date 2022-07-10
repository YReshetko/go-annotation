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

	for _, n := range node.InnerNodes() {
		var ann []Mapping
		for _, annotation := range n.Annotations() {
			m, ok := annotation.(Mapping)
			if !ok {
				continue
			}
			ann = append(ann, m)
		}
		fmt.Println("Annotations:", ann)
		fs, err := pkg.FunctionSignature(n.GoNode())
		if err != nil {
			return err
		}

		resBuilders := p.buildersByFields("res", fs.Results, node)
		paramBuilders := p.buildersByFields("val", fs.Params, node)
		fmt.Println("================ Result ==================")
		for _, v := range resBuilders {
			for _, builder := range v {
				builder.print("")
			}
		}
		fmt.Println("================ Fields ==================")
		for _, v := range paramBuilders {
			for _, builder := range v {
				builder.print("")
			}
		}

		fmt.Println("================ Search ===================")
		for _, mapping := range ann {
			for k, builders := range resBuilders {
				fmt.Println(builders[0].find(k + "." + mapping.Target))
			}
			for k, builders := range paramBuilders {
				fmt.Println(builders[0].find(k + "." + mapping.Source))
			}
		}
		fmt.Println("================ Lines ===================")
		bc := BuildContext{
			imports: map[string]string{},
			lines:   []string{},
		}
		for _, v := range resBuilders {
			for _, builder := range v {
				builder.mapping("", &bc, paramBuilders, ann)
			}
		}
		for a, p := range bc.imports {
			fmt.Println(p, a)
		}
		for _, line := range bc.lines {
			fmt.Println(line)
		}

	}
	return nil
}

func (p *Processor) buildersByFields(namePrefix string, fields []pkg.NodeField, node pkg.Node) map[string][]ResultBuilder {
	builders := make(map[string][]ResultBuilder)
	for i, field := range fields {
		if field.FieldType != pkg.StructureFieldType && field.FieldType != pkg.SelectorFieldType {
			// TODO extend to other types
			continue
		}
		name := field.Name
		if len(name) == 0 {
			field.Name = fmt.Sprintf("%s_%d", namePrefix, i)
			name = field.Name
		}

		builders[name] = p.processParam(field, node)
	}
	return builders
}

func (p *Processor) fieldsList(n pkg.Node) []ResultBuilder {
	//out := []string{}
	out := []ResultBuilder{}
	ast.Inspect(n.GoNode(), func(node ast.Node) bool {
		switch v := node.(type) {
		case *ast.Field:
			params := pkg.Param(v)
			for _, param := range params {
				out = append(out, p.processParam(param, n)...)
			}
		}
		return true
	})
	return out
}

func (p *Processor) processParam(param pkg.NodeField, n pkg.Node) []ResultBuilder {
	switch param.FieldType {
	case pkg.BasicFieldType, pkg.InterfaceFieldType:
		return []ResultBuilder{
			{
				typeName:    param.TypeName,
				fieldName:   param.Name,
				builderType: Value,
				isPointer:   param.IsPointer,
			},
		}
	case pkg.SelectorFieldType:
		rb := ResultBuilder{
			pkg:         getParamImport(param, n),
			selector:    param.Selector,
			typeName:    param.TypeName,
			fieldName:   param.Name,
			builderType: Instance,
			isPointer:   param.IsPointer,
		}
		s := pkg.Selector{
			TypeName:      rb.typeName,
			PackageImport: rb.pkg,
		}
		externalNode := p.lookup(n, s)
		if externalNode.NodeType() != pkg.Structure {
			return []ResultBuilder{rb}
		}
		extFields := p.fieldsList(externalNode)
		rb.children = extFields
		return []ResultBuilder{rb}
	case pkg.ArrayFieldType:
		return []ResultBuilder{
			{
				builderType: Array,
				isPointer:   param.IsPointer,
				fieldName:   param.Name,
				children:    p.processParam(*param.Value, n),
			},
		}
	case pkg.MapFieldType:
		return []ResultBuilder{
			{
				builderType: Map,
				isPointer:   param.IsPointer,
				fieldName:   param.Name,
				children:    p.processParam(*param.Value, n),
			},
		}
	}
	return nil
}

func getParamImport(field pkg.NodeField, n pkg.Node) string {
	if field.Selector != "" {
		return pkg.FindImport(n.FileSpec(), field.Selector)
	}
	// TODO Move the package calculation to annotation tool Lookup to get rid of metadata exported to processors
	return n.ModuleName() + "/" + n.Dir()

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
