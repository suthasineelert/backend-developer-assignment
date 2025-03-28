package controllers

import "backend-developer-assignment/app/services"

type Controller struct {
	AuthController AuthController
	// Add more controllers if needed (e.g., TransactionController, AccountController)
}

func InitController(service *services.Service) *Controller {
	return &Controller{
		AuthController: *NewAuthController(service.UserService),
	}
}
