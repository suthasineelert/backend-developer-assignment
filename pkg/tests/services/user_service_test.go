package services

import (
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/app/services"
	mocks "backend-developer-assignment/pkg/mocks/repositories"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// UserServiceTestSuite defines the test suite
type UserServiceTestSuite struct {
	suite.Suite
	userRepository         *mocks.UserRepository
	userGreetingRepository *mocks.UserGreetingRepository
	service                services.UserService
}

// SetupTest runs before each test
func (s *UserServiceTestSuite) SetupTest() {
	s.userRepository = new(mocks.UserRepository)
	s.userGreetingRepository = new(mocks.UserGreetingRepository)
	s.service = services.NewUserService(s.userRepository, s.userGreetingRepository)
}

// TestGetUserByID tests the GetUserByID function
func (s *UserServiceTestSuite) TestGetUserByID() {
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
				UserID: "000018b0e1a211ef95a30242ac180002",
				Name:   "Test User",
				PIN:    "hashedpin123",
			},
			mockError:     nil,
			expectedUser:  &models.User{UserID: "000018b0e1a211ef95a30242ac180002", Name: "Test User", PIN: "hashedpin123"},
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
		s.Run(tc.name, func() {
			// Mock the repository method
			s.userRepository.On("GetByID", tc.userID).Return(tc.mockUser, tc.mockError)

			// Call the service method
			user, err := s.service.GetUserByID(tc.userID)

			// Assert results
			assertUserResponse(s.T(), user, err, tc.expectedUser, tc.expectedError)

			// Verify expected method calls
			s.userRepository.AssertExpectations(s.T())
		})
	}
}

// TestUpdateUser tests the UpdateUser function
func (s *UserServiceTestSuite) TestUpdateUser() {
	testCases := []struct {
		name          string
		user          *models.User
		mockError     error
		expectedError error
	}{
		{
			name: "Success - Valid Update",
			user: &models.User{
				UserID: "000018b0e1a211ef95a30242ac180002",
				Name:   "Updated User Name",
				PIN:    "newhashedpin456",
			},
			mockError:     nil,
			expectedError: nil,
		},
		{
			name: "Failure - User Not Found",
			user: &models.User{
				UserID: "000018b0e1a211ef95a30242ac180001",
				Name:   "Non-existent User",
			},
			mockError:     errors.New("user not found"),
			expectedError: errors.New("user not found"),
		},
		{
			name: "Failure - Database Error",
			user: &models.User{
				UserID: "000018b0e1a211ef95a30242ac180002",
				Name:   "Test User",
			},
			mockError:     errors.New("database connection failed"),
			expectedError: errors.New("database connection failed"),
		},
		{
			name: "Failure - Invalid User Data",
			user: &models.User{
				UserID: "",
				Name:   "Invalid User",
			},
			mockError:     errors.New("invalid user data"),
			expectedError: errors.New("invalid user data"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Mock the repository method
			s.userRepository.On("Update", tc.user).Return(tc.mockError)

			// Call the service method
			err := s.service.UpdateUser(tc.user)

			// Assert results
			if tc.expectedError != nil {
				assert.Error(s.T(), err)
				assert.Equal(s.T(), tc.expectedError.Error(), err.Error(), "Expected error: %v, Got: %v", tc.expectedError, err)
			} else {
				assert.NoError(s.T(), err)
			}

			// Verify expected method calls
			s.userRepository.AssertExpectations(s.T())
		})
	}
}

// Helper function for assertions
func assertUserResponse(t *testing.T, actualUser *models.User, actualErr error, expectedUser *models.User, expectedErr error) {
	if expectedErr != nil {
		assert.Error(t, actualErr)
		assert.Equal(t, expectedErr.Error(), actualErr.Error(), "Expected error: %v, Got: %v", expectedErr, actualErr)
	} else {
		assert.NoError(t, actualErr)
		assert.Equal(t, expectedUser, actualUser, "Expected user: %+v, Got: %+v", expectedUser, actualUser)
	}
}

// TestGetUserGreetingByID tests the GetUserGreetingByID function
func (s *UserServiceTestSuite) TestGetUserGreetingByID() {
	testCases := []struct {
		name             string
		userID           string
		mockUserGreeting *models.UserGreeting
		mockError        error
		expectedGreeting *models.UserGreeting
		expectedError    error
	}{
		{
			name:   "Success - Valid User Greeting",
			userID: "000018b0e1a211ef95a30242ac180002",
			mockUserGreeting: &models.UserGreeting{
				UserID:   "000018b0e1a211ef95a30242ac180002",
				Greeting: "Hello, welcome back!",
			},
			mockError: nil,
			expectedGreeting: &models.UserGreeting{
				UserID:   "000018b0e1a211ef95a30242ac180002",
				Greeting: "Hello, welcome back!",
			},
			expectedError: nil,
		},
		{
			name:             "Failure - User Greeting Not Found",
			userID:           "000018b0e1a211ef95a30242ac180001",
			mockUserGreeting: nil,
			mockError:        errors.New("user greeting not found"),
			expectedGreeting: nil,
			expectedError:    errors.New("user greeting not found"),
		},
		{
			name:             "Failure - Database Error",
			userID:           "invalid-user-id",
			mockUserGreeting: nil,
			mockError:        errors.New("database connection failed"),
			expectedGreeting: nil,
			expectedError:    errors.New("database connection failed"),
		},
		{
			name:             "Failure - Empty UserID",
			userID:           "",
			mockUserGreeting: nil,
			mockError:        errors.New("invalid user ID"),
			expectedGreeting: nil,
			expectedError:    errors.New("invalid user ID"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Mock the repository method
			s.userGreetingRepository.On("GetByID", tc.userID).Return(tc.mockUserGreeting, tc.mockError)

			// Call the service method
			greeting, err := s.service.GetUserGreetingByID(tc.userID)

			// Assert results
			assertUserGreetingResponse(s.T(), greeting, err, tc.expectedGreeting, tc.expectedError)

			// Verify expected method calls
			s.userGreetingRepository.AssertExpectations(s.T())
		})
	}
}

// TestUpdateUserGreeting tests the UpdateUserGreeting function
func (s *UserServiceTestSuite) TestUpdateUserGreeting() {
	testCases := []struct {
		name          string
		greeting      *models.UserGreeting
		mockError     error
		expectedError error
	}{
		{
			name: "Success - Valid Update",
			greeting: &models.UserGreeting{
				UserID:   "000018b0e1a211ef95a30242ac180002",
				Greeting: "Updated greeting message!",
			},
			mockError:     nil,
			expectedError: nil,
		},
		{
			name: "Failure - User Greeting Not Found",
			greeting: &models.UserGreeting{
				UserID:   "000018b0e1a211ef95a30242ac180001",
				Greeting: "Non-existent greeting",
			},
			mockError:     errors.New("user greeting not found"),
			expectedError: errors.New("user greeting not found"),
		},
		{
			name: "Failure - Database Error",
			greeting: &models.UserGreeting{
				UserID:   "000018b0e1a211ef95a30242ac180002",
				Greeting: "Test greeting",
			},
			mockError:     errors.New("database connection failed"),
			expectedError: errors.New("database connection failed"),
		},
		{
			name: "Failure - Invalid Greeting Data",
			greeting: &models.UserGreeting{
				UserID:   "",
				Greeting: "Invalid greeting",
			},
			mockError:     errors.New("invalid greeting data"),
			expectedError: errors.New("invalid greeting data"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Mock the repository method
			s.userGreetingRepository.On("Update", tc.greeting).Return(tc.mockError)

			// Call the service method
			err := s.service.UpdateUserGreeting(tc.greeting)

			// Assert results
			if tc.expectedError != nil {
				assert.Error(s.T(), err)
				assert.Equal(s.T(), tc.expectedError.Error(), err.Error(), "Expected error: %v, Got: %v", tc.expectedError, err)
			} else {
				assert.NoError(s.T(), err)
			}

			// Verify expected method calls
			s.userGreetingRepository.AssertExpectations(s.T())
		})
	}
}

// Helper function for greeting assertions
func assertUserGreetingResponse(t *testing.T, actualGreeting *models.UserGreeting, actualErr error, expectedGreeting *models.UserGreeting, expectedErr error) {
	if expectedErr != nil {
		assert.Error(t, actualErr)
		assert.Equal(t, expectedErr.Error(), actualErr.Error(), "Expected error: %v, Got: %v", expectedErr, actualErr)
	} else {
		assert.NoError(t, actualErr)
		assert.Equal(t, expectedGreeting, actualGreeting, "Expected greeting: %+v, Got: %+v", expectedGreeting, actualGreeting)
	}
}

// Run the test suite
func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
