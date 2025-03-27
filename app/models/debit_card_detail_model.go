package models

import (
	"github.com/google/uuid"
)

// DebitCardDetail represents the debit_card_details table
type DebitCardDetail struct {
	CardID     uuid.UUID `db:"card_id" json:"card_id" validate:"required"`
	UserID     uuid.UUID `db:"user_id" json:"user_id" validate:"required"`
	Issuer     string    `db:"issuer" json:"issuer"`
	Number     string    `db:"number" json:"number"`
	DummyCol10 string    `db:"dummy_col_10" json:"dummy_col_10"`
}