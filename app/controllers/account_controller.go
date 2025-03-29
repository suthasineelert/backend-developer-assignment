package controllers

import (
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/app/services"

	fiber "github.com/gofiber/fiber/v2"
)

// AccountController handles HTTP requests related to accounts
type AccountController struct {
	accountService services.AccountService
}

// NewAccountController creates a new instance of AccountController
func NewAccountController(service services.AccountService) *AccountController {
	return &AccountController{
		accountService: service,
	}
}

// ListAccounts retrieves all accounts for a user
// @Summary List all accounts for a user
// @Description Get all accounts with details for a specific user
// @Tags accounts
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {array} models.AccountWithDetails
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{user_id}/accounts [get]
func (c *AccountController) ListAccounts(ctx *fiber.Ctx) error {
	userID := ctx.Params("user_id")
	if userID == "" {
		return ErrorResponse(ctx, fiber.StatusBadRequest, "user_id is required")
	}

	accounts, err := c.accountService.GetAccountsWithDetailsByUserID(userID)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(accounts)
}

// GetAccount retrieves a specific account by ID
// @Summary Get account details
// @Description Get detailed information about a specific account
// @Tags accounts
// @Accept json
// @Produce json
// @Param account_id path string true "Account ID"
// @Success 200 {object} models.AccountWithDetails
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /accounts/{account_id} [get]
func (c *AccountController) GetAccount(ctx *fiber.Ctx) error {
	accountID := ctx.Params("account_id")
	if accountID == "" {
		return ErrorResponse(ctx, fiber.StatusBadRequest, "account_id is required")
	}

	account, err := c.accountService.GetAccountWithDetails(accountID)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusNotFound, "account not found")
	}

	return ctx.JSON(account)
}

// CreateAccount creates a new account
// @Summary Create a new account
// @Description Create a new account with all related details
// @Tags accounts
// @Accept json
// @Produce json
// @Param account body models.AccountWithDetails true "Account details"
// @Success 201 {object} models.AccountWithDetails
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /accounts [post]
func (c *AccountController) CreateAccount(ctx *fiber.Ctx) error {
	var accountDetails models.AccountWithDetails
	if err := ctx.BodyParser(&accountDetails); err != nil {
		return ErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	// Validate required fields
	if accountDetails.UserID == "" {
		return ErrorResponse(ctx, fiber.StatusBadRequest, "user_id is required")
	}

	if err := c.accountService.CreateAccountWithDetails(&accountDetails); err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusCreated).JSON(accountDetails)
}

// UpdateAccount updates an existing account
// @Summary Update account details
// @Description Update basic information about an account
// @Tags accounts
// @Accept json
// @Produce json
// @Param account_id path string true "Account ID"
// @Param account body models.Account true "Account details"
// @Success 200 {object} models.Account
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /accounts/{account_id} [put]
func (c *AccountController) UpdateAccount(ctx *fiber.Ctx) error {
	accountID := ctx.Params("account_id")
	if accountID == "" {
		return ErrorResponse(ctx, fiber.StatusBadRequest, "account_id is required")
	}

	// First check if the account exists
	existingAccount, err := c.accountService.GetAccountByID(accountID)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusNotFound, "account not found")
	}

	// Bind the request body to the account model
	var account models.Account
	if err := ctx.BodyParser(&account); err != nil {
		return ErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	// Ensure the account ID in the path matches the one in the body
	account.AccountID = accountID

	// Preserve the user ID from the existing account
	account.UserID = existingAccount.UserID

	if err := c.accountService.UpdateAccount(&account); err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(account)
}

// DepositRequest represents a deposit request
type DepositRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

// Deposit adds funds to an account
// @Summary Deposit funds
// @Description Add funds to an account
// @Tags accounts
// @Accept json
// @Produce json
// @Param account_id path string true "Account ID"
// @Param deposit body DepositRequest true "Deposit details"
// @Success 200 {object} models.AccountWithDetails
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /accounts/{account_id}/deposit [post]
func (c *AccountController) Deposit(ctx *fiber.Ctx) error {
	accountID := ctx.Params("account_id")
	if accountID == "" {
		return ErrorResponse(ctx, fiber.StatusBadRequest, "account_id is required")
	}

	// Get the current account details
	account, err := c.accountService.GetAccountWithDetails(accountID)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusNotFound, "account not found")
	}

	// Parse the deposit request
	var req DepositRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	// Validate amount is greater than 0
	if req.Amount <= 0 {
		return ErrorResponse(ctx, fiber.StatusBadRequest, "amount must be greater than 0")
	}

	// Calculate the new balance
	newBalance := account.Amount + req.Amount

	// Update the account balance
	if err := c.accountService.UpdateAccountBalance(accountID, newBalance); err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	// Get the updated account details
	updatedAccount, err := c.accountService.GetAccountWithDetails(accountID)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to retrieve updated account")
	}

	return ctx.JSON(updatedAccount)
}

// WithdrawRequest represents a withdrawal request
type WithdrawRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

// Withdraw removes funds from an account
// @Summary Withdraw funds
// @Description Remove funds from an account
// @Tags accounts
// @Accept json
// @Produce json
// @Param account_id path string true "Account ID"
// @Param withdraw body WithdrawRequest true "Withdrawal details"
// @Success 200 {object} models.AccountWithDetails
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /accounts/{account_id}/withdraw [post]
func (c *AccountController) Withdraw(ctx *fiber.Ctx) error {
	accountID := ctx.Params("account_id")
	if accountID == "" {
		return ErrorResponse(ctx, fiber.StatusBadRequest, "account_id is required")
	}

	// Get the current account details
	account, err := c.accountService.GetAccountWithDetails(accountID)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusNotFound, "account not found")
	}

	// Parse the withdrawal request
	var req WithdrawRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	// Validate amount is greater than 0
	if req.Amount <= 0 {
		return ErrorResponse(ctx, fiber.StatusBadRequest, "amount must be greater than 0")
	}

	// Check if there are sufficient funds
	if account.Amount < req.Amount {
		return ErrorResponse(ctx, fiber.StatusBadRequest, "insufficient funds")
	}

	// Calculate the new balance
	newBalance := account.Amount - req.Amount

	// Update the account balance
	if err := c.accountService.UpdateAccountBalance(accountID, newBalance); err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}

	// Get the updated account details
	updatedAccount, err := c.accountService.GetAccountWithDetails(accountID)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to retrieve updated account")
	}

	return ctx.JSON(updatedAccount)
}

// TransferRequest represents a transfer request
type TransferRequest struct {
	ToAccountID string  `json:"to_account_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
}

// Transfer moves funds from one account to another
// @Summary Transfer funds
// @Description Move funds from one account to another
// @Tags accounts
// @Accept json
// @Produce json
// @Param account_id path string true "Source Account ID"
// @Param transfer body TransferRequest true "Transfer details"
// @Success 200 {object} models.AccountWithDetails
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /accounts/{account_id}/transfer [post]
func (c *AccountController) Transfer(ctx *fiber.Ctx) error {
	fromAccountID := ctx.Params("account_id")
	if fromAccountID == "" {
		return ErrorResponse(ctx, fiber.StatusBadRequest, "source account_id is required")
	}

	// Parse the transfer request
	var req TransferRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	// Validate amount is greater than 0
	if req.Amount <= 0 {
		return ErrorResponse(ctx, fiber.StatusBadRequest, "amount must be greater than 0")
	}

	// Check if source and destination accounts are different
	if fromAccountID == req.ToAccountID {
		return ErrorResponse(ctx, fiber.StatusBadRequest, "source and destination accounts must be different")
	}

	// Get the source account details
	sourceAccount, err := c.accountService.GetAccountWithDetails(fromAccountID)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusNotFound, "source account not found")
	}

	// Get the destination account details
	destAccount, err := c.accountService.GetAccountWithDetails(req.ToAccountID)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusNotFound, "destination account not found")
	}

	// Check if there are sufficient funds in the source account
	if sourceAccount.Amount < req.Amount {
		return ErrorResponse(ctx, fiber.StatusBadRequest, "insufficient funds")
	}

	// Calculate the new balances
	newSourceBalance := sourceAccount.Amount - req.Amount
	newDestBalance := destAccount.Amount + req.Amount

	// Update the source account balance
	if err := c.accountService.UpdateAccountBalance(fromAccountID, newSourceBalance); err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to update source account")
	}

	// Update the destination account balance
	if err := c.accountService.UpdateAccountBalance(req.ToAccountID, newDestBalance); err != nil {
		// If this fails, we should try to revert the source account update
		// This is a simplified error handling - in a real system, you'd use a transaction
		_ = c.accountService.UpdateAccountBalance(fromAccountID, sourceAccount.Amount)
		return ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to update destination account")
	}

	// Get the updated source account details
	updatedSourceAccount, err := c.accountService.GetAccountWithDetails(fromAccountID)
	if err != nil {
		return ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to retrieve updated source account")
	}

	return ctx.JSON(updatedSourceAccount)
}
