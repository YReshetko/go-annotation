package pkg

import (
	"github.com/YReshetko/go-annotation/internal/environment"
	"github.com/YReshetko/go-annotation/internal/nodes"
	"github.com/YReshetko/go-annotation/internal/output"
)

func Process() {
	args := environment.LoadArguments()
	annotatedNodes := panicOnErr(nodes.ReadProject(args.ProjectPath))

	nodes := make([]node, len(annotatedNodes))
	for i, annotatedNode := range annotatedNodes {
		nodes[i] = newNode(annotatedNode)
	}

	usedProcessors := processNode(nodes)

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

func processNode(nodes []node) map[AnnotationProcessor]struct{} {
	usedProcessors := make(map[AnnotationProcessor]struct{})
	for _, n := range nodes {
		for _, annotation := range n.Annotations() {
			p, ok := processor(annotation)
			if !ok {
				continue
			}
			usedProcessors[p] = struct{}{}
			err := p.Process(annotation, n)
			if err != nil {
				panic(err)
			}
		}
		for p, _ := range processNode(n.inner) {
			usedProcessors[p] = struct{}{}
		}
	}
	return usedProcessors
}

func panicOnErr[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
