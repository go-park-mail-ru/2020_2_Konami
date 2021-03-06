// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package session is a generated GoMock package.
package session

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockRepository is a mock of Repository interface
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// GetUserId mocks base method
func (m *MockRepository) GetUserId(token string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserId", token)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserId indicates an expected call of GetUserId
func (mr *MockRepositoryMockRecorder) GetUserId(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserId", reflect.TypeOf((*MockRepository)(nil).GetUserId), token)
}

// CreateSession mocks base method
func (m *MockRepository) CreateSession(userId int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", userId)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSession indicates an expected call of CreateSession
func (mr *MockRepositoryMockRecorder) CreateSession(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockRepository)(nil).CreateSession), userId)
}

// RemoveSession mocks base method
func (m *MockRepository) RemoveSession(token string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveSession", token)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveSession indicates an expected call of RemoveSession
func (mr *MockRepositoryMockRecorder) RemoveSession(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveSession", reflect.TypeOf((*MockRepository)(nil).RemoveSession), token)
}
