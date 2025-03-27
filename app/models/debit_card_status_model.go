package models

import (
	"github.com/google/uuid"
)

// DebitCardStatus represents the debit_card_status table
type DebitCardStatus struct {
	CardID    uuid.UUID `db:"card_id" json:"card_id" validate:"required"`
	UserID    uuid.UUID `db:"user_id" json:"user_id" validate:"required"`
	Status    string    `db:"status" json:"status"`
	DummyCol8 string    `db:"dummy_col_8" json:"dummy_col_8"`
}