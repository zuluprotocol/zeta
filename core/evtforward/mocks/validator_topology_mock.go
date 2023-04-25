// Code generated by MockGen. DO NOT EDIT.
// Source: zuluprotocol/zeta/core/evtforward (interfaces: ValidatorTopology)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockValidatorTopology is a mock of ValidatorTopology interface.
type MockValidatorTopology struct {
	ctrl     *gomock.Controller
	recorder *MockValidatorTopologyMockRecorder
}

// MockValidatorTopologyMockRecorder is the mock recorder for MockValidatorTopology.
type MockValidatorTopologyMockRecorder struct {
	mock *MockValidatorTopology
}

// NewMockValidatorTopology creates a new mock instance.
func NewMockValidatorTopology(ctrl *gomock.Controller) *MockValidatorTopology {
	mock := &MockValidatorTopology{ctrl: ctrl}
	mock.recorder = &MockValidatorTopologyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockValidatorTopology) EXPECT() *MockValidatorTopologyMockRecorder {
	return m.recorder
}

// AllNodeIDs mocks base method.
func (m *MockValidatorTopology) AllNodeIDs() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllNodeIDs")
	ret0, _ := ret[0].([]string)
	return ret0
}

// AllNodeIDs indicates an expected call of AllNodeIDs.
func (mr *MockValidatorTopologyMockRecorder) AllNodeIDs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllNodeIDs", reflect.TypeOf((*MockValidatorTopology)(nil).AllNodeIDs))
}

// SelfNodeID mocks base method.
func (m *MockValidatorTopology) SelfNodeID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelfNodeID")
	ret0, _ := ret[0].(string)
	return ret0
}

// SelfNodeID indicates an expected call of SelfNodeID.
func (mr *MockValidatorTopologyMockRecorder) SelfNodeID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelfNodeID", reflect.TypeOf((*MockValidatorTopology)(nil).SelfNodeID))
}
