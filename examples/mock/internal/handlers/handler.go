// @Mock(name="HandlersMock")
package handlers

import "net/http"

// InterfaceReference test interface reference for mocking
// Annotations:
//		@Mock
type InterfaceReference SomeHandler

// EmbeddedInterface test interface embedding for mocking
// Annotations:
//		@Mock
type EmbeddedInterface interface {
	CopySomeInternalHandler
	SomeAnotherHandler
}

// SomeHandler test interface for mocking
// Annotations:
//		@Mock
type SomeHandler interface {
	Get(int) float64
	Set(float32) int64
}

// SomeAnotherHandler test interface for mocking
// Annotations:
//		@Mock(name="OverrodeHandlerMockName")
type SomeAnotherHandler interface {
	Handle(r http.Request) (http.Response, error)
	Proxy() SomeHandler
	SetProxy(SomeHandler) error
}

type (
	// SomeInternalHandler test interface for mocking
	// Annotations:
	//		@Mock(name="{{ .TypeName }}OverrodeMock", sub="mocks_2")
	SomeInternalHandler interface {
		Handle(r http.Request) (http.Response, error)
		Proxy() SomeHandler
		SetProxy(SomeHandler) error
	}

	// CopySomeInternalHandler test interface for mocking
	// Annotations:
	//		@Mock(name="OverrodeMock", sub="mocks_2")
	CopySomeInternalHandler interface {
		Handle(r http.Request) (http.Response, error)
		Proxy() SomeHandler
		SetProxy(SomeHandler) error
	}
)

// FunctionForTestingItsTool test interface for mocking
// Annotations:
//		@Mock(sub="mocks_2")
type FunctionForTestingItsTool func(SomeInternalHandler, http.Response) http.Request

// Testing package interface
func DoSomething(r http.Request) {

}
func DoSomeAnother(response http.Response) {

}
