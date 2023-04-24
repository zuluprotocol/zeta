// Code generated by MockGen. DO NOT EDIT.
// Source: zuluprotocol/zeta/zeta/core/types (interfaces: StateProvider)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	types "zuluprotocol/zeta/zeta/core/types"
	gomock "github.com/golang/mock/gomock"
)

// MockStateProvider is a mock of StateProvider interface.
type MockStateProvider struct {
	ctrl     *gomock.Controller
	recorder *MockStateProviderMockRecorder
}

// MockStateProviderMockRecorder is the mock recorder for MockStateProvider.
type MockStateProviderMockRecorder struct {
	mock *MockStateProvider
}

// NewMockStateProvider creates a new mock instance.
func NewMockStateProvider(ctrl *gomock.Controller) *MockStateProvider {
	mock := &MockStateProvider{ctrl: ctrl}
	mock.recorder = &MockStateProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStateProvider) EXPECT() *MockStateProviderMockRecorder {
	return m.recorder
}

// GetState mocks base method.
func (m *MockStateProvider) GetState(arg0 string) ([]byte, []types.StateProvider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetState", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].([]types.StateProvider)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetState indicates an expected call of GetState.
func (mr *MockStateProviderMockRecorder) GetState(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetState", reflect.TypeOf((*MockStateProvider)(nil).GetState), arg0)
}

// Keys mocks base method.
func (m *MockStateProvider) Keys() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Keys")
	ret0, _ := ret[0].([]string)
	return ret0
}

// Keys indicates an expected call of Keys.
func (mr *MockStateProviderMockRecorder) Keys() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Keys", reflect.TypeOf((*MockStateProvider)(nil).Keys))
}

// LoadState mocks base method.
func (m *MockStateProvider) LoadState(arg0 context.Context, arg1 *types.Payload) ([]types.StateProvider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadState", arg0, arg1)
	ret0, _ := ret[0].([]types.StateProvider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadState indicates an expected call of LoadState.
func (mr *MockStateProviderMockRecorder) LoadState(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadState", reflect.TypeOf((*MockStateProvider)(nil).LoadState), arg0, arg1)
}

// Namespace mocks base method.
func (m *MockStateProvider) Namespace() types.SnapshotNamespace {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Namespace")
	ret0, _ := ret[0].(types.SnapshotNamespace)
	return ret0
}

// Namespace indicates an expected call of Namespace.
func (mr *MockStateProviderMockRecorder) Namespace() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Namespace", reflect.TypeOf((*MockStateProvider)(nil).Namespace))
}

// Stopped mocks base method.
func (m *MockStateProvider) Stopped() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stopped")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Stopped indicates an expected call of Stopped.
func (mr *MockStateProviderMockRecorder) Stopped() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stopped", reflect.TypeOf((*MockStateProvider)(nil).Stopped))
}
