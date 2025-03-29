package controllers

import (
	"backend-developer-assignment/app/controllers"
	"backend-developer-assignment/app/models"
	mocks "backend-developer-assignment/pkg/mocks/services"
	"backend-developer-assignment/pkg/utils"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// AuthControllerTestSuite defines the test suite for AuthController
type AuthControllerTestSuite struct {
	suite.Suite
	app         *fiber.App
	mockService *mocks.UserService
	userID      string
	tokens      *utils.Tokens
}

// SetupSuite runs once before all tests
func (s *AuthControllerTestSuite) SetupSuite() {
	// Set up environment for JWT
	os.Setenv("JWT_SECRET_KEY", "test-secret-key")
	os.Setenv("JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT", "15")
	os.Setenv("JWT_REFRESH_KEY", "test-refresh-key")
	os.Setenv("JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT", "720")

	// Generate a user ID for testing
	s.userID = uuid.New().String()
}

// SetupTest runs before each test case
func (s *AuthControllerTestSuite) SetupTest() {
	s.app = fiber.New()
	s.mockService = new(mocks.UserService)

	authController := controllers.NewAuthController(s.mockService)
	s.app.Post("/verify-pin", authController.VerifyPin)
	s.app.Post("/token/renew", authController.RenewTokens)

	// Generate tokens for testing
	var err error
	s.tokens, err = utils.GenerateNewTokens(s.userID)
	s.Require().NoError(err)
}

// Helper function to test response
func (s *AuthControllerTestSuite) testResponse(resp *http.Response, expectedCode int, expectedBody map[string]interface{}) {
	s.Equal(expectedCode, resp.StatusCode)

	var respBody map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respBody)

	if expectedCode == fiber.StatusOK {
		s.Contains(respBody, "tokens")
		tokens, ok := respBody["tokens"].(map[string]interface{})
		s.True(ok)
		s.Contains(tokens, "access")
		s.Contains(tokens, "refresh")
	} else {
		s.Equal(expectedBody, respBody)
	}
}

// TestVerifyPin_Success checks if PIN verification works
func (s *AuthControllerTestSuite) TestVerifyPin_Success() {
	// Hash a PIN for testing
	pin := "123456"
	hashPin, _ := utils.HashPIN(pin)
	userID := uuid.New().String()

	// Setup mock expectations
	s.mockService.On("GetUserByID", userID).Return(&models.User{
		UserID: userID,
		Name:   "Test User",
		PIN:    hashPin,
	}, nil)

	// Create request
	reqBody, _ := json.Marshal(map[string]string{
		"user_id": userID,
		"pin":     pin,
	})
	req := httptest.NewRequest(http.MethodPost, "/verify-pin", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Test the endpoint
	resp, err := s.app.Test(req)
	s.NoError(err)

	// Expected response
	expectedBody := map[string]interface{}{
		"tokens": map[string]interface{}{
			"access":  "",
			"refresh": "",
		},
	}
	s.testResponse(resp, fiber.StatusOK, expectedBody)

	// Verify expected method calls
	s.mockService.AssertExpectations(s.T())
}

// TestVerifyPin_UserNotFound checks if a missing user returns a 404
func (s *AuthControllerTestSuite) TestVerifyPin_UserNotFound() {
	// Generate a user ID
	userID := uuid.New().String()

	// Setup mock expectations
	s.mockService.On("GetUserByID", userID).Return(nil, errors.New("user not found"))

	// Create request
	reqBody, _ := json.Marshal(map[string]string{
		"user_id": userID,
		"pin":     "123456",
	})
	req := httptest.NewRequest(http.MethodPost, "/verify-pin", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Test the endpoint
	resp, err := s.app.Test(req)
	s.NoError(err)

	// Expected response
	expectedBody := map[string]interface{}{
		"message": "User does not exist",
	}
	s.testResponse(resp, fiber.StatusNotFound, expectedBody)

	// Verify expected method calls
	s.mockService.AssertExpectations(s.T())
}

// TestVerifyPin_InvalidPIN checks if incorrect PIN returns an unauthorized error
func (s *AuthControllerTestSuite) TestVerifyPin_InvalidPIN() {
	// Hash a PIN for testing
	pin := "123456"
	hashPin, _ := utils.HashPIN(pin)
	userID := uuid.New().String()

	// Setup mock expectations
	s.mockService.On("GetUserByID", userID).Return(&models.User{
		UserID: userID,
		Name:   "Test User",
		PIN:    hashPin,
	}, nil)

	// Create request with wrong PIN
	reqBody, _ := json.Marshal(map[string]string{
		"user_id": userID,
		"pin":     "wrong-pin",
	})
	req := httptest.NewRequest(http.MethodPost, "/verify-pin", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Test the endpoint
	resp, err := s.app.Test(req)
	s.NoError(err)

	// Expected response
	expectedBody := map[string]interface{}{
		"message": "Invalid PIN",
	}
	s.testResponse(resp, fiber.StatusUnauthorized, expectedBody)

	// Verify expected method calls
	s.mockService.AssertExpectations(s.T())
}

// TestRenewTokens_Success tests successful token renewal
func (s *AuthControllerTestSuite) TestRenewTokens_Success() {
	// Setup mock expectations
	s.mockService.On("GetUserByID", s.userID).Return(&models.User{
		UserID: s.userID,
		Name:   "Test User",
	}, nil)

	// Create request body
	reqBody, _ := json.Marshal(models.Renew{
		RefreshToken: s.tokens.Refresh,
	})

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/token/renew", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.tokens.Access)

	// Test the endpoint
	resp, err := s.app.Test(req)
	s.NoError(err)

	// Check status code
	s.Equal(fiber.StatusOK, resp.StatusCode)

	// Parse response body
	var respBody map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respBody)

	// Check response structure
	s.Contains(respBody, "tokens")
	tokens, ok := respBody["tokens"].(map[string]interface{})
	s.True(ok)
	s.Contains(tokens, "access")
	s.Contains(tokens, "refresh")

	// Verify expected method calls
	s.mockService.AssertExpectations(s.T())
}

// TestRenewTokens_InvalidRefreshToken tests renewal with invalid refresh token
func (s *AuthControllerTestSuite) TestRenewTokens_InvalidRefreshToken() {
	// Create request body with invalid refresh token
	reqBody, _ := json.Marshal(models.Renew{
		RefreshToken: "invalid-refresh-token",
	})

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/token/renew", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.tokens.Access)

	// Test the endpoint
	resp, err := s.app.Test(req)
	s.NoError(err)

	// Check status code - should be bad request for invalid token
	s.Equal(fiber.StatusBadRequest, resp.StatusCode)
}

// TestRenewTokens_MissingAccessToken tests renewal without access token
func (s *AuthControllerTestSuite) TestRenewTokens_MissingAccessToken() {
	// Create request body
	reqBody, _ := json.Marshal(models.Renew{
		RefreshToken: s.tokens.Refresh,
	})

	// Create request without Authorization header
	req := httptest.NewRequest(http.MethodPost, "/token/renew", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Test the endpoint
	resp, err := s.app.Test(req)
	s.NoError(err)

	// Check status code - should be unauthorized or internal server error
	s.True(resp.StatusCode == fiber.StatusUnauthorized || resp.StatusCode == fiber.StatusInternalServerError)
}

// TestRenewTokens_UserNotFound tests renewal when user is not found
func (s *AuthControllerTestSuite) TestRenewTokens_UserNotFound() {
	// Setup mock expectations - user not found
	s.mockService.On("GetUserByID", s.userID).Return(nil, assert.AnError)

	// Create request body
	reqBody, _ := json.Marshal(models.Renew{
		RefreshToken: s.tokens.Refresh,
	})

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/token/renew", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.tokens.Access)

	// Test the endpoint
	resp, err := s.app.Test(req)
	s.NoError(err)

	// Check status code
	s.Equal(fiber.StatusNotFound, resp.StatusCode)

	// Verify expected method calls
	s.mockService.AssertExpectations(s.T())
}

// Run the test suite
func TestAuthControllerTestSuite(t *testing.T) {
	suite.Run(t, new(AuthControllerTestSuite))
}
