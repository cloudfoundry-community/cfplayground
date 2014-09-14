// This file was generated by counterfeiter
package fakes

import (
	"net/http"
	"sync"

	"github.com/cloudfoundry-community/cfplayground/server"
)

type FakeServerHandlers struct {
	InitSessionStub        func(http.ResponseWriter, *http.Request)
	initSessionMutex       sync.RWMutex
	initSessionArgsForCall []struct {
		arg1 http.ResponseWriter
		arg2 *http.Request
	}
	RedirectBaseStub        func(http.ResponseWriter, *http.Request)
	redirectBaseMutex       sync.RWMutex
	redirectBaseArgsForCall []struct {
		arg1 http.ResponseWriter
		arg2 *http.Request
	}
	UploadHandlerStub        func(http.ResponseWriter, *http.Request)
	uploadHandlerMutex       sync.RWMutex
	uploadHandlerArgsForCall []struct {
		arg1 http.ResponseWriter
		arg2 *http.Request
	}
	DeleteHandlerStub        func(http.ResponseWriter, *http.Request)
	deleteHandlerMutex       sync.RWMutex
	deleteHandlerArgsForCall []struct {
		arg1 http.ResponseWriter
		arg2 *http.Request
	}
	BasePathStub        func() string
	basePathMutex       sync.RWMutex
	basePathArgsForCall []struct{}
	basePathReturns     struct {
		result1 string
	}
}

func (fake *FakeServerHandlers) InitSession(arg1 http.ResponseWriter, arg2 *http.Request) {
	fake.initSessionMutex.Lock()
	defer fake.initSessionMutex.Unlock()
	fake.initSessionArgsForCall = append(fake.initSessionArgsForCall, struct {
		arg1 http.ResponseWriter
		arg2 *http.Request
	}{arg1, arg2})
	if fake.InitSessionStub != nil {
		fake.InitSessionStub(arg1, arg2)
	}
}

func (fake *FakeServerHandlers) InitSessionCallCount() int {
	fake.initSessionMutex.RLock()
	defer fake.initSessionMutex.RUnlock()
	return len(fake.initSessionArgsForCall)
}

func (fake *FakeServerHandlers) InitSessionArgsForCall(i int) (http.ResponseWriter, *http.Request) {
	fake.initSessionMutex.RLock()
	defer fake.initSessionMutex.RUnlock()
	return fake.initSessionArgsForCall[i].arg1, fake.initSessionArgsForCall[i].arg2
}

func (fake *FakeServerHandlers) RedirectBase(arg1 http.ResponseWriter, arg2 *http.Request) {
	fake.redirectBaseMutex.Lock()
	defer fake.redirectBaseMutex.Unlock()
	fake.redirectBaseArgsForCall = append(fake.redirectBaseArgsForCall, struct {
		arg1 http.ResponseWriter
		arg2 *http.Request
	}{arg1, arg2})
	if fake.RedirectBaseStub != nil {
		fake.RedirectBaseStub(arg1, arg2)
	}
}

func (fake *FakeServerHandlers) RedirectBaseCallCount() int {
	fake.redirectBaseMutex.RLock()
	defer fake.redirectBaseMutex.RUnlock()
	return len(fake.redirectBaseArgsForCall)
}

func (fake *FakeServerHandlers) RedirectBaseArgsForCall(i int) (http.ResponseWriter, *http.Request) {
	fake.redirectBaseMutex.RLock()
	defer fake.redirectBaseMutex.RUnlock()
	return fake.redirectBaseArgsForCall[i].arg1, fake.redirectBaseArgsForCall[i].arg2
}

func (fake *FakeServerHandlers) UploadHandler(arg1 http.ResponseWriter, arg2 *http.Request) {
	fake.uploadHandlerMutex.Lock()
	defer fake.uploadHandlerMutex.Unlock()
	fake.uploadHandlerArgsForCall = append(fake.uploadHandlerArgsForCall, struct {
		arg1 http.ResponseWriter
		arg2 *http.Request
	}{arg1, arg2})
	if fake.UploadHandlerStub != nil {
		fake.UploadHandlerStub(arg1, arg2)
	}
}

func (fake *FakeServerHandlers) UploadHandlerCallCount() int {
	fake.uploadHandlerMutex.RLock()
	defer fake.uploadHandlerMutex.RUnlock()
	return len(fake.uploadHandlerArgsForCall)
}

func (fake *FakeServerHandlers) UploadHandlerArgsForCall(i int) (http.ResponseWriter, *http.Request) {
	fake.uploadHandlerMutex.RLock()
	defer fake.uploadHandlerMutex.RUnlock()
	return fake.uploadHandlerArgsForCall[i].arg1, fake.uploadHandlerArgsForCall[i].arg2
}

func (fake *FakeServerHandlers) DeleteHandler(arg1 http.ResponseWriter, arg2 *http.Request) {
	fake.deleteHandlerMutex.Lock()
	defer fake.deleteHandlerMutex.Unlock()
	fake.deleteHandlerArgsForCall = append(fake.deleteHandlerArgsForCall, struct {
		arg1 http.ResponseWriter
		arg2 *http.Request
	}{arg1, arg2})
	if fake.DeleteHandlerStub != nil {
		fake.DeleteHandlerStub(arg1, arg2)
	}
}

func (fake *FakeServerHandlers) DeleteHandlerCallCount() int {
	fake.deleteHandlerMutex.RLock()
	defer fake.deleteHandlerMutex.RUnlock()
	return len(fake.deleteHandlerArgsForCall)
}

func (fake *FakeServerHandlers) DeleteHandlerArgsForCall(i int) (http.ResponseWriter, *http.Request) {
	fake.deleteHandlerMutex.RLock()
	defer fake.deleteHandlerMutex.RUnlock()
	return fake.deleteHandlerArgsForCall[i].arg1, fake.deleteHandlerArgsForCall[i].arg2
}

func (fake *FakeServerHandlers) BasePath() string {
	fake.basePathMutex.Lock()
	defer fake.basePathMutex.Unlock()
	fake.basePathArgsForCall = append(fake.basePathArgsForCall, struct{}{})
	if fake.BasePathStub != nil {
		return fake.BasePathStub()
	} else {
		return fake.basePathReturns.result1
	}
}

func (fake *FakeServerHandlers) BasePathCallCount() int {
	fake.basePathMutex.RLock()
	defer fake.basePathMutex.RUnlock()
	return len(fake.basePathArgsForCall)
}

func (fake *FakeServerHandlers) BasePathReturns(result1 string) {
	fake.BasePathStub = nil
	fake.basePathReturns = struct {
		result1 string
	}{result1}
}

var _ server.ServerHandlers = new(FakeServerHandlers)
