package models

import (
	"github.com/google/uuid"
)

// Transaction represents the transactions table
type Transaction struct {
	TransactionID uuid.UUID `db:"transaction_id" json:"transaction_id" validate:"required"`
	UserID        uuid.UUID `db:"user_id" json:"user_id" validate:"required"`
	Name          string    `db:"name" json:"name"`
	Image         string    `db:"image" json:"image"`
	IsBank        bool      `db:"isBank" json:"is_bank"`
	DummyCol6     string    `db:"dummy_col_6" json:"dummy_col_6"`
}