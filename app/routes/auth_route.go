package routes

import (
	"backend-developer-assignment/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func AuthRoute(route fiber.Router, controller *controllers.Controller) {
	route.Post("/auth/verify-pin", controller.AuthController.VerifyPin)
}
