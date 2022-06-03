package pkg

func CastAnnotation[T Annotation](a Annotation) T {
	t, ok := a.(T)
	if !ok {
		panic("unable parse annotation")
	}
	return t
}
