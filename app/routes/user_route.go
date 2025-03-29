package routes

import (
	"backend-developer-assignment/app/controllers"
	"backend-developer-assignment/pkg/middleware"

	fiber "github.com/gofiber/fiber/v2"
)

func UserRoute(route fiber.Router, controller *controllers.Controller) {
	route.Post("/greeting", middleware.JWTProtected(), controller.AuthController.VerifyPin)
}
