package middleware

import (
	"log"
	"os"

	fiberzap "github.com/gofiber/contrib/fiberzap/v2"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
)

// FiberMiddleware provide Fiber's built-in middlewares.
func FiberMiddleware(a *fiber.App) {
	// Initialize the global logger
	InitLogger()

	a.Use(
		// Add CORS to each route.
		cors.New(),
		// logger
		fiberzap.New(fiberzap.Config{
			Logger: GetLogger(),
		}),
		// Compression
		compress.New(compress.Config{
			Level: compress.LevelBestSpeed,
		}),
	)
}

// Global logger instance
var globalLogger *zap.Logger

// InitLogger initializes the global Zap logger
func InitLogger() {
	var logConfig zap.Config
	// Get log level from environment variable or use production as default
	if os.Getenv("APP_ENV") == "dev" {
		logConfig = zap.NewDevelopmentConfig()
	} else {
		logConfig = zap.NewDevelopmentConfig()
	}

	if os.Getenv("LOG_LEVEL") == "DEBUG" {
		logConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	var err error
	globalLogger, err = logConfig.Build()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
}

// GetLogger returns the global logger instance
func GetLogger() *zap.Logger {
	if globalLogger == nil {
		InitLogger()
	}
	return globalLogger
}
