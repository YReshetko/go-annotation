package common

import "github.com/google/uuid"

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

type WithUUID struct {
	UID *uuid.UUID
}
type StringForUUID struct {
	UID string
}
