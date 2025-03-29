package controllers

import (
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/app/services"
	"backend-developer-assignment/pkg/base"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// TransactionController holds the services related to transactions.
type TransactionController struct {
	TransactionService services.TransactionService
}

// NewTransactionController creates a new TransactionController.
func NewTransactionController(transactionService services.TransactionService) *TransactionController {
	return &TransactionController{
		TransactionService: transactionService,
	}
}

// ListTransactions retrieves all transactions for a user.
// @Description Retrieves all transactions for a user.
// @Summary List transactions
// @Tags Transactions
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Success 200 {object} object{transactions=[]models.Transaction,total=int} "List of transactions"
// @Failure 400 {object} base.ErrorResponse "Invalid input format"
// @Failure 401 {object} base.ErrorResponse "Unauthorized"
// @Failure 500 {object} base.ErrorResponse "Failed to retrieve transactions"
// @Router /transactions [get]
func (c *TransactionController) ListTransactions(ctx *fiber.Ctx) error {
	type listTransactionResponse struct {
		Transactions []*models.Transaction `json:"transactions"`
		Total        int                   `json:"total"`
	}
	pageQuery := ctx.Query("page", "1")
	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		logger.Warn("Cannot parse page query to int, default to 1", zap.String("page", pageQuery), zap.Error(err))
		page = 1
	}

	userID := ctx.Locals("userID").(string)

	transactions, total, err := c.TransactionService.GetTransactionsByUserID(userID, page)
	if err != nil {
		logger.Error("Failed to get transactions", zap.String("user_id", userID), zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).JSON(base.ErrorResponse{
			Message: "Failed to retrieve transactions",
		})
	}

	return ctx.JSON(listTransactionResponse{
		Transactions: transactions,
		Total:        total,
	})
}
