// @title My API
// @version 1.0
// @description This is a sample API
// @host localhost:8080
// @BasePath /api/v1
package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"

	"backend-developer-assignment/pkg/configs"
	"backend-developer-assignment/pkg/middleware"
	"backend-developer-assignment/pkg/routes"
	"backend-developer-assignment/pkg/utils"
	"backend-developer-assignment/platform/database"

	_ "github.com/joho/godotenv/autoload" // load .env file automatically
)

// @title Backend Developer Assignment API
// @version 1.0
// @description This is an auto-generated API Docs.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email suthasinee.ler@mail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Define Fiber config.
	config := configs.FiberConfig()

	// Define a new Fiber app with config.
	app := fiber.New(config)

	// Middlewares.
	middleware.FiberMiddleware(app) // Register Fiber's middleware for app.

	// Migrate database
	err := database.Migrate()
	if err != nil {
		log.Fatal(err)
	}

	// Routes.
	routes.SwaggerRoute(app)  // Register a route for API Docs (Swagger).
	routes.PublicRoutes(app)  // Register a public routes for app.
	routes.PrivateRoutes(app) // Register a private routes for app.
	routes.NotFoundRoute(app) // Register route for 404 Error.

	// Start server (with or without graceful shutdown).
	if os.Getenv("APP_ENV") == "dev" {
		utils.StartServer(app)
	} else {
		utils.StartServerWithGracefulShutdown(app)
	}
}
