package services_test

import (
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/app/services"
	mocks "backend-developer-assignment/pkg/mocks/repositories"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// BannerServiceTestSuite is a test suite for BannerService
type BannerServiceTestSuite struct {
	suite.Suite
	bannerRepository *mocks.BannerRepository
	service          services.BannerService
}

// SetupTest sets up the test suite
func (s *BannerServiceTestSuite) SetupTest() {
	s.bannerRepository = new(mocks.BannerRepository)
	s.service = services.NewBannerService(s.bannerRepository)
}

// TestGetBannerByID tests the GetBannerByID function
func (s *BannerServiceTestSuite) TestGetBannerByID() {
	bannerID := "banner-123"

	testCases := []struct {
		name           string
		mockBanner     *models.Banner
		mockError      error
		expectedBanner *models.Banner
		expectedError  error
	}{
		{
			name: "Success - Banner Found",
			mockBanner: &models.Banner{
				BaseModel:   &models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()},
				BannerID:    bannerID,
				UserID:      "user-123",
				Title:       "Test Banner",
				Description: "This is a test banner",
				Image:       "test-image.jpg",
			},
			mockError: nil,
			expectedBanner: &models.Banner{
				BaseModel:   &models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()},
				BannerID:    bannerID,
				UserID:      "user-123",
				Title:       "Test Banner",
				Description: "This is a test banner",
				Image:       "test-image.jpg",
			},
			expectedError: nil,
		},
		{
			name:           "Success - Banner Not Found",
			mockBanner:     nil,
			mockError:      nil,
			expectedBanner: nil,
			expectedError:  nil,
		},
		{
			name:           "Failure - Database Error",
			mockBanner:     nil,
			mockError:      errors.New("database connection failed"),
			expectedBanner: nil,
			expectedError:  errors.New("database connection failed"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Mock the repository method
			s.bannerRepository.On("GetBannerByID", bannerID).Return(tc.mockBanner, tc.mockError).Once()

			// Call the service method
			banner, err := s.service.GetBannerByID(bannerID)

			// Assert results
			if tc.expectedError != nil {
				assert.Error(s.T(), err)
				assert.Equal(s.T(), tc.expectedError.Error(), err.Error())
				assert.Nil(s.T(), banner)
			} else {
				assert.NoError(s.T(), err)
				if tc.expectedBanner == nil {
					assert.Nil(s.T(), banner)
				} else {
					assert.Equal(s.T(), tc.expectedBanner.BannerID, banner.BannerID)
					assert.Equal(s.T(), tc.expectedBanner.UserID, banner.UserID)
					assert.Equal(s.T(), tc.expectedBanner.Title, banner.Title)
					assert.Equal(s.T(), tc.expectedBanner.Description, banner.Description)
					assert.Equal(s.T(), tc.expectedBanner.Image, banner.Image)
				}
			}

			// Verify expected method calls
			s.bannerRepository.AssertExpectations(s.T())
		})
	}
}

// TestGetBannersByUserID tests the GetBannersByUserID function
func (s *BannerServiceTestSuite) TestGetBannersByUserID() {
	userID := "user-123"

	testCases := []struct {
		name            string
		mockBanners     []*models.Banner
		mockError       error
		expectedBanners []*models.Banner
		expectedError   error
	}{
		{
			name: "Success - Banners Found",
			mockBanners: []*models.Banner{
				{
					BaseModel:   &models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()},
					BannerID:    "banner-123",
					UserID:      userID,
					Title:       "Test Banner 1",
					Description: "This is test banner 1",
					Image:       "test-image-1.jpg",
				},
				{
					BaseModel:   &models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()},
					BannerID:    "banner-456",
					UserID:      userID,
					Title:       "Test Banner 2",
					Description: "This is test banner 2",
					Image:       "test-image-2.jpg",
				},
			},
			mockError: nil,
			expectedBanners: []*models.Banner{
				{
					BaseModel:   &models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()},
					BannerID:    "banner-123",
					UserID:      userID,
					Title:       "Test Banner 1",
					Description: "This is test banner 1",
					Image:       "test-image-1.jpg",
				},
				{
					BaseModel:   &models.BaseModel{CreatedAt: time.Now(), UpdatedAt: time.Now()},
					BannerID:    "banner-456",
					UserID:      userID,
					Title:       "Test Banner 2",
					Description: "This is test banner 2",
					Image:       "test-image-2.jpg",
				},
			},
			expectedError: nil,
		},
		{
			name:            "Success - No Banners Found",
			mockBanners:     []*models.Banner{},
			mockError:       nil,
			expectedBanners: []*models.Banner{},
			expectedError:   nil,
		},
		{
			name:            "Failure - Database Error",
			mockBanners:     nil,
			mockError:       errors.New("database connection failed"),
			expectedBanners: nil,
			expectedError:   errors.New("database connection failed"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Mock the repository method
			s.bannerRepository.On("GetBannersByUserID", userID).Return(tc.mockBanners, tc.mockError).Once()

			// Call the service method
			banners, err := s.service.GetBannersByUserID(userID)

			// Assert results
			if tc.expectedError != nil {
				assert.Error(s.T(), err)
				assert.Equal(s.T(), tc.expectedError.Error(), err.Error())
				assert.Nil(s.T(), banners)
			} else {
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), len(tc.expectedBanners), len(banners))

				// If we have banners, check their properties
				if len(tc.expectedBanners) > 0 {
					for i, expectedBanner := range tc.expectedBanners {
						assert.Equal(s.T(), expectedBanner.BannerID, banners[i].BannerID)
						assert.Equal(s.T(), expectedBanner.UserID, banners[i].UserID)
						assert.Equal(s.T(), expectedBanner.Title, banners[i].Title)
						assert.Equal(s.T(), expectedBanner.Description, banners[i].Description)
						assert.Equal(s.T(), expectedBanner.Image, banners[i].Image)
					}
				}
			}

			// Verify expected method calls
			s.bannerRepository.AssertExpectations(s.T())
		})
	}
}

// TestMain runs the test suite
func TestBannerService(t *testing.T) {
	suite.Run(t, new(BannerServiceTestSuite))
}
