// Code generated by MockGen. DO NOT EDIT.
// Source: zuluprotocol/zeta/zeta/datanode/candlesv2 (interfaces: CandleStore)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	entities "zuluprotocol/zeta/zeta/datanode/entities"
	gomock "github.com/golang/mock/gomock"
)

// MockCandleStore is a mock of CandleStore interface.
type MockCandleStore struct {
	ctrl     *gomock.Controller
	recorder *MockCandleStoreMockRecorder
}

// MockCandleStoreMockRecorder is the mock recorder for MockCandleStore.
type MockCandleStoreMockRecorder struct {
	mock *MockCandleStore
}

// NewMockCandleStore creates a new mock instance.
func NewMockCandleStore(ctrl *gomock.Controller) *MockCandleStore {
	mock := &MockCandleStore{ctrl: ctrl}
	mock.recorder = &MockCandleStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCandleStore) EXPECT() *MockCandleStoreMockRecorder {
	return m.recorder
}

// CandleExists mocks base method.
func (m *MockCandleStore) CandleExists(arg0 context.Context, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CandleExists", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CandleExists indicates an expected call of CandleExists.
func (mr *MockCandleStoreMockRecorder) CandleExists(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CandleExists", reflect.TypeOf((*MockCandleStore)(nil).CandleExists), arg0, arg1)
}

// GetCandleDataForTimeSpan mocks base method.
func (m *MockCandleStore) GetCandleDataForTimeSpan(arg0 context.Context, arg1 string, arg2, arg3 *time.Time, arg4 entities.CursorPagination) ([]entities.Candle, entities.PageInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCandleDataForTimeSpan", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].([]entities.Candle)
	ret1, _ := ret[1].(entities.PageInfo)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetCandleDataForTimeSpan indicates an expected call of GetCandleDataForTimeSpan.
func (mr *MockCandleStoreMockRecorder) GetCandleDataForTimeSpan(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCandleDataForTimeSpan", reflect.TypeOf((*MockCandleStore)(nil).GetCandleDataForTimeSpan), arg0, arg1, arg2, arg3, arg4)
}

// GetCandleIDForIntervalAndMarket mocks base method.
func (m *MockCandleStore) GetCandleIDForIntervalAndMarket(arg0 context.Context, arg1, arg2 string) (bool, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCandleIDForIntervalAndMarket", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetCandleIDForIntervalAndMarket indicates an expected call of GetCandleIDForIntervalAndMarket.
func (mr *MockCandleStoreMockRecorder) GetCandleIDForIntervalAndMarket(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCandleIDForIntervalAndMarket", reflect.TypeOf((*MockCandleStore)(nil).GetCandleIDForIntervalAndMarket), arg0, arg1, arg2)
}

// GetCandlesForMarket mocks base method.
func (m *MockCandleStore) GetCandlesForMarket(arg0 context.Context, arg1 string) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCandlesForMarket", arg0, arg1)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCandlesForMarket indicates an expected call of GetCandlesForMarket.
func (mr *MockCandleStoreMockRecorder) GetCandlesForMarket(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCandlesForMarket", reflect.TypeOf((*MockCandleStore)(nil).GetCandlesForMarket), arg0, arg1)
}
