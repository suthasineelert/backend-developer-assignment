package controllers

import (
	"backend-developer-assignment/app/controllers"
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/app/queries"
	"backend-developer-assignment/pkg/utils"
	"backend-developer-assignment/platform/database"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserQueries is a mock implementation of the UserQueries interface
type MockUserQueries struct {
	mock.Mock
}

func (m *MockUserQueries) GetUserByID(id string) (models.User, error) {
	args := m.Called(id)
	return args.Get(0).(models.User), args.Error(1)
}

func TestVerifyPin(t *testing.T) {
	// Set up environment for JWT
	os.Setenv("JWT_SECRET_KEY", "test-secret-key")

	// Create a new Fiber app
	app := fiber.New()

	// Register the VerifyPin handler
	app.Post("/api/auth/verify-pin", controllers.VerifyPin)

	hashPin, _ := utils.HashPIN("123456")

	// Test cases
	tests := []struct {
		description  string
		setupMock    func(*MockUserQueries)
		requestBody  map[string]string
		expectedCode int
		expectedBody map[string]interface{}
	}{
		{
			description: "Valid PIN verification",
			setupMock: func(mockQueries *MockUserQueries) {
				// Create a user with a hashed PIN (123456)
				userID := uuid.New().String()
				mockQueries.On("GetUserByID", userID).Return(models.User{
					UserID: uuid.MustParse(userID),
					Name:   "Test User",
					PIN:    hashPin,
				}, nil)
			},
			requestBody: map[string]string{
				"user_id": uuid.New().String(),
				"pin":     "123456",
			},
			expectedCode: fiber.StatusOK,
			expectedBody: map[string]interface{}{
				"token": "", // We'll just check if it exists
			},
		},
		{
			description: "User not found",
			setupMock: func(mockQueries *MockUserQueries) {
				mockQueries.On("GetUserByID", mock.Anything).Return(models.User{}, fiber.NewError(fiber.StatusNotFound, "User not found"))
			},
			requestBody: map[string]string{
				"user_id": uuid.New().String(),
				"pin":     "123456",
			},
			expectedCode: fiber.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": true,
				"msg":   "User does not exist",
			},
		},
		{
			description: "Invalid PIN",
			setupMock: func(mockQueries *MockUserQueries) {
				userID := uuid.New().String()
				mockQueries.On("GetUserByID", userID).Return(models.User{
					UserID: uuid.MustParse(userID),
					Name:   "Test User",
					PIN:    hashPin,
				}, nil)
			},
			requestBody: map[string]string{
				"user_id": uuid.New().String(),
				"pin":     "wrong-pin",
			},
			expectedCode: fiber.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": true,
				"msg":   "Invalid PIN",
			},
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			// Create the mock object
			mockUserQueries := new(MockUserQueries)

			// Setup the mock expectations
			test.setupMock(mockUserQueries)

			// Create a wrapper that satisfies the queries.UserQueries interface
			userQueriesWrapper := &queries.UserQueries{
				DB: nil, // We don't need a real DB for the tests
			}

			// Replace the real DB function with our mock
			originalDB := database.GetDB
			database.GetDB = func() *database.Queries {
				return &database.Queries{
					UserQueries: userQueriesWrapper,
					DB:          nil,
				}
			}

			// Override the GetUserByID method to use our mock
			originalGetUserByID := userQueriesWrapper.GetUserByID
			userQueriesWrapper.GetUserByID = func(id string) (models.User, error) {
				return mockUserQueries.GetUserByID(id)
			}

			// Restore original functions after test
			defer func() {
				database.GetDB = originalDB
				if originalGetUserByID != nil {
					userQueriesWrapper.GetUserByID = originalGetUserByID
				}
			}()

			// Create request
			reqBody, _ := json.Marshal(test.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/auth/verify-pin", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			// Test the endpoint
			resp, err := app.Test(req)
			assert.NoError(t, err)

			// Check status code
			assert.Equal(t, test.expectedCode, resp.StatusCode)

			// Parse response body
			var respBody map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&respBody)

			// For successful responses, just check if token exists
			if test.expectedCode == fiber.StatusOK {
				assert.Contains(t, respBody, "token")
				assert.NotEmpty(t, respBody["token"])
			} else {
				// For error responses, check the exact response
				assert.Equal(t, test.expectedBody, respBody)
			}

			// Verify that all expected mock calls were made
			mockUserQueries.AssertExpectations(t)
		})
	}
}
