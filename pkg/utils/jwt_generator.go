package utils

import (
	"os"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

// GenerateNewTokens func for generate a new Access & Refresh tokens.
func GenerateNewTokens(id string) (*string, error) {
	// Generate JWT Access token.
	accessToken, err := generateNewAccessToken(id)
	if err != nil {
		// Return token generation error.
		return nil, err
	}

	return &accessToken, nil
}

func generateNewAccessToken(id string) (string, error) {
	// Set secret key from .env file.
	secret := os.Getenv("JWT_SECRET_KEY")

	// Set expires minutes count for secret key from .env file.
	minutesCount, _ := strconv.Atoi(os.Getenv("JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT"))

	// Create a new claims.
	claims := jwt.MapClaims{}

	// Set public claims:
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(minutesCount)).Unix()

	// Create a new JWT access token with claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate token.
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		// Return error, it JWT token generation failed.
		return "", err
	}

	return t, nil
}
