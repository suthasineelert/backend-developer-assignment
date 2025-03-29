package controllers

import (
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/app/services"
	"backend-developer-assignment/pkg/base"
	"backend-developer-assignment/pkg/utils"

	fiber "github.com/gofiber/fiber/v2"
)

// UserController holds the services related to users.
type UserController struct {
	UserService services.UserService
}

// NewUserController creates a new UserController.
func NewUserController(userService services.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

// Greeting get user's greeting message
func (c *UserController) Greeting(ctx *fiber.Ctx) error {
	tokenData, err := utils.ExtractTokenMetadata(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(base.ErrorResponse{
			Message: "Cannot get user info: " + err.Error(),
		})
	}

	var greeting *models.UserGreeting

	greeting, err = c.UserService.GetUserGreetingByID(tokenData.UserID.String())
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(base.ErrorResponse{
			Message: "User greeting not found",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": greeting.Greeting,
	})
}
