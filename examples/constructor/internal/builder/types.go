// Invalid example but that would be nice to explore the case deeply

package builder

import (
	"fmt"
	"net/http"
)

type CommonInterface interface {
	Do()
}

// StructureForExample1 builder example
// @Builder(type="pointer")
type StructureForExample1 struct {
	Field1 http.Client
	Field2 int
	Field3 string
}

// StructureForExample2 builder example
// @Builder(type="pointer")
type StructureForExample2 struct {
	Field1 http.Client
	Field2 int
	Field3 string
}

func (s StructureForExample1) Do() {
	fmt.Println("StructureForExample1 is doing something, values:", s.Field2, s.Field3)
}

func (s StructureForExample2) Do() {
	fmt.Println("StructureForExample2 is doing something, values:", s.Field2, s.Field3)
}
