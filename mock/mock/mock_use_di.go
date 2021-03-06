// Code generated by MockGen. DO NOT EDIT.
// Source: use_di.go

// Package mock_mock is a generated GoMock package.
package mock_mock

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockDiA is a mock of DiA interface
type MockDiA struct {
	ctrl     *gomock.Controller
	recorder *MockDiAMockRecorder
}

// MockDiAMockRecorder is the mock recorder for MockDiA
type MockDiAMockRecorder struct {
	mock *MockDiA
}

// NewMockDiA creates a new mock instance
func NewMockDiA(ctrl *gomock.Controller) *MockDiA {
	mock := &MockDiA{ctrl: ctrl}
	mock.recorder = &MockDiAMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDiA) EXPECT() *MockDiAMockRecorder {
	return m.recorder
}

// SomeMethodDiA mocks base method
func (m *MockDiA) SomeMethodDiA() (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SomeMethodDiA")
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SomeMethodDiA indicates an expected call of SomeMethodDiA
func (mr *MockDiAMockRecorder) SomeMethodDiA() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SomeMethodDiA", reflect.TypeOf((*MockDiA)(nil).SomeMethodDiA))
}
