package constructor

import (
	"fmt"
	"github.com/YReshetko/go-annotation/annotations/constructor/annotations"
	"github.com/YReshetko/go-annotation/annotations/constructor/generators"
	"go/ast"
	"path/filepath"

	annotation "github.com/YReshetko/go-annotation/pkg"
	"github.com/hashicorp/go-multierror"
)

func init() {
	p := &Processor{cache: newCache()}
	annotation.Register[annotations.Constructor](p)
	annotation.Register[annotations.Optional](p)
	annotation.Register[annotations.Builder](p)
	annotation.Register[annotations.PostConstruct](p)
	annotation.RegisterNoop[annotations.Exclude]()
	annotation.RegisterNoop[annotations.Init]()
}

var _ annotation.AnnotationProcessor = (*Processor)(nil)

type Processor struct {
	cache *cache
}

type generator interface {
	Generate([]generators.PostConstructValues) ([]byte, []generators.Import, error)
	Name() string
}

func (p *Processor) Process(node annotation.Node) error {
	return multierror.Append(
		addAnnotatedTypeSpec[annotations.Constructor](p, node, newConstructorGenerator),
		addAnnotatedTypeSpec[annotations.Optional](p, node, newOptionalGenerator),
		addAnnotatedTypeSpec[annotations.Builder](p, node, newBuilderGenerator),
		p.addPostConstruct(node),
	).ErrorOrNil()
}

func (p *Processor) addPostConstruct(node annotation.Node) error {
	typeName, pcv, err := generators.PostConstructReceiverName(node)
	if err != nil {
		return fmt.Errorf("unable to build PostConstruct: %w", err)
	}

	if len(typeName) > 0 {
		p.cache.addPostConstruct(node.Dir(), node.PackageName(), typeName, pcv)
	}
	return nil
}

func addAnnotatedTypeSpec[T any](p *Processor, node annotation.Node, builder func(*ast.TypeSpec, T, annotation.Node) generator) error {
	a, ts, ok, err := findAnnotatedTypeSpec[T](node)
	if err != nil {
		return err
	}

	if ok {
		p.cache.addGenerator(node.Dir(), node.PackageName(), ts.Name.Name, builder(ts, a, node))
	}
	return nil
}

func findAnnotatedTypeSpec[T any](node annotation.Node) (T, *ast.TypeSpec, bool, error) {
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

func newConstructorGenerator(ts *ast.TypeSpec, a annotations.Constructor, n annotation.Node) generator {
	return generators.NewConstructorGenerator(ts, a, n)
}

func newOptionalGenerator(ts *ast.TypeSpec, a annotations.Optional, n annotation.Node) generator {
	return generators.NewOptionalGenerator(ts, a, n)
}

func newBuilderGenerator(ts *ast.TypeSpec, a annotations.Builder, n annotation.Node) generator {
	return generators.NewBuilderGenerator(ts, a, n)
}

func (p *Processor) Output() map[string][]byte {
	out := map[string][]byte{}
	data, err := p.cache.generate()
	if err != nil {
		panic(err)
	}
	for k, gd := range data {
		if len(gd) == 0 {
			continue
		}
		out[filepath.Join(k.dir, "constructor.gen.go")] = generators.Generate(k.pkg, gd)
	}

	return out
}

func (p *Processor) Version() string {
	return "0.0.3"
}

func (p *Processor) Name() string {
	return "Constructor"
}
