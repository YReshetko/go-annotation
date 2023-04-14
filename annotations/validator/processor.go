package validator

import (
	"errors"
	"fmt"
	"github.com/YReshetko/go-annotation/annotations/validator/annotations"
	"github.com/YReshetko/go-annotation/annotations/validator/generators"
	annotation "github.com/YReshetko/go-annotation/pkg"
	"go/ast"
	"path/filepath"
)

func init() {
	p := &Processor{cache: map[string]*pkgOutput{}}
	annotation.Register[annotations.Validator](p)
}

var _ annotation.AnnotationProcessor = (*Processor)(nil)

type Processor struct {
	cache map[string]*pkgOutput
}
type pkgOutput struct {
	name string
	o    []generators.Output
}

func (p Processor) Process(node annotation.Node) error {
	validators := annotation.FindAnnotations[annotations.Validator](node.Annotations())
	if len(validators) == 0 || len(validators) > 1 {
		return errors.New("validator annotation should be single per type")
	}
	validator := validators[0]
	ts, ok := annotation.CastNode[*ast.TypeSpec](node)
	if !ok {
		return errors.New("validator annotation is applicable to ast.TypeSpec only")
	}

	_, ok = ts.Type.(*ast.StructType)
	if !ok {
		return errors.New("validator should be declared for ast.StructType only")
	}

	g := generators.NewGenerator(node, ts, validator)
	o := generators.NewOutput(g.Validatable())
	pkgO, ok := p.cache[node.Meta().Dir()]
	if !ok {
		pkgO = &pkgOutput{
			name: node.Meta().PackageName(),
			o:    make([]generators.Output, 0),
		}
		p.cache[node.Meta().Dir()] = pkgO
	}

	pkgO.o = append(pkgO.o, o)

	return nil
}

func (p Processor) Output() map[string][]byte {
	out := make(map[string][]byte)
	for dir, outputs := range p.cache {
		d := []byte(fmt.Sprintf("package %s\n\n", outputs.name))
		for _, output := range outputs.o {
			data := output.Generate()
			d = append(d, data...)
		}

		out[filepath.Join(dir, "validators.gen.go")] = d
	}
	return out
}

func (p Processor) Version() string {
	return "0.0.1"
}

func (p Processor) Name() string {
	return "Validator"
}
