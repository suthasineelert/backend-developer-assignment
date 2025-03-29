package controllers

import (
	"backend-developer-assignment/app/controllers"
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/pkg/middleware"
	mocks "backend-developer-assignment/pkg/mocks/services"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// UserControllerTestSuite defines the test suite for UserController
type UserControllerTestSuite struct {
	suite.Suite
	app         *fiber.App
	mockService *mocks.UserService
	testToken   string
}

// SetupSuite runs once before all tests
func (s *UserControllerTestSuite) SetupSuite() {
	// Set up JWT environment for testing
	os.Setenv("JWT_SECRET_KEY", "test-secret-key")

	// Generate a test token
	s.testToken = s.generateTestToken()
}

// SetupTest runs before each test case
func (s *UserControllerTestSuite) SetupTest() {
	s.app = fiber.New()
	s.mockService = new(mocks.UserService)

	userController := controllers.NewUserController(s.mockService)
	s.app.Get("/users/greeting", middleware.JWTProtected(), userController.Greeting)
}

// Helper function to generate a test JWT token
func (s *UserControllerTestSuite) generateTestToken() string {
	// Create token claims
	claims := jwt.MapClaims{
		"id":  uuid.New().String(),
		"exp": time.Now().Add(time.Hour * 24).Unix(), // 24 hours
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	s.NoError(err)

	return tokenString
}

// Helper function to test response
func (s *UserControllerTestSuite) testResponse(resp *http.Response, expectedCode int, expectedMessage string) {
	s.Equal(expectedCode, resp.StatusCode)

	var respBody map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respBody)

	s.Equal(expectedMessage, respBody["message"])
}

// TestGreeting_Success checks if greeting retrieval works
func (s *UserControllerTestSuite) TestGreeting_Success() {
	greeting := &models.UserGreeting{
		UserID:   uuid.New(),
		Greeting: "Hello, welcome back!",
	}
	// Setup mock expectations
	s.mockService.On("GetUserGreetingByID", mock.Anything).Return(greeting, nil)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/users/greeting", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.testToken)

	// Test the endpoint
	resp, err := s.app.Test(req)
	s.NoError(err)

	// Check response
	s.testResponse(resp, fiber.StatusOK, greeting.Greeting)

	// Verify expected method calls
	s.mockService.AssertExpectations(s.T())
}

// TestGreeting_NotFound checks if missing greeting returns 404
func (s *UserControllerTestSuite) TestGreeting_NotFound() {
	// Setup mock expectations
	s.mockService.On("GetUserGreetingByID", mock.Anything).Return(nil, errors.New("user greeting not found"))

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/users/greeting", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.testToken)

	// Test the endpoint
	resp, err := s.app.Test(req)
	s.NoError(err)

	// Check response
	s.testResponse(resp, fiber.StatusNotFound, "User greeting not found")

	// Verify expected method calls
	s.mockService.AssertExpectations(s.T())
}

// TestGreeting_Unauthorized checks if missing token returns 401
func (s *UserControllerTestSuite) TestGreeting_Unauthorized() {
	// Create request without token
	req := httptest.NewRequest(http.MethodGet, "/users/greeting", nil)
	req.Header.Set("Content-Type", "application/json")
	// No Authorization header

	// Test the endpoint
	resp, err := s.app.Test(req)
	s.NoError(err)

	// Check response - assuming your app returns 401 for missing token
	s.Equal(fiber.StatusUnauthorized, resp.StatusCode)
}

// Run the test suite
func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}
