// Code generated by MockGen. DO NOT EDIT.
// Source: zuluprotocol/zeta/zeta/datanode/networkhistory (interfaces: NetworkHistory)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	networkhistory "zuluprotocol/zeta/zeta/datanode/networkhistory"
	snapshot "zuluprotocol/zeta/zeta/datanode/networkhistory/snapshot"
	sqlstore "zuluprotocol/zeta/zeta/datanode/sqlstore"
	v2 "zuluprotocol/zeta/zeta/protos/data-node/api/v2"
	gomock "github.com/golang/mock/gomock"
)

// MockNetworkHistory is a mock of NetworkHistory interface.
type MockNetworkHistory struct {
	ctrl     *gomock.Controller
	recorder *MockNetworkHistoryMockRecorder
}

// MockNetworkHistoryMockRecorder is the mock recorder for MockNetworkHistory.
type MockNetworkHistoryMockRecorder struct {
	mock *MockNetworkHistory
}

// NewMockNetworkHistory creates a new mock instance.
func NewMockNetworkHistory(ctrl *gomock.Controller) *MockNetworkHistory {
	mock := &MockNetworkHistory{ctrl: ctrl}
	mock.recorder = &MockNetworkHistoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNetworkHistory) EXPECT() *MockNetworkHistoryMockRecorder {
	return m.recorder
}

// FetchHistorySegment mocks base method.
func (m *MockNetworkHistory) FetchHistorySegment(arg0 context.Context, arg1 string) (networkhistory.Segment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchHistorySegment", arg0, arg1)
	ret0, _ := ret[0].(networkhistory.Segment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchHistorySegment indicates an expected call of FetchHistorySegment.
func (mr *MockNetworkHistoryMockRecorder) FetchHistorySegment(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchHistorySegment", reflect.TypeOf((*MockNetworkHistory)(nil).FetchHistorySegment), arg0, arg1)
}

// GetDatanodeBlockSpan mocks base method.
func (m *MockNetworkHistory) GetDatanodeBlockSpan(arg0 context.Context) (sqlstore.DatanodeBlockSpan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDatanodeBlockSpan", arg0)
	ret0, _ := ret[0].(sqlstore.DatanodeBlockSpan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDatanodeBlockSpan indicates an expected call of GetDatanodeBlockSpan.
func (mr *MockNetworkHistoryMockRecorder) GetDatanodeBlockSpan(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDatanodeBlockSpan", reflect.TypeOf((*MockNetworkHistory)(nil).GetDatanodeBlockSpan), arg0)
}

// GetMostRecentHistorySegmentFromPeers mocks base method.
func (m *MockNetworkHistory) GetMostRecentHistorySegmentFromPeers(arg0 context.Context, arg1 []int) (*networkhistory.PeerResponse, map[string]*v2.GetMostRecentNetworkHistorySegmentResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMostRecentHistorySegmentFromPeers", arg0, arg1)
	ret0, _ := ret[0].(*networkhistory.PeerResponse)
	ret1, _ := ret[1].(map[string]*v2.GetMostRecentNetworkHistorySegmentResponse)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetMostRecentHistorySegmentFromPeers indicates an expected call of GetMostRecentHistorySegmentFromPeers.
func (mr *MockNetworkHistoryMockRecorder) GetMostRecentHistorySegmentFromPeers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMostRecentHistorySegmentFromPeers", reflect.TypeOf((*MockNetworkHistory)(nil).GetMostRecentHistorySegmentFromPeers), arg0, arg1)
}

// ListAllHistorySegments mocks base method.
func (m *MockNetworkHistory) ListAllHistorySegments() ([]networkhistory.Segment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAllHistorySegments")
	ret0, _ := ret[0].([]networkhistory.Segment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAllHistorySegments indicates an expected call of ListAllHistorySegments.
func (mr *MockNetworkHistoryMockRecorder) ListAllHistorySegments() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAllHistorySegments", reflect.TypeOf((*MockNetworkHistory)(nil).ListAllHistorySegments))
}

// LoadNetworkHistoryIntoDatanode mocks base method.
func (m *MockNetworkHistory) LoadNetworkHistoryIntoDatanode(arg0 context.Context, arg1 networkhistory.ContiguousHistory, arg2 sqlstore.ConnectionConfig, arg3, arg4 bool) (snapshot.LoadResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadNetworkHistoryIntoDatanode", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(snapshot.LoadResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadNetworkHistoryIntoDatanode indicates an expected call of LoadNetworkHistoryIntoDatanode.
func (mr *MockNetworkHistoryMockRecorder) LoadNetworkHistoryIntoDatanode(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadNetworkHistoryIntoDatanode", reflect.TypeOf((*MockNetworkHistory)(nil).LoadNetworkHistoryIntoDatanode), arg0, arg1, arg2, arg3, arg4)
}
