package annotation

import (
	"fmt"
	goAST "go/ast"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"

	"github.com/YReshetko/go-annotation/internal/ast"
	"github.com/YReshetko/go-annotation/internal/environment"
	"github.com/YReshetko/go-annotation/internal/module"
	"github.com/YReshetko/go-annotation/internal/output"
	"github.com/YReshetko/go-annotation/internal/parser"
	"github.com/YReshetko/go-annotation/internal/tag"
	. "github.com/YReshetko/go-annotation/internal/utils/stream"
)

// Process is entry point for parsing project annotations
func Process() {
	root := environment.ProjectRoot()
	m := panicOnError(module.Load(root))

	run(m, m.Files())

	if len(rerunable) > 0 {
		f := filesToReRun()
		for ; len(f) > 0; f = filesToReRun() {
			run(m, f)
		}
	}
}

func filesToReRun() []string {
	return FlatMap(Map(OfMap(rerunable), ExtractVal2[string, Rerunable]()), filesToRerun).
		Filter(isNotEmpty).ToSlice()
}

func isNotEmpty(s string) bool {
	return len(s) > 0
}
func filesToRerun(t Rerunable) []string {
	defer t.Clear()
	out := t.ToRerun()
	return out
}

func run(m module.Module, files []string) {
	// Execute processors for annotations
	err := MapPair(OfSlice(files).Map(joinPath(m.Root())), toAstFile).RangeErr(moduleNodeProcessor(m))
	if err != nil {
		panic(err)
	}

	// Persist required data for annotation processors
	Map(OfMap(processors), ExtractVal2[string, AnnotationProcessor]()).
		Filter(DistinctBy(annotationProcessorDistinct)).Range(storeData)
}

func moduleNodeProcessor(m module.Module) func(Pair[string, *goAST.File]) error {
	return func(p Pair[string, *goAST.File]) error {
		f := p.Val2
		path := p.Val1
		var err error
		// TODO set parents no framework node
		ast.Walk(f, func(n goAST.Node, parents []goAST.Node) bool {
			if err != nil {
				return false
			}
			a, ok := annotationsByNode(n)
			if !ok {
				return true
			}
			internalNode := newNode(m, path, f, n, parents, filledAnnotations(a))

			err = Map(Map(OfSlice(a), toParsedAnnotationName).
				Filter(Distinct[string]()), processorByAnnotationName).
				Filter(NonNil[AnnotationProcessor]()).
				Filter(DistinctBy(annotationProcessorDistinct)).
				RangeErr(processNode(internalNode))

			return err == nil
		})
		if err != nil {
			printer.Fprint(os.Stdout, &token.FileSet{}, f)
		}
		return err
	}
}

func storeData(a AnnotationProcessor) {
	po := a.Output()
	if len(po) == 0 {
		return
	}

	meta := map[string]string{
		"annotation_name": a.Name(),
		"annotation_ver":  a.Version(),
	}

	for path, data := range po {
		if err := output.Store(path, data, meta); err != nil {
			fmt.Println(string(data))
			panic(err)
		}
	}
}

func annotationsByNode(n goAST.Node) ([]parser.Annotation, bool) {
	c, ok := ast.Comment(n)
	if !ok {
		return nil, false
	}

	a := panicOnError(parser.Parse(c))
	if len(a) == 0 {
		return nil, false
	}
	return a, true
}

func filledAnnotations(a []parser.Annotation) []Annotation {
	return Map(MapPair(OfSlice(a), annotationTypeByParsedAnnotation).
		Filter(NonNilPair[parser.Annotation, Annotation]()),
		fillAnnotation).ToSlice()
}

func processNode(n Node) func(AnnotationProcessor) error {
	return func(annotationProcessor AnnotationProcessor) error {
		return annotationProcessor.Process(n)
	}
}

func fillAnnotation(p Pair[parser.Annotation, Annotation]) Annotation {
	return panicOnError(tag.Parse(p.Val2, p.Val1.Params()))
}

func annotationTypeByParsedAnnotation(a parser.Annotation) Annotation {
	return annotations[a.Name()]
}

func annotationProcessorDistinct(a AnnotationProcessor) string {
	return a.Name() + "_" + a.Version()
}

func toParsedAnnotationName(a parser.Annotation) string {
	return a.Name()
}

func joinPath(root string) func(string) string {
	return func(path string) string {
		return filepath.Join(root, path)
	}
}

func toAstFile(p string) *goAST.File {
	return panicOnError(ast.LoadFileAST(p))
}

func panicOnError[T any](v T, e error) T {
	if e != nil {
		panic(e)
	}
	return v
}
