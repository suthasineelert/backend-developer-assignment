package models

import (
	"github.com/google/uuid"
	"time"
)

// AccountFlag represents the account_flags table
type AccountFlag struct {
	FlagID    int       `db:"flag_id" json:"flag_id"`
	AccountID uuid.UUID `db:"account_id" json:"account_id" validate:"required"`
	UserID    uuid.UUID `db:"user_id" json:"user_id" validate:"required"`
	FlagType  string    `db:"flag_type" json:"flag_type" validate:"required"`
	FlagValue string    `db:"flag_value" json:"flag_value" validate:"required"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}