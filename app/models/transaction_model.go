package models

type TransactionType string

const (
	Deposit    TransactionType = "deposit"
	Withdrawal TransactionType = "withdrawal"
	Transfer   TransactionType = "transfer"
)

// Transaction represents the transactions table
type Transaction struct {
	*BaseModel
	TransactionID   string  `db:"transaction_id" json:"transaction_id" validate:"required"`
	AccountID       string  `db:"account_id" json:"account_id" validate:"required"`
	UserID          string  `db:"user_id" json:"user_id" validate:"required"`
	Name            string  `db:"name" json:"name"`
	Image           string  `db:"image" json:"image"`
	IsBank          bool    `db:"isBank" json:"is_bank"`
	Amount          float64 `db:"amount" json:"amount" validate:"required"`
	TransactionType string  `db:"transaction_type" json:"transaction_type" validate:"required"` // deposit, withdrawal, transfer
}
