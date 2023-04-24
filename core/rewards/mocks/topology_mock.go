// Code generated by MockGen. DO NOT EDIT.
// Source: zuluprotocol/zeta/zeta/core/rewards (interfaces: Topology)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	types "zuluprotocol/zeta/zeta/core/types"
	gomock "github.com/golang/mock/gomock"
)

// MockTopology is a mock of Topology interface.
type MockTopology struct {
	ctrl     *gomock.Controller
	recorder *MockTopologyMockRecorder
}

// MockTopologyMockRecorder is the mock recorder for MockTopology.
type MockTopologyMockRecorder struct {
	mock *MockTopology
}

// NewMockTopology creates a new mock instance.
func NewMockTopology(ctrl *gomock.Controller) *MockTopology {
	mock := &MockTopology{ctrl: ctrl}
	mock.recorder = &MockTopologyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTopology) EXPECT() *MockTopologyMockRecorder {
	return m.recorder
}

// GetRewardsScores mocks base method.
func (m *MockTopology) GetRewardsScores(arg0 context.Context, arg1 string, arg2 []*types.ValidatorData, arg3 types.StakeScoreParams) (*types.ScoreData, *types.ScoreData) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRewardsScores", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*types.ScoreData)
	ret1, _ := ret[1].(*types.ScoreData)
	return ret0, ret1
}

// GetRewardsScores indicates an expected call of GetRewardsScores.
func (mr *MockTopologyMockRecorder) GetRewardsScores(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRewardsScores", reflect.TypeOf((*MockTopology)(nil).GetRewardsScores), arg0, arg1, arg2, arg3)
}

// RecalcValidatorSet mocks base method.
func (m *MockTopology) RecalcValidatorSet(arg0 context.Context, arg1 string, arg2 []*types.ValidatorData, arg3 types.StakeScoreParams) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RecalcValidatorSet", arg0, arg1, arg2, arg3)
}

// RecalcValidatorSet indicates an expected call of RecalcValidatorSet.
func (mr *MockTopologyMockRecorder) RecalcValidatorSet(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecalcValidatorSet", reflect.TypeOf((*MockTopology)(nil).RecalcValidatorSet), arg0, arg1, arg2, arg3)
}
