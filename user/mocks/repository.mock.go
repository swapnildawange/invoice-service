// Code generated by MockGen. DO NOT EDIT.
// Source: invoice_service/user/repository (interfaces: Repository)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	model "invoice_service/model"
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

// CreateUser mocks base method
func (m *MockRepository) CreateUser(arg0 context.Context, arg1 model.CreateUserRequest) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser
func (mr *MockRepositoryMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockRepository)(nil).CreateUser), arg0, arg1)
}

// DeleteUser mocks base method
func (m *MockRepository) DeleteUser(arg0 context.Context, arg1 model.DeleteUserReq) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser
func (mr *MockRepositoryMockRecorder) DeleteUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockRepository)(nil).DeleteUser), arg0, arg1)
}

// EditUser mocks base method
func (m *MockRepository) EditUser(arg0 context.Context, arg1 model.EditUserRequest) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditUser", arg0, arg1)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditUser indicates an expected call of EditUser
func (mr *MockRepositoryMockRecorder) EditUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditUser", reflect.TypeOf((*MockRepository)(nil).EditUser), arg0, arg1)
}

// GetUser mocks base method
func (m *MockRepository) GetUser(arg0 context.Context, arg1 int) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser
func (mr *MockRepositoryMockRecorder) GetUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockRepository)(nil).GetUser), arg0, arg1)
}

// GetUserFromAuth mocks base method
func (m *MockRepository) GetUserFromAuth(arg0 context.Context, arg1 string) (int, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserFromAuth", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetUserFromAuth indicates an expected call of GetUserFromAuth
func (mr *MockRepositoryMockRecorder) GetUserFromAuth(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserFromAuth", reflect.TypeOf((*MockRepository)(nil).GetUserFromAuth), arg0, arg1)
}

// ListUsers mocks base method
func (m *MockRepository) ListUsers(arg0 context.Context, arg1 model.UserFilter) ([]model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUsers", arg0, arg1)
	ret0, _ := ret[0].([]model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUsers indicates an expected call of ListUsers
func (mr *MockRepositoryMockRecorder) ListUsers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUsers", reflect.TypeOf((*MockRepository)(nil).ListUsers), arg0, arg1)
}