package routes

import (
	"backend-developer-assignment/app/controllers"

	fiber "github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App, controller *controllers.Controller) {
	route := app.Group("/api/v1")

	AuthRoute(route, controller)
	UserRoute(route, controller)
	AccountRoute(route, controller)
	TransactionRoute(route, controller)
	DebitCardRoute(route, controller)
	BannerRoute(route, controller)

	SwaggerRoute(app)  // Register a route for API Docs (Swagger).
	NotFoundRoute(app) // Register route for 404 Error.
}
