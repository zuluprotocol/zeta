// Code generated by MockGen. DO NOT EDIT.
// Source: code.zetaprotocol.io/vega/libs/jsonrpc (interfaces: Command)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	jsonrpc "code.zetaprotocol.io/vega/libs/jsonrpc"
	gomock "github.com/golang/mock/gomock"
)

// MockCommand is a mock of Command interface.
type MockCommand struct {
	ctrl     *gomock.Controller
	recorder *MockCommandMockRecorder
}

// MockCommandMockRecorder is the mock recorder for MockCommand.
type MockCommandMockRecorder struct {
	mock *MockCommand
}

// NewMockCommand creates a new mock instance.
func NewMockCommand(ctrl *gomock.Controller) *MockCommand {
	mock := &MockCommand{ctrl: ctrl}
	mock.recorder = &MockCommandMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommand) EXPECT() *MockCommandMockRecorder {
	return m.recorder
}

// Handle mocks base method.
func (m *MockCommand) Handle(arg0 context.Context, arg1 jsonrpc.Params) (jsonrpc.Result, *jsonrpc.ErrorDetails) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handle", arg0, arg1)
	ret0, _ := ret[0].(jsonrpc.Result)
	ret1, _ := ret[1].(*jsonrpc.ErrorDetails)
	return ret0, ret1
}

// Handle indicates an expected call of Handle.
func (mr *MockCommandMockRecorder) Handle(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockCommand)(nil).Handle), arg0, arg1)
}
