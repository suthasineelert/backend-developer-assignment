package routes

import (
	"backend-developer-assignment/app/controllers"
	"backend-developer-assignment/pkg/middleware"

	fiber "github.com/gofiber/fiber/v2"
)

func AuthRoute(route fiber.Router, controller *controllers.Controller) {
	route.Post("/auth/verify-pin", controller.AuthController.VerifyPin)
	route.Post("/token/renew", middleware.JWTProtected(), controller.AuthController.RenewTokens)
}
