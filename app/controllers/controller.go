package controllers

import "backend-developer-assignment/app/services"

type Controller struct {
	AuthController AuthController
	UserController UserController
}

func InitController(service *services.Service) *Controller {
	return &Controller{
		AuthController: *NewAuthController(service.UserService),
		UserController: *NewUserController(service.UserService),
	}
}
