package models

// DebitCardDetail represents the debit_card_details table
type DebitCardDetail struct {
	CardID string `db:"card_id" json:"card_id" validate:"required"`
	UserID string `db:"user_id" json:"user_id" validate:"required"`
	Issuer string `db:"issuer" json:"issuer"`
	Number string `db:"number" json:"number"`
}
