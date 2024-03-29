// Code generated by MockGen. DO NOT EDIT.
// Source: internal/taskrpg/deps.go

// Package mock_taskrpg is a generated GoMock package.
package taskrpg

import (
	reflect "reflect"
	models "rpgMonster/models"

	gomock "github.com/golang/mock/gomock"
)

// MockIoservice is a mock of Ioservice interface.
type MockIoservice struct {
	ctrl     *gomock.Controller
	recorder *MockIoserviceMockRecorder
}

// MockIoserviceMockRecorder is the mock recorder for MockIoservice.
type MockIoserviceMockRecorder struct {
	mock *MockIoservice
}

// NewMockIoservice creates a new mock instance.
func NewMockIoservice(ctrl *gomock.Controller) *MockIoservice {
	mock := &MockIoservice{ctrl: ctrl}
	mock.recorder = &MockIoserviceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIoservice) EXPECT() *MockIoserviceMockRecorder {
	return m.recorder
}

// GetTopics mocks base method.
func (m *MockIoservice) GetTopics(file string) []models.TopicDTO {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTopics", file)
	ret0, _ := ret[0].([]models.TopicDTO)
	return ret0
}

// GetTopics indicates an expected call of GetTopics.
func (mr *MockIoserviceMockRecorder) GetTopics(file interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTopics", reflect.TypeOf((*MockIoservice)(nil).GetTopics), file)
}

// LoadPlayers mocks base method.
func (m *MockIoservice) LoadPlayers(file string) []models.PlayerDTO {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadPlayers", file)
	ret0, _ := ret[0].([]models.PlayerDTO)
	return ret0
}

// LoadPlayers indicates an expected call of LoadPlayers.
func (mr *MockIoserviceMockRecorder) LoadPlayers(file interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadPlayers", reflect.TypeOf((*MockIoservice)(nil).LoadPlayers), file)
}

// SavePlayers mocks base method.
func (m *MockIoservice) SavePlayers(file string, players [][]string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SavePlayers", file, players)
}

// SavePlayers indicates an expected call of SavePlayers.
func (mr *MockIoserviceMockRecorder) SavePlayers(file, players interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SavePlayers", reflect.TypeOf((*MockIoservice)(nil).SavePlayers), file, players)
}

// SaveTopics mocks base method.
func (m *MockIoservice) SaveTopics(file string, topics [][]string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveTopics", file, topics)
}

// SaveTopics indicates an expected call of SaveTopics.
func (mr *MockIoserviceMockRecorder) SaveTopics(file, topics interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveTopics", reflect.TypeOf((*MockIoservice)(nil).SaveTopics), file, topics)
}
