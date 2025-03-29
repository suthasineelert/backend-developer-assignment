package models

import "backend-developer-assignment/pkg/base"

// AccountBalance represents the account_balances table
type AccountBalance struct {
	*base.BaseModel
	AccountID string  `db:"account_id" json:"account_id" validate:"required"`
	UserID    string  `db:"user_id" json:"user_id" validate:"required"`
	Amount    float64 `db:"amount" json:"amount"`
}
