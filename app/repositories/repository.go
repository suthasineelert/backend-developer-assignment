package repositories

import (
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	UserRepository          UserRepository
	UserGreetingsRepository UserGreetingRepository
}

func InitRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserRepository:          NewUserRepository(db),
		UserGreetingsRepository: NewUserGreetingsRepository(db),
	}
}
