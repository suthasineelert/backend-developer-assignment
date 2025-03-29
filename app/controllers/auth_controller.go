package controllers

import (
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/app/services"
	"backend-developer-assignment/pkg/base"
	"backend-developer-assignment/pkg/utils"
	"fmt"
	"log"
	"time"

	fiber "github.com/gofiber/fiber/v2"
)

// AuthController holds the services related to users.
type AuthController struct {
	UserService services.UserService
}

// NewAuthController creates a new AuthController.
func NewAuthController(userService services.UserService) *AuthController {
	return &AuthController{
		UserService: userService,
	}
}

// VerifyPin method for user PIN verification.
// @Description Verify user PIN and return JWT token
// @Summary Verify PIN and get JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body controllers.VerifyPin.verifyPinRequest true "PIN verification request"
// @Success 200 {object} object{tokens=object{access=string,refresh=string}} "JWT tokens"
// @Failure 400 {object} base.ErrorResponse "Invalid input format"
// @Failure 401 {object} base.ErrorResponse "Invalid PIN"
// @Failure 404 {object} base.ErrorResponse "User does not exist"
// @Failure 500 {object} base.ErrorResponse "Failed to generate token"
// @Router /auth/verify-pin [post]
func (c *AuthController) VerifyPin(ctx *fiber.Ctx) error {
	type verifyPinRequest struct {
		UserID string `json:"user_id"`
		PIN    string `json:"pin"`
	}
	var request verifyPinRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(base.ErrorResponse{
			Message: "Invalid input format: " + err.Error(),
		})
	}

	// Fetch stored PIN hash
	var user *models.User
	user, err := c.UserService.GetUserByID(request.UserID)
	if err != nil {
		log.Printf("Failed to get user by ID: %s. %s", request.UserID, err.Error())
		return ctx.Status(fiber.StatusNotFound).JSON(base.ErrorResponse{
			Message: "User does not exist",
		})
	}

	// Verify PIN
	if !utils.VerifyPIN(user.PIN, request.PIN) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(base.ErrorResponse{
			Message: "Invalid PIN",
		})
	}

	// Generate JWT Token
	token, err := utils.GenerateNewTokens(user.UserID)
	log.Printf("user id: %s", user.UserID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(base.ErrorResponse{
			Message: fmt.Sprintf("Failed to generate token for user: %s. %s", user.UserID, err.Error()),
		})
	}

	return ctx.JSON(fiber.Map{
		"tokens": fiber.Map{
			"access":  token.Access,
			"refresh": token.Refresh,
		},
	})
}

// RenewTokens method for renew access and refresh tokens.
// @Description Renew access and refresh tokens.
// @Summary renew access and refresh tokens
// @Tags Token
// @Accept json
// @Produce json
// @Param refresh_token body models.Renew true "Refresh token"
// @Success 200 {string} status "ok"
// @Security ApiKeyAuth
// @Router /token/renew [post]
func (c *AuthController) RenewTokens(ctx *fiber.Ctx) error {
	// Get now time.
	now := time.Now().Unix()

	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(ctx)
	if err != nil {
		// Return status 500 and JWT parse error.
		return ctx.Status(fiber.StatusInternalServerError).JSON(base.ErrorResponse{
			Message: err.Error(),
		})
	}

	// Set expiration time from JWT data of current user.
	expiresAccessToken := claims.Expires

	// Checking, if now time greather than Access token expiration time.
	if now > expiresAccessToken {
		// Return status 401 and unauthorized error message.
		return ctx.Status(fiber.StatusUnauthorized).JSON(base.ErrorResponse{
			Message: "unauthorized, check expiration time of your token",
		})
	}

	// Create a new renew refresh token struct.
	renew := &models.Renew{}

	// Checking received data from JSON body.
	if err := ctx.BodyParser(renew); err != nil {
		// Return, if JSON data is not correct.
		return ctx.Status(fiber.StatusBadRequest).JSON(base.ErrorResponse{
			Message: err.Error(),
		})
	}

	// Set expiration time from Refresh token of current user.
	expiresRefreshToken, err := utils.ParseRefreshToken(renew.RefreshToken)
	if err != nil {
		// Return status 400 and error message.
		return ctx.Status(fiber.StatusBadRequest).JSON(base.ErrorResponse{
			Message: err.Error(),
		})
	}

	// Checking, if now time greather than Refresh token expiration time.
	if now < expiresRefreshToken {
		// Define user ID.
		userID := claims.UserID

		// Get user by ID.
		_, err = c.UserService.GetUserByID(userID)
		if err != nil {
			log.Printf("Failed to get user by ID: %s. %s", userID, err.Error())
			// Return, if user not found.
			return ctx.Status(fiber.StatusNotFound).JSON(base.ErrorResponse{
				Message: "user not found",
			})
		}

		// Generate JWT Access & Refresh tokens.
		tokens, err := utils.GenerateNewTokens(userID)
		if err != nil {
			// Return status 500 and token generation error.
			return ctx.Status(fiber.StatusInternalServerError).JSON(base.ErrorResponse{
				Message: err.Error(),
			})
		}

		return ctx.JSON(fiber.Map{
			"tokens": fiber.Map{
				"access":  tokens.Access,
				"refresh": tokens.Refresh,
			},
		})
	} else {
		// Return status 401 and unauthorized error message.
		return ctx.Status(fiber.StatusUnauthorized).JSON(base.ErrorResponse{
			Message: "unauthorized, your session was ended earlier",
		})
	}
}
