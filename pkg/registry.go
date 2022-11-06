package annotation

import (
	"go/ast"
	"reflect"
)

type Annotation any

type Node interface {
	// Base API

	Annotations() []Annotation
	ASTNode() ast.Node
	AnnotatedNode(ast.Node) Node

	// Metadata API

	Root() string
	Dir() string
	FileName() string
	PackageName() string
	Imports() []*ast.ImportSpec

	// Lookup API

	// FindImportByAlias returns related import for alias in current ast.File.
	// For example:
	// import "github.com/YReshetko/go-annotation/internal/tag"
	// ...
	// tag.Parse(...)
	// FindImportByAlias("tag") returns "github.com/YReshetko/go-annotation/internal/tag", true
	FindImportByAlias(alias string) (string, bool)

	// FindNodeByAlias returns related Node by alias and a type/function name from related module
	// if alias is empty, then the search will go in current directory of ast.File
	FindNodeByAlias(alias, nodeName string) (Node, error)
}

type AnnotationProcessor interface {
	Process(node Node) error
	Output() map[string][]byte
	Version() string
	Name() string
}

type Rerunable interface {
	ToRerun() []string
	Clear()
}

var _ AnnotationProcessor = (*noopProcessor)(nil)

type noopProcessor struct{}

func (n noopProcessor) Process(Node) error        { return nil }
func (n noopProcessor) Output() map[string][]byte { return nil }
func (n noopProcessor) Version() string           { return "" }
func (n noopProcessor) Name() string              { return "noop" }

var processors = map[string]AnnotationProcessor{}
var rerunable = map[string]Rerunable{}
var annotations = map[string]Annotation{}

func Register[T Annotation](processor AnnotationProcessor) {
	var v T
	if reflect.TypeOf(v).Kind() != reflect.Struct {
		panic("unable to register non-struct annotation")
	}
	annotationName := reflect.TypeOf(v).Name()
	processors[annotationName] = processor
	annotations[annotationName] = v

	if r, ok := processor.(Rerunable); ok {
		rerunable[processorKey(processor)] = r
	}
}

func RegisterNoop[T Annotation]() {
	Register[T](noopProcessor{})
}

func processorByAnnotationName(n string) AnnotationProcessor {
	return processors[n]
}

func processorKey(p AnnotationProcessor) string {
	return p.Name() + p.Version()
}

func processor[T Annotation]() (AnnotationProcessor, bool) {
	var v T
	if reflect.TypeOf(v).Kind() != reflect.Struct {
		panic("unable to retrieve non-struct annotation")
	}
	annotationName := reflect.TypeOf(v).Name()
	a, ok := processors[annotationName]
	return a, ok
}
