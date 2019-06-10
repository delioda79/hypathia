// Code generated by counterfeiter. DO NOT EDIT.
package searchfakes

import (
	"sync"

	"github.com/taxibeat/hypatia/search"
)

type FakeDocument struct {
	ContentStub        func() (interface{}, error)
	contentMutex       sync.RWMutex
	contentArgsForCall []struct {
	}
	contentReturns struct {
		result1 interface{}
		result2 error
	}
	contentReturnsOnCall map[int]struct {
		result1 interface{}
		result2 error
	}
	GetIDStub        func() string
	getIDMutex       sync.RWMutex
	getIDArgsForCall []struct {
	}
	getIDReturns struct {
		result1 string
	}
	getIDReturnsOnCall map[int]struct {
		result1 string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeDocument) Content() (interface{}, error) {
	fake.contentMutex.Lock()
	ret, specificReturn := fake.contentReturnsOnCall[len(fake.contentArgsForCall)]
	fake.contentArgsForCall = append(fake.contentArgsForCall, struct {
	}{})
	fake.recordInvocation("Content", []interface{}{})
	fake.contentMutex.Unlock()
	if fake.ContentStub != nil {
		return fake.ContentStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.contentReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeDocument) ContentCallCount() int {
	fake.contentMutex.RLock()
	defer fake.contentMutex.RUnlock()
	return len(fake.contentArgsForCall)
}

func (fake *FakeDocument) ContentCalls(stub func() (interface{}, error)) {
	fake.contentMutex.Lock()
	defer fake.contentMutex.Unlock()
	fake.ContentStub = stub
}

func (fake *FakeDocument) ContentReturns(result1 interface{}, result2 error) {
	fake.contentMutex.Lock()
	defer fake.contentMutex.Unlock()
	fake.ContentStub = nil
	fake.contentReturns = struct {
		result1 interface{}
		result2 error
	}{result1, result2}
}

func (fake *FakeDocument) ContentReturnsOnCall(i int, result1 interface{}, result2 error) {
	fake.contentMutex.Lock()
	defer fake.contentMutex.Unlock()
	fake.ContentStub = nil
	if fake.contentReturnsOnCall == nil {
		fake.contentReturnsOnCall = make(map[int]struct {
			result1 interface{}
			result2 error
		})
	}
	fake.contentReturnsOnCall[i] = struct {
		result1 interface{}
		result2 error
	}{result1, result2}
}

func (fake *FakeDocument) GetID() string {
	fake.getIDMutex.Lock()
	ret, specificReturn := fake.getIDReturnsOnCall[len(fake.getIDArgsForCall)]
	fake.getIDArgsForCall = append(fake.getIDArgsForCall, struct {
	}{})
	fake.recordInvocation("GetID", []interface{}{})
	fake.getIDMutex.Unlock()
	if fake.GetIDStub != nil {
		return fake.GetIDStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.getIDReturns
	return fakeReturns.result1
}

func (fake *FakeDocument) GetIDCallCount() int {
	fake.getIDMutex.RLock()
	defer fake.getIDMutex.RUnlock()
	return len(fake.getIDArgsForCall)
}

func (fake *FakeDocument) GetIDCalls(stub func() string) {
	fake.getIDMutex.Lock()
	defer fake.getIDMutex.Unlock()
	fake.GetIDStub = stub
}

func (fake *FakeDocument) GetIDReturns(result1 string) {
	fake.getIDMutex.Lock()
	defer fake.getIDMutex.Unlock()
	fake.GetIDStub = nil
	fake.getIDReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeDocument) GetIDReturnsOnCall(i int, result1 string) {
	fake.getIDMutex.Lock()
	defer fake.getIDMutex.Unlock()
	fake.GetIDStub = nil
	if fake.getIDReturnsOnCall == nil {
		fake.getIDReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.getIDReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeDocument) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.contentMutex.RLock()
	defer fake.contentMutex.RUnlock()
	fake.getIDMutex.RLock()
	defer fake.getIDMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeDocument) recordInvocation(key string, args []interface{}) {
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

var _ search.Document = new(FakeDocument)
