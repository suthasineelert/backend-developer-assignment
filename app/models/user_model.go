package models

import (
	"github.com/google/uuid"
)

// User struct to describe User object.
type User struct {
	UserID uuid.UUID `db:"user_id" json:"user_id" validate:"required"`
	Name   string    `db:"name" json:"name" validate:"required"`
	PIN    string    `db:"pin" json:"-" validate:"required"`
}
