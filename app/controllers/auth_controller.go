package controllers

import (
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/app/services"
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
// @Param request body object{user_id=string,pin=string} true "PIN verification request"
// @Success 200 {object} object{token=string} "JWT token"
// @Failure 400 {object} object{error=boolean,msg=string} "Invalid input format"
// @Failure 404 {object} object{error=boolean,msg=string} "User does not exist"
// @Failure 401 {object} object{error=boolean,msg=string} "Invalid PIN"
// @Failure 500 {object} object{error=boolean,msg=string} "Failed to generate token for user"
// @Router /api/auth/verify-pin [post]
func (c *AuthController) VerifyPin(ctx *fiber.Ctx) error {
	type verifyPinRequest struct {
		UserID string `json:"user_id"`
		PIN    string `json:"pin"`
	}
	var request verifyPinRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "msg": "Invalid input format: " + err.Error()})
	}

	// Fetch stored PIN hash
	var user *models.User
	user, err := c.UserService.GetUserByID(request.UserID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": true, "msg": "User does not exist"})
	}

	// Verify PIN
	if !utils.VerifyPIN(user.PIN, request.PIN) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": true, "msg": "Invalid PIN"})
	}

	// Generate JWT Token
	token, err := utils.GenerateNewTokens(user.UserID.String())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "msg": fmt.Sprintf("Failed to generate token for user: %s. %s", user.UserID.String(), err.Error())})
	}

	return ctx.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"tokens": fiber.Map{
			"access":  token.Access,
			"refresh": token.Refresh,
		},
	})
}
