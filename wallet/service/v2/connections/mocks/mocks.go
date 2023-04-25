// Code generated by MockGen. DO NOT EDIT.
// Source: zuluprotocol/zeta/wallet/service/v2/connections (interfaces: WalletStore,TimeService,TokenStore)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	connections "zuluprotocol/zeta/wallet/service/v2/connections"
	wallet "zuluprotocol/zeta/wallet/wallet"
	gomock "github.com/golang/mock/gomock"
)

// MockWalletStore is a mock of WalletStore interface.
type MockWalletStore struct {
	ctrl     *gomock.Controller
	recorder *MockWalletStoreMockRecorder
}

// MockWalletStoreMockRecorder is the mock recorder for MockWalletStore.
type MockWalletStoreMockRecorder struct {
	mock *MockWalletStore
}

// NewMockWalletStore creates a new mock instance.
func NewMockWalletStore(ctrl *gomock.Controller) *MockWalletStore {
	mock := &MockWalletStore{ctrl: ctrl}
	mock.recorder = &MockWalletStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWalletStore) EXPECT() *MockWalletStoreMockRecorder {
	return m.recorder
}

// GetWallet mocks base method.
func (m *MockWalletStore) GetWallet(arg0 context.Context, arg1 string) (wallet.Wallet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWallet", arg0, arg1)
	ret0, _ := ret[0].(wallet.Wallet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWallet indicates an expected call of GetWallet.
func (mr *MockWalletStoreMockRecorder) GetWallet(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWallet", reflect.TypeOf((*MockWalletStore)(nil).GetWallet), arg0, arg1)
}

// OnUpdate mocks base method.
func (m *MockWalletStore) OnUpdate(arg0 func(context.Context, wallet.Event)) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnUpdate", arg0)
}

// OnUpdate indicates an expected call of OnUpdate.
func (mr *MockWalletStoreMockRecorder) OnUpdate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnUpdate", reflect.TypeOf((*MockWalletStore)(nil).OnUpdate), arg0)
}

// UnlockWallet mocks base method.
func (m *MockWalletStore) UnlockWallet(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnlockWallet", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnlockWallet indicates an expected call of UnlockWallet.
func (mr *MockWalletStoreMockRecorder) UnlockWallet(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnlockWallet", reflect.TypeOf((*MockWalletStore)(nil).UnlockWallet), arg0, arg1, arg2)
}

// MockTimeService is a mock of TimeService interface.
type MockTimeService struct {
	ctrl     *gomock.Controller
	recorder *MockTimeServiceMockRecorder
}

// MockTimeServiceMockRecorder is the mock recorder for MockTimeService.
type MockTimeServiceMockRecorder struct {
	mock *MockTimeService
}

// NewMockTimeService creates a new mock instance.
func NewMockTimeService(ctrl *gomock.Controller) *MockTimeService {
	mock := &MockTimeService{ctrl: ctrl}
	mock.recorder = &MockTimeServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTimeService) EXPECT() *MockTimeServiceMockRecorder {
	return m.recorder
}

// Now mocks base method.
func (m *MockTimeService) Now() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Now")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// Now indicates an expected call of Now.
func (mr *MockTimeServiceMockRecorder) Now() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Now", reflect.TypeOf((*MockTimeService)(nil).Now))
}

// MockTokenStore is a mock of TokenStore interface.
type MockTokenStore struct {
	ctrl     *gomock.Controller
	recorder *MockTokenStoreMockRecorder
}

// MockTokenStoreMockRecorder is the mock recorder for MockTokenStore.
type MockTokenStoreMockRecorder struct {
	mock *MockTokenStore
}

// NewMockTokenStore creates a new mock instance.
func NewMockTokenStore(ctrl *gomock.Controller) *MockTokenStore {
	mock := &MockTokenStore{ctrl: ctrl}
	mock.recorder = &MockTokenStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenStore) EXPECT() *MockTokenStoreMockRecorder {
	return m.recorder
}

// DeleteToken mocks base method.
func (m *MockTokenStore) DeleteToken(arg0 connections.Token) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteToken", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteToken indicates an expected call of DeleteToken.
func (mr *MockTokenStoreMockRecorder) DeleteToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteToken", reflect.TypeOf((*MockTokenStore)(nil).DeleteToken), arg0)
}

// DescribeToken mocks base method.
func (m *MockTokenStore) DescribeToken(arg0 connections.Token) (connections.TokenDescription, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DescribeToken", arg0)
	ret0, _ := ret[0].(connections.TokenDescription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeToken indicates an expected call of DescribeToken.
func (mr *MockTokenStoreMockRecorder) DescribeToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeToken", reflect.TypeOf((*MockTokenStore)(nil).DescribeToken), arg0)
}

// ListTokens mocks base method.
func (m *MockTokenStore) ListTokens() ([]connections.TokenSummary, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTokens")
	ret0, _ := ret[0].([]connections.TokenSummary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTokens indicates an expected call of ListTokens.
func (mr *MockTokenStoreMockRecorder) ListTokens() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTokens", reflect.TypeOf((*MockTokenStore)(nil).ListTokens))
}

// OnUpdate mocks base method.
func (m *MockTokenStore) OnUpdate(arg0 func(context.Context, ...connections.TokenDescription)) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnUpdate", arg0)
}

// OnUpdate indicates an expected call of OnUpdate.
func (mr *MockTokenStoreMockRecorder) OnUpdate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnUpdate", reflect.TypeOf((*MockTokenStore)(nil).OnUpdate), arg0)
}

// SaveToken mocks base method.
func (m *MockTokenStore) SaveToken(arg0 connections.TokenDescription) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveToken", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveToken indicates an expected call of SaveToken.
func (mr *MockTokenStoreMockRecorder) SaveToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveToken", reflect.TypeOf((*MockTokenStore)(nil).SaveToken), arg0)
}

// TokenExists mocks base method.
func (m *MockTokenStore) TokenExists(arg0 connections.Token) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TokenExists", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TokenExists indicates an expected call of TokenExists.
func (mr *MockTokenStoreMockRecorder) TokenExists(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TokenExists", reflect.TypeOf((*MockTokenStore)(nil).TokenExists), arg0)
}
