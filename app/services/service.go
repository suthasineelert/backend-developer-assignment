package services

import "backend-developer-assignment/app/repositories"

type Service struct {
	UserService UserService
	// Add more services if needed (e.g., TransactionService, AccountService)
}

func InitService(repo *repositories.Repository) *Service {
	return &Service{
		UserService: NewUserService(repo.UserRepo),
	}
}
