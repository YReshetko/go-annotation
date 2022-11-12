package common

type SliceValue struct {
}

type MapKey struct {
	Field1 *string
	Field2 string
	Field3 int
}

type Common struct {
	Field1 *string
	Field2 string
	Field3 int
	Slice  []SliceValue
}

type Common2 struct {
	Field1 *string
	Field2 string
	Field3 int
	Slice  []SliceValue
}
