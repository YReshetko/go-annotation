package parametrized

import "net/http"

// GenericInterfaceMock test interface for mocking generics
// Annotations:
//		@Mock(name="GenericInterfaceMock")
type GenericInterfaceMock GenericInterface[int, float32]

type GenericInterface[T any, V comparable] interface {
	Process(T, []V) (chan T, V)
	Get() V
}

func Process[T any, V comparable](g GenericInterface[T, V], t T, v []V) bool {
	ch, _ := g.Process(t, v)

	ind := 0
	for _ = range ch {
		ind++
	}
	return ind%5 == 3
}

// ParametrizedFuncMock test function for mocking generics
// Annotations:
//		@Mock(name="ParametrizedFuncMock")
type ParametrizedFuncMock ParametrizedFunc[*http.Request]

type ParametrizedFunc[T comparable] func(t T) []T
