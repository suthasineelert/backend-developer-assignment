package controllers

import (
	"backend-developer-assignment/app/controllers"
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/pkg/middleware"
	mocks "backend-developer-assignment/pkg/mocks/services"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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
	route := s.app.Group("/users", middleware.AuthProtected()...)
	route.Get("/greeting", userController.GetUserGreeting)
	route.Put("/greeting", userController.UpdateUserGreeting)
	route.Get("/profile", userController.GetUser)
	route.Patch("/profile", userController.UpdateUser)
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
		UserID:   uuid.New().String(),
		Greeting: "Hello, welcome back!",
	}
	// Setup mock expectations
	s.mockService.On("GetUserGreetingByID", mock.Anything).Return(greeting, nil)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/users/greeting", http.NoBody)
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
	req := httptest.NewRequest(http.MethodGet, "/users/greeting", http.NoBody)
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
	req := httptest.NewRequest(http.MethodGet, "/users/greeting", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	// No Authorization header

	// Test the endpoint
	resp, err := s.app.Test(req)
	s.NoError(err)

	// Check response - assuming your app returns 401 for missing token
	s.Equal(fiber.StatusUnauthorized, resp.StatusCode)
}

// TestUpdateUserGreeting_Success checks if greeting update works
func (s *UserControllerTestSuite) TestUpdateUserGreeting_Success() {
	// Setup mock expectations
	s.mockService.On("UpdateUserGreeting", mock.Anything).Return(nil)

	// Create request body
	requestBody := map[string]string{
		"message": "Updated greeting message",
	}
	requestJSON, _ := json.Marshal(requestBody)

	// Create request
	req := httptest.NewRequest(http.MethodPut, "/users/greeting",
		bytes.NewReader(requestJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.testToken)

	// Test the endpoint
	resp, err := s.app.Test(req)
	s.NoError(err)

	// Check response
	s.Equal(fiber.StatusOK, resp.StatusCode)

	var respBody map[string]string
	json.NewDecoder(resp.Body).Decode(&respBody)
	s.Equal("Updated greeting message", respBody["message"])

	// Verify expected method calls
	s.mockService.AssertExpectations(s.T())
}

// TestUpdateUserGreeting_BadRequest checks if invalid request body returns 400
func (s *UserControllerTestSuite) TestUpdateUserGreeting_BadRequest() {
	// Create request with invalid JSON
	req := httptest.NewRequest(http.MethodPut, "/users/greeting",
		strings.NewReader("{invalid json}"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.testToken)

	// Test the endpoint
	resp, err := s.app.Test(req)
	s.NoError(err)

	// Check response
	s.Equal(fiber.StatusBadRequest, resp.StatusCode)

	var respBody map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respBody)
	s.Contains(respBody["message"].(string), "Invalid input format")
}

// TestUpdateUserGreeting_UpdateFailed checks if service error returns 404
func (s *UserControllerTestSuite) TestUpdateUserGreeting_UpdateFailed() {
	// Setup mock expectations
	s.mockService.On("UpdateUserGreeting", mock.Anything).Return(errors.New("update failed"))

	// Create request body
	requestBody := map[string]string{
		"message": "Updated greeting message",
	}
	requestJSON, _ := json.Marshal(requestBody)

	// Create request
	req := httptest.NewRequest(http.MethodPut, "/users/greeting",
		bytes.NewReader(requestJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.testToken)

	// Test the endpoint
	resp, err := s.app.Test(req)
	s.NoError(err)

	// Check response
	s.Equal(fiber.StatusNotFound, resp.StatusCode)

	var respBody map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respBody)
	s.Equal("Fail to update user greeting", respBody["message"])

	// Verify expected method calls
	s.mockService.AssertExpectations(s.T())
}

// TestGetUser_Success checks if user retrieval works
func (s *UserControllerTestSuite) TestGetUser_Success() {
	// Create a test user
	testUser := &models.User{
		UserID: uuid.New().String(),
		Name:   "Test User",
		PIN:    "123456",
	}

	// Setup mock expectations
	s.mockService.On("GetUserByID", mock.Anything).Return(testUser, nil)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/users/profile", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.testToken)

	// Test the endpoint
	resp, err := s.app.Test(req)
	s.NoError(err)

	// Check response
	s.Equal(fiber.StatusOK, resp.StatusCode)

	var respUser models.User
	err = json.NewDecoder(resp.Body).Decode(&respUser)
	s.NoError(err)

	s.Equal(testUser.UserID, respUser.UserID)
	s.Equal(testUser.Name, respUser.Name)
	s.Empty(respUser.PIN) // should not show pin

	// Verify expected method calls
	s.mockService.AssertExpectations(s.T())
}

// TestGetUser_NotFound checks if missing user returns 404
func (s *UserControllerTestSuite) TestGetUser_NotFound() {
	// Setup mock expectations
	s.mockService.On("GetUserByID", mock.Anything).Return(nil, errors.New("user not found"))

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/users/profile", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.testToken)

	// Test the endpoint
	resp, err := s.app.Test(req)
	s.NoError(err)

	// Check response
	s.Equal(fiber.StatusNotFound, resp.StatusCode)

	var respBody map[string]string
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	s.NoError(err)

	s.Equal("User not found", respBody["message"])

	// Verify expected method calls
	s.mockService.AssertExpectations(s.T())
}

// TestGetUser_Unauthorized checks if missing token returns 401
func (s *UserControllerTestSuite) TestGetUser_Unauthorized() {
	// Create request without token
	req := httptest.NewRequest(http.MethodGet, "/users/profile", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	// No Authorization header

	// Test the endpoint
	resp, err := s.app.Test(req)
	s.NoError(err)

	// Check response
	s.Equal(fiber.StatusUnauthorized, resp.StatusCode)
}

// TestUpdateUser_Success checks if user update works correctly
func (s *UserControllerTestSuite) TestUpdateUser_Success() {
	// Setup mock expectations
	s.mockService.On("UpdateUser", mock.Anything).Return(nil)

	// Create request body
	requestBody := map[string]string{
		"name": "Updated User Name",
	}
	requestJSON, _ := json.Marshal(requestBody)

	// Create request
	req := httptest.NewRequest(http.MethodPatch, "/users/profile", bytes.NewReader(requestJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.testToken)

	// Test the endpoint
	resp, err := s.app.Test(req)
	s.NoError(err)

	// Check response
	s.Equal(fiber.StatusOK, resp.StatusCode)

	var respUser models.User
	err = json.NewDecoder(resp.Body).Decode(&respUser)
	s.NoError(err)

	s.Equal("Updated User Name", respUser.Name)

	// Verify expected method calls
	s.mockService.AssertExpectations(s.T())
}

// TestUpdateUser_BadRequest checks if invalid request body returns 400
func (s *UserControllerTestSuite) TestUpdateUser_BadRequest() {
	// Create request with invalid JSON
	req := httptest.NewRequest(http.MethodPatch, "/users/profile", strings.NewReader("{invalid json}"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.testToken)

	// Test the endpoint
	resp, err := s.app.Test(req)
	s.NoError(err)

	// Check response
	s.Equal(fiber.StatusBadRequest, resp.StatusCode)

	var respBody map[string]string
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	s.NoError(err)

	s.Contains(respBody["message"], "Invalid input format")
}

// TestUpdateUser_NotFound checks if service error returns 404
func (s *UserControllerTestSuite) TestUpdateUser_NotFound() {
	// Setup mock expectations
	s.mockService.On("UpdateUser", mock.Anything).Return(errors.New("user not found"))

	// Create request body
	requestBody := map[string]string{
		"name": "Updated User Name",
	}
	requestJSON, _ := json.Marshal(requestBody)

	// Create request
	req := httptest.NewRequest(http.MethodPatch, "/users/profile", bytes.NewReader(requestJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.testToken)

	// Test the endpoint
	resp, err := s.app.Test(req)
	s.NoError(err)

	// Check response
	s.Equal(fiber.StatusNotFound, resp.StatusCode)

	var respBody map[string]string
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	s.NoError(err)

	s.Equal("User not found", respBody["message"])

	// Verify expected method calls
	s.mockService.AssertExpectations(s.T())
}

// Run the test suite
func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}
