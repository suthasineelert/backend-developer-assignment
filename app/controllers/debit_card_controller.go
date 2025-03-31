package controllers

import (
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/app/services"
	"backend-developer-assignment/pkg/configs"
	"backend-developer-assignment/pkg/utils"

	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// DebitCardController handles HTTP requests related to debit cards
type DebitCardController struct {
	debitCardService services.DebitCardService
}

// NewDebitCardController creates a new instance of DebitCardController
func NewDebitCardController(service services.DebitCardService) *DebitCardController {
	return &DebitCardController{
		debitCardService: service,
	}
}

// ListDebitCards returns all debit cards for a user
func (c *DebitCardController) ListDebitCards(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(string)

	// Get cards from service
	cards, err := c.debitCardService.GetCardWithDetailByUserID(userID)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(cards)
}

// GetDebitCard returns a specific debit card by ID
func (c *DebitCardController) GetDebitCard(ctx *fiber.Ctx) error {
	// Get card_id from path parameters
	cardID := ctx.Params("id")
	if cardID == "" {
		return ErrorResponse(ctx, fiber.StatusBadRequest, "card_id is required")
	}

	// Get card from service
	card, err := c.debitCardService.GetCardWithDetailByID(cardID)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusNotFound, "Debit card not found")
	}

	return ctx.Status(fiber.StatusOK).JSON(card)
}

// CreateDebitCard creates a new debit card with all its details
func (c *DebitCardController) CreateDebitCard(ctx *fiber.Ctx) error {
	type createDebitCardRequest struct {
		Name        string `json:"name" validate:"required,alpha"`
		Issuer      string `json:"issuer" validate:"required,alpha"`
		Color       string `json:"color" validate:"iscolor"`
		BorderColor string `json:"border_color" validate:"iscolor"`
	}
	// Parse request body
	var request createDebitCardRequest

	if err := ctx.BodyParser(&request); err != nil {
		return ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request body")
	}

	validate := utils.NewValidator()
	if err := validate.Struct(request); err != nil {
		logger.Info("Validation error", zap.Error(err))
		return ErrorResponse(ctx, fiber.StatusBadRequest, utils.ValidatorErrors(err))
	}

	if request.Color == "" {
		// set default color
		request.Color = configs.DEFAULT_DEBIT_CARD_COLOR
	}
	if request.BorderColor == "" {
		// set default border color
		request.BorderColor = configs.DEFAULT_DEBIT_CARD_BORDER_COLOR
	}

	userID := ctx.Locals("userID").(string)

	// Create the main card
	card := &models.DebitCardWithDetails{
		UserID:      userID,
		Name:        request.Name,
		Issuer:      request.Issuer,
		Color:       request.Color,
		BorderColor: request.BorderColor,
	}

	// Create the card with all its details
	err := c.debitCardService.CreateCardWithDetails(card)
	if err != nil {
		logger.Error("Failed to create card", zap.Error(err))
		return ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	// Retrieve the created card with all its details
	createdCard, err := c.debitCardService.GetCardWithDetailByID(card.CardID)
	if err != nil {
		logger.Error("Failed to retrieve card details after create", zap.Error(err))
		return ErrorResponse(ctx, fiber.StatusInternalServerError, "Card created but failed to retrieve details")
	}

	return ctx.Status(fiber.StatusCreated).JSON(createdCard)
}

// UpdateDebitCard updates an existing debit card
func (c *DebitCardController) UpdateDebitCard(ctx *fiber.Ctx) error {
	type updateDebitCardRequest struct {
		Name        string `json:"name" validate:"alpha"`
		Color       string `json:"color" validate:"iscolor"`
		BorderColor string `json:"border_color" validate:"iscolor"`
	}
	// Parse request body
	var request updateDebitCardRequest

	// Get card_id from path parameters
	cardID := ctx.Params("id")
	if cardID == "" {
		return ErrorResponse(ctx, fiber.StatusBadRequest, "card_id is required")
	}

	// Check if the card exists
	existingCard, err := c.debitCardService.GetCardByID(cardID)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusNotFound, "Debit card not found")
	}

	if err := ctx.BodyParser(&request); err != nil {
		return ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request body")
	}

	validate := utils.NewValidator()
	if err := validate.Struct(request); err != nil {
		logger.Info("Validation error", zap.Error(err))
		return ErrorResponse(ctx, fiber.StatusBadRequest, utils.ValidatorErrors(err))
	}

	if err := c.debitCardService.UpdateCard(existingCard, request.Name, request.Color, request.BorderColor); err != nil {
		logger.Error("Failed to update card", zap.String("card_id", cardID), zap.Error(err))
		return ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to update card: "+err.Error())
	}

	// Retrieve the updated card with all its details
	updatedCard, err := c.debitCardService.GetCardWithDetailByID(cardID)
	if err != nil {
		logger.Error("Card updated but failed to retrieve details", zap.String("card_id", cardID), zap.Error(err))
		return ErrorResponse(ctx, fiber.StatusInternalServerError, "Card updated but failed to retrieve details")
	}

	return ctx.Status(fiber.StatusOK).JSON(updatedCard)
}

// DeleteDebitCard deletes a debit card
func (c *DebitCardController) DeleteDebitCard(ctx *fiber.Ctx) error {
	// Get card_id from path parameters
	cardID := ctx.Params("id")
	if cardID == "" {
		return ErrorResponse(ctx, fiber.StatusBadRequest, "card_id is required")
	}

	// Check if the card exists
	_, err := c.debitCardService.GetCardByID(cardID)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusNotFound, "Debit card not found")
	}

	// Delete the card (soft delete)
	if err := c.debitCardService.DeleteCard(cardID); err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to delete card: "+err.Error())
	}

	return ctx.Status(fiber.StatusNoContent).Send(nil)
}
