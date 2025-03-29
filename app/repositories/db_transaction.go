package repositories

import "github.com/jmoiron/sqlx"

// DBTransaction defines the interface for database transactions
type DBTransaction interface {
	Commit() error
	Rollback() error
	Exec(query string, args ...interface{}) (interface{}, error)
}

// sqlxTransaction is a wrapper around sqlx.Tx that implements DBTransaction
type sqlxTransaction struct {
	tx *sqlx.Tx
}

func (t *sqlxTransaction) Commit() error {
	return t.tx.Commit()
}

func (t *sqlxTransaction) Rollback() error {
	return t.tx.Rollback()
}

func (t *sqlxTransaction) Exec(query string, args ...interface{}) (interface{}, error) {
	return t.tx.Exec(query, args...)
}
