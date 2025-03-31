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
	testUserID  string
}

// SetupSuite runs once before all tests
func (s *UserControllerTestSuite) SetupSuite() {
	// Set up JWT environment for testing
	os.Setenv("JWT_SECRET_KEY", "test-secret-key")

	// Generate a test user ID
	s.testUserID = uuid.New().String()

	// Generate a test token with the test user ID
	s.testToken = s.generateTestToken(s.testUserID)
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
func (s *UserControllerTestSuite) generateTestToken(userID string) string {
	// Create token claims
	claims := jwt.MapClaims{
		"id":  userID,
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
		UserID:   s.testUserID,
		Greeting: "Hello, welcome back!",
	}
	// Setup mock expectations
	s.mockService.On("GetUserGreetingByID", s.testUserID).Return(greeting, nil)

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
	s.mockService.On("GetUserGreetingByID", s.testUserID).Return(nil, errors.New("user greeting not found"))

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
	s.mockService.On("UpdateUserGreeting", mock.MatchedBy(func(greeting *models.UserGreeting) bool {
		return greeting.UserID == s.testUserID && greeting.Greeting == "Updated greeting message"
	})).Return(nil)

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
	s.mockService.On("UpdateUserGreeting", mock.MatchedBy(func(greeting *models.UserGreeting) bool {
		return greeting.UserID == s.testUserID && greeting.Greeting == "Updated greeting message"
	})).Return(errors.New("update failed"))

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
		UserID: s.testUserID,
		Name:   "Test User",
		PIN:    "123456",
	}

	// Setup mock expectations
	s.mockService.On("GetUserByID", s.testUserID).Return(testUser, nil)

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
	s.mockService.On("GetUserByID", s.testUserID).Return(nil, errors.New("user not found"))

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
	s.mockService.On("UpdateUser", mock.MatchedBy(func(user *models.User) bool {
		return user.UserID == s.testUserID && user.Name == "Updated User Name"
	})).Return(nil)

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
	s.Equal(s.testUserID, respUser.UserID)

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
	s.mockService.On("UpdateUser", mock.MatchedBy(func(user *models.User) bool {
		return user.UserID == s.testUserID && user.Name == "Updated User Name"
	})).Return(errors.New("user not found"))

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

// TestUpdateUser_EmptyName checks if empty name validation works
func (s *UserControllerTestSuite) TestUpdateUser_EmptyName() {
	// Create request body with empty name
	requestBody := map[string]string{
		"name": "",
	}
	requestJSON, _ := json.Marshal(requestBody)

	// Create request
	req := httptest.NewRequest(http.MethodPatch, "/users/profile", bytes.NewReader(requestJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.testToken)

	// Setup mock expectations - we expect the service to be called even with empty name
	s.mockService.On("UpdateUser", mock.MatchedBy(func(user *models.User) bool {
		return user.UserID == s.testUserID && user.Name == ""
	})).Return(nil)

	// Test the endpoint
	resp, err := s.app.Test(req)
	s.NoError(err)

	// Check response
	s.Equal(fiber.StatusOK, resp.StatusCode)

	// Verify expected method calls
	s.mockService.AssertExpectations(s.T())
}

// Run the test suite
func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}
