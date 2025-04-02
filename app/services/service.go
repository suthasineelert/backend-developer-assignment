package services

import (
	"backend-developer-assignment/app/repositories"
	"backend-developer-assignment/pkg/middleware"
	"backend-developer-assignment/pkg/types"
)

type Service struct {
	UserService        UserService
	TransactionService TransactionService
	DebitCardService   DebitCardService
	AccountService     AccountService
	BannerService      BannerService
}

var logger = middleware.GetLogger()

func InitService(repo *repositories.Repository, txProvider repositories.TxProvider, redisClient types.CacheClient) *Service {
	return &Service{
		UserService:        NewUserService(repo.UserRepository, repo.UserGreetingsRepository),
		TransactionService: NewTransactionService(repo.TransactionRepository, redisClient),
		DebitCardService:   NewDebitCardService(repo.DebitCardRepository),
		AccountService:     NewAccountService(repo.AccountRepository, repo.TransactionRepository, txProvider),
		BannerService:      NewBannerService(repo.BannerRepository),
	}
}
