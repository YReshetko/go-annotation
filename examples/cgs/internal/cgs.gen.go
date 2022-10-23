// Code generated by CGS annotation processor. DO NOT EDIT.
// versions:
//		go-annotation: 0.0.2-alpha
//		CGS: 0.0.1

package internal

import (
	bytes "bytes"
	constraints "golang.org/x/exp/constraints"
)

func NewSomeSomeStructureThisIsMyTemplate(a int, b float64, c *bool, d **complex128) SomeStructure {
	return SomeStructure{
		a: a,
		b: b,
		c: c,
		d: d,
	}
}

func NewAnotherStructOverride(fn func(**SomeStructure) AnotherStruct, buff bytes.Buffer, a SomeStructure, b *SomeStructure, c int, d int) *AnotherStruct {
	return &AnotherStruct{
		fn:   fn,
		buff: buff,
		a:    a,
		b:    b,
		c:    c,
		d:    d,
	}
}

func NewTheThirdStruct(a SomeStructure, b *SomeStructure, c int, d int, fn func(**SomeStructure) AnotherStruct) TheThirdStruct {
	return TheThirdStruct{
		a:  a,
		b:  b,
		c:  c,
		d:  d,
		fn: fn,
	}
}

func NewStackStruct[T stack[T]](a stack[T], q queue[stack[T]], fn func(**SomeStructure) AnotherStruct) StackStruct[T] {
	return StackStruct[T]{
		a:  a,
		q:  q,
		fn: fn,
	}
}

func NewStackQueueStruct[T comparable, V constraints.Integer](fn func(**SomeStructure) AnotherStruct, buff bytes.Buffer, a stack[T], q queue[V]) StackQueueStruct[T, V] {
	return StackQueueStruct[T, V]{
		fn:   fn,
		buff: buff,
		a:    a,
		q:    q,
	}
}
