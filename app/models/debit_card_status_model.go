package models

type CardStatus string

const (
	CardStatusActive     CardStatus = "active"
	CardStatusInactive   CardStatus = "inactive"
	CardStatusInprogress CardStatus = "in-progress"
	CardStatusBlocked    CardStatus = "blocked"
)

// DebitCardStatus represents the debit_card_status table
type DebitCardStatus struct {
	*BaseModel
	CardID string `db:"card_id" json:"card_id" validate:"required"`
	UserID string `db:"user_id" json:"user_id" validate:"required"`
	Status string `db:"status" json:"status"` // active, inactive, in-progress, blocked
}
