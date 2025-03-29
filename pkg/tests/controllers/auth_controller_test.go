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
	"github.com/stretchr/testify/suite"
)

// AuthControllerTestSuite defines the test suite for AuthController
type AuthControllerTestSuite struct {
	suite.Suite
	app         *fiber.App
	mockService *mocks.UserService
}

// SetupTest runs before each test case
func (s *AuthControllerTestSuite) SetupTest() {
	s.app = fiber.New()
	s.mockService = new(mocks.UserService)

	authController := controllers.NewAuthController(s.mockService)
	s.app.Post("/verify-pin", authController.VerifyPin)

	// Set up environment variable for JWT
	os.Setenv("JWT_SECRET_KEY", "test-secret-key")
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
		UserID: uuid.MustParse(userID),
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
		UserID: uuid.MustParse(userID),
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

// Run the test suite
func TestAuthControllerTestSuite(t *testing.T) {
	suite.Run(t, new(AuthControllerTestSuite))
}
