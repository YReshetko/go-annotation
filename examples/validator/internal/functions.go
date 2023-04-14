package internal

// SomeStructWithFunc @Validator
type SomeStructWithFunc struct {
	fn  *func(bool2 bool) string
	fn2 func(bool2 bool) string
}
