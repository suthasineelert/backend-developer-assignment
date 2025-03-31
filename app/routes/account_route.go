package routes

import (
	"backend-developer-assignment/app/controllers"
	"backend-developer-assignment/pkg/middleware"

	fiber "github.com/gofiber/fiber/v2"
)

func AccountRoute(route fiber.Router, controller *controllers.Controller) {
	// Group user routes with JWT protection
	accountRoutes := route.Group("/accounts", middleware.AuthProtected()...)
	accountRoutes.Get("/", controller.AccountController.ListAccounts)
	accountRoutes.Get("/:id", controller.AccountController.GetAccount)
	accountRoutes.Patch("/:id", controller.AccountController.UpdateAccount)
	accountRoutes.Post("", controller.AccountController.CreateAccount)
	accountRoutes.Post("/:id/main", controller.AccountController.SetMainAccount)
	accountRoutes.Post("/:id/deposit", controller.AccountController.Deposit)
	accountRoutes.Post("/:id/withdraw", controller.AccountController.Withdraw)
	accountRoutes.Post("/:id/transfer", controller.AccountController.Transfer)

}
