package models

import "time"

// AccountWithDetailsFlat represents a flattened view of an account with all its related information
type AccountWithDetails struct {
	// Account fields
	AccountID     string    `json:"account_id"`
	UserID        string    `json:"user_id"`
	Type          string    `json:"type"` // saving-account, credit-loan, goal-driven-saving
	Currency      string    `json:"currency"`
	AccountNumber string    `json:"account_number"`
	Issuer        string    `json:"issuer"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// AccountDetail fields
	Color         string `json:"color"`
	IsMainAccount bool   `json:"is_main_account"`
	Progress      int    `json:"progress"`

	// AccountBalance fields
	Amount float64 `json:"amount"`

	// AccountFlags
	Flags []*AccountFlag `json:"flags"`
}
