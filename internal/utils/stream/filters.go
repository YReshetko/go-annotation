package stream

import "reflect"

func Distinct[T comparable]() func(T) bool {
	s := map[T]struct{}{}
	return func(t T) bool {
		_, ok := s[t]
		if !ok {
			s[t] = struct{}{}
			return true
		}
		return false
	}
}

func DistinctBy[T any, B comparable](extractor func(T) B) func(T) bool {
	s := map[B]struct{}{}
	return func(t T) bool {
		v := extractor(t)
		_, ok := s[v]
		if !ok {
			s[v] = struct{}{}
			return true
		}
		return false
	}
}

func NonNil[T any]() func(v T) bool {
	return func(v T) bool {
		return !reflect.ValueOf(&v).Elem().IsZero()
	}
}

func NonNilPair[A, B any]() func(Pair[A, B]) bool {
	return func(p Pair[A, B]) bool {
		return !reflect.ValueOf(&p.Val1).Elem().IsZero() && !reflect.ValueOf(&p.Val2).Elem().IsZero()
	}
}

func ExtractVal1[A, B any]() func(Pair[A, B]) A {
	return func(p Pair[A, B]) A {
		return p.Val1
	}
}

func ExtractVal2[A, B any]() func(Pair[A, B]) B {
	return func(p Pair[A, B]) B {
		return p.Val2
	}
}

type emtyable interface {
	string | []any | map[any]any
}

func IsEmpty[T emtyable](t T) bool {
	return len(t) == 0
}

func IsNotEmpty[T emtyable](t T) bool {
	return len(t) != 0
}
