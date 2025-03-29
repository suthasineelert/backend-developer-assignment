package models

// AccountBalance represents the account_balances table
type AccountBalance struct {
	AccountID string  `db:"account_id" json:"account_id" validate:"required"`
	UserID    string  `db:"user_id" json:"user_id" validate:"required"`
	Amount    float64 `db:"amount" json:"amount"`
	DummyCol4 string  `db:"dummy_col_4" json:"dummy_col_4"`
}
