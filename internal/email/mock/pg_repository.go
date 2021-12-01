// Code generated by MockGen. DO NOT EDIT.
// Source: pg_repository.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	models "github.com/email-microservice/internal/models"
	utils "github.com/email-microservice/pkg/utils"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	reflect "reflect"
)

// MockEmailsRepository is a mock of EmailsRepository interface
type MockEmailsRepository struct {
	ctrl     *gomock.Controller
	recorder *MockEmailsRepositoryMockRecorder
}

// MockEmailsRepositoryMockRecorder is the mock recorder for MockEmailsRepository
type MockEmailsRepositoryMockRecorder struct {
	mock *MockEmailsRepository
}

// NewMockEmailsRepository creates a new mock instance
func NewMockEmailsRepository(ctrl *gomock.Controller) *MockEmailsRepository {
	mock := &MockEmailsRepository{ctrl: ctrl}
	mock.recorder = &MockEmailsRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockEmailsRepository) EXPECT() *MockEmailsRepositoryMockRecorder {
	return m.recorder
}

// CreateEmail mocks base method
func (m *MockEmailsRepository) CreateEmail(ctx context.Context, email *models.Email) (*models.Email, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEmail", ctx, email)
	ret0, _ := ret[0].(*models.Email)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateEmail indicates an expected call of CreateEmail
func (mr *MockEmailsRepositoryMockRecorder) CreateEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEmail", reflect.TypeOf((*MockEmailsRepository)(nil).CreateEmail), ctx, email)
}

// FindEmailById mocks base method
func (m *MockEmailsRepository) FindEmailById(ctx context.Context, emailID uuid.UUID) (*models.Email, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindEmailById", ctx, emailID)
	ret0, _ := ret[0].(*models.Email)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindEmailById indicates an expected call of FindEmailById
func (mr *MockEmailsRepositoryMockRecorder) FindEmailById(ctx, emailID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindEmailById", reflect.TypeOf((*MockEmailsRepository)(nil).FindEmailById), ctx, emailID)
}

// FindEmailsByReceiver mocks base method
func (m *MockEmailsRepository) FindEmailsByReceiver(ctx context.Context, to string, query *utils.PaginationQuery) (*models.EmailsList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindEmailsByReceiver", ctx, to, query)
	ret0, _ := ret[0].(*models.EmailsList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindEmailsByReceiver indicates an expected call of FindEmailsByReceiver
func (mr *MockEmailsRepositoryMockRecorder) FindEmailsByReceiver(ctx, to, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindEmailsByReceiver", reflect.TypeOf((*MockEmailsRepository)(nil).FindEmailsByReceiver), ctx, to, query)
}
