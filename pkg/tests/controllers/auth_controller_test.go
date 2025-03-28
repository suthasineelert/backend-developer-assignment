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
)

// setupTest creates a new Fiber app, mock service, and controller for testing
func setupTest() (*fiber.App, *mocks.MockUserService, *controllers.AuthController) {
	app := fiber.New()
	mockService := new(mocks.MockUserService)
	authController := controllers.NewAuthController(mockService)
	app.Post("/verify-pin", authController.VerifyPin)
	return app, mockService, authController
}

// testResponse tests the response from the API
func testResponse(t *testing.T, resp *http.Response, expectedCode int, expectedBody map[string]interface{}) {
	assert.Equal(t, expectedCode, resp.StatusCode)

	var respBody map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respBody)

	if expectedCode == fiber.StatusOK {
		assert.Contains(t, respBody, "tokens")
		tokens, ok := respBody["tokens"].(map[string]interface{})
		assert.True(t, ok)
		assert.Contains(t, tokens, "access")
		assert.Contains(t, tokens, "refresh")
	} else {
		assert.Equal(t, expectedBody, respBody)
	}
}

func TestVerifyPin_Success(t *testing.T) {
	// Set up environment for JWT
	os.Setenv("JWT_SECRET_KEY", "test-secret-key")

	// Setup test
	app, mockService, _ := setupTest()

	// Hash a PIN for testing
	pin := "123456"
	hashPin, _ := utils.HashPIN(pin)
	userID := uuid.New().String()

	// Setup mock expectations
	mockService.On("GetUserByID", userID).Return(&models.User{
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
	resp, err := app.Test(req)
	assert.NoError(t, err)

	// Test response
	expectedBody := map[string]interface{}{
		"tokens": map[string]interface{}{
			"access":  "",
			"refresh": "",
		},
	}
	testResponse(t, resp, fiber.StatusOK, expectedBody)

	// Verify that all expected calls were made
	mockService.AssertExpectations(t)
}

func TestVerifyPin_UserNotFound(t *testing.T) {
	// Setup test
	app, mockService, _ := setupTest()

	// Generate a user ID
	userID := uuid.New().String()

	// Setup mock expectations
	mockService.On("GetUserByID", userID).Return(nil, errors.New("user not found"))

	// Create request
	reqBody, _ := json.Marshal(map[string]string{
		"user_id": userID,
		"pin":     "123456",
	})
	req := httptest.NewRequest(http.MethodPost, "/verify-pin", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Test the endpoint
	resp, err := app.Test(req)
	assert.NoError(t, err)

	// Test response
	expectedBody := map[string]interface{}{
		"message": "User does not exist",
	}
	testResponse(t, resp, fiber.StatusNotFound, expectedBody)

	// Verify that all expected calls were made
	mockService.AssertExpectations(t)
}

func TestVerifyPin_InvalidPIN(t *testing.T) {
	// Setup test
	app, mockService, _ := setupTest()

	// Hash a PIN for testing
	pin := "123456"
	hashPin, _ := utils.HashPIN(pin)
	userID := uuid.New().String()

	// Setup mock expectations
	mockService.On("GetUserByID", userID).Return(&models.User{
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
	resp, err := app.Test(req)
	assert.NoError(t, err)

	// Test response
	expectedBody := map[string]interface{}{
		"message": "Invalid PIN",
	}
	testResponse(t, resp, fiber.StatusUnauthorized, expectedBody)

	// Verify that all expected calls were made
	mockService.AssertExpectations(t)
}
