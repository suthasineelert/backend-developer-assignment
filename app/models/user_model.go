package models

// User struct to describe User object.
type User struct {
	*BaseModel
	UserID string `db:"user_id" json:"user_id" validate:"required"`
	Name   string `db:"name" json:"name" validate:"required"`
	PIN    string `db:"pin" json:"-"`
}
