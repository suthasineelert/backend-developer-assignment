package routes

import (
	"backend-developer-assignment/app/controllers"
	"backend-developer-assignment/pkg/middleware"

	fiber "github.com/gofiber/fiber/v2"
)

func TransactionRoute(route fiber.Router, controller *controllers.Controller) {
	// Group user routes with JWT protection
	transactionRoutes := route.Group("/transaction", middleware.AuthProtected()...)
	transactionRoutes.Get("", controller.TransactionController.ListTransactions)
}
