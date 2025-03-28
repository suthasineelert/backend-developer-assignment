package mocks

import (
	"backend-developer-assignment/app/models"

	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of the UserRepository
type MockUserRepository struct {
	mock.Mock
}

// GetByID mocks the GetByID method of UserRepository
func (m *MockUserRepository) GetByID(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// Add any other methods required by the repositories.UserRepository interface
func (m *MockUserRepository) GetByName(name string) (*models.User, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}
