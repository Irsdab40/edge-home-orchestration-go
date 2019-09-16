/*******************************************************************************
 * Copyright 2019 Samsung Electronics All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *******************************************************************************/

// Code generated by MockGen. DO NOT EDIT.
// Source: discovery.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockDiscovery is a mock of Discovery interface
type MockDiscovery struct {
	ctrl     *gomock.Controller
	recorder *MockDiscoveryMockRecorder
}

// MockDiscoveryMockRecorder is the mock recorder for MockDiscovery
type MockDiscoveryMockRecorder struct {
	mock *MockDiscovery
}

// NewMockDiscovery creates a new mock instance
func NewMockDiscovery(ctrl *gomock.Controller) *MockDiscovery {
	mock := &MockDiscovery{ctrl: ctrl}
	mock.recorder = &MockDiscoveryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDiscovery) EXPECT() *MockDiscoveryMockRecorder {
	return m.recorder
}

// StartDiscovery mocks base method
func (m *MockDiscovery) StartDiscovery(UUIDpath, platform, executionType string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartDiscovery", UUIDpath, platform, executionType)
	ret0, _ := ret[0].(error)
	return ret0
}

// StartDiscovery indicates an expected call of StartDiscovery
func (mr *MockDiscoveryMockRecorder) StartDiscovery(UUIDpath, platform, executionType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartDiscovery", reflect.TypeOf((*MockDiscovery)(nil).StartDiscovery), UUIDpath, platform, executionType)
}

// StopDiscovery mocks base method
func (m *MockDiscovery) StopDiscovery() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StopDiscovery")
}

// StopDiscovery indicates an expected call of StopDiscovery
func (mr *MockDiscoveryMockRecorder) StopDiscovery() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopDiscovery", reflect.TypeOf((*MockDiscovery)(nil).StopDiscovery))
}

// DeleteDeviceWithIP mocks base method
func (m *MockDiscovery) DeleteDeviceWithIP(targetIP string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteDeviceWithIP", targetIP)
}

// DeleteDeviceWithIP indicates an expected call of DeleteDeviceWithIP
func (mr *MockDiscoveryMockRecorder) DeleteDeviceWithIP(targetIP interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteDeviceWithIP", reflect.TypeOf((*MockDiscovery)(nil).DeleteDeviceWithIP), targetIP)
}

// DeleteDeviceWithID mocks base method
func (m *MockDiscovery) DeleteDeviceWithID(ID string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteDeviceWithID", ID)
}

// DeleteDeviceWithID indicates an expected call of DeleteDeviceWithID
func (mr *MockDiscoveryMockRecorder) DeleteDeviceWithID(ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteDeviceWithID", reflect.TypeOf((*MockDiscovery)(nil).DeleteDeviceWithID), ID)
}

// AddNewServiceName mocks base method
func (m *MockDiscovery) AddNewServiceName(serviceName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNewServiceName", serviceName)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddNewServiceName indicates an expected call of AddNewServiceName
func (mr *MockDiscoveryMockRecorder) AddNewServiceName(serviceName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNewServiceName", reflect.TypeOf((*MockDiscovery)(nil).AddNewServiceName), serviceName)
}

// RemoveServiceName mocks base method
func (m *MockDiscovery) RemoveServiceName(serviceName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveServiceName", serviceName)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveServiceName indicates an expected call of RemoveServiceName
func (mr *MockDiscoveryMockRecorder) RemoveServiceName(serviceName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveServiceName", reflect.TypeOf((*MockDiscovery)(nil).RemoveServiceName), serviceName)
}

// ResetServiceName mocks base method
func (m *MockDiscovery) ResetServiceName() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ResetServiceName")
}

// ResetServiceName indicates an expected call of ResetServiceName
func (mr *MockDiscoveryMockRecorder) ResetServiceName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetServiceName", reflect.TypeOf((*MockDiscovery)(nil).ResetServiceName))
}
