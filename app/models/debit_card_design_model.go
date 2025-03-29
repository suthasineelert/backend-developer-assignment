package models

// DebitCardDesign represents the debit_card_design table
type DebitCardDesign struct {
	CardID      string `db:"card_id" json:"card_id" validate:"required"`
	UserID      string `db:"user_id" json:"user_id" validate:"required"`
	Color       string `db:"color" json:"color"`
	BorderColor string `db:"border_color" json:"border_color"`
}
