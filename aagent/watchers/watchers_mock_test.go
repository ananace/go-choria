// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/choria-io/go-choria/aagent/watchers (interfaces: Machine,Watcher,WatcherConstructor)

// mockgen -package watchers -destination watchers_mock_test.go github.com/choria-io/go-choria/aagent/watchers Machine,Watcher,WatcherConstructor  

// Package watchers is a generated GoMock package.
package watchers

import (
	context "context"
	reflect "reflect"
	sync "sync"
	time "time"

	watcher "github.com/choria-io/go-choria/aagent/watchers/watcher"
	gomock "github.com/golang/mock/gomock"
	jsm_go "github.com/nats-io/jsm.go"
)

// MockMachine is a mock of Machine interface.
type MockMachine struct {
	ctrl     *gomock.Controller
	recorder *MockMachineMockRecorder
}

// MockMachineMockRecorder is the mock recorder for MockMachine.
type MockMachineMockRecorder struct {
	mock *MockMachine
}

// NewMockMachine creates a new mock instance.
func NewMockMachine(ctrl *gomock.Controller) *MockMachine {
	mock := &MockMachine{ctrl: ctrl}
	mock.recorder = &MockMachineMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMachine) EXPECT() *MockMachineMockRecorder {
	return m.recorder
}

// ChoriaStatusFile mocks base method.
func (m *MockMachine) ChoriaStatusFile() (string, int) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChoriaStatusFile")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(int)
	return ret0, ret1
}

// ChoriaStatusFile indicates an expected call of ChoriaStatusFile.
func (mr *MockMachineMockRecorder) ChoriaStatusFile() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChoriaStatusFile", reflect.TypeOf((*MockMachine)(nil).ChoriaStatusFile))
}

// Debugf mocks base method.
func (m *MockMachine) Debugf(arg0, arg1 string, arg2 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Debugf", varargs...)
}

// Debugf indicates an expected call of Debugf.
func (mr *MockMachineMockRecorder) Debugf(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debugf", reflect.TypeOf((*MockMachine)(nil).Debugf), varargs...)
}

// Directory mocks base method.
func (m *MockMachine) Directory() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Directory")
	ret0, _ := ret[0].(string)
	return ret0
}

// Directory indicates an expected call of Directory.
func (mr *MockMachineMockRecorder) Directory() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Directory", reflect.TypeOf((*MockMachine)(nil).Directory))
}

// Errorf mocks base method.
func (m *MockMachine) Errorf(arg0, arg1 string, arg2 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Errorf", varargs...)
}

// Errorf indicates an expected call of Errorf.
func (mr *MockMachineMockRecorder) Errorf(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Errorf", reflect.TypeOf((*MockMachine)(nil).Errorf), varargs...)
}

// Identity mocks base method.
func (m *MockMachine) Identity() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Identity")
	ret0, _ := ret[0].(string)
	return ret0
}

// Identity indicates an expected call of Identity.
func (mr *MockMachineMockRecorder) Identity() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Identity", reflect.TypeOf((*MockMachine)(nil).Identity))
}

// Infof mocks base method.
func (m *MockMachine) Infof(arg0, arg1 string, arg2 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Infof", varargs...)
}

// Infof indicates an expected call of Infof.
func (mr *MockMachineMockRecorder) Infof(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Infof", reflect.TypeOf((*MockMachine)(nil).Infof), varargs...)
}

// InstanceID mocks base method.
func (m *MockMachine) InstanceID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InstanceID")
	ret0, _ := ret[0].(string)
	return ret0
}

// InstanceID indicates an expected call of InstanceID.
func (mr *MockMachineMockRecorder) InstanceID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InstanceID", reflect.TypeOf((*MockMachine)(nil).InstanceID))
}

// JetStreamConnection mocks base method.
func (m *MockMachine) JetStreamConnection() (*jsm_go.Manager, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "JetStreamConnection")
	ret0, _ := ret[0].(*jsm_go.Manager)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// JetStreamConnection indicates an expected call of JetStreamConnection.
func (mr *MockMachineMockRecorder) JetStreamConnection() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "JetStreamConnection", reflect.TypeOf((*MockMachine)(nil).JetStreamConnection))
}

// MainCollective mocks base method.
func (m *MockMachine) MainCollective() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MainCollective")
	ret0, _ := ret[0].(string)
	return ret0
}

// MainCollective indicates an expected call of MainCollective.
func (mr *MockMachineMockRecorder) MainCollective() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MainCollective", reflect.TypeOf((*MockMachine)(nil).MainCollective))
}

// Name mocks base method.
func (m *MockMachine) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockMachineMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockMachine)(nil).Name))
}

// NotifyWatcherState mocks base method.
func (m *MockMachine) NotifyWatcherState(arg0 string, arg1 interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "NotifyWatcherState", arg0, arg1)
}

// NotifyWatcherState indicates an expected call of NotifyWatcherState.
func (mr *MockMachineMockRecorder) NotifyWatcherState(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NotifyWatcherState", reflect.TypeOf((*MockMachine)(nil).NotifyWatcherState), arg0, arg1)
}

// OverrideData mocks base method.
func (m *MockMachine) OverrideData() ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OverrideData")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OverrideData indicates an expected call of OverrideData.
func (mr *MockMachineMockRecorder) OverrideData() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OverrideData", reflect.TypeOf((*MockMachine)(nil).OverrideData))
}

// State mocks base method.
func (m *MockMachine) State() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "State")
	ret0, _ := ret[0].(string)
	return ret0
}

// State indicates an expected call of State.
func (mr *MockMachineMockRecorder) State() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "State", reflect.TypeOf((*MockMachine)(nil).State))
}

// TextFileDirectory mocks base method.
func (m *MockMachine) TextFileDirectory() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TextFileDirectory")
	ret0, _ := ret[0].(string)
	return ret0
}

// TextFileDirectory indicates an expected call of TextFileDirectory.
func (mr *MockMachineMockRecorder) TextFileDirectory() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TextFileDirectory", reflect.TypeOf((*MockMachine)(nil).TextFileDirectory))
}

// TimeStampSeconds mocks base method.
func (m *MockMachine) TimeStampSeconds() int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TimeStampSeconds")
	ret0, _ := ret[0].(int64)
	return ret0
}

// TimeStampSeconds indicates an expected call of TimeStampSeconds.
func (mr *MockMachineMockRecorder) TimeStampSeconds() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TimeStampSeconds", reflect.TypeOf((*MockMachine)(nil).TimeStampSeconds))
}

// Transition mocks base method.
func (m *MockMachine) Transition(arg0 string, arg1 ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Transition", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Transition indicates an expected call of Transition.
func (mr *MockMachineMockRecorder) Transition(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Transition", reflect.TypeOf((*MockMachine)(nil).Transition), varargs...)
}

// Version mocks base method.
func (m *MockMachine) Version() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Version")
	ret0, _ := ret[0].(string)
	return ret0
}

// Version indicates an expected call of Version.
func (mr *MockMachineMockRecorder) Version() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Version", reflect.TypeOf((*MockMachine)(nil).Version))
}

// Warnf mocks base method.
func (m *MockMachine) Warnf(arg0, arg1 string, arg2 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Warnf", varargs...)
}

// Warnf indicates an expected call of Warnf.
func (mr *MockMachineMockRecorder) Warnf(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warnf", reflect.TypeOf((*MockMachine)(nil).Warnf), varargs...)
}

// Watchers mocks base method.
func (m *MockMachine) Watchers() []*WatcherDef {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Watchers")
	ret0, _ := ret[0].([]*WatcherDef)
	return ret0
}

// Watchers indicates an expected call of Watchers.
func (mr *MockMachineMockRecorder) Watchers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Watchers", reflect.TypeOf((*MockMachine)(nil).Watchers))
}

// MockWatcher is a mock of Watcher interface.
type MockWatcher struct {
	ctrl     *gomock.Controller
	recorder *MockWatcherMockRecorder
}

// MockWatcherMockRecorder is the mock recorder for MockWatcher.
type MockWatcherMockRecorder struct {
	mock *MockWatcher
}

// NewMockWatcher creates a new mock instance.
func NewMockWatcher(ctrl *gomock.Controller) *MockWatcher {
	mock := &MockWatcher{ctrl: ctrl}
	mock.recorder = &MockWatcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWatcher) EXPECT() *MockWatcherMockRecorder {
	return m.recorder
}

// AnnounceInterval mocks base method.
func (m *MockWatcher) AnnounceInterval() time.Duration {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AnnounceInterval")
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

// AnnounceInterval indicates an expected call of AnnounceInterval.
func (mr *MockWatcherMockRecorder) AnnounceInterval() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AnnounceInterval", reflect.TypeOf((*MockWatcher)(nil).AnnounceInterval))
}

// CurrentState mocks base method.
func (m *MockWatcher) CurrentState() interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CurrentState")
	ret0, _ := ret[0].(interface{})
	return ret0
}

// CurrentState indicates an expected call of CurrentState.
func (mr *MockWatcherMockRecorder) CurrentState() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CurrentState", reflect.TypeOf((*MockWatcher)(nil).CurrentState))
}

// Delete mocks base method.
func (m *MockWatcher) Delete() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Delete")
}

// Delete indicates an expected call of Delete.
func (mr *MockWatcherMockRecorder) Delete() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockWatcher)(nil).Delete))
}

// Name mocks base method.
func (m *MockWatcher) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockWatcherMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockWatcher)(nil).Name))
}

// NotifyStateChance mocks base method.
func (m *MockWatcher) NotifyStateChance() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "NotifyStateChance")
}

// NotifyStateChance indicates an expected call of NotifyStateChance.
func (mr *MockWatcherMockRecorder) NotifyStateChance() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NotifyStateChance", reflect.TypeOf((*MockWatcher)(nil).NotifyStateChance))
}

// Run mocks base method.
func (m *MockWatcher) Run(arg0 context.Context, arg1 *sync.WaitGroup) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Run", arg0, arg1)
}

// Run indicates an expected call of Run.
func (mr *MockWatcherMockRecorder) Run(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockWatcher)(nil).Run), arg0, arg1)
}

// Type mocks base method.
func (m *MockWatcher) Type() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Type")
	ret0, _ := ret[0].(string)
	return ret0
}

// Type indicates an expected call of Type.
func (mr *MockWatcherMockRecorder) Type() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Type", reflect.TypeOf((*MockWatcher)(nil).Type))
}

// MockWatcherConstructor is a mock of WatcherConstructor interface.
type MockWatcherConstructor struct {
	ctrl     *gomock.Controller
	recorder *MockWatcherConstructorMockRecorder
}

// MockWatcherConstructorMockRecorder is the mock recorder for MockWatcherConstructor.
type MockWatcherConstructorMockRecorder struct {
	mock *MockWatcherConstructor
}

// NewMockWatcherConstructor creates a new mock instance.
func NewMockWatcherConstructor(ctrl *gomock.Controller) *MockWatcherConstructor {
	mock := &MockWatcherConstructor{ctrl: ctrl}
	mock.recorder = &MockWatcherConstructorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWatcherConstructor) EXPECT() *MockWatcherConstructorMockRecorder {
	return m.recorder
}

// EventType mocks base method.
func (m *MockWatcherConstructor) EventType() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EventType")
	ret0, _ := ret[0].(string)
	return ret0
}

// EventType indicates an expected call of EventType.
func (mr *MockWatcherConstructorMockRecorder) EventType() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EventType", reflect.TypeOf((*MockWatcherConstructor)(nil).EventType))
}

// New mocks base method.
func (m *MockWatcherConstructor) New(arg0 watcher.Machine, arg1 string, arg2 []string, arg3, arg4, arg5 string, arg6 time.Duration, arg7 map[string]interface{}) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "New", arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// New indicates an expected call of New.
func (mr *MockWatcherConstructorMockRecorder) New(arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "New", reflect.TypeOf((*MockWatcherConstructor)(nil).New), arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7)
}

// Type mocks base method.
func (m *MockWatcherConstructor) Type() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Type")
	ret0, _ := ret[0].(string)
	return ret0
}

// Type indicates an expected call of Type.
func (mr *MockWatcherConstructorMockRecorder) Type() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Type", reflect.TypeOf((*MockWatcherConstructor)(nil).Type))
}

// UnmarshalNotification mocks base method.
func (m *MockWatcherConstructor) UnmarshalNotification(arg0 []byte) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnmarshalNotification", arg0)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UnmarshalNotification indicates an expected call of UnmarshalNotification.
func (mr *MockWatcherConstructorMockRecorder) UnmarshalNotification(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnmarshalNotification", reflect.TypeOf((*MockWatcherConstructor)(nil).UnmarshalNotification), arg0)
}
