package repositories

import (
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	UserRepo UserRepository
	// Add more repositories here (e.g., TransactionRepo, AccountRepo)
}

func InitRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserRepo: NewUserRepository(db),
	}
}
