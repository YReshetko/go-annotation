package constructor

import (
	"fmt"
	"github.com/YReshetko/go-annotation/annotations/constructor/annotations"
	"github.com/YReshetko/go-annotation/annotations/constructor/generators"
	"go/ast"
	"path/filepath"

	annotation "github.com/YReshetko/go-annotation/pkg"
)

func init() {
	p := &Processor{
		gen:            make(map[fileKey][]generator),
		postConstructs: make(map[fileKey]map[string][]generators.PostConstructValues),
	}
	annotation.Register[annotations.Constructor](p)
	annotation.Register[annotations.Optional](p)
	annotation.Register[annotations.Builder](p)
	annotation.Register[annotations.PostConstruct](p)
	annotation.RegisterNoop[annotations.Exclude]()
	annotation.RegisterNoop[annotations.Init]()
}

var _ annotation.AnnotationProcessor = (*Processor)(nil)

type fileKey struct {
	dir string
	pkg string
}

type Processor struct {
	gen            map[fileKey][]generator
	postConstructs map[fileKey]map[string][]generators.PostConstructValues
}

type generator interface {
	Generate(map[string][]generators.PostConstructValues) ([]byte, []generators.Import, error)
}

func (p *Processor) Process(node annotation.Node) error {
	key := fileKey{
		dir: node.Dir(),
		pkg: node.PackageName(),
	}
	c, ts, ok, err := findAnnotation[annotations.Constructor](node)
	if err != nil {
		return err
	}

	if ok {
		p.gen[key] = append(p.gen[key], generators.NewConstructorGenerator(ts, c, node))
	}

	o, ts, ok, err := findAnnotation[annotations.Optional](node)
	if err != nil {
		return err
	}

	if ok {
		p.gen[key] = append(p.gen[key], generators.NewOptionalGenerator(ts, o, node))
	}

	b, ts, ok, err := findAnnotation[annotations.Builder](node)
	if err != nil {
		return err
	}

	if ok {
		p.gen[key] = append(p.gen[key], generators.NewBuilderGenerator(ts, b, node))
	}

	typeName, pcv, err := generators.PostConstructReceiverName(node)
	if err != nil {
		return fmt.Errorf("unable to build PostConstruct: %w", err)
	}

	if len(typeName) > 0 {
		pcs, ok := p.postConstructs[key]
		if !ok {
			pcs = map[string][]generators.PostConstructValues{}
			p.postConstructs[key] = pcs
		}
		pcvs := append(pcs[typeName], pcv)
		pcs[typeName] = pcvs
		fmt.Println(p.postConstructs)
	}

	return nil
}

func findAnnotation[T any](node annotation.Node) (T, *ast.TypeSpec, bool, error) {
	var a T
	ans := annotation.FindAnnotations[T](node.Annotations())
	if len(ans) == 0 {
		return a, nil, false, nil
	}

	if len(ans) > 1 {
		return a, nil, false, fmt.Errorf("expected 1 %T annotation, but got: %d", ans[0], len(ans))
	}

	ts, ok := annotation.CastNode[*ast.TypeSpec](node)
	if !ok {
		return a, nil, false, fmt.Errorf("unable to create constructor for %t: should be ast.TypeSpec", node.Node())
	}
	return ans[0], ts, true, nil
}

func (p *Processor) Output() map[string][]byte {
	out := map[string][]byte{}
	for k, gs := range p.gen {
		if len(gs) == 0 {
			continue
		}
		data := make([]generators.Data, len(gs))
		for i, g := range gs {
			d, im, _ := g.Generate(p.postConstructs[k])
			data[i] = generators.Data{
				Data:    d,
				Imports: im,
			}
		}

		out[filepath.Join(k.dir, "constructor.gen.go")] = generators.Generate(k.pkg, data)
	}

	return out
}

func (p *Processor) Version() string {
	return "0.0.2"
}

func (p *Processor) Name() string {
	return "Constructor"
}
