package models

// UserGreeting represents the user_greetings table
type UserGreeting struct {
	UserID    string `db:"user_id" json:"user_id" validate:"required"`
	Greeting  string `db:"greeting" json:"greeting"`
	DummyCol2 string `db:"dummy_col_2" json:"dummy_col_2"`
}
