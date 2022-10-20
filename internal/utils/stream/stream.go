// NOTE: Take a look at https://github.com/mariomac/gostream

package stream

// NOTE: depending on this buffer and number of original array the functions OfSlice, OfMap, Map can be blockable or non blockable
// You always must use the functions that finalize stream: ForEach, Stream.Value
const buffer = 3

type Stream[T any] chan T

type Pair[A, B any] struct {
	Val1 A
	Val2 B
}

func (s Stream[T]) ToSlice() []T {
	var out []T
	for t := range s {
		//fmt.Println("ToSlice")
		out = append(out, t)
	}
	return out
}

func (s Stream[T]) Filter(f func(T) bool) Stream[T] {
	out := OfSlice[T](nil)
	go func() {
		defer close(out)
		for t := range s {
			if f(t) {
				//fmt.Println("Filter")
				out <- t
			}
		}
	}()
	return out
}

func (s Stream[T]) ForEach(f func(T)) {
	for t := range s {
		//fmt.Println("ForEach")
		f(t)
	}
}

func (s Stream[T]) ForEachErr(f func(T) error) error {
	var err error
	for t := range s {
		//fmt.Println("ForEachErr")
		if err == nil {
			err = f(t)
		}
	}
	return err
}

func Map[T, E any](s Stream[T], f func(T) E) Stream[E] {
	out := OfSlice[E](nil)
	go func() {
		defer close(out)
		for t := range s {
			//fmt.Println("Map")
			out <- f(t)
		}
	}()
	return out
}

func MapPair[T, E any](s Stream[T], f func(T) E) Stream[Pair[T, E]] {
	out := OfSlice[Pair[T, E]](nil)
	go func() {
		defer close(out)
		for t := range s {
			//fmt.Println("MapPair")
			e := f(t)
			out <- Pair[T, E]{Val1: t, Val2: e}
		}
	}()
	return out
}

func FlatMap[T, E any](s Stream[T], f func(T) []E) Stream[E] {
	out := OfSlice[E](nil)
	go func() {
		defer close(out)
		for t := range s {
			for _, v := range f(t) {
				//fmt.Println("FlatMap")
				out <- v
			}
		}
	}()
	return out
}

func OfSlice[T any](v []T) Stream[T] {
	s := make(chan T, buffer)
	if len(v) != 0 {
		go func() {
			defer close(s)
			for _, t := range v {
				//fmt.Println("OfSlice")
				s <- t
			}
		}()
	}
	return s
}

func OfMap[K comparable, V any](m map[K]V) Stream[Pair[K, V]] {
	s := OfSlice[Pair[K, V]](nil)
	if len(m) != 0 {
		go func() {
			defer close(s)
			for k, v := range m {
				//fmt.Println("OfSlice")
				s <- Pair[K, V]{Val1: k, Val2: v}
			}
		}()
	}
	return s
}
