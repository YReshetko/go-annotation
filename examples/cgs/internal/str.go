package internal

import (
	"bytes"
	"golang.org/x/exp/constraints"
)

type stack[T any] []T
type queue[B any] []B

// @Constructor(name="NewSome{{.TypeName}}ThisIsMyTemplate")
type SomeStructure struct {
	a int
	b float64
	c *bool
	d **complex128
}

// @Constructor(name="NewAnotherStructOverride", type="pointer")
type AnotherStruct struct {
	a    SomeStructure
	b    *SomeStructure
	c, d int
	fn   func(**SomeStructure) AnotherStruct
	buff bytes.Buffer
}

// @Constructor
type TheThirdStruct struct {
	a    SomeStructure
	b    *SomeStructure
	c, d int
	fn   func(**SomeStructure) AnotherStruct
}

// @Constructor
type StackStruct[T stack[T]] struct {
	a  stack[T]
	q  queue[stack[T]]
	fn func(**SomeStructure) AnotherStruct
}

// @Constructor
type StackQueueStruct[T comparable, V constraints.Integer] struct {
	a    stack[T]
	q    queue[V]
	fn   func(**SomeStructure) AnotherStruct
	buff bytes.Buffer
}
