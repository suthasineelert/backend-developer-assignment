package services

import (
	"backend-developer-assignment/app/repositories"
	"backend-developer-assignment/pkg/middleware"
)

type Service struct {
	UserService        UserService
	TransactionService TransactionService
	DebitCardService   DebitCardService
	AccountService     AccountService
}

var logger = middleware.GetLogger()

func InitService(repo *repositories.Repository, txProvider repositories.TxProvider) *Service {
	return &Service{
		UserService:        NewUserService(repo.UserRepository, repo.UserGreetingsRepository),
		TransactionService: NewTransactionService(repo.TransactionRepository),
		DebitCardService:   NewDebitCardService(repo.DebitCardRepository),
		AccountService:     NewAccountService(repo.AccountRepository, repo.TransactionRepository, txProvider),
	}
}
