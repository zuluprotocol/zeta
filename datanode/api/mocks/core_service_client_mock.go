// Code generated by MockGen. DO NOT EDIT.
// Source: zuluprotocol/zeta/datanode/api (interfaces: CoreServiceClient)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	v1 "zuluprotocol/zeta/protos/zeta/api/v1"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockCoreServiceClient is a mock of CoreServiceClient interface.
type MockCoreServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockCoreServiceClientMockRecorder
}

// MockCoreServiceClientMockRecorder is the mock recorder for MockCoreServiceClient.
type MockCoreServiceClientMockRecorder struct {
	mock *MockCoreServiceClient
}

// NewMockCoreServiceClient creates a new mock instance.
func NewMockCoreServiceClient(ctrl *gomock.Controller) *MockCoreServiceClient {
	mock := &MockCoreServiceClient{ctrl: ctrl}
	mock.recorder = &MockCoreServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCoreServiceClient) EXPECT() *MockCoreServiceClientMockRecorder {
	return m.recorder
}

// CheckRawTransaction mocks base method.
func (m *MockCoreServiceClient) CheckRawTransaction(arg0 context.Context, arg1 *v1.CheckRawTransactionRequest, arg2 ...grpc.CallOption) (*v1.CheckRawTransactionResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CheckRawTransaction", varargs...)
	ret0, _ := ret[0].(*v1.CheckRawTransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckRawTransaction indicates an expected call of CheckRawTransaction.
func (mr *MockCoreServiceClientMockRecorder) CheckRawTransaction(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckRawTransaction", reflect.TypeOf((*MockCoreServiceClient)(nil).CheckRawTransaction), varargs...)
}

// CheckTransaction mocks base method.
func (m *MockCoreServiceClient) CheckTransaction(arg0 context.Context, arg1 *v1.CheckTransactionRequest, arg2 ...grpc.CallOption) (*v1.CheckTransactionResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CheckTransaction", varargs...)
	ret0, _ := ret[0].(*v1.CheckTransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckTransaction indicates an expected call of CheckTransaction.
func (mr *MockCoreServiceClientMockRecorder) CheckTransaction(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckTransaction", reflect.TypeOf((*MockCoreServiceClient)(nil).CheckTransaction), varargs...)
}

// GetSpamStatistics mocks base method.
func (m *MockCoreServiceClient) GetSpamStatistics(arg0 context.Context, arg1 *v1.GetSpamStatisticsRequest, arg2 ...grpc.CallOption) (*v1.GetSpamStatisticsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetSpamStatistics", varargs...)
	ret0, _ := ret[0].(*v1.GetSpamStatisticsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSpamStatistics indicates an expected call of GetSpamStatistics.
func (mr *MockCoreServiceClientMockRecorder) GetSpamStatistics(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSpamStatistics", reflect.TypeOf((*MockCoreServiceClient)(nil).GetSpamStatistics), varargs...)
}

// GetZetaTime mocks base method.
func (m *MockCoreServiceClient) GetZetaTime(arg0 context.Context, arg1 *v1.GetZetaTimeRequest, arg2 ...grpc.CallOption) (*v1.GetZetaTimeResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetZetaTime", varargs...)
	ret0, _ := ret[0].(*v1.GetZetaTimeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetZetaTime indicates an expected call of GetZetaTime.
func (mr *MockCoreServiceClientMockRecorder) GetZetaTime(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetZetaTime", reflect.TypeOf((*MockCoreServiceClient)(nil).GetZetaTime), varargs...)
}

// LastBlockHeight mocks base method.
func (m *MockCoreServiceClient) LastBlockHeight(arg0 context.Context, arg1 *v1.LastBlockHeightRequest, arg2 ...grpc.CallOption) (*v1.LastBlockHeightResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "LastBlockHeight", varargs...)
	ret0, _ := ret[0].(*v1.LastBlockHeightResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LastBlockHeight indicates an expected call of LastBlockHeight.
func (mr *MockCoreServiceClientMockRecorder) LastBlockHeight(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastBlockHeight", reflect.TypeOf((*MockCoreServiceClient)(nil).LastBlockHeight), varargs...)
}

// ObserveEventBus mocks base method.
func (m *MockCoreServiceClient) ObserveEventBus(arg0 context.Context, arg1 ...grpc.CallOption) (v1.CoreService_ObserveEventBusClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ObserveEventBus", varargs...)
	ret0, _ := ret[0].(v1.CoreService_ObserveEventBusClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ObserveEventBus indicates an expected call of ObserveEventBus.
func (mr *MockCoreServiceClientMockRecorder) ObserveEventBus(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ObserveEventBus", reflect.TypeOf((*MockCoreServiceClient)(nil).ObserveEventBus), varargs...)
}

// PropagateChainEvent mocks base method.
func (m *MockCoreServiceClient) PropagateChainEvent(arg0 context.Context, arg1 *v1.PropagateChainEventRequest, arg2 ...grpc.CallOption) (*v1.PropagateChainEventResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PropagateChainEvent", varargs...)
	ret0, _ := ret[0].(*v1.PropagateChainEventResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PropagateChainEvent indicates an expected call of PropagateChainEvent.
func (mr *MockCoreServiceClientMockRecorder) PropagateChainEvent(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PropagateChainEvent", reflect.TypeOf((*MockCoreServiceClient)(nil).PropagateChainEvent), varargs...)
}

// Statistics mocks base method.
func (m *MockCoreServiceClient) Statistics(arg0 context.Context, arg1 *v1.StatisticsRequest, arg2 ...grpc.CallOption) (*v1.StatisticsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Statistics", varargs...)
	ret0, _ := ret[0].(*v1.StatisticsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Statistics indicates an expected call of Statistics.
func (mr *MockCoreServiceClientMockRecorder) Statistics(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Statistics", reflect.TypeOf((*MockCoreServiceClient)(nil).Statistics), varargs...)
}

// SubmitRawTransaction mocks base method.
func (m *MockCoreServiceClient) SubmitRawTransaction(arg0 context.Context, arg1 *v1.SubmitRawTransactionRequest, arg2 ...grpc.CallOption) (*v1.SubmitRawTransactionResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SubmitRawTransaction", varargs...)
	ret0, _ := ret[0].(*v1.SubmitRawTransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubmitRawTransaction indicates an expected call of SubmitRawTransaction.
func (mr *MockCoreServiceClientMockRecorder) SubmitRawTransaction(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubmitRawTransaction", reflect.TypeOf((*MockCoreServiceClient)(nil).SubmitRawTransaction), varargs...)
}

// SubmitTransaction mocks base method.
func (m *MockCoreServiceClient) SubmitTransaction(arg0 context.Context, arg1 *v1.SubmitTransactionRequest, arg2 ...grpc.CallOption) (*v1.SubmitTransactionResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SubmitTransaction", varargs...)
	ret0, _ := ret[0].(*v1.SubmitTransactionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubmitTransaction indicates an expected call of SubmitTransaction.
func (mr *MockCoreServiceClientMockRecorder) SubmitTransaction(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubmitTransaction", reflect.TypeOf((*MockCoreServiceClient)(nil).SubmitTransaction), varargs...)
}
