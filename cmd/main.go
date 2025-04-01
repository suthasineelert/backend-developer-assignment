// @title My API
// @version 1.0
// @description This is a sample API
// @host localhost:8080
// @BasePath /api/v1
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	fiber "github.com/gofiber/fiber/v2"

	"backend-developer-assignment/app/controllers"
	"backend-developer-assignment/app/repositories"
	"backend-developer-assignment/app/routes"
	"backend-developer-assignment/app/services"
	"backend-developer-assignment/pkg/configs"
	"backend-developer-assignment/pkg/middleware"
	"backend-developer-assignment/pkg/utils"
	"backend-developer-assignment/platform/database"

	_ "backend-developer-assignment/docs" // to import docs generated by Swag CLI

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
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Open the DB connection
	db, err := database.MysqlConnection()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	// Create a channel to signal when shutdown is complete
	shutdownComplete := make(chan struct{})

	// Set up graceful shutdown handler
	go func() {
		// Channel to listen for termination signals
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		// Wait for termination signal
		<-sigChan

		// Close database connection after server shutdown
		defer func() {
			log.Println("Closing database connection...")
			if err := db.Close(); err != nil {
				log.Printf("Error closing database connection: %v", err)
			} else {
				log.Println("Database connection closed successfully")
			}
			close(shutdownComplete)
		}()
	}()

	// Migrate database
	err = database.Migrate(db)
	if err != nil {
		log.Fatal(err)
	}

	// Define Fiber config.
	config := configs.FiberConfig()

	// Define a new Fiber app with config.
	app := fiber.New(config)

	redisClient := configs.RedisConnection()

	// Middlewares.
	middleware.FiberMiddleware(app) // Register Fiber's middleware for app.

	// Initialize repoList, services, and controllers
	txProvider := repositories.NewTransactionProvider(db)
	repoList := repositories.InitRepository(db)
	serviceList := services.InitService(repoList, txProvider, redisClient)
	controllerList := controllers.InitController(serviceList)
	// Routes
	routes.InitRoutes(app, controllerList)

	utils.StartServerWithGracefulShutdown(app, redisClient)
}
