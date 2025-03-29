package routes

import (
	"backend-developer-assignment/app/controllers"
	"backend-developer-assignment/pkg/middleware"

	fiber "github.com/gofiber/fiber/v2"
)

func UserRoute(route fiber.Router, controller *controllers.Controller) {
	// Group user routes with JWT protection
	userRoutes := route.Group("/user", middleware.AuthProtected()...)
	userRoutes.Get("/greeting", controller.UserController.GetUserGreeting)
	userRoutes.Put("/greeting", controller.UserController.UpdateUserGreeting)
	userRoutes.Get("/profile", controller.UserController.GetUser)
	userRoutes.Patch("/profile", controller.UserController.UpdateUser)
}
