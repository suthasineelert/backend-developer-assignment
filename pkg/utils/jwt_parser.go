package utils

import (
	"os"
	"strings"

	fiber "github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
)

// TokenMetadata struct to describe metadata in JWT.
type TokenMetadata struct {
	UserID  string
	Expires int64
}

// ExtractTokenMetadata func to extract metadata from JWT.
func ExtractTokenMetadata(ctx *fiber.Ctx) (*TokenMetadata, error) {
	token, err := verifyToken(ctx)
	if err != nil {
		return nil, err
	}

	// Setting and checking token and credentials.
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		// User ID.
		userID, ok := claims["id"].(string)
		if !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid token: missing user ID")
		}

		// Expires time.
		expiresFloat, ok := claims["exp"].(float64)
		if !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid token: missing expiration")
		}
		expires := int64(expiresFloat)

		return &TokenMetadata{
			UserID:  userID,
			Expires: expires,
		}, nil
	}

	return nil, err
}

func extractToken(c *fiber.Ctx) string {
	bearToken := c.Get("Authorization")

	// Normally Authorization HTTP header.
	onlyToken := strings.Split(bearToken, " ")
	if len(onlyToken) == 2 {
		return onlyToken[1]
	}

	return ""
}

func verifyToken(c *fiber.Ctx) (*jwt.Token, error) {
	tokenString := extractToken(c)

	token, err := jwt.Parse(tokenString, jwtKeyFunc)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(os.Getenv("JWT_SECRET_KEY")), nil
}
