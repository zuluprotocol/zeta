// Code generated by MockGen. DO NOT EDIT.
// Source: zuluprotocol/zeta/zeta/core/bridges (interfaces: ETHClient)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	big "math/big"
	reflect "reflect"

	ethereum "github.com/ethereum/go-ethereum"
	common "github.com/ethereum/go-ethereum/common"
	types "github.com/ethereum/go-ethereum/core/types"
	gomock "github.com/golang/mock/gomock"
)

// MockETHClient is a mock of ETHClient interface.
type MockETHClient struct {
	ctrl     *gomock.Controller
	recorder *MockETHClientMockRecorder
}

// MockETHClientMockRecorder is the mock recorder for MockETHClient.
type MockETHClientMockRecorder struct {
	mock *MockETHClient
}

// NewMockETHClient creates a new mock instance.
func NewMockETHClient(ctrl *gomock.Controller) *MockETHClient {
	mock := &MockETHClient{ctrl: ctrl}
	mock.recorder = &MockETHClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockETHClient) EXPECT() *MockETHClientMockRecorder {
	return m.recorder
}

// CallContract mocks base method.
func (m *MockETHClient) CallContract(arg0 context.Context, arg1 ethereum.CallMsg, arg2 *big.Int) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CallContract", arg0, arg1, arg2)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CallContract indicates an expected call of CallContract.
func (mr *MockETHClientMockRecorder) CallContract(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CallContract", reflect.TypeOf((*MockETHClient)(nil).CallContract), arg0, arg1, arg2)
}

// CodeAt mocks base method.
func (m *MockETHClient) CodeAt(arg0 context.Context, arg1 common.Address, arg2 *big.Int) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CodeAt", arg0, arg1, arg2)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CodeAt indicates an expected call of CodeAt.
func (mr *MockETHClientMockRecorder) CodeAt(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CodeAt", reflect.TypeOf((*MockETHClient)(nil).CodeAt), arg0, arg1, arg2)
}

// CollateralBridgeAddress mocks base method.
func (m *MockETHClient) CollateralBridgeAddress() common.Address {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CollateralBridgeAddress")
	ret0, _ := ret[0].(common.Address)
	return ret0
}

// CollateralBridgeAddress indicates an expected call of CollateralBridgeAddress.
func (mr *MockETHClientMockRecorder) CollateralBridgeAddress() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CollateralBridgeAddress", reflect.TypeOf((*MockETHClient)(nil).CollateralBridgeAddress))
}

// ConfirmationsRequired mocks base method.
func (m *MockETHClient) ConfirmationsRequired() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConfirmationsRequired")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// ConfirmationsRequired indicates an expected call of ConfirmationsRequired.
func (mr *MockETHClientMockRecorder) ConfirmationsRequired() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConfirmationsRequired", reflect.TypeOf((*MockETHClient)(nil).ConfirmationsRequired))
}

// CurrentHeight mocks base method.
func (m *MockETHClient) CurrentHeight(arg0 context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CurrentHeight", arg0)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CurrentHeight indicates an expected call of CurrentHeight.
func (mr *MockETHClientMockRecorder) CurrentHeight(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CurrentHeight", reflect.TypeOf((*MockETHClient)(nil).CurrentHeight), arg0)
}

// EstimateGas mocks base method.
func (m *MockETHClient) EstimateGas(arg0 context.Context, arg1 ethereum.CallMsg) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EstimateGas", arg0, arg1)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EstimateGas indicates an expected call of EstimateGas.
func (mr *MockETHClientMockRecorder) EstimateGas(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EstimateGas", reflect.TypeOf((*MockETHClient)(nil).EstimateGas), arg0, arg1)
}

// FilterLogs mocks base method.
func (m *MockETHClient) FilterLogs(arg0 context.Context, arg1 ethereum.FilterQuery) ([]types.Log, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FilterLogs", arg0, arg1)
	ret0, _ := ret[0].([]types.Log)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FilterLogs indicates an expected call of FilterLogs.
func (mr *MockETHClientMockRecorder) FilterLogs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FilterLogs", reflect.TypeOf((*MockETHClient)(nil).FilterLogs), arg0, arg1)
}

// HeaderByNumber mocks base method.
func (m *MockETHClient) HeaderByNumber(arg0 context.Context, arg1 *big.Int) (*types.Header, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HeaderByNumber", arg0, arg1)
	ret0, _ := ret[0].(*types.Header)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HeaderByNumber indicates an expected call of HeaderByNumber.
func (mr *MockETHClientMockRecorder) HeaderByNumber(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HeaderByNumber", reflect.TypeOf((*MockETHClient)(nil).HeaderByNumber), arg0, arg1)
}

// PendingCodeAt mocks base method.
func (m *MockETHClient) PendingCodeAt(arg0 context.Context, arg1 common.Address) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PendingCodeAt", arg0, arg1)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PendingCodeAt indicates an expected call of PendingCodeAt.
func (mr *MockETHClientMockRecorder) PendingCodeAt(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PendingCodeAt", reflect.TypeOf((*MockETHClient)(nil).PendingCodeAt), arg0, arg1)
}

// PendingNonceAt mocks base method.
func (m *MockETHClient) PendingNonceAt(arg0 context.Context, arg1 common.Address) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PendingNonceAt", arg0, arg1)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PendingNonceAt indicates an expected call of PendingNonceAt.
func (mr *MockETHClientMockRecorder) PendingNonceAt(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PendingNonceAt", reflect.TypeOf((*MockETHClient)(nil).PendingNonceAt), arg0, arg1)
}

// SendTransaction mocks base method.
func (m *MockETHClient) SendTransaction(arg0 context.Context, arg1 *types.Transaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendTransaction", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendTransaction indicates an expected call of SendTransaction.
func (mr *MockETHClientMockRecorder) SendTransaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendTransaction", reflect.TypeOf((*MockETHClient)(nil).SendTransaction), arg0, arg1)
}

// SubscribeFilterLogs mocks base method.
func (m *MockETHClient) SubscribeFilterLogs(arg0 context.Context, arg1 ethereum.FilterQuery, arg2 chan<- types.Log) (ethereum.Subscription, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeFilterLogs", arg0, arg1, arg2)
	ret0, _ := ret[0].(ethereum.Subscription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubscribeFilterLogs indicates an expected call of SubscribeFilterLogs.
func (mr *MockETHClientMockRecorder) SubscribeFilterLogs(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeFilterLogs", reflect.TypeOf((*MockETHClient)(nil).SubscribeFilterLogs), arg0, arg1, arg2)
}

// SuggestGasPrice mocks base method.
func (m *MockETHClient) SuggestGasPrice(arg0 context.Context) (*big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SuggestGasPrice", arg0)
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SuggestGasPrice indicates an expected call of SuggestGasPrice.
func (mr *MockETHClientMockRecorder) SuggestGasPrice(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SuggestGasPrice", reflect.TypeOf((*MockETHClient)(nil).SuggestGasPrice), arg0)
}

// SuggestGasTipCap mocks base method.
func (m *MockETHClient) SuggestGasTipCap(arg0 context.Context) (*big.Int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SuggestGasTipCap", arg0)
	ret0, _ := ret[0].(*big.Int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SuggestGasTipCap indicates an expected call of SuggestGasTipCap.
func (mr *MockETHClientMockRecorder) SuggestGasTipCap(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SuggestGasTipCap", reflect.TypeOf((*MockETHClient)(nil).SuggestGasTipCap), arg0)
}
