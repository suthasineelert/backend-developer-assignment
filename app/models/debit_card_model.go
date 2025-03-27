package models

import (
	"github.com/google/uuid"
)

// DebitCard represents the debit_cards table
type DebitCard struct {
	CardID    uuid.UUID `db:"card_id" json:"card_id" validate:"required"`
	UserID    uuid.UUID `db:"user_id" json:"user_id" validate:"required"`
	Name      string    `db:"name" json:"name"`
	DummyCol7 string    `db:"dummy_col_7" json:"dummy_col_7"`
}