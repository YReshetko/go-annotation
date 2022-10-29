package pkg

import (
	"go/ast"
	"reflect"
)

type Annotation any

type Node interface {
	// Base API

	Annotations() []Annotation
	Node() ast.Node
	AnnotatedNode(ast.Node) Node

	// Metadata API

	Root() string
	Dir() string
	FileName() string
	PackageName() string
	Imports() []*ast.ImportSpec

	// Lookup API

	FindImportByAlias(alias string) (string, bool)
}

type AnnotationProcessor interface {
	Process(node Node) error
	Output() map[string][]byte
	Version() string
	Name() string
}

var processors = map[string]AnnotationProcessor{}
var annotations = map[string]Annotation{}

func Register[T Annotation](processor AnnotationProcessor) {
	var v T
	if reflect.TypeOf(v).Kind() != reflect.Struct {
		panic("unable to register non-struct annotation")
	}
	annotationName := reflect.TypeOf(v).Name()
	processors[annotationName] = processor
	annotations[annotationName] = v
}

func processorByAnnotationName(n string) AnnotationProcessor {
	return processors[n]
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
