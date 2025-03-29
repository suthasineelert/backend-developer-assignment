package controllers

import (
	"backend-developer-assignment/app/services"
	"backend-developer-assignment/pkg/middleware"
)

type Controller struct {
	AuthController        AuthController
	UserController        UserController
	TransactionController TransactionController
	DebitCardController   DebitCardController
}

var logger = middleware.GetLogger()

func InitController(service *services.Service) *Controller {
	return &Controller{
		AuthController:        *NewAuthController(service.UserService),
		UserController:        *NewUserController(service.UserService),
		TransactionController: *NewTransactionController(service.TransactionService),
		DebitCardController:   *NewDebitCardController(service.DebitCardService),
	}
}
