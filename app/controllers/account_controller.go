package controllers

import (
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/app/services"
	"backend-developer-assignment/pkg/configs"
	"backend-developer-assignment/pkg/utils"
	"errors"

	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// AccountController handles account-related HTTP requests
type AccountController struct {
	accountService services.AccountService
}

// NewAccountController creates a new AccountController
func NewAccountController(accountService services.AccountService) *AccountController {
	return &AccountController{accountService: accountService}
}

// ListAccounts retrieves all accounts for a user
func (ac *AccountController) ListAccounts(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(string)

	// Get accounts from service
	accounts, err := ac.accountService.GetAccountsWithDetailByUserID(userID)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(accounts)
}

// GetAccount retrieves a single account by ID
func (ac *AccountController) GetAccount(ctx *fiber.Ctx) error {
	// Get account_id from path parameters
	accountID := ctx.Params("id")
	if accountID == "" {
		return ErrorResponse(ctx, fiber.StatusBadRequest, "account_id is required")
	}

	// Get account info
	account, err := ac.accountService.GetAccountWithDetailByID(accountID)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusNotFound, "Account not found")
	}

	return ctx.Status(fiber.StatusOK).JSON(account)
}

// CreateAccount handles account creation
func (ac *AccountController) CreateAccount(ctx *fiber.Ctx) error {
	type createAccountRequest struct {
		Type          string  `json:"type" validate:"required,oneof=saving-account credit-loan goal-driven-saving"`
		Currency      string  `json:"currency" validate:"required,alpha"`
		AccountNumber string  `json:"account_number" validate:"required"`
		Issuer        string  `json:"issuer" validate:"required,alpha"`
		Color         string  `json:"color" validate:"iscolor"`
		IsMainAccount bool    `json:"is_main_account"`
		Amount        float64 `json:"amount"`
	}

	// Parse request body
	var request createAccountRequest

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
		request.Color = configs.DEFAULT_ACCOUNT_COLOR
	}

	userID := ctx.Locals("userID").(string)

	// Create the account with details
	account := &models.AccountWithDetails{
		UserID:        userID,
		Type:          request.Type,
		Currency:      request.Currency,
		AccountNumber: request.AccountNumber,
		Issuer:        request.Issuer,
		Color:         request.Color,
		IsMainAccount: request.IsMainAccount,
		Progress:      0,
		Amount:        request.Amount,
	}

	// Create the account with all its details
	err := ac.accountService.CreateAccountWithDetails(account)
	if err != nil {
		logger.Error("Failed to create account", zap.Error(err))
		return ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	// Retrieve the created account with all details
	createdAccount, err := ac.accountService.GetAccountWithDetailByID(account.AccountID)
	if err != nil {
		logger.Error("Failed to retrieve account details after create", zap.Error(err))
		return ErrorResponse(ctx, fiber.StatusInternalServerError, "Account created but failed to retrieve details")
	}

	return ctx.Status(fiber.StatusCreated).JSON(createdAccount)
}

// UpdateAccount updates an existing account
func (ac *AccountController) UpdateAccount(ctx *fiber.Ctx) error {
	type updateAccountRequest struct {
		Type          string `json:"type" validate:"omitempty,oneof=saving-account credit-loan goal-driven-saving"`
		Currency      string `json:"currency" validate:"omitempty,alpha"`
		AccountNumber string `json:"account_number" validate:"omitempty"`
		Issuer        string `json:"issuer" validate:"omitempty,alpha"`
		Color         string `json:"color" validate:"omitempty,iscolor"`
		Progress      int    `json:"progress" validate:"omitempty,min=0,max=100"`
	}

	// Parse request body
	var request updateAccountRequest

	// Get account_id from path parameters
	accountID := ctx.Params("id")
	if accountID == "" {
		return ErrorResponse(ctx, fiber.StatusBadRequest, "account_id is required")
	}

	// Check if the account exists
	existingAccount, err := ac.accountService.GetAccountByID(accountID)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusNotFound, "Account not found")
	}

	if err := ctx.BodyParser(&request); err != nil {
		return ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request body")
	}

	validate := utils.NewValidator()
	if err := validate.Struct(request); err != nil {
		logger.Info("Validation error", zap.Error(err))
		return ErrorResponse(ctx, fiber.StatusBadRequest, utils.ValidatorErrors(err))
	}

	account := &models.AccountWithDetails{
		AccountID:     existingAccount.AccountID,
		UserID:        existingAccount.UserID,
		Type:          request.Type,
		Currency:      request.Currency,
		AccountNumber: request.AccountNumber,
		Issuer:        request.Issuer,
		Color:         request.Color,
		Progress:      request.Progress,
	}

	if err := ac.accountService.UpdateAccount(account); err != nil {
		logger.Error("Failed to update account", zap.String("account_id", accountID), zap.Error(err))
		return ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to update account: "+err.Error())
	}

	// Retrieve the updated account with all its details
	updatedAccount, err := ac.accountService.GetAccountWithDetailByID(accountID)
	if err != nil {
		logger.Error("Account updated but failed to retrieve details", zap.String("account_id", accountID), zap.Error(err))
		return ErrorResponse(ctx, fiber.StatusInternalServerError, "Account updated but failed to retrieve details")
	}

	return ctx.Status(fiber.StatusOK).JSON(updatedAccount)
}

// SetMainAccount sets an account as the main account
func (ac *AccountController) SetMainAccount(ctx *fiber.Ctx) error {
	// Get account_id from path parameters
	accountID := ctx.Params("id")
	if accountID == "" {
		return ErrorResponse(ctx, fiber.StatusBadRequest, "account_id is required")
	}

	// Check if the account exists
	existingAccount, err := ac.accountService.GetAccountByID(accountID)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusNotFound, "Account not found")
	}

	if err := ac.accountService.SetMainAccount(existingAccount); err != nil {
		logger.Error("Failed to set main account", zap.String("account_id", accountID), zap.Error(err))
		return ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to set main account: "+err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Main account set successfully"})
}

// Withdraw handles withdrawing money from an account
func (ac *AccountController) Withdraw(ctx *fiber.Ctx) error {
	type withdrawRequest struct {
		Amount float64 `json:"amount" validate:"required,gt=0"`
	}

	// Get account_id from path parameters
	accountID := ctx.Params("id")
	if accountID == "" {
		return ErrorResponse(ctx, fiber.StatusBadRequest, "account_id is required")
	}

	var request withdrawRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request body")
	}

	validate := utils.NewValidator()
	if err := validate.Struct(request); err != nil {
		logger.Info("Validation error", zap.Error(err))
		return ErrorResponse(ctx, fiber.StatusBadRequest, utils.ValidatorErrors(err))
	}

	// Get the account to ensure it exists
	account, err := ac.accountService.GetAccountWithDetailByID(accountID)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusNotFound, "Account not found")
	}

	// Check if account has sufficient funds
	if account.Amount < request.Amount {
		return ErrorResponse(ctx, fiber.StatusBadRequest, "Insufficient funds")
	}

	// Use a transaction to handle the withdrawal with proper locking
	updatedBalance, err := ac.accountService.WithdrawFromAccount(accountID, request.Amount)
	if err != nil {
		if err.Error() == "insufficient funds" {
			return ErrorResponse(ctx, fiber.StatusBadRequest, "Insufficient funds")
		}
		logger.Error("Failed to withdraw from account", zap.String("account_id", accountID), zap.Float64("amount", request.Amount), zap.Error(err))
		return ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to process withdrawal")
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Withdrawal successful",
		"amount":  request.Amount,
		"balance": updatedBalance,
	})
}

// Deposit handles depositing money to an account
func (ac *AccountController) Deposit(ctx *fiber.Ctx) error {
	type depositRequest struct {
		Amount float64 `json:"amount" validate:"required,gt=0"`
	}

	// Get account_id from path parameters
	accountID := ctx.Params("id")
	if accountID == "" {
		return ErrorResponse(ctx, fiber.StatusBadRequest, "account_id is required")
	}

	var request depositRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request body")
	}

	validate := utils.NewValidator()
	if err := validate.Struct(request); err != nil {
		logger.Info("Validation error", zap.Error(err))
		return ErrorResponse(ctx, fiber.StatusBadRequest, utils.ValidatorErrors(err))
	}

	// Use a transaction to handle the deposit with proper locking
	updatedBalance, err := ac.accountService.DepositToAccount(accountID, request.Amount)
	if err != nil {
		logger.Error("Failed to deposit to account", zap.String("account_id", accountID), zap.Float64("amount", request.Amount), zap.Error(err))
		return ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to process deposit")
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Deposit successful",
		"amount":  request.Amount,
		"balance": updatedBalance,
	})
}

// Transfer handles transferring money between accounts
func (ac *AccountController) Transfer(ctx *fiber.Ctx) error {
	type transferRequest struct {
		FromAccountID string  `json:"from_account_id" validate:"required"`
		ToAccountID   string  `json:"to_account_id" validate:"required"`
		Amount        float64 `json:"amount" validate:"required,gt=0"`
	}

	var request transferRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid request body")
	}

	validate := utils.NewValidator()
	if err := validate.Struct(request); err != nil {
		logger.Info("Validation error", zap.Error(err))
		return ErrorResponse(ctx, fiber.StatusBadRequest, utils.ValidatorErrors(err))
	}

	// Use a transaction to handle the transfer with proper locking
	result, err := ac.accountService.TransferBetweenAccounts(
		request.FromAccountID,
		request.ToAccountID,
		request.Amount,
	)

	if err != nil {
		if errors.Is(err, services.ErrInsufficientFunds) {
			return ErrorResponse(ctx, fiber.StatusBadRequest, "Insufficient funds in source account")
		}
		logger.Error("Failed to transfer between accounts",
			zap.String("from_account_id", request.FromAccountID),
			zap.String("to_account_id", request.ToAccountID),
			zap.Float64("amount", request.Amount),
			zap.Error(err))
		return ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to process transfer")
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":             "Transfer successful",
		"amount":              request.Amount,
		"from_account":        request.FromAccountID,
		"to_account":          request.ToAccountID,
		"source_balance":      result.SourceBalance,
		"destination_balance": result.DestinationBalance,
	})
}
