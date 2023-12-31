// Code generated by MockGen. DO NOT EDIT.
// Source: NotificationService/internal/repo (interfaces: MessageRepoInterface)

// Package repo is a generated GoMock package.
package repo

import (
	dto "NotificationService/internal/dto"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
)

// MockMessageRepoInterface is a mock of MessageRepoInterface interface.
type MockMessageRepoInterface struct {
	ctrl     *gomock.Controller
	recorder *MockMessageRepoInterfaceMockRecorder
}

// MockMessageRepoInterfaceMockRecorder is the mock recorder for MockMessageRepoInterface.
type MockMessageRepoInterfaceMockRecorder struct {
	mock *MockMessageRepoInterface
}

// NewMockMessageRepoInterface creates a new mock instance.
func NewMockMessageRepoInterface(ctrl *gomock.Controller) *MockMessageRepoInterface {
	mock := &MockMessageRepoInterface{ctrl: ctrl}
	mock.recorder = &MockMessageRepoInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMessageRepoInterface) EXPECT() *MockMessageRepoInterfaceMockRecorder {
	return m.recorder
}

// CreateFailedMessage mocks base method.
func (m *MockMessageRepoInterface) CreateFailedMessage(arg0 *dto.Message) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFailedMessage", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// CreateFailedMessage indicates an expected call of CreateFailedMessage.
func (mr *MockMessageRepoInterfaceMockRecorder) CreateFailedMessage(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFailedMessage", reflect.TypeOf((*MockMessageRepoInterface)(nil).CreateFailedMessage), arg0)
}

// CreateSuccessMessage mocks base method.
func (m *MockMessageRepoInterface) CreateSuccessMessage(arg0 *dto.Message) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSuccessMessage", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// CreateSuccessMessage indicates an expected call of CreateSuccessMessage.
func (mr *MockMessageRepoInterfaceMockRecorder) CreateSuccessMessage(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSuccessMessage", reflect.TypeOf((*MockMessageRepoInterface)(nil).CreateSuccessMessage), arg0)
}

// GetStatus mocks base method.
func (m *MockMessageRepoInterface) GetStatus(arg0, arg1 primitive.ObjectID) dto.MessageStatus {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStatus", arg0, arg1)
	ret0, _ := ret[0].(dto.MessageStatus)
	return ret0
}

// GetStatus indicates an expected call of GetStatus.
func (mr *MockMessageRepoInterfaceMockRecorder) GetStatus(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatus", reflect.TypeOf((*MockMessageRepoInterface)(nil).GetStatus), arg0, arg1)
}
