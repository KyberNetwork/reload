// Code generated by MockGen. DO NOT EDIT.
// Source: dmm-aggregator-backend/pkg/reload (interfaces: Reloader)

// Package reload is a generated GoMock package.
package reload

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockReloader is a mock of Reloader interface.
type MockReloader struct {
	ctrl     *gomock.Controller
	recorder *MockReloaderMockRecorder
}

// MockReloaderMockRecorder is the mock recorder for MockReloader.
type MockReloaderMockRecorder struct {
	mock *MockReloader
}

// NewMockReloader creates a new mock instance.
func NewMockReloader(ctrl *gomock.Controller) *MockReloader {
	mock := &MockReloader{ctrl: ctrl}
	mock.recorder = &MockReloaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReloader) EXPECT() *MockReloaderMockRecorder {
	return m.recorder
}

// Reload mocks base method.
func (m *MockReloader) Reload(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reload", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Reload indicates an expected call of Reload.
func (mr *MockReloaderMockRecorder) Reload(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reload", reflect.TypeOf((*MockReloader)(nil).Reload), arg0, arg1)
}