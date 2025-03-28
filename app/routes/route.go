package routes

import (
	"backend-developer-assignment/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App, controller *controllers.Controller) {
	route := app.Group("/api/v1")

	AuthRoute(route, controller)

	// route.Post("/book", middleware.JWTProtected(), controllers.CreateBook)

	SwaggerRoute(app)  // Register a route for API Docs (Swagger).
	NotFoundRoute(app) // Register route for 404 Error.
}
