package middleware

import (
	"backend-developer-assignment/pkg/base"
	"backend-developer-assignment/pkg/utils"
	"os"

	fiber "github.com/gofiber/fiber/v2"

	jwtMiddleware "github.com/gofiber/contrib/jwt"
)

// AuthProtected combines JWT protection and user ID extraction in a single middleware chain
func AuthProtected() []fiber.Handler {
	return []fiber.Handler{
		JWTProtected(),
		ExtractJwtClaim(),
	}
}

// JWTProtected func for specify routes group with JWT authentication.
// See: https://github.com/gofiber/contrib/jwt
func JWTProtected() func(*fiber.Ctx) error {
	// Create config for JWT authentication middleware.
	config := jwtMiddleware.Config{
		SigningKey:   jwtMiddleware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET_KEY"))},
		ContextKey:   "jwt", // used in private routes
		ErrorHandler: jwtError,
	}

	return jwtMiddleware.New(config)
}

// ExtractJwtClaim middleware extracts the user ID from JWT token and stores it in context
// This middleware should be used after JWTProtected middleware
func ExtractJwtClaim() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Use the utils package to verify and extract token metadata
		tokenMetadata, err := utils.ExtractTokenMetadata(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(base.ErrorResponse{
				Message: "Invalid token: " + err.Error(),
			})
		}

		// Store user ID in context for later use in controllers
		c.Locals("userID", tokenMetadata.UserID)

		// Continue to the next middleware/handler
		return c.Next()
	}
}

func jwtError(c *fiber.Ctx, err error) error {
	// Return status 401 and failed authentication error.
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).JSON(base.ErrorResponse{
			Message: err.Error(),
		})
	}

	// Return status 401 and failed authentication error.
	return c.Status(fiber.StatusUnauthorized).JSON(base.ErrorResponse{
		Message: err.Error(),
	})
}
