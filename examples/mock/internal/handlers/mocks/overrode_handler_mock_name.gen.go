// Code generated by Mock annotation processor. DO NOT EDIT.
// versions:
//		go-annotation: 0.0.2-alpha
//		Mock: 0.0.1
// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"net/http"
	"sync"

	"github.com/YReshetko/go-annotation/examples/mock/internal/handlers"
)

type OverrodeHandlerMockName struct {
	HandleStub        func(http.Request) (http.Response, error)
	handleMutex       sync.RWMutex
	handleArgsForCall []struct {
		arg1 http.Request
	}
	handleReturns struct {
		result1 http.Response
		result2 error
	}
	handleReturnsOnCall map[int]struct {
		result1 http.Response
		result2 error
	}
	ProxyStub        func() handlers.SomeHandler
	proxyMutex       sync.RWMutex
	proxyArgsForCall []struct {
	}
	proxyReturns struct {
		result1 handlers.SomeHandler
	}
	proxyReturnsOnCall map[int]struct {
		result1 handlers.SomeHandler
	}
	SetProxyStub        func(handlers.SomeHandler) error
	setProxyMutex       sync.RWMutex
	setProxyArgsForCall []struct {
		arg1 handlers.SomeHandler
	}
	setProxyReturns struct {
		result1 error
	}
	setProxyReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *OverrodeHandlerMockName) Handle(arg1 http.Request) (http.Response, error) {
	fake.handleMutex.Lock()
	ret, specificReturn := fake.handleReturnsOnCall[len(fake.handleArgsForCall)]
	fake.handleArgsForCall = append(fake.handleArgsForCall, struct {
		arg1 http.Request
	}{arg1})
	stub := fake.HandleStub
	fakeReturns := fake.handleReturns
	fake.recordInvocation("Handle", []interface{}{arg1})
	fake.handleMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *OverrodeHandlerMockName) HandleCallCount() int {
	fake.handleMutex.RLock()
	defer fake.handleMutex.RUnlock()
	return len(fake.handleArgsForCall)
}

func (fake *OverrodeHandlerMockName) HandleCalls(stub func(http.Request) (http.Response, error)) {
	fake.handleMutex.Lock()
	defer fake.handleMutex.Unlock()
	fake.HandleStub = stub
}

func (fake *OverrodeHandlerMockName) HandleArgsForCall(i int) http.Request {
	fake.handleMutex.RLock()
	defer fake.handleMutex.RUnlock()
	argsForCall := fake.handleArgsForCall[i]
	return argsForCall.arg1
}

func (fake *OverrodeHandlerMockName) HandleReturns(result1 http.Response, result2 error) {
	fake.handleMutex.Lock()
	defer fake.handleMutex.Unlock()
	fake.HandleStub = nil
	fake.handleReturns = struct {
		result1 http.Response
		result2 error
	}{result1, result2}
}

func (fake *OverrodeHandlerMockName) HandleReturnsOnCall(i int, result1 http.Response, result2 error) {
	fake.handleMutex.Lock()
	defer fake.handleMutex.Unlock()
	fake.HandleStub = nil
	if fake.handleReturnsOnCall == nil {
		fake.handleReturnsOnCall = make(map[int]struct {
			result1 http.Response
			result2 error
		})
	}
	fake.handleReturnsOnCall[i] = struct {
		result1 http.Response
		result2 error
	}{result1, result2}
}

func (fake *OverrodeHandlerMockName) Proxy() handlers.SomeHandler {
	fake.proxyMutex.Lock()
	ret, specificReturn := fake.proxyReturnsOnCall[len(fake.proxyArgsForCall)]
	fake.proxyArgsForCall = append(fake.proxyArgsForCall, struct {
	}{})
	stub := fake.ProxyStub
	fakeReturns := fake.proxyReturns
	fake.recordInvocation("Proxy", []interface{}{})
	fake.proxyMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *OverrodeHandlerMockName) ProxyCallCount() int {
	fake.proxyMutex.RLock()
	defer fake.proxyMutex.RUnlock()
	return len(fake.proxyArgsForCall)
}

func (fake *OverrodeHandlerMockName) ProxyCalls(stub func() handlers.SomeHandler) {
	fake.proxyMutex.Lock()
	defer fake.proxyMutex.Unlock()
	fake.ProxyStub = stub
}

func (fake *OverrodeHandlerMockName) ProxyReturns(result1 handlers.SomeHandler) {
	fake.proxyMutex.Lock()
	defer fake.proxyMutex.Unlock()
	fake.ProxyStub = nil
	fake.proxyReturns = struct {
		result1 handlers.SomeHandler
	}{result1}
}

func (fake *OverrodeHandlerMockName) ProxyReturnsOnCall(i int, result1 handlers.SomeHandler) {
	fake.proxyMutex.Lock()
	defer fake.proxyMutex.Unlock()
	fake.ProxyStub = nil
	if fake.proxyReturnsOnCall == nil {
		fake.proxyReturnsOnCall = make(map[int]struct {
			result1 handlers.SomeHandler
		})
	}
	fake.proxyReturnsOnCall[i] = struct {
		result1 handlers.SomeHandler
	}{result1}
}

func (fake *OverrodeHandlerMockName) SetProxy(arg1 handlers.SomeHandler) error {
	fake.setProxyMutex.Lock()
	ret, specificReturn := fake.setProxyReturnsOnCall[len(fake.setProxyArgsForCall)]
	fake.setProxyArgsForCall = append(fake.setProxyArgsForCall, struct {
		arg1 handlers.SomeHandler
	}{arg1})
	stub := fake.SetProxyStub
	fakeReturns := fake.setProxyReturns
	fake.recordInvocation("SetProxy", []interface{}{arg1})
	fake.setProxyMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *OverrodeHandlerMockName) SetProxyCallCount() int {
	fake.setProxyMutex.RLock()
	defer fake.setProxyMutex.RUnlock()
	return len(fake.setProxyArgsForCall)
}

func (fake *OverrodeHandlerMockName) SetProxyCalls(stub func(handlers.SomeHandler) error) {
	fake.setProxyMutex.Lock()
	defer fake.setProxyMutex.Unlock()
	fake.SetProxyStub = stub
}

func (fake *OverrodeHandlerMockName) SetProxyArgsForCall(i int) handlers.SomeHandler {
	fake.setProxyMutex.RLock()
	defer fake.setProxyMutex.RUnlock()
	argsForCall := fake.setProxyArgsForCall[i]
	return argsForCall.arg1
}

func (fake *OverrodeHandlerMockName) SetProxyReturns(result1 error) {
	fake.setProxyMutex.Lock()
	defer fake.setProxyMutex.Unlock()
	fake.SetProxyStub = nil
	fake.setProxyReturns = struct {
		result1 error
	}{result1}
}

func (fake *OverrodeHandlerMockName) SetProxyReturnsOnCall(i int, result1 error) {
	fake.setProxyMutex.Lock()
	defer fake.setProxyMutex.Unlock()
	fake.SetProxyStub = nil
	if fake.setProxyReturnsOnCall == nil {
		fake.setProxyReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.setProxyReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *OverrodeHandlerMockName) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.handleMutex.RLock()
	defer fake.handleMutex.RUnlock()
	fake.proxyMutex.RLock()
	defer fake.proxyMutex.RUnlock()
	fake.setProxyMutex.RLock()
	defer fake.setProxyMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *OverrodeHandlerMockName) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ handlers.SomeAnotherHandler = new(OverrodeHandlerMockName)