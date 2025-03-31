package controllers

import (
	"backend-developer-assignment/app/services"
	"backend-developer-assignment/pkg/middleware"

	fiber "github.com/gofiber/fiber/v2"
)

type Controller struct {
	AuthController        AuthController
	UserController        UserController
	TransactionController TransactionController
	DebitCardController   DebitCardController
	AccountController     AccountController
}

var logger = middleware.GetLogger()

func InitController(service *services.Service) *Controller {
	return &Controller{
		AuthController:        *NewAuthController(service.UserService),
		UserController:        *NewUserController(service.UserService),
		TransactionController: *NewTransactionController(service.TransactionService),
		DebitCardController:   *NewDebitCardController(service.DebitCardService),
		AccountController:     *NewAccountController(service.AccountService),
	}
}

func ErrorResponse(ctx *fiber.Ctx, statusCode int, message string) error {
	return ctx.Status(statusCode).JSON(fiber.Map{
		"code":    string(statusCode),
		"message": message,
	})
}
