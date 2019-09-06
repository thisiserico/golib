// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/thisiserico/golib/log (interfaces: Logger)

// Package log is a generated GoMock package.
package log

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockLogger is a mock of Logger interface
type MockLogger struct {
	ctrl     *gomock.Controller
	recorder *MockLoggerMockRecorder
}

// MockLoggerMockRecorder is the mock recorder for MockLogger
type MockLoggerMockRecorder struct {
	mock *MockLogger
}

// NewMockLogger creates a new mock instance
func NewMockLogger(ctrl *gomock.Controller) *MockLogger {
	mock := &MockLogger{ctrl: ctrl}
	mock.recorder = &MockLoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLogger) EXPECT() *MockLoggerMockRecorder {
	return m.recorder
}

// Error mocks base method
func (m *MockLogger) Error(arg0 context.Context, arg1 error, arg2 Tags) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Error", arg0, arg1, arg2)
}

// Error indicates an expected call of Error
func (mr *MockLoggerMockRecorder) Error(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockLogger)(nil).Error), arg0, arg1, arg2)
}

// Fatal mocks base method
func (m *MockLogger) Fatal(arg0 context.Context, arg1 error, arg2 Tags) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Fatal", arg0, arg1, arg2)
}

// Fatal indicates an expected call of Fatal
func (mr *MockLoggerMockRecorder) Fatal(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fatal", reflect.TypeOf((*MockLogger)(nil).Fatal), arg0, arg1, arg2)
}

// Info mocks base method
func (m *MockLogger) Info(arg0 context.Context, arg1 string, arg2 Tags) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Info", arg0, arg1, arg2)
}

// Info indicates an expected call of Info
func (mr *MockLoggerMockRecorder) Info(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockLogger)(nil).Info), arg0, arg1, arg2)
}
