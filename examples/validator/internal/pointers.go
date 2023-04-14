package internal

type SomeInternalStruct struct {
	m  map[string]struct{}
	m2 *map[string]struct{}
	s  []string
	s2 *[]string
	s3 *[7]string
	d  int
	b  bool
	b2 *bool
	i  interface{}
	i2 *interface{}
	i3 any
	i4 *any
}

// PointerValidation @Validator
type PointerValidation struct {
	sis  *SomeInternalStruct
	sis2 SomeInternalStruct
	s    []string
	d    int
	m    map[string]struct{}
}
