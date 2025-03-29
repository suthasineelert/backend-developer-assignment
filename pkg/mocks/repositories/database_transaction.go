package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockTx is a mock implementation of repositories.DBTransaction
type MockTx struct {
	mock.Mock
}

// Commit mocks the Commit method
func (m *MockTx) Commit() error {
	args := m.Called()
	return args.Error(0)
}

// Rollback mocks the Rollback method
func (m *MockTx) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

// Exec mocks the Exec method
func (m *MockTx) Exec(query string, args ...interface{}) (interface{}, error) {
	mockArgs := m.Called(append([]interface{}{query}, args...)...)
	return mockArgs.Get(0), mockArgs.Error(1)
}
