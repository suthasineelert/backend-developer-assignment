package models

import "backend-developer-assignment/pkg/base"

// Banner represents the banners table
type Banner struct {
	*base.BaseModel
	BannerID    string `db:"banner_id" json:"banner_id" validate:"required"`
	UserID      string `db:"user_id" json:"user_id" validate:"required"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	Image       string `db:"image" json:"image"`
}
