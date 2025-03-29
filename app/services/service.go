package services

import "backend-developer-assignment/app/repositories"

type Service struct {
	UserService        UserService
	TransactionService TransactionService
	DebitCardService   DebitCardService
}

func InitService(repo *repositories.Repository) *Service {
	return &Service{
		UserService:        NewUserService(repo.UserRepository, repo.UserGreetingsRepository),
		TransactionService: NewTransactionService(repo.TransactionRepository),
		DebitCardService:   NewDebitCardService(repo.DebitCardRepository),
	}
}
