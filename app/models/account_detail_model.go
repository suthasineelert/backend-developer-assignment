package models

import "backend-developer-assignment/pkg/base"

// AccountDetail represents the account_details table
type AccountDetail struct {
	*base.BaseModel
	AccountID     string `db:"account_id" json:"account_id" validate:"required"`
	UserID        string `db:"user_id" json:"user_id" validate:"required"`
	Color         string `db:"color" json:"color"`
	IsMainAccount bool   `db:"is_main_account" json:"is_main_account"`
	Progress      int    `db:"progress" json:"progress"`
}
