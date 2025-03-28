package services

import (
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/app/services"
	mocks "backend-developer-assignment/pkg/mocks/repositories"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByID(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name          string
		userID        string
		mockUser      *models.User
		mockError     error
		expectedUser  *models.User
		expectedError error
	}{
		{
			name:   "Success - Valid User",
			userID: "000018b0e1a211ef95a30242ac180002",
			mockUser: &models.User{
				UserID: uuid.MustParse("000018b0e1a211ef95a30242ac180002"),
				Name:   "Test User",
				PIN:    "hashedpin123",
			},
			mockError:     nil,
			expectedUser:  &models.User{UserID: uuid.MustParse("000018b0e1a211ef95a30242ac180002"), Name: "Test User", PIN: "hashedpin123"},
			expectedError: nil,
		},
		{
			name:          "Failure - User Not Found",
			userID:        "000018b0e1a211ef95a30242ac180001",
			mockUser:      nil,
			mockError:     errors.New("user not found"),
			expectedUser:  nil,
			expectedError: errors.New("user not found"),
		},
		{
			name:          "Failure - Database Error",
			userID:        "invalid-user-id",
			mockUser:      nil,
			mockError:     errors.New("database connection failed"),
			expectedUser:  nil,
			expectedError: errors.New("database connection failed"),
		},
		{
			name:          "Failure - Empty UserID",
			userID:        "",
			mockUser:      nil,
			mockError:     errors.New("invalid user ID"),
			expectedUser:  nil,
			expectedError: errors.New("invalid user ID"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new mock repository
			mockRepo := new(mocks.MockUserRepository)

			// Set mock expectations
			mockRepo.On("GetByID", tc.userID).Return(tc.mockUser, tc.mockError)

			// Create the service with the mock repository
			service := services.NewUserService(mockRepo)

			// Call the function being tested
			user, err := service.GetUserByID(tc.userID)

			// Assert results using helper function
			assertUserResponse(t, user, err, tc.expectedUser, tc.expectedError)

			// Verify expected method calls
			mockRepo.AssertExpectations(t)
		})
	}
}

// Helper function for assertions
func assertUserResponse(t *testing.T, actualUser *models.User, actualErr error, expectedUser *models.User, expectedErr error) {
	if expectedErr != nil {
		assert.Error(t, actualErr)
		assert.Equal(t, actualErr.Error(), expectedErr.Error(), "Expected error: %v, Got: %v", expectedErr, actualErr)
	} else {
		assert.NoError(t, actualErr)
		assert.True(t, assert.ObjectsAreEqual(expectedUser, actualUser), "Expected user: %+v, Got: %+v", expectedUser, actualUser)
	}
}
