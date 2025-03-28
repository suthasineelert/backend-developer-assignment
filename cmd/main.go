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

	"backend-developer-assignment/app/controllers"
	"backend-developer-assignment/app/repositories"
	"backend-developer-assignment/app/routes"
	"backend-developer-assignment/app/services"
	"backend-developer-assignment/pkg/configs"
	"backend-developer-assignment/pkg/middleware"
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
	// Open the DB connection
	db, err := database.MysqlConnection()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	// Ensure database connection is closed when application exits
	defer db.Close()

	// Migrate database
	err = database.Migrate(db)
	if err != nil {
		log.Fatal(err)
	}

	// Define Fiber config.
	config := configs.FiberConfig()

	// Define a new Fiber app with config.
	app := fiber.New(config)

	// Middlewares.
	middleware.FiberMiddleware(app) // Register Fiber's middleware for app.

	// Initialize repositories, services, and controllers
	repositories := repositories.InitRepository(db)
	services := services.InitService(repositories)
	controllers := controllers.InitController(services)
	// Routes
	routes.InitRoutes(app, controllers)

	// Start server (with or without graceful shutdown).
	if os.Getenv("APP_ENV") == "dev" {
		utils.StartServer(app)
	} else {
		utils.StartServerWithGracefulShutdown(app)
	}
}
