// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/user/repository.go
//
// Generated by this command:
//
//	mockgen -source internal/repository/user/repository.go -destination internal/repository/user/repository_mock.go
//

// Package mock_userrepository is a generated GoMock package.
package userrepository

import (
	model "avito/internal/model"
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// CloseConnection mocks base method.
func (m *MockRepository) CloseConnection() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseConnection")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseConnection indicates an expected call of CloseConnection.
func (mr *MockRepositoryMockRecorder) CloseConnection() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseConnection", reflect.TypeOf((*MockRepository)(nil).CloseConnection))
}

// Save mocks base method.
func (m *MockRepository) Save(ctx context.Context, user model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockRepositoryMockRecorder) Save(ctx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockRepository)(nil).Save), ctx, user)
}

// UserByEmail mocks base method.
func (m *MockRepository) UserByEmail(ctx context.Context, email string) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserByEmail", ctx, email)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserByEmail indicates an expected call of UserByEmail.
func (mr *MockRepositoryMockRecorder) UserByEmail(ctx, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserByEmail", reflect.TypeOf((*MockRepository)(nil).UserByEmail), ctx, email)
}
