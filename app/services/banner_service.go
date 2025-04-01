package services

import (
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/app/repositories"

	"go.uber.org/zap"
)

// BannerService defines the interface for banner operations
type BannerService interface {
	GetBannerByID(bannerID string) (*models.Banner, error)
	GetBannersByUserID(userID string) ([]*models.Banner, error)
}

// BannerServiceImpl implements BannerService
type BannerServiceImpl struct {
	bannerRepository repositories.BannerRepository
}

// NewBannerService creates a new banner service
func NewBannerService(bannerRepository repositories.BannerRepository) BannerService {
	return &BannerServiceImpl{
		bannerRepository: bannerRepository,
	}
}

// GetBannerByID retrieves a banner by its ID
func (s *BannerServiceImpl) GetBannerByID(bannerID string) (*models.Banner, error) {
	banner, err := s.bannerRepository.GetBannerByID(bannerID)
	if err != nil {
		logger.Error("Failed to get banner by ID", zap.String("banner_id", bannerID), zap.Error(err))
		return nil, err
	}
	
	if banner == nil {
		logger.Info("Banner not found", zap.String("banner_id", bannerID))
		return nil, nil
	}
	
	return banner, nil
}

// GetBannersByUserID retrieves all banners for a specific user
func (s *BannerServiceImpl) GetBannersByUserID(userID string) ([]*models.Banner, error) {
	banners, err := s.bannerRepository.GetBannersByUserID(userID)
	if err != nil {
		logger.Error("Failed to get banners by user ID", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}
	
	return banners, nil
}