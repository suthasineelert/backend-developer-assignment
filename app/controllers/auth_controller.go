package controllers

import (
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/app/services"
	"backend-developer-assignment/pkg/base"
	"backend-developer-assignment/pkg/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
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

// VerifyPin verifies a PIN against a stored hash and get a JWT token
// @Description Verify user PIN and return JWT token
// @Summary Verify PIN and get JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body controllers.verifyPinRequest true "PIN verification request"
// @Success 200 {object} object{tokens=object{access=string,refresh=string}} "JWT tokens"
// @Failure 400 {object} base.ErrorResponse "Invalid input format"
// @Failure 401 {object} base.ErrorResponse "Invalid PIN"
// @Failure 404 {object} base.ErrorResponse "User does not exist"
// @Failure 500 {object} base.ErrorResponse "Failed to generate token"
// @Router /api/auth/verify-pin [post]
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
	token, err := utils.GenerateNewTokens(user.UserID.String())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(base.ErrorResponse{
			Message: fmt.Sprintf("Failed to generate token for user: %s. %s", user.UserID.String(), err.Error()),
		})
	}

	return ctx.JSON(fiber.Map{
		"tokens": fiber.Map{
			"access":  token.Access,
			"refresh": token.Refresh,
		},
	})
}
