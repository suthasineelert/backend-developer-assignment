package mocks

import (
	"backend-developer-assignment/app/models"

	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of the UserService
type MockUserService struct {
	mock.Mock
}

// GetUserByID mocks the GetUserByID method
func (m *MockUserService) GetUserByID(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}
