package repositories

import (
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	UserRepository          UserRepository
	UserGreetingsRepository UserGreetingRepository
	TransactionRepository   TransactionRepository
	DebitCardRepository     DebitCardRepository
	AccountRepository       AccountRepository
}

func InitRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserRepository:          NewUserRepository(db),
		UserGreetingsRepository: NewUserGreetingsRepository(db),
		TransactionRepository:   NewTransactionRepository(db),
		DebitCardRepository:     NewDebitCardRepository(db),
		AccountRepository:       NewAccountRepository(db),
	}
}
