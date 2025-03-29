package sqlx

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/mock"
)

// MockTx is a mock implementation of sqlx.Tx
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
func (m *MockTx) Exec(query string, args ...interface{}) (sql.Result, error) {
	mockArgs := m.Called(append([]interface{}{query}, args...)...)
	result := mockArgs.Get(0)
	if result == nil {
		return nil, mockArgs.Error(1)
	}
	return result.(sql.Result), mockArgs.Error(1)
}

// Query mocks the Query method
func (m *MockTx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	mockArgs := m.Called(append([]interface{}{query}, args...)...)
	rows := mockArgs.Get(0)
	if rows == nil {
		return nil, mockArgs.Error(1)
	}
	return rows.(*sql.Rows), mockArgs.Error(1)
}

// Queryx mocks the Queryx method
func (m *MockTx) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	mockArgs := m.Called(append([]interface{}{query}, args...)...)
	rows := mockArgs.Get(0)
	if rows == nil {
		return nil, mockArgs.Error(1)
	}
	return rows.(*sqlx.Rows), mockArgs.Error(1)
}

// QueryRowx mocks the QueryRowx method
func (m *MockTx) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	mockArgs := m.Called(append([]interface{}{query}, args...)...)
	return mockArgs.Get(0).(*sqlx.Row)
}

// Select mocks the Select method
func (m *MockTx) Select(dest interface{}, query string, args ...interface{}) error {
	mockArgs := m.Called(append([]interface{}{dest, query}, args...)...)
	return mockArgs.Error(0)
}

// Get mocks the Get method
func (m *MockTx) Get(dest interface{}, query string, args ...interface{}) error {
	mockArgs := m.Called(append([]interface{}{dest, query}, args...)...)
	return mockArgs.Error(0)
}

// Prepare mocks the Prepare method
func (m *MockTx) Prepare(query string) (*sql.Stmt, error) {
	args := m.Called(query)
	stmt := args.Get(0)
	if stmt == nil {
		return nil, args.Error(1)
	}
	return stmt.(*sql.Stmt), args.Error(1)
}

// Stmt mocks the Stmt method
func (m *MockTx) Stmt(stmt *sql.Stmt) *sql.Stmt {
	args := m.Called(stmt)
	return args.Get(0).(*sql.Stmt)
}

// DriverName mocks the DriverName method
func (m *MockTx) DriverName() string {
	args := m.Called()
	return args.String(0)
}

// Rebind mocks the Rebind method
func (m *MockTx) Rebind(query string) string {
	args := m.Called(query)
	return args.String(0)
}

// BindNamed mocks the BindNamed method
func (m *MockTx) BindNamed(query string, arg interface{}) (string, []interface{}, error) {
	args := m.Called(query, arg)
	return args.String(0), args.Get(1).([]interface{}), args.Error(2)
}

// MustExec mocks the MustExec method
func (m *MockTx) MustExec(query string, args ...interface{}) sql.Result {
	mockArgs := m.Called(append([]interface{}{query}, args...)...)
	return mockArgs.Get(0).(sql.Result)
}

// Preparex mocks the Preparex method
func (m *MockTx) Preparex(query string) (*sqlx.Stmt, error) {
	args := m.Called(query)
	stmt := args.Get(0)
	if stmt == nil {
		return nil, args.Error(1)
	}
	return stmt.(*sqlx.Stmt), args.Error(1)
}

// Stmtx mocks the Stmtx method
func (m *MockTx) Stmtx(stmt interface{}) *sqlx.Stmt {
	args := m.Called(stmt)
	return args.Get(0).(*sqlx.Stmt)
}

// NamedExec mocks the NamedExec method
func (m *MockTx) NamedExec(query string, arg interface{}) (sql.Result, error) {
	args := m.Called(query, arg)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(sql.Result), args.Error(1)
}

// NamedQuery mocks the NamedQuery method
func (m *MockTx) NamedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	args := m.Called(query, arg)
	rows := args.Get(0)
	if rows == nil {
		return nil, args.Error(1)
	}
	return rows.(*sqlx.Rows), args.Error(1)
}

// PrepareNamed mocks the PrepareNamed method
func (m *MockTx) PrepareNamed(query string) (*sqlx.NamedStmt, error) {
	args := m.Called(query)
	stmt := args.Get(0)
	if stmt == nil {
		return nil, args.Error(1)
	}
	return stmt.(*sqlx.NamedStmt), args.Error(1)
}

// NamedStmt mocks the NamedStmt method
func (m *MockTx) NamedStmt(stmt *sqlx.NamedStmt) *sqlx.NamedStmt {
	args := m.Called(stmt)
	return args.Get(0).(*sqlx.NamedStmt)
}