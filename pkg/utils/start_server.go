package utils

import (
	"backend-developer-assignment/platform/cache"
	"log"
	"os"
	"os/signal"
	"syscall"

	fiber "github.com/gofiber/fiber/v2"
)

// StartServerWithGracefulShutdown function for starting server with a graceful shutdown.
func StartServerWithGracefulShutdown(app *fiber.App, redisClient *cache.RedisClient) {
	// Build Fiber connection URL
	fiberConnURL, _ := ConnectionURLBuilder("fiber")

	// Listen from a different goroutine
	go func() {
		if err := app.Listen(fiberConnURL); err != nil {
			log.Fatalf("Server is not running! Reason: %v", err)
		}
	}()

	sigChannel := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	sig := <-sigChannel // This blocks the main thread until an interrupt is received
	log.Printf("Received signal: %v. Shutting down gracefully...", sig)

	// Shutdown Fiber server
	if err := app.Shutdown(); err != nil {
		// Error from closing listeners, or context timeout:
		log.Fatalf("Server is not shutting down! Reason: %v", err)
	}

	log.Println("Running cleanup tasks...")
	redisClient.Close()

	log.Println("Fiber was successful shutdown.")
}
