package models

import (
	"github.com/google/uuid"
)

// Banner represents the banners table
type Banner struct {
	BannerID    uuid.UUID `db:"banner_id" json:"banner_id" validate:"required"`
	UserID      uuid.UUID `db:"user_id" json:"user_id" validate:"required"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Image       string    `db:"image" json:"image"`
	DummyCol11  string    `db:"dummy_col_11" json:"dummy_col_11"`
}