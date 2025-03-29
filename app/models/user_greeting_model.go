package models

import "backend-developer-assignment/pkg/base"

// UserGreeting represents the user_greetings table
type UserGreeting struct {
	*base.BaseModel
	UserID   string `db:"user_id" json:"user_id" validate:"required"`
	Greeting string `db:"greeting" json:"greeting"`
}
