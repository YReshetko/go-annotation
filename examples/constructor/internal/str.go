package internal

import (
	"bytes"
	"net/http"

	"golang.org/x/exp/constraints"
)

type stack[T any] []T
type queue[B any] []B

// @Constructor(name="NewSome{{.TypeName}}ThisIsMyTemplate")
// @Optional(constructor="NewSome{{ .TypeName }}ThisIsMyTemplateOptional")
type SomeStructure struct {
	// @Exclude
	a int
	// @Exclude
	b float64
	// @Init(len="5", cap="10")
	slice []map[chan int]string
	// @Init
	slice2 []map[chan int]string
	// @Init(cap="5")
	maps map[chan []int]struct{ A http.Request }
	// @Init
	chanals chan []struct{ A http.Request }
	// @Init(cap="5")
	chanalsCap chan []struct{ A http.Request }
	c          *bool
	d          **complex128
}

// @Constructor(name="NewAnotherStructOverride", type="pointer")
// @Optional(constructor="New{{ .TypeName }}Optional")
type AnotherStruct struct {
	a SomeStructure
	b *SomeStructure
	// @Exclude
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
// @Optional(constructor="New{{ .TypeName }}Optional", type="pointer", with="WithSQS{{ .FieldName }}")
type StackQueueStruct[T comparable, V constraints.Integer] struct {
	a    stack[T]
	q    queue[V]
	fn   func(**SomeStructure) AnotherStruct
	buff bytes.Buffer
	// @Init
	str chan map[T][]V
}

func validation() {
	_ = NewStackQueueStructOptional[int, int](
		WithSQSA[int, int]([]int{}),
		WithSQSBuff[int, int](bytes.Buffer{}),
		WithSQSFn[int, int](func(i **SomeStructure) AnotherStruct {
			return AnotherStruct{}
		}),
		WithSQSQ[int, int]([]int{}),
		WithSQSStr[int, int](make(chan map[int][]int)),
	)
}
