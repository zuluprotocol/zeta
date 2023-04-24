// Code generated by MockGen. DO NOT EDIT.
// Source: zuluprotocol/zeta/zeta/core/settlement (interfaces: MarketPosition)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	num "zuluprotocol/zeta/zeta/libs/num"
	gomock "github.com/golang/mock/gomock"
)

// MockMarketPosition is a mock of MarketPosition interface.
type MockMarketPosition struct {
	ctrl     *gomock.Controller
	recorder *MockMarketPositionMockRecorder
}

// MockMarketPositionMockRecorder is the mock recorder for MockMarketPosition.
type MockMarketPositionMockRecorder struct {
	mock *MockMarketPosition
}

// NewMockMarketPosition creates a new mock instance.
func NewMockMarketPosition(ctrl *gomock.Controller) *MockMarketPosition {
	mock := &MockMarketPosition{ctrl: ctrl}
	mock.recorder = &MockMarketPositionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMarketPosition) EXPECT() *MockMarketPositionMockRecorder {
	return m.recorder
}

// Buy mocks base method.
func (m *MockMarketPosition) Buy() int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Buy")
	ret0, _ := ret[0].(int64)
	return ret0
}

// Buy indicates an expected call of Buy.
func (mr *MockMarketPositionMockRecorder) Buy() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Buy", reflect.TypeOf((*MockMarketPosition)(nil).Buy))
}

// ClearPotentials mocks base method.
func (m *MockMarketPosition) ClearPotentials() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ClearPotentials")
}

// ClearPotentials indicates an expected call of ClearPotentials.
func (mr *MockMarketPositionMockRecorder) ClearPotentials() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearPotentials", reflect.TypeOf((*MockMarketPosition)(nil).ClearPotentials))
}

// Party mocks base method.
func (m *MockMarketPosition) Party() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Party")
	ret0, _ := ret[0].(string)
	return ret0
}

// Party indicates an expected call of Party.
func (mr *MockMarketPositionMockRecorder) Party() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Party", reflect.TypeOf((*MockMarketPosition)(nil).Party))
}

// Price mocks base method.
func (m *MockMarketPosition) Price() *num.Uint {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Price")
	ret0, _ := ret[0].(*num.Uint)
	return ret0
}

// Price indicates an expected call of Price.
func (mr *MockMarketPositionMockRecorder) Price() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Price", reflect.TypeOf((*MockMarketPosition)(nil).Price))
}

// Sell mocks base method.
func (m *MockMarketPosition) Sell() int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sell")
	ret0, _ := ret[0].(int64)
	return ret0
}

// Sell indicates an expected call of Sell.
func (mr *MockMarketPositionMockRecorder) Sell() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sell", reflect.TypeOf((*MockMarketPosition)(nil).Sell))
}

// Size mocks base method.
func (m *MockMarketPosition) Size() int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Size")
	ret0, _ := ret[0].(int64)
	return ret0
}

// Size indicates an expected call of Size.
func (mr *MockMarketPositionMockRecorder) Size() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Size", reflect.TypeOf((*MockMarketPosition)(nil).Size))
}

// VWBuy mocks base method.
func (m *MockMarketPosition) VWBuy() *num.Uint {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VWBuy")
	ret0, _ := ret[0].(*num.Uint)
	return ret0
}

// VWBuy indicates an expected call of VWBuy.
func (mr *MockMarketPositionMockRecorder) VWBuy() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VWBuy", reflect.TypeOf((*MockMarketPosition)(nil).VWBuy))
}

// VWSell mocks base method.
func (m *MockMarketPosition) VWSell() *num.Uint {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VWSell")
	ret0, _ := ret[0].(*num.Uint)
	return ret0
}

// VWSell indicates an expected call of VWSell.
func (mr *MockMarketPositionMockRecorder) VWSell() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VWSell", reflect.TypeOf((*MockMarketPosition)(nil).VWSell))
}
