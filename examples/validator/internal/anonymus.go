package internal

// SomeStructWithAnonymous @Validator
type SomeStructWithAnonymous struct {
	an struct {
		f     float32
		i     int
		b     bool
		s     string
		inter any
		fn    func()
	}
}
