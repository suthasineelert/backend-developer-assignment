package controllers

import (
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/app/services"
	"backend-developer-assignment/pkg/base"

	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
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

// GetUserGreeting get user's greeting message
// @Summary Get user's greeting message
// @Description Retrieves the greeting message for the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} controllers.GetUserGreeting.getUserGreetingResponse "Returns the greeting message"
// @Failure 401 {object} base.ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure 404 {object} base.ErrorResponse "User greeting not found"
// @Router /user/greeting [get]
func (c *UserController) GetUserGreeting(ctx *fiber.Ctx) error {
	type getUserGreetingResponse struct {
		Message string `json:"message"`
	}

	userID := ctx.Locals("userID").(string)

	logger.Info("Get user greeting", zap.String("user_id", userID))

	var greeting *models.UserGreeting

	greeting, err := c.UserService.GetUserGreetingByID(userID)
	if err != nil {
		logger.Info("Failed to get user greeting", zap.String("user_id", userID), zap.Error(err))
		return ctx.Status(fiber.StatusNotFound).JSON(base.ErrorResponse{
			Message: "User greeting not found",
		})
	}

	return ctx.JSON(getUserGreetingResponse{
		Message: greeting.Greeting,
	})
}

// UpdateUserGreeting update user's greeting message
// @Summary Update user's greeting message
// @Description Update the greeting message of the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body controllers.UpdateUserGreeting.updateUserGreetingRequest true "Request body"
// @Success 200 {object} controllers.UpdateUserGreeting.updateUserGreetingResponse "Returns the updated greeting message"
// @Failure 400 {object} base.ErrorResponse "Bad request - Invalid input format"
// @Failure 401 {object} base.ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure 404 {object} base.ErrorResponse "User not found"
// @Failure 500 {object} base.ErrorResponse "Internal server error"
// @Router /user/greeting [put]
func (c *UserController) UpdateUserGreeting(ctx *fiber.Ctx) error {
	type updateUserGreetingRequest struct {
		Message string `json:"message"`
	}
	type updateUserGreetingResponse struct {
		Message string `json:"message"`
	}

	userID := ctx.Locals("userID").(string)

	logger.Info("Update user greeting", zap.String("user_id", userID))

	var request updateUserGreetingRequest
	if err := ctx.BodyParser(&request); err != nil {
		logger.Info("Failed to parse request body", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(base.ErrorResponse{
			Message: "Invalid input format: " + err.Error(),
		})
	}

	var greeting models.UserGreeting
	greeting.UserID = userID
	greeting.Greeting = request.Message

	err := c.UserService.UpdateUserGreeting(&greeting)
	if err != nil {
		logger.Error("Failed to update user greeting", zap.Error(err))
		return ctx.Status(fiber.StatusNotFound).JSON(base.ErrorResponse{
			Message: "Fail to update user greeting",
		})
	}

	return ctx.JSON(updateUserGreetingResponse{
		Message: greeting.Greeting,
	})
}

// GetUser get user's information
// @Summary Get user's information
// @Description Retrieves the information of the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.User "Returns the user information"
// @Failure 401 {object} base.ErrorResponse "Unauthorized - Invalid or missing token"
// @Failure 404 {object} base.ErrorResponse "User not found"
// @Router /user [get]
func (c *UserController) GetUser(ctx *fiber.Ctx) error {
	var user *models.User

	userID := ctx.Locals("userID").(string)

	user, err := c.UserService.GetUserByID(userID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(base.ErrorResponse{
			Message: "User not found",
		})
	}

	return ctx.JSON(user)
}

func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
	type updateUserRequest struct {
		Name string `json:"name"`
	}

	userID := ctx.Locals("userID").(string)

	var request updateUserRequest
	if err := ctx.BodyParser(&request); err != nil {
		logger.Info("Failed to parse request body", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(base.ErrorResponse{
			Message: "Invalid input format: " + err.Error(),
		})
	}

	logger.Info("Update user", zap.String("user_id", userID), zap.String("name", request.Name))

	var user models.User
	user.UserID = userID
	user.Name = request.Name

	err := c.UserService.UpdateUser(&user)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(base.ErrorResponse{
			Message: "User not found",
		})
	}

	return ctx.JSON(user)
}
