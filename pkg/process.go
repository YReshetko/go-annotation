package pkg

import (
	"github.com/YReshetko/go-annotation/internal/debug"
	"github.com/YReshetko/go-annotation/internal/environment"
	"github.com/YReshetko/go-annotation/internal/lookup"
	"github.com/YReshetko/go-annotation/internal/module"
	"github.com/YReshetko/go-annotation/internal/output"
)

func Process() {
	args := environment.LoadArguments()
	loader, rootSelector, err := module.NewLoader(args, lookup.ModuleStructure)
	if err != nil {
		debug.Critical("unable to init project loader %w", err)
	}

	for _, annotationProcessor := range processors {
		annotationProcessor.SetLookup(func(n Node, selector Selector) Node {
			no, ok := n.(node)
			if !ok {
				debug.Critical("unable to cast node for selector %v", selector)
			}
			out, err := loader.FindNode(no.n.Selector, selector.PackageImport, selector.TypeName)
			if err != nil {
				debug.Critical("unable to load node for selector %v: %w", selector, err)
			}
			return newNode(out)
		})
	}

	annotatedNodes, err := loader.AllAnnotatedNodes(rootSelector)
	if err != nil {
		debug.Critical("unable to get root annotated nodes %w", err)
	}

	nodes := make([]node, len(annotatedNodes))
	for i, annotatedNode := range annotatedNodes {
		nodes[i] = newNode(annotatedNode)
	}

	usedProcessors := processNodes(nodes)

	for processor, _ := range usedProcessors {
		pOut := processor.Output()
		o := map[string][]byte{}
		// TODO validate file name duplications in the same package
		for path, data := range pOut {
			o[string(path)] = data
		}
		if err := output.Save(args, o); err != nil {
			panic(err)
		}
	}
}

func processNodes(nodes []node) map[AnnotationProcessor]struct{} {
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
		for p, _ := range processNodes(n.inner) {
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
