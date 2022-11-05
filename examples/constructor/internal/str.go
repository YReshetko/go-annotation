package internal

import (
	"bytes"
	"fmt"
	"net/http"

	"golang.org/x/exp/constraints"
)

type stack[T any] []T
type queue[B any] []B

// SomeStructure for testing base constructor options
// Annotations:
// 		@Constructor(name="NewSome{{.TypeName}}ThisIsMyTemplate")
// 		@Optional(constructor="NewSome{{ .TypeName }}ThisIsMyTemplateOptional")
//		@Builder(name="My{{.TypeName}}Builder", build="Build{{.FieldName}}Field", type="pointer")
type SomeStructure struct {
	a          int                                     // @Exclude
	b          float64                                 // @Exclude
	slice      []map[chan int]string                   // @Init(len="5", cap="15")
	slice2     []map[chan int]string                   // @Init
	maps       map[chan []int]struct{ A http.Request } // @Init(cap="5")
	chanals    chan []struct{ A http.Request }         // @Init
	chanalsCap chan []struct{ A http.Request }         // @Init(cap="5")
	c          *bool
	d          **complex128
}

// AnotherStruct for testing base constructor options
// Annotations:
// 		@Constructor(name="NewAnotherStructOverride", type="pointer")
// 		@Optional(constructor="New{{ .TypeName }}Optional")
type AnotherStruct struct {
	a    SomeStructure
	b    *SomeStructure
	c, d int // @Exclude
	fn   func(**SomeStructure) AnotherStruct
	buff bytes.Buffer
}

// TheThirdStruct for testing base constructor options
// Annotations:
// 		@Constructor
type TheThirdStruct struct {
	a    SomeStructure
	b    *SomeStructure
	c, d int
	fn   func(**SomeStructure) AnotherStruct
}

// StackStruct for testing parametrized constructor options
// Annotations:
// 		@Constructor
type StackStruct[T stack[T]] struct {
	a  stack[T]
	q  queue[stack[T]]
	fn func(**SomeStructure) AnotherStruct
}

// StackQueueStruct for testing parametrized constructor options
// Annotations:
// 		@Constructor(type="pointer")
// 		@Optional(constructor="New{{ .TypeName }}Optional", type="pointer", with="WithSQS{{ .FieldName }}")
//		@Builder(name="My{{.TypeName}}Builder", build="Build{{.FieldName}}Field", type="pointer")
type StackQueueStruct[T comparable, V constraints.Integer] struct {
	a    stack[T]
	q    queue[V]
	simp T
	vimp V
	fn   func(**SomeStructure) AnotherStruct
	buff bytes.Buffer
	str  chan map[T][]V // @Init
}

// @PostConstruct(priority="7")
func (s StackQueueStruct[T, V]) postConstruct3() {
	fmt.Println(s, "-3")
}

// @PostConstruct(priority="5")
func (s StackQueueStruct[T, V]) postConstruct1() {
	fmt.Println(s, "-1")
}

// @PostConstruct(priority="6")
func (s StackQueueStruct[T, V]) postConstruct2() {
	fmt.Println(s, "-2")
}

func Validation() {
	_ = NewStackQueueStructOptional[int, int](
		WithSQSA[int, int]([]int{}),
		WithSQSBuff[int, int](bytes.Buffer{}),
		WithSQSFn[int, int](func(i **SomeStructure) AnotherStruct {
			return AnotherStruct{}
		}),
		WithSQSQ[int, int]([]int{}),
		WithSQSStr[int, int](make(chan map[int][]int)),
	)

	nsb := NewStackQueueStructBuilder[bool, int]()
	_ = nsb.BuildSimpField(false).BuildAField(stack[bool]{}).Build()

	fp := false
	c := complex128(10)
	cp := &c
	b := NewSomeStructureBuilder()
	_ = b.BuildCField(&fp).
		BuildChanalsField(nil).
		BuildChanalsCapField(nil).
		BuildMapsField(nil).
		BuildDField(&cp).Build()

}
