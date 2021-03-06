// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go

// Package session is a generated GoMock package.
package session

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockUseCase is a mock of UseCase interface
type MockUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUseCaseMockRecorder
}

// MockUseCaseMockRecorder is the mock recorder for MockUseCase
type MockUseCaseMockRecorder struct {
	mock *MockUseCase
}

// NewMockUseCase creates a new mock instance
func NewMockUseCase(ctrl *gomock.Controller) *MockUseCase {
	mock := &MockUseCase{ctrl: ctrl}
	mock.recorder = &MockUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUseCase) EXPECT() *MockUseCaseMockRecorder {
	return m.recorder
}

// GetUserId mocks base method
func (m *MockUseCase) GetUserId(token string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserId", token)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserId indicates an expected call of GetUserId
func (mr *MockUseCaseMockRecorder) GetUserId(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserId", reflect.TypeOf((*MockUseCase)(nil).GetUserId), token)
}

// CreateSession mocks base method
func (m *MockUseCase) CreateSession(userId int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", userId)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSession indicates an expected call of CreateSession
func (mr *MockUseCaseMockRecorder) CreateSession(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockUseCase)(nil).CreateSession), userId)
}

// RemoveSession mocks base method
func (m *MockUseCase) RemoveSession(token string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveSession", token)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveSession indicates an expected call of RemoveSession
func (mr *MockUseCaseMockRecorder) RemoveSession(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveSession", reflect.TypeOf((*MockUseCase)(nil).RemoveSession), token)
}
