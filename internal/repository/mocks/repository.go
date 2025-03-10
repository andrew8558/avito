// Code generated by MockGen. DO NOT EDIT.
// Source: ./repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	model "Avito/internal/model"
	repository "Avito/internal/repository"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	pgx "github.com/jackc/pgx/v4"
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

// AddSendCoinEvent mocks base method.
func (m *MockRepository) AddSendCoinEvent(ctx context.Context, tx pgx.Tx, fromUser, toUser string, amount int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSendCoinEvent", ctx, tx, fromUser, toUser, amount)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddSendCoinEvent indicates an expected call of AddSendCoinEvent.
func (mr *MockRepositoryMockRecorder) AddSendCoinEvent(ctx, tx, fromUser, toUser, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSendCoinEvent", reflect.TypeOf((*MockRepository)(nil).AddSendCoinEvent), ctx, tx, fromUser, toUser, amount)
}

// AddUser mocks base method.
func (m *MockRepository) AddUser(ctx context.Context, tx pgx.Tx, login, password string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", ctx, tx, login, password)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUser indicates an expected call of AddUser.
func (mr *MockRepositoryMockRecorder) AddUser(ctx, tx, login, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockRepository)(nil).AddUser), ctx, tx, login, password)
}

// BeginTransaction mocks base method.
func (m *MockRepository) BeginTransaction(ctx context.Context, options *pgx.TxOptions) (pgx.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeginTransaction", ctx, options)
	ret0, _ := ret[0].(pgx.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BeginTransaction indicates an expected call of BeginTransaction.
func (mr *MockRepositoryMockRecorder) BeginTransaction(ctx, options interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeginTransaction", reflect.TypeOf((*MockRepository)(nil).BeginTransaction), ctx, options)
}

// CommitTx mocks base method.
func (m *MockRepository) CommitTx(ctx context.Context, tx pgx.Tx) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CommitTx", ctx, tx)
}

// CommitTx indicates an expected call of CommitTx.
func (mr *MockRepositoryMockRecorder) CommitTx(ctx, tx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommitTx", reflect.TypeOf((*MockRepository)(nil).CommitTx), ctx, tx)
}

// GetReceiveCoinEvents mocks base method.
func (m *MockRepository) GetReceiveCoinEvents(ctx context.Context, tx pgx.Tx, login string) ([]model.ReceiveCoinEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReceiveCoinEvents", ctx, tx, login)
	ret0, _ := ret[0].([]model.ReceiveCoinEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReceiveCoinEvents indicates an expected call of GetReceiveCoinEvents.
func (mr *MockRepositoryMockRecorder) GetReceiveCoinEvents(ctx, tx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReceiveCoinEvents", reflect.TypeOf((*MockRepository)(nil).GetReceiveCoinEvents), ctx, tx, login)
}

// GetSendCoinEvents mocks base method.
func (m *MockRepository) GetSendCoinEvents(ctx context.Context, tx pgx.Tx, login string) ([]model.SendCoinEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSendCoinEvents", ctx, tx, login)
	ret0, _ := ret[0].([]model.SendCoinEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSendCoinEvents indicates an expected call of GetSendCoinEvents.
func (mr *MockRepositoryMockRecorder) GetSendCoinEvents(ctx, tx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSendCoinEvents", reflect.TypeOf((*MockRepository)(nil).GetSendCoinEvents), ctx, tx, login)
}

// GetUser mocks base method.
func (m *MockRepository) GetUser(ctx context.Context, tx pgx.Tx, login string) (*repository.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", ctx, tx, login)
	ret0, _ := ret[0].(*repository.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockRepositoryMockRecorder) GetUser(ctx, tx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockRepository)(nil).GetUser), ctx, tx, login)
}

// GetUserItems mocks base method.
func (m *MockRepository) GetUserItems(ctx context.Context, tx pgx.Tx, login string) ([]model.Item, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserItems", ctx, tx, login)
	ret0, _ := ret[0].([]model.Item)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserItems indicates an expected call of GetUserItems.
func (mr *MockRepositoryMockRecorder) GetUserItems(ctx, tx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserItems", reflect.TypeOf((*MockRepository)(nil).GetUserItems), ctx, tx, login)
}

// RollbackTx mocks base method.
func (m *MockRepository) RollbackTx(ctx context.Context, tx pgx.Tx) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RollbackTx", ctx, tx)
}

// RollbackTx indicates an expected call of RollbackTx.
func (mr *MockRepositoryMockRecorder) RollbackTx(ctx, tx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RollbackTx", reflect.TypeOf((*MockRepository)(nil).RollbackTx), ctx, tx)
}

// UpdateItemPurchaseCount mocks base method.
func (m *MockRepository) UpdateItemPurchaseCount(ctx context.Context, tx pgx.Tx, login, item string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateItemPurchaseCount", ctx, tx, login, item)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateItemPurchaseCount indicates an expected call of UpdateItemPurchaseCount.
func (mr *MockRepositoryMockRecorder) UpdateItemPurchaseCount(ctx, tx, login, item interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateItemPurchaseCount", reflect.TypeOf((*MockRepository)(nil).UpdateItemPurchaseCount), ctx, tx, login, item)
}

// UpdateUserBalance mocks base method.
func (m *MockRepository) UpdateUserBalance(ctx context.Context, tx pgx.Tx, login string, amount int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserBalance", ctx, tx, login, amount)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserBalance indicates an expected call of UpdateUserBalance.
func (mr *MockRepositoryMockRecorder) UpdateUserBalance(ctx, tx, login, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserBalance", reflect.TypeOf((*MockRepository)(nil).UpdateUserBalance), ctx, tx, login, amount)
}
