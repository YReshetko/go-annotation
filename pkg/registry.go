package pkg

import (
	"go/ast"
	"reflect"
)

type Annotation any

// TODO Make node for processing
type Node interface {
	Annotations() []Annotation
	Node() ast.Node

	AnnotatedNode(ast.Node) Node

	// TODO Define all required methods
	Dir() string
	FileName() string
	PackageName() string
	Imports() []*ast.ImportSpec
}

type AnnotationProcessor interface {
	Process(node Node) error
	Output() map[string][]byte
	Version() string
	Name() string
}

var processors = map[string]AnnotationProcessor{}
var annotations = map[string]Annotation{}

func Register(annotation Annotation, processor AnnotationProcessor) {
	if reflect.TypeOf(annotation).Kind() != reflect.Struct {
		panic("unable to register non-struct annotation")
	}
	annotationName := reflect.TypeOf(annotation).Name()
	processors[annotationName] = processor
	annotations[annotationName] = annotation
}

func processorByAnnotationName(n string) AnnotationProcessor {
	return processors[n]
}

func processor(annotation Annotation) (AnnotationProcessor, bool) {
	if reflect.TypeOf(annotation).Kind() != reflect.Struct {
		panic("unable to retrieve non-struct annotation")
	}
	annotationName := reflect.TypeOf(annotation).Name()
	a, ok := processors[annotationName]
	return a, ok
}
