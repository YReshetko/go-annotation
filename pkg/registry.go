package pkg

import "reflect"

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

func processor(annotation Annotation) (AnnotationProcessor, bool) {
	if reflect.TypeOf(annotation).Kind() != reflect.Struct {
		panic("unable to register non-struct annotation")
	}
	annotationName := reflect.TypeOf(annotation).Name()
	a, ok := processors[annotationName]
	return a, ok
}

/*func Processors() map[string]AnnotationProcessor {
	return processors
}

func Annotations() map[string]Annotation {
	return annotations
}
*/
