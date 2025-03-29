package models

import (
	"backend-developer-assignment/pkg/base"
)

// DebitCard represents the debit_cards table
type DebitCard struct {
	*base.BaseModel
	CardID string `db:"card_id" json:"card_id" validate:"required"`
	UserID string `db:"user_id" json:"user_id" validate:"required"`
	Name   string `db:"name" json:"name"`
}
