package repositories

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type DB interface {
	Query(string, ...any) (*sql.Rows, error)
	Get(dest any, query string, args ...any) error
	Exec(query string, args ...any) (sql.Result, error)
	Select(dest any, query string, args ...any) error
}

type Adapters struct {
	AccountRepository     AccountRepository
	TransactionRepository TransactionRepository
}

type TxProvider interface {
	Transact(txFunc func(adapters Adapters) error) error
}

type TransactionProvider struct {
	db *sqlx.DB
}

func NewTransactionProvider(db *sqlx.DB) *TransactionProvider {
	return &TransactionProvider{
		db: db,
	}
}

func (p *TransactionProvider) Transact(txFunc func(adapters Adapters) error) error {
	return runInTx(p.db, func(tx *sqlx.Tx) error {
		adapters := Adapters{
			AccountRepository:     NewAccountRepository(tx),
			TransactionRepository: NewTransactionRepository(tx),
		}

		return txFunc(adapters)
	})
}

func runInTx(db interface{}, fn func(tx *sqlx.Tx) error) error {
	var tx *sqlx.Tx
	var err error
	var isTx bool = false

	// Check if db is *sql.DB or *sql.Tx
	//nolint:all // This type switch is intentional for handling different DB types
	switch d := db.(type) {
	case *sqlx.DB:
		// Begin a new transaction
		tx, err = d.Beginx()
		if err != nil {
			return err
		}
	case *sqlx.Tx:
		// Use the existing transaction
		tx = d
		isTx = true
	default:
		return errors.New("invalid db type, only *sql.DB and *sql.Tx are supported")
	}

	err = fn(tx)
	if err == nil {
		if isTx {
			return nil // no need to commit if it's a transaction, since it might be nested in an outer tx
		}
		return tx.Commit()
	}

	rollbackErr := tx.Rollback()
	if rollbackErr != nil {
		return errors.Join(err, rollbackErr)
	}

	return err
}
