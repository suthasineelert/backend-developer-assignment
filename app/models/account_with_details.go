package models

import "time"

// AccountWithDetailsFlat represents a flattened view of an account with all its related information
type AccountWithDetails struct {
	// Account fields
	AccountID     string     `json:"account_id" db:"account_id"`
	UserID        string     `json:"user_id" db:"user_id"`
	Type          string     `json:"type" db:"type"` // saving-account, credit-loan, goal-driven-saving
	Currency      string     `json:"currency" db:"currency"`
	AccountNumber string     `json:"account_number" db:"account_number"`
	Issuer        string     `json:"issuer" db:"issuer"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`

	// AccountDetail fields
	Color         string `json:"color" db:"color"`
	IsMainAccount bool   `json:"is_main_account" db:"is_main_account"`
	Progress      int    `json:"progress" db:"progress"`

	// AccountBalance fields
	Amount float64 `json:"amount" db:"amount"`

	// AccountFlags
	Flags []*AccountFlag `json:"flags" db:"-"` // Using db:"-" to indicate this field is not directly mapped from DB
}
