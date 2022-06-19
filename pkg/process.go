package pkg

import (
	"github.com/YReshetko/go-annotation/internal/annotation/tag"
	"github.com/YReshetko/go-annotation/internal/environment"
	"github.com/YReshetko/go-annotation/internal/nodes"
	"github.com/YReshetko/go-annotation/internal/output"
)

func Process() {
	args := environment.LoadArguments()
	annotatedNodes := panicOnErr(nodes.ReadProject(args.ProjectPath))

	usedProcessors := map[AnnotationProcessor]struct{}{}

	for _, node := range annotatedNodes {
		intNode := newInternalNode(node)

		for _, annotation := range node.Annotations {
			p, ok := Processors()[annotation.Name()]
			if !ok {
				continue
			}
			a, ok := Annotations()[annotation.Name()]
			if !ok {
				continue
			}

			usedProcessors[p] = struct{}{}
			err := p.Process(tag.Parse(a, annotation), intNode)
			if err != nil {
				panic(err)
			}
		}
	}
	for processor, _ := range usedProcessors {
		pOut := processor.Output()
		o := map[string][]byte{}
		// TODO validate duplications
		for path, data := range pOut {
			o[string(path)] = data
		}
		if err := output.Save(args, o); err != nil {
			panic(err)
		}

	}
}

func panicOnErr[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
