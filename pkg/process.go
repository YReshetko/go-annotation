package pkg

import (
	"github.com/YReshetko/go-annotation/internal/annotation/tag"
	"github.com/YReshetko/go-annotation/internal/environment"
	"github.com/YReshetko/go-annotation/internal/nodes"
)

func Process() {
	args := environment.LoadArguments()
	nodes := panicOnErr(nodes.ReadProject(args.ProjectPath))

	for _, node := range nodes {
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

			err := p.Process(tag.Parse(a, annotation), intNode)
			if err != nil {
				panic(err)
			}
		}
	}
}

func panicOnErr[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
