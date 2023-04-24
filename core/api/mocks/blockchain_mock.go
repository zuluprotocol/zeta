// Code generated by MockGen. DO NOT EDIT.
// Source: code.zetaprotocol.io/vega/core/api (interfaces: Blockchain)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	v1 "code.zetaprotocol.io/vega/protos/vega/commands/v1"
	gomock "github.com/golang/mock/gomock"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
)

// MockBlockchain is a mock of Blockchain interface.
type MockBlockchain struct {
	ctrl     *gomock.Controller
	recorder *MockBlockchainMockRecorder
}

// MockBlockchainMockRecorder is the mock recorder for MockBlockchain.
type MockBlockchainMockRecorder struct {
	mock *MockBlockchain
}

// NewMockBlockchain creates a new mock instance.
func NewMockBlockchain(ctrl *gomock.Controller) *MockBlockchain {
	mock := &MockBlockchain{ctrl: ctrl}
	mock.recorder = &MockBlockchainMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBlockchain) EXPECT() *MockBlockchainMockRecorder {
	return m.recorder
}

// CheckRawTransaction mocks base method.
func (m *MockBlockchain) CheckRawTransaction(arg0 context.Context, arg1 []byte) (*coretypes.ResultCheckTx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckRawTransaction", arg0, arg1)
	ret0, _ := ret[0].(*coretypes.ResultCheckTx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckRawTransaction indicates an expected call of CheckRawTransaction.
func (mr *MockBlockchainMockRecorder) CheckRawTransaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckRawTransaction", reflect.TypeOf((*MockBlockchain)(nil).CheckRawTransaction), arg0, arg1)
}

// CheckTransaction mocks base method.
func (m *MockBlockchain) CheckTransaction(arg0 context.Context, arg1 *v1.Transaction) (*coretypes.ResultCheckTx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckTransaction", arg0, arg1)
	ret0, _ := ret[0].(*coretypes.ResultCheckTx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckTransaction indicates an expected call of CheckTransaction.
func (mr *MockBlockchainMockRecorder) CheckTransaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckTransaction", reflect.TypeOf((*MockBlockchain)(nil).CheckTransaction), arg0, arg1)
}

// GetChainID mocks base method.
func (m *MockBlockchain) GetChainID(arg0 context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChainID", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChainID indicates an expected call of GetChainID.
func (mr *MockBlockchainMockRecorder) GetChainID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChainID", reflect.TypeOf((*MockBlockchain)(nil).GetChainID), arg0)
}

// GetGenesisTime mocks base method.
func (m *MockBlockchain) GetGenesisTime(arg0 context.Context) (time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGenesisTime", arg0)
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGenesisTime indicates an expected call of GetGenesisTime.
func (mr *MockBlockchainMockRecorder) GetGenesisTime(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGenesisTime", reflect.TypeOf((*MockBlockchain)(nil).GetGenesisTime), arg0)
}

// GetNetworkInfo mocks base method.
func (m *MockBlockchain) GetNetworkInfo(arg0 context.Context) (*coretypes.ResultNetInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNetworkInfo", arg0)
	ret0, _ := ret[0].(*coretypes.ResultNetInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNetworkInfo indicates an expected call of GetNetworkInfo.
func (mr *MockBlockchainMockRecorder) GetNetworkInfo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNetworkInfo", reflect.TypeOf((*MockBlockchain)(nil).GetNetworkInfo), arg0)
}

// GetStatus mocks base method.
func (m *MockBlockchain) GetStatus(arg0 context.Context) (*coretypes.ResultStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStatus", arg0)
	ret0, _ := ret[0].(*coretypes.ResultStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStatus indicates an expected call of GetStatus.
func (mr *MockBlockchainMockRecorder) GetStatus(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatus", reflect.TypeOf((*MockBlockchain)(nil).GetStatus), arg0)
}

// GetUnconfirmedTxCount mocks base method.
func (m *MockBlockchain) GetUnconfirmedTxCount(arg0 context.Context) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnconfirmedTxCount", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnconfirmedTxCount indicates an expected call of GetUnconfirmedTxCount.
func (mr *MockBlockchainMockRecorder) GetUnconfirmedTxCount(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnconfirmedTxCount", reflect.TypeOf((*MockBlockchain)(nil).GetUnconfirmedTxCount), arg0)
}

// Health mocks base method.
func (m *MockBlockchain) Health() (*coretypes.ResultHealth, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Health")
	ret0, _ := ret[0].(*coretypes.ResultHealth)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Health indicates an expected call of Health.
func (mr *MockBlockchainMockRecorder) Health() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Health", reflect.TypeOf((*MockBlockchain)(nil).Health))
}

// SubmitRawTransactionAsync mocks base method.
func (m *MockBlockchain) SubmitRawTransactionAsync(arg0 context.Context, arg1 []byte) (*coretypes.ResultBroadcastTx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubmitRawTransactionAsync", arg0, arg1)
	ret0, _ := ret[0].(*coretypes.ResultBroadcastTx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubmitRawTransactionAsync indicates an expected call of SubmitRawTransactionAsync.
func (mr *MockBlockchainMockRecorder) SubmitRawTransactionAsync(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubmitRawTransactionAsync", reflect.TypeOf((*MockBlockchain)(nil).SubmitRawTransactionAsync), arg0, arg1)
}

// SubmitRawTransactionCommit mocks base method.
func (m *MockBlockchain) SubmitRawTransactionCommit(arg0 context.Context, arg1 []byte) (*coretypes.ResultBroadcastTxCommit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubmitRawTransactionCommit", arg0, arg1)
	ret0, _ := ret[0].(*coretypes.ResultBroadcastTxCommit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubmitRawTransactionCommit indicates an expected call of SubmitRawTransactionCommit.
func (mr *MockBlockchainMockRecorder) SubmitRawTransactionCommit(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubmitRawTransactionCommit", reflect.TypeOf((*MockBlockchain)(nil).SubmitRawTransactionCommit), arg0, arg1)
}

// SubmitRawTransactionSync mocks base method.
func (m *MockBlockchain) SubmitRawTransactionSync(arg0 context.Context, arg1 []byte) (*coretypes.ResultBroadcastTx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubmitRawTransactionSync", arg0, arg1)
	ret0, _ := ret[0].(*coretypes.ResultBroadcastTx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubmitRawTransactionSync indicates an expected call of SubmitRawTransactionSync.
func (mr *MockBlockchainMockRecorder) SubmitRawTransactionSync(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubmitRawTransactionSync", reflect.TypeOf((*MockBlockchain)(nil).SubmitRawTransactionSync), arg0, arg1)
}

// SubmitTransactionAsync mocks base method.
func (m *MockBlockchain) SubmitTransactionAsync(arg0 context.Context, arg1 *v1.Transaction) (*coretypes.ResultBroadcastTx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubmitTransactionAsync", arg0, arg1)
	ret0, _ := ret[0].(*coretypes.ResultBroadcastTx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubmitTransactionAsync indicates an expected call of SubmitTransactionAsync.
func (mr *MockBlockchainMockRecorder) SubmitTransactionAsync(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubmitTransactionAsync", reflect.TypeOf((*MockBlockchain)(nil).SubmitTransactionAsync), arg0, arg1)
}

// SubmitTransactionCommit mocks base method.
func (m *MockBlockchain) SubmitTransactionCommit(arg0 context.Context, arg1 *v1.Transaction) (*coretypes.ResultBroadcastTxCommit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubmitTransactionCommit", arg0, arg1)
	ret0, _ := ret[0].(*coretypes.ResultBroadcastTxCommit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubmitTransactionCommit indicates an expected call of SubmitTransactionCommit.
func (mr *MockBlockchainMockRecorder) SubmitTransactionCommit(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubmitTransactionCommit", reflect.TypeOf((*MockBlockchain)(nil).SubmitTransactionCommit), arg0, arg1)
}

// SubmitTransactionSync mocks base method.
func (m *MockBlockchain) SubmitTransactionSync(arg0 context.Context, arg1 *v1.Transaction) (*coretypes.ResultBroadcastTx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubmitTransactionSync", arg0, arg1)
	ret0, _ := ret[0].(*coretypes.ResultBroadcastTx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubmitTransactionSync indicates an expected call of SubmitTransactionSync.
func (mr *MockBlockchainMockRecorder) SubmitTransactionSync(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubmitTransactionSync", reflect.TypeOf((*MockBlockchain)(nil).SubmitTransactionSync), arg0, arg1)
}
