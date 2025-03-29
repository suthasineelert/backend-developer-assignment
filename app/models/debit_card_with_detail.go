package models

import "time"

// DebitCardWithDetails represents a joined view of all debit card related tables
type DebitCardWithDetails struct {
	// DebitCard fields
	CardID    string    `db:"card_id" json:"card_id"`
	UserID    string    `db:"user_id" json:"user_id"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// DebitCardDetail fields
	Issuer string `db:"issuer" json:"issuer"`
	Number string `db:"number" json:"number"`

	// DebitCardDesign fields
	Color       string `db:"color" json:"color"`
	BorderColor string `db:"border_color" json:"border_color"`

	// DebitCardStatus fields
	Status string `db:"status" json:"status"`
}
