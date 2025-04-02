package repositories

import (
	"backend-developer-assignment/app/models"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// BannerRepository defines the interface for banner operations
type BannerRepository interface {
	GetBannerByID(bannerID string) (*models.Banner, error)
	GetBannersByUserID(userID string) ([]*models.Banner, error)
}

// BannerRepositoryImpl implements BannerRepository
type BannerRepositoryImpl struct {
	db *sqlx.DB
}

// NewBannerRepository creates a new banner repository
func NewBannerRepository(db *sqlx.DB) BannerRepository {
	return &BannerRepositoryImpl{
		db: db,
	}
}

// GetBannerByID retrieves a banner by its ID
func (r *BannerRepositoryImpl) GetBannerByID(bannerID string) (*models.Banner, error) {
	banner := &models.Banner{}
	query := `SELECT banner_id, user_id, title, description, image, created_at, updated_at FROM banners WHERE banner_id = ?`

	err := r.db.Get(banner, query, bannerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil, nil when no banner is found
		}
		return nil, err
	}

	return banner, nil
}

// GetBannersByUserID retrieves all banners for a specific user
func (r *BannerRepositoryImpl) GetBannersByUserID(userID string) ([]*models.Banner, error) {
	banners := []*models.Banner{}
	query := `SELECT banner_id, user_id, title, description, image, created_at, updated_at FROM banners WHERE user_id = ? ORDER BY created_at DESC`

	err := r.db.Select(&banners, query, userID)
	if err != nil {
		return nil, err
	}

	return banners, nil
}
