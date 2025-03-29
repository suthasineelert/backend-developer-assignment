package services

import (
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/app/repositories"
)

type UserService interface {
	GetUserByID(id string) (*models.User, error)
	GetUserGreetingByID(id string) (*models.UserGreeting, error)
}

// UserService contains business logic related to users.
type UserServiceImpl struct {
	UserRepository         repositories.UserRepository
	UserGreetingRepository repositories.UserGreetingRepository
}

// NewUserService creates a new UserService.
func NewUserService(userRepository repositories.UserRepository, userGreetingRepository repositories.UserGreetingRepository) UserService {
	return &UserServiceImpl{
		UserRepository:         userRepository,
		UserGreetingRepository: userGreetingRepository,
	}
}

// GetUserByID retrieves a user by ID.
func (s *UserServiceImpl) GetUserByID(id string) (*models.User, error) {
	return s.UserRepository.GetByID(id)
}

// GetUserByID retrieves a user greeting by ID.
func (s *UserServiceImpl) GetUserGreetingByID(id string) (*models.UserGreeting, error) {
	return s.UserGreetingRepository.GetByID(id)
}
