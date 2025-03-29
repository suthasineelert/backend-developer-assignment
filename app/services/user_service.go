package services

import (
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/app/repositories"
)

type UserService interface {
	GetUserByID(id string) (*models.User, error)
	GetUserGreetingByID(id string) (*models.UserGreeting, error)
	UpdateUserGreeting(greeting *models.UserGreeting) error
	UpdateUser(user *models.User) error
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
func (s *UserServiceImpl) UpdateUser(user *models.User) error {
	return s.UserRepository.Update(user)
}

// GetUserByID retrieves a user greeting by ID.
func (s *UserServiceImpl) GetUserGreetingByID(id string) (*models.UserGreeting, error) {
	return s.UserGreetingRepository.GetByID(id)
}

// UpdateUserGreeting updates a user greeting.
func (s *UserServiceImpl) UpdateUserGreeting(greeting *models.UserGreeting) error {
	return s.UserGreetingRepository.Update(greeting)
}
