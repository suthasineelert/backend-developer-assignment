package models

import (
	"github.com/google/uuid"
)

// DebitCardDesign represents the debit_card_design table
type DebitCardDesign struct {
	CardID      uuid.UUID `db:"card_id" json:"card_id" validate:"required"`
	UserID      uuid.UUID `db:"user_id" json:"user_id" validate:"required"`
	Color       string    `db:"color" json:"color"`
	BorderColor string    `db:"border_color" json:"border_color"`
	DummyCol9   string    `db:"dummy_col_9" json:"dummy_col_9"`
}