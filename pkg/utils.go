package pkg

import . "github.com/YReshetko/go-annotation/internal/utils/stream"

func FindAnnotations[T any](a []Annotation) []T {
	return Map(OfSlice(a).Filter(func(a Annotation) bool {
		_, ok := a.(T)
		return ok
	}), func(t Annotation) T {
		return t.(T)
	}).ToSlice()
}
