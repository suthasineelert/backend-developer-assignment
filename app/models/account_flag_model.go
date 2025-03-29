package models

import "backend-developer-assignment/pkg/base"

// AccountFlag represents the account_flags table
type AccountFlag struct {
	*base.BaseModel
	FlagID    int    `db:"flag_id" json:"flag_id"`
	AccountID string `db:"account_id" json:"account_id" validate:"required"`
	UserID    string `db:"user_id" json:"user_id" validate:"required"`
	FlagType  string `db:"flag_type" json:"flag_type" validate:"required"`
	FlagValue string `db:"flag_value" json:"flag_value" validate:"required"`
}
