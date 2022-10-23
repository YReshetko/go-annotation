package cgs

import (
	"fmt"
	"go/ast"
	"path/filepath"

	annotation "github.com/YReshetko/go-annotation/pkg"
)

func init() {
	p := &Processor{gen: make(map[string][]toGenerate)}
	annotation.Register[Constructor](p)
}

var _ annotation.AnnotationProcessor = (*Processor)(nil)

type Processor struct {
	gen map[string][]toGenerate
}

type toGenerate struct {
	node        *ast.TypeSpec
	annotation  Constructor
	packageName string
	an          annotation.Node
}

func (p *Processor) Process(node annotation.Node) error {
	annotations := annotation.FindAnnotations[Constructor](node.Annotations())
	if len(annotations) == 0 {
		return nil
	}

	if len(annotations) > 1 {
		return fmt.Errorf("expected 1 Constructor annotation, but got: %d", len(annotations))
	}

	n, ok := annotation.CastNode[*ast.TypeSpec](node)
	if !ok {
		return fmt.Errorf("unable to create constructor fot %t: should be ast.TypeSpec", node.Node())
	}

	p.gen[node.Dir()] = append(p.gen[node.Dir()], toGenerate{
		node:        n,
		annotation:  annotations[0],
		packageName: node.PackageName(),
		an:          node,
	})

	return nil
}

func (p *Processor) Output() map[string][]byte {
	out := map[string][]byte{}
	for path, generates := range p.gen {
		if len(generates) == 0 {
			continue
		}
		out[filepath.Join(path, "cgs.gen.go")] = generateConstructors(generates)
	}

	return out
}

func (p *Processor) Version() string {
	return "0.0.1"
}

func (p *Processor) Name() string {
	return "CGS"
}
