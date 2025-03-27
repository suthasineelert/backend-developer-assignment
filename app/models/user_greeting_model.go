package models

import (
	"github.com/google/uuid"
)

// UserGreeting represents the user_greetings table
type UserGreeting struct {
	UserID    uuid.UUID `db:"user_id" json:"user_id" validate:"required"`
	Greeting  string    `db:"greeting" json:"greeting"`
	DummyCol2 string    `db:"dummy_col_2" json:"dummy_col_2"`
}