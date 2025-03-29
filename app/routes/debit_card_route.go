package routes

import (
	"backend-developer-assignment/app/controllers"
	"backend-developer-assignment/pkg/middleware"

	fiber "github.com/gofiber/fiber/v2"
)

func DebitCardRoute(route fiber.Router, controller *controllers.Controller) {
	// Group user routes with JWT protection
	debitCardRoutes := route.Group("/debit-cards", middleware.AuthProtected()...)
	debitCardRoutes.Get("", controller.DebitCardController.ListDebitCards)
	debitCardRoutes.Get("/:id", controller.DebitCardController.GetDebitCard)
	debitCardRoutes.Post("", controller.DebitCardController.CreateDebitCard)
	debitCardRoutes.Put("/:id", controller.DebitCardController.UpdateDebitCard)
	debitCardRoutes.Delete("/:id", controller.DebitCardController.DeleteDebitCard)
}
