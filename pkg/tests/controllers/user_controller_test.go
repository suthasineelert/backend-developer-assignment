package controllers

import (
	"backend-developer-assignment/app/controllers"
	"backend-developer-assignment/app/models"
	mocks "backend-developer-assignment/pkg/mocks/services"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

// UserControllerTestSuite defines the test suite for UserController
type UserControllerTestSuite struct {
	suite.Suite
	app         *fiber.App
	mockService *mocks.UserService
}

// SetupTest runs before each test case
func (s *UserControllerTestSuite) SetupTest() {
	s.app = fiber.New()
	s.mockService = new(mocks.UserService)

	userController := controllers.NewUserController(s.mockService)
	s.app.Get("/users/greeting", userController.Greeting)
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
	// Create test data
	userID := uuid.New().String()
	greeting := &models.UserGreeting{
		UserID:   uuid.MustParse(userID),
		Greeting: "Hello, welcome back!",
	}

	// Setup mock expectations
	s.mockService.On("GetUserGreetingByID", userID).Return(greeting, nil)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/users/"+userID+"/greeting", http.NoBody)
	req.Header.Set("Content-Type", "application/json")

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
	// Create test data
	userID := uuid.New().String()

	// Setup mock expectations
	s.mockService.On("GetUserGreetingByID", userID).Return(nil, errors.New("user greeting not found"))

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/users/"+userID+"/greeting", http.NoBody)
	req.Header.Set("Content-Type", "application/json")

	// Test the endpoint
	resp, err := s.app.Test(req)
	s.NoError(err)

	// Check response
	s.testResponse(resp, fiber.StatusNotFound, "User greeting not found")

	// Verify expected method calls
	s.mockService.AssertExpectations(s.T())
}

// Run the test suite
func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}
