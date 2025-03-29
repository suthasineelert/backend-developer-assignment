package models

// DebitCard represents the debit_cards table
type DebitCard struct {
	*BaseModel
	CardID string `db:"card_id" json:"card_id" validate:"required"`
	UserID string `db:"user_id" json:"user_id" validate:"required"`
	Name   string `db:"name" json:"name"`
}
