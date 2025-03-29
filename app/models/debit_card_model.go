package models

import "time"

// DebitCard represents the debit_cards table
type DebitCard struct {
	CardID    string     `db:"card_id" json:"card_id" validate:"required"`
	UserID    string     `db:"user_id" json:"user_id" validate:"required"`
	Name      string     `db:"name" json:"name"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"` // for soft delete
}
