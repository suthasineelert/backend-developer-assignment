package controllers_test

import (
	"backend-developer-assignment/app/controllers"
	"backend-developer-assignment/app/models"
	mocks "backend-developer-assignment/pkg/mocks/services"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// BannerControllerTestSuite defines the test suite
type BannerControllerTestSuite struct {
	suite.Suite
	app            *fiber.App
	bannerService  *mocks.BannerService
	controller     *controllers.BannerController
	testUserID     string
	testBannerID   string
	testBannerData *models.Banner
}

// SetupTest runs before each test
func (s *BannerControllerTestSuite) SetupTest() {
	s.app = fiber.New()
	s.bannerService = new(mocks.BannerService)
	s.controller = controllers.NewBannerController(s.bannerService)
	s.testUserID = "test-user-id"
	s.testBannerID = "test-banner-id"

	now := time.Now()
	s.testBannerData = &models.Banner{
		BaseModel:   &models.BaseModel{CreatedAt: now, UpdatedAt: now},
		BannerID:    s.testBannerID,
		UserID:      s.testUserID,
		Title:       "Test Banner",
		Description: "This is a test banner",
		Image:       "test-image.jpg",
	}

	// Setup routes
	s.app.Get("/banners", func(c *fiber.Ctx) error {
		c.Locals("userID", s.testUserID)
		return s.controller.ListBanners(c)
	})

	s.app.Get("/banners/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", s.testUserID)
		return s.controller.GetBanner(c)
	})
}

// TestListBanners tests the ListBanners controller method
func (s *BannerControllerTestSuite) TestListBanners() {
	// Test case: successful retrieval of banners
	s.bannerService.On("GetBannersByUserID", s.testUserID).Return([]*models.Banner{s.testBannerData}, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/banners", http.NoBody)
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

	var banners []*models.Banner
	err = json.NewDecoder(resp.Body).Decode(&banners)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), banners, 1)
	assert.Equal(s.T(), s.testBannerID, banners[0].BannerID)

	// Test case: error retrieving banners
	s.bannerService.On("GetBannersByUserID", s.testUserID).Return(nil, errors.New("database error")).Once()

	req = httptest.NewRequest(http.MethodGet, "/banners", http.NoBody)
	resp, err = s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusInternalServerError, resp.StatusCode)

	s.bannerService.AssertExpectations(s.T())
}

// TestGetBanner tests the GetBanner controller method
func (s *BannerControllerTestSuite) TestGetBanner() {
	// Test case: successful retrieval of banner
	s.bannerService.On("GetBannerByID", s.testBannerID).Return(s.testBannerData, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/banners/"+s.testBannerID, http.NoBody)
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

	var banner models.Banner
	err = json.NewDecoder(resp.Body).Decode(&banner)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), s.testBannerID, banner.BannerID)

	// Test case: banner not found
	s.bannerService.On("GetBannerByID", "nonexistent-id").Return(nil, nil).Once()

	req = httptest.NewRequest(http.MethodGet, "/banners/nonexistent-id", http.NoBody)
	resp, err = s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusNotFound, resp.StatusCode)

	// Test case: error retrieving banner
	s.bannerService.On("GetBannerByID", "error-id").Return(nil, errors.New("database error")).Once()

	req = httptest.NewRequest(http.MethodGet, "/banners/error-id", http.NoBody)
	resp, err = s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusInternalServerError, resp.StatusCode)

	s.bannerService.AssertExpectations(s.T())
}

// TestBannerControllerSuite runs the test suite
func TestBannerControllerSuite(t *testing.T) {
	suite.Run(t, new(BannerControllerTestSuite))
}
