package models

// Account represents the accounts table
type Account struct {
	AccountID     string `db:"account_id" json:"account_id" validate:"required"`
	UserID        string `db:"user_id" json:"user_id" validate:"required"`
	Type          string `db:"type" json:"type"`
	Currency      string `db:"currency" json:"currency"`
	AccountNumber string `db:"account_number" json:"account_number"`
	Issuer        string `db:"issuer" json:"issuer"`
	DummyCol3     string `db:"dummy_col_3" json:"dummy_col_3"`
}
