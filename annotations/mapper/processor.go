package mapper

import (
	"fmt"
	"github.com/YReshetko/go-annotation/annotations/mapper/annotations"
	"github.com/YReshetko/go-annotation/annotations/mapper/generators"
	"github.com/YReshetko/go-annotation/annotations/mapper/templates"
	annotation "github.com/YReshetko/go-annotation/pkg"
	"go/ast"
	"path/filepath"
)

func init() {
	p := &Processor{cache: map[key][]mapperData{}}
	annotation.Register[annotations.Mapper](p)
	annotation.RegisterNoop[annotations.Mapping]()
	annotation.RegisterNoop[annotations.SliceMapping]()
	annotation.RegisterNoop[annotations.IgnoreDefaultMapping]()
}

var _ annotation.AnnotationProcessor = (*Processor)(nil)

type key struct {
	dir string
	pkg string
}

type Processor struct {
	cache map[key][]mapperData
}

type mapperData struct {
	data    []byte
	imports []generators.Import
}

func (p *Processor) Process(node annotation.Node) error {
	a, ts, err := validateAndGetMapperWithTypeSpec(node)
	if err != nil {
		return err
	}

	if ts == nil {
		return nil
	}

	mapperName, err := a.BuildName(ts.Name.String())
	if err != nil {
		return fmt.Errorf("unable to build mapper name: %w", err)
	}

	mapperGenerator := generators.NewMapperGeneratorBuilder().
		Node(node).
		IntName(ts.Name.String()).
		IntType(ts.Type.(*ast.InterfaceType)).
		StructName(mapperName).
		Build()

	data, imports, err := mapperGenerator.Generate()
	if err != nil {
		return fmt.Errorf("unable to generate mapper for %s: %w", ts.Name.String(), err)
	}

	k := key{
		dir: node.Dir(),
		pkg: node.PackageName(),
	}

	p.cache[k] = append(p.cache[k], mapperData{
		data:    data,
		imports: imports,
	})
	return nil
}

func (p *Processor) Output() map[string][]byte {
	out := map[string][]byte{}

	for k, data := range p.cache {
		var rd []byte
		distinctImports := map[generators.Import]struct{}{}
		for _, d := range data {
			rd = append(rd, d.data...)
			for _, g := range d.imports {
				distinctImports[g] = struct{}{}
			}
		}
		importsSlice := make([]generators.Import, len(distinctImports))
		var ind int
		for imp, _ := range distinctImports {
			importsSlice[ind] = imp
			ind++
		}

		fileData, err := templates.Execute(templates.FileTemplate, map[string]interface{}{
			"PackageName": k.pkg,
			"Data":        string(rd),
			"HasImports":  len(importsSlice) != 0,
			"Imports":     importsSlice,
		})
		if err != nil {
			panic(err)
		}
		out[filepath.Join(k.dir, "mappers.gen.go")] = fileData
	}

	return out
}

func (p *Processor) Version() string {
	return "0.0.1-alpha"
}

func (p *Processor) Name() string {
	return "Mapper"
}
