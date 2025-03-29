package models

import "backend-developer-assignment/pkg/base"

// DebitCardStatus represents the debit_card_status table
type DebitCardStatus struct {
	*base.BaseModel
	CardID string `db:"card_id" json:"card_id" validate:"required"`
	UserID string `db:"user_id" json:"user_id" validate:"required"`
	Status string `db:"status" json:"status"`
}
