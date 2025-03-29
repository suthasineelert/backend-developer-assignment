package models

type AccountType string

const (
	SavingAccount AccountType = "saving-account"
	CreditLoan    AccountType = "credit-loan"
	GoalDriven    AccountType = "goal-driven-saving"
)

// Account represents the accounts table
type Account struct {
	*BaseModel
	AccountID     string `db:"account_id" json:"account_id" validate:"required"`
	UserID        string `db:"user_id" json:"user_id" validate:"required"`
	Type          string `db:"type" json:"type"` // saving-account, credit-loan, goal-driven-saving
	Currency      string `db:"currency" json:"currency"`
	AccountNumber string `db:"account_number" json:"account_number"`
	Issuer        string `db:"issuer" json:"issuer"`
}
