package models

// UserGreeting represents the user_greetings table
type UserGreeting struct {
	*BaseModel
	UserID   string `db:"user_id" json:"user_id" validate:"required"`
	Greeting string `db:"greeting" json:"greeting"`
}
