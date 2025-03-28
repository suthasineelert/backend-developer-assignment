package services

import (
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/app/repositories"
)

type UserService interface {
	GetUserByID(id string) (*models.User, error)
}

// UserService contains business logic related to users.
type UserServiceImpl struct {
	UserRepository repositories.UserRepository
}

// NewUserService creates a new UserService.
func NewUserService(userRepository repositories.UserRepository) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
	}
}

// GetUserByID retrieves a user by ID.
func (s *UserServiceImpl) GetUserByID(id string) (*models.User, error) {
	return s.UserRepository.GetByID(id)
}
