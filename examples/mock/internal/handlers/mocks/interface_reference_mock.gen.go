// Code generated by Mock annotation processor. DO NOT EDIT.
// versions:
//
//	go: go1.22.4
//	go-annotation: 0.1.0
//	Mock: 0.0.1
//
// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"

	"github.com/YReshetko/go-annotation/examples/mock/internal/handlers"
)

type InterfaceReferenceMock struct {
	GetStub        func(int) float64
	getMutex       sync.RWMutex
	getArgsForCall []struct {
		arg1 int
	}
	getReturns struct {
		result1 float64
	}
	getReturnsOnCall map[int]struct {
		result1 float64
	}
	SetStub        func(float32) int64
	setMutex       sync.RWMutex
	setArgsForCall []struct {
		arg1 float32
	}
	setReturns struct {
		result1 int64
	}
	setReturnsOnCall map[int]struct {
		result1 int64
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *InterfaceReferenceMock) Get(arg1 int) float64 {
	fake.getMutex.Lock()
	ret, specificReturn := fake.getReturnsOnCall[len(fake.getArgsForCall)]
	fake.getArgsForCall = append(fake.getArgsForCall, struct {
		arg1 int
	}{arg1})
	stub := fake.GetStub
	fakeReturns := fake.getReturns
	fake.recordInvocation("Get", []interface{}{arg1})
	fake.getMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *InterfaceReferenceMock) GetCallCount() int {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return len(fake.getArgsForCall)
}

func (fake *InterfaceReferenceMock) GetCalls(stub func(int) float64) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = stub
}

func (fake *InterfaceReferenceMock) GetArgsForCall(i int) int {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	argsForCall := fake.getArgsForCall[i]
	return argsForCall.arg1
}

func (fake *InterfaceReferenceMock) GetReturns(result1 float64) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = nil
	fake.getReturns = struct {
		result1 float64
	}{result1}
}

func (fake *InterfaceReferenceMock) GetReturnsOnCall(i int, result1 float64) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = nil
	if fake.getReturnsOnCall == nil {
		fake.getReturnsOnCall = make(map[int]struct {
			result1 float64
		})
	}
	fake.getReturnsOnCall[i] = struct {
		result1 float64
	}{result1}
}

func (fake *InterfaceReferenceMock) Set(arg1 float32) int64 {
	fake.setMutex.Lock()
	ret, specificReturn := fake.setReturnsOnCall[len(fake.setArgsForCall)]
	fake.setArgsForCall = append(fake.setArgsForCall, struct {
		arg1 float32
	}{arg1})
	stub := fake.SetStub
	fakeReturns := fake.setReturns
	fake.recordInvocation("Set", []interface{}{arg1})
	fake.setMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *InterfaceReferenceMock) SetCallCount() int {
	fake.setMutex.RLock()
	defer fake.setMutex.RUnlock()
	return len(fake.setArgsForCall)
}

func (fake *InterfaceReferenceMock) SetCalls(stub func(float32) int64) {
	fake.setMutex.Lock()
	defer fake.setMutex.Unlock()
	fake.SetStub = stub
}

func (fake *InterfaceReferenceMock) SetArgsForCall(i int) float32 {
	fake.setMutex.RLock()
	defer fake.setMutex.RUnlock()
	argsForCall := fake.setArgsForCall[i]
	return argsForCall.arg1
}

func (fake *InterfaceReferenceMock) SetReturns(result1 int64) {
	fake.setMutex.Lock()
	defer fake.setMutex.Unlock()
	fake.SetStub = nil
	fake.setReturns = struct {
		result1 int64
	}{result1}
}

func (fake *InterfaceReferenceMock) SetReturnsOnCall(i int, result1 int64) {
	fake.setMutex.Lock()
	defer fake.setMutex.Unlock()
	fake.SetStub = nil
	if fake.setReturnsOnCall == nil {
		fake.setReturnsOnCall = make(map[int]struct {
			result1 int64
		})
	}
	fake.setReturnsOnCall[i] = struct {
		result1 int64
	}{result1}
}

func (fake *InterfaceReferenceMock) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	fake.setMutex.RLock()
	defer fake.setMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *InterfaceReferenceMock) recordInvocation(key string, args []interface{}) {
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

var _ handlers.InterfaceReference = new(InterfaceReferenceMock)
