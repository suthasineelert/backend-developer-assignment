package models

// AccountDetail represents the account_details table
type AccountDetail struct {
	AccountID     string `db:"account_id" json:"account_id" validate:"required"`
	UserID        string `db:"user_id" json:"user_id" validate:"required"`
	Color         string `db:"color" json:"color"`
	IsMainAccount bool   `db:"is_main_account" json:"is_main_account"`
	Progress      int    `db:"progress" json:"progress"`
	DummyCol5     string `db:"dummy_col_5" json:"dummy_col_5"`
}
