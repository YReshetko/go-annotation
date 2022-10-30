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
	p := &Processor{gen: make(map[fileKey][]generator)}
	annotation.Register[annotations.Constructor](p)
	annotation.Register[annotations.Optional](p)
	annotation.RegisterNoop[annotations.Exclude]()
	annotation.RegisterNoop[annotations.Init]()
}

var _ annotation.AnnotationProcessor = (*Processor)(nil)

type Processor struct {
	gen map[fileKey][]generator
}

type generator interface {
	Generate() ([]byte, []generators.Import, error)
}

type fileKey struct {
	dir string
	pkg string
}
type toGenerate struct {
	packageName string
	generators  []generator
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
			d, im, _ := g.Generate()
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
