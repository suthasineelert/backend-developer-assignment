package controllers_test

import (
	"backend-developer-assignment/app/controllers"
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/app/services"
	mocks "backend-developer-assignment/pkg/mocks/services"
	"backend-developer-assignment/pkg/types"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// AccountControllerTestSuite defines the test suite
type AccountControllerTestSuite struct {
	suite.Suite
	app             *fiber.App
	accountService  *mocks.AccountService
	controller      *controllers.AccountController
	testUserID      string
	testAccountID   string
	testAccountData *models.AccountWithDetails
}

// SetupTest runs before each test
func (s *AccountControllerTestSuite) SetupTest() {
	s.app = fiber.New()
	s.accountService = new(mocks.AccountService)
	s.controller = controllers.NewAccountController(s.accountService)
	s.testUserID = "test-user-id"
	s.testAccountID = "test-account-id"

	now := time.Now()
	s.testAccountData = &models.AccountWithDetails{
		AccountID:     s.testAccountID,
		UserID:        s.testUserID,
		Type:          "saving-account",
		Currency:      "USD",
		AccountNumber: "123456789",
		Issuer:        "TestBank",
		Color:         "#FF0000",
		Progress:      0,
		Amount:        1000.0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Setup routes
	s.app.Get("/accounts", func(c *fiber.Ctx) error {
		c.Locals("userID", s.testUserID)
		return s.controller.ListAccounts(c)
	})

	s.app.Get("/accounts/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", s.testUserID)
		return s.controller.GetAccount(c)
	})

	s.app.Post("/accounts", func(c *fiber.Ctx) error {
		c.Locals("userID", s.testUserID)
		return s.controller.CreateAccount(c)
	})

	s.app.Put("/accounts/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", s.testUserID)
		return s.controller.UpdateAccount(c)
	})

	s.app.Put("/accounts/:id/main", func(c *fiber.Ctx) error {
		c.Locals("userID", s.testUserID)
		return s.controller.SetMainAccount(c)
	})

	s.app.Post("/accounts/:id/withdraw", func(c *fiber.Ctx) error {
		c.Locals("userID", s.testUserID)
		return s.controller.Withdraw(c)
	})

	s.app.Post("/accounts/:id/deposit", func(c *fiber.Ctx) error {
		c.Locals("userID", s.testUserID)
		return s.controller.Deposit(c)
	})

	s.app.Post("/accounts/transfer", func(c *fiber.Ctx) error {
		c.Locals("userID", s.testUserID)
		return s.controller.Transfer(c)
	})
}

// TestListAccounts tests the ListAccounts controller method
func (s *AccountControllerTestSuite) TestListAccounts() {
	// Test case: successful retrieval of accounts
	s.accountService.On("GetAccountsWithDetailByUserID", s.testUserID).Return([]*models.AccountWithDetails{s.testAccountData}, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/accounts", http.NoBody)
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

	var accounts []*models.AccountWithDetails
	err = json.NewDecoder(resp.Body).Decode(&accounts)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), accounts, 1)
	assert.Equal(s.T(), s.testAccountID, accounts[0].AccountID)

	// Test case: error retrieving accounts
	s.accountService.On("GetAccountsWithDetailByUserID", s.testUserID).Return(nil, errors.New("database error")).Once()

	req = httptest.NewRequest(http.MethodGet, "/accounts", http.NoBody)
	resp, err = s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusInternalServerError, resp.StatusCode)

	s.accountService.AssertExpectations(s.T())
}

// TestGetAccount tests the GetAccount controller method
func (s *AccountControllerTestSuite) TestGetAccount() {
	// Test case: successful retrieval of account
	s.accountService.On("GetAccountWithDetailByID", s.testAccountID).Return(s.testAccountData, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/accounts/"+s.testAccountID, http.NoBody)
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

	var account models.AccountWithDetails
	err = json.NewDecoder(resp.Body).Decode(&account)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), s.testAccountID, account.AccountID)

	// Test case: account not found
	s.accountService.On("GetAccountWithDetailByID", "nonexistent-id").Return(nil, errors.New("account not found")).Once()

	req = httptest.NewRequest(http.MethodGet, "/accounts/nonexistent-id", http.NoBody)
	resp, err = s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusNotFound, resp.StatusCode)
}

// TestCreateAccount tests the CreateAccount controller method
func (s *AccountControllerTestSuite) TestCreateAccount() {
	// Test case: successful account creation
	createRequest := map[string]interface{}{
		"type":            "saving-account",
		"currency":        "USD",
		"account_number":  "123456789",
		"issuer":          "TestBank",
		"color":           "#FF0000",
		"is_main_account": false,
		"amount":          1000.0,
	}

	requestBody, _ := json.Marshal(createRequest)

	s.accountService.On("CreateAccountWithDetails", mock.AnythingOfType("*models.AccountWithDetails")).Return(nil).Once()
	s.accountService.On("GetAccountWithDetailByID", mock.AnythingOfType("string")).Return(s.testAccountData, nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusCreated, resp.StatusCode)

	var account models.AccountWithDetails
	err = json.NewDecoder(resp.Body).Decode(&account)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), s.testAccountID, account.AccountID)

	// Test case: validation error
	invalidRequest := map[string]interface{}{
		"type":           "invalid-type", // Invalid account type
		"currency":       "USD",
		"account_number": "123456789",
		"issuer":         "TestBank",
	}

	requestBody, _ = json.Marshal(invalidRequest)

	req = httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)

	// Test case: service error
	s.accountService.On("CreateAccountWithDetails", mock.AnythingOfType("*models.AccountWithDetails")).Return(errors.New("database error")).Once()

	requestBody, _ = json.Marshal(createRequest)

	req = httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusInternalServerError, resp.StatusCode)

	s.accountService.AssertExpectations(s.T())
}

// TestUpdateAccount tests the UpdateAccount controller method
func (s *AccountControllerTestSuite) TestUpdateAccount() {
	// Test case: successful account update
	updateRequest := map[string]interface{}{
		"type":           "credit-loan",
		"currency":       "EUR",
		"account_number": "987654321",
		"issuer":         "NewBank",
		"color":          "#00FF00",
		"progress":       50,
	}

	requestBody, _ := json.Marshal(updateRequest)

	s.accountService.On("GetAccountByID", s.testAccountID).Return(&models.Account{
		AccountID: s.testAccountID,
		UserID:    s.testUserID,
	}, nil).Once()

	s.accountService.On("UpdateAccount", mock.AnythingOfType("*models.AccountWithDetails")).Return(nil).Once()

	updatedAccount := *s.testAccountData
	updatedAccount.Type = "credit-loan"
	updatedAccount.Currency = "EUR"
	updatedAccount.AccountNumber = "987654321"
	updatedAccount.Issuer = "NewBank"
	updatedAccount.Color = "#00FF00"
	updatedAccount.Progress = 50

	s.accountService.On("GetAccountWithDetailByID", s.testAccountID).Return(&updatedAccount, nil).Once()

	req := httptest.NewRequest(http.MethodPut, "/accounts/"+s.testAccountID, bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

	var account models.AccountWithDetails
	err = json.NewDecoder(resp.Body).Decode(&account)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "credit-loan", account.Type)
	assert.Equal(s.T(), "EUR", account.Currency)

	// Test case: account not found
	s.accountService.On("GetAccountByID", "nonexistent-id").Return(nil, errors.New("account not found")).Once()

	req = httptest.NewRequest(http.MethodPut, "/accounts/nonexistent-id", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusNotFound, resp.StatusCode)

	s.accountService.AssertExpectations(s.T())
}

// TestSetMainAccount tests the SetMainAccount controller method
func (s *AccountControllerTestSuite) TestSetMainAccount() {
	// Test case: successful setting of main account
	s.accountService.On("GetAccountByID", s.testAccountID).Return(&models.Account{
		AccountID: s.testAccountID,
		UserID:    s.testUserID,
	}, nil).Once()

	s.accountService.On("SetMainAccount", mock.AnythingOfType("*models.Account")).Return(nil).Once()

	req := httptest.NewRequest(http.MethodPut, "/accounts/"+s.testAccountID+"/main", http.NoBody)
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

	// Test case: account not found
	s.accountService.On("GetAccountByID", "nonexistent-id").Return(nil, errors.New("account not found")).Once()

	req = httptest.NewRequest(http.MethodPut, "/accounts/nonexistent-id/main", http.NoBody)
	resp, err = s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusNotFound, resp.StatusCode)

	// Test case: service error
	s.accountService.On("GetAccountByID", "error-id").Return(&models.Account{
		AccountID: "error-id",
		UserID:    s.testUserID,
	}, nil).Once()

	s.accountService.On("SetMainAccount", mock.AnythingOfType("*models.Account")).Return(errors.New("database error")).Once()

	req = httptest.NewRequest(http.MethodPut, "/accounts/error-id/main", http.NoBody)
	resp, err = s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusInternalServerError, resp.StatusCode)

	s.accountService.AssertExpectations(s.T())
}

// TestWithdraw tests the Withdraw controller method
func (s *AccountControllerTestSuite) TestWithdraw() {
	// Test case: successful withdrawal
	withdrawRequest := map[string]interface{}{
		"amount": 500.0,
	}

	requestBody, _ := json.Marshal(withdrawRequest)

	s.accountService.On("GetAccountWithDetailByID", s.testAccountID).Return(s.testAccountData, nil).Once()
	s.accountService.On("WithdrawFromAccount", s.testAccountID, 500.0).Return(500.0, nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/accounts/"+s.testAccountID+"/withdraw", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "Withdrawal successful", response["message"])
	assert.Equal(s.T(), 500.0, response["amount"])
	assert.Equal(s.T(), 500.0, response["balance"])

	// Test case: insufficient funds
	s.accountService.On("GetAccountWithDetailByID", "low-balance-id").Return(&models.AccountWithDetails{
		AccountID: "low-balance-id",
		UserID:    s.testUserID,
		Amount:    100.0,
	}, nil).Once()

	req = httptest.NewRequest(http.MethodPost, "/accounts/low-balance-id/withdraw", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)

	// Test case: service error
	s.accountService.On("GetAccountWithDetailByID", "error-id").Return(s.testAccountData, nil).Once()
	s.accountService.On("WithdrawFromAccount", "error-id", 500.0).Return(0.0, errors.New("database error")).Once()

	req = httptest.NewRequest(http.MethodPost, "/accounts/error-id/withdraw", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusInternalServerError, resp.StatusCode)

	s.accountService.AssertExpectations(s.T())
}

// TestDeposit tests the Deposit controller method
func (s *AccountControllerTestSuite) TestDeposit() {
	// Test case: successful deposit
	depositRequest := map[string]interface{}{
		"amount": 500.0,
	}

	requestBody, _ := json.Marshal(depositRequest)

	s.accountService.On("DepositToAccount", s.testAccountID, 500.0).Return(1500.0, nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/accounts/"+s.testAccountID+"/deposit", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "Deposit successful", response["message"])
	assert.Equal(s.T(), 500.0, response["amount"])
	assert.Equal(s.T(), 1500.0, response["balance"])

	// Test case: service error
	s.accountService.On("DepositToAccount", "error-id", 500.0).Return(0.0, errors.New("database error")).Once()

	req = httptest.NewRequest(http.MethodPost, "/accounts/error-id/deposit", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusInternalServerError, resp.StatusCode)

	s.accountService.AssertExpectations(s.T())
}

// TestTransfer tests the Transfer controller method
func (s *AccountControllerTestSuite) TestTransfer() {
	// Test case: successful transfer
	transferRequest := map[string]interface{}{
		"from_account_id": "source-account-id",
		"to_account_id":   "dest-account-id",
		"amount":          500.0,
	}

	requestBody, _ := json.Marshal(transferRequest)

	transferResult := &types.TransferResult{
		SourceBalance:      500.0,
		DestinationBalance: 1500.0,
	}

	s.accountService.On("TransferBetweenAccounts", "source-account-id", "dest-account-id", 500.0).Return(transferResult, nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/accounts/transfer", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "Transfer successful", response["message"])
	assert.Equal(s.T(), 500.0, response["amount"])
	assert.Equal(s.T(), "source-account-id", response["from_account"])
	assert.Equal(s.T(), "dest-account-id", response["to_account"])
	assert.Equal(s.T(), 500.0, response["source_balance"])
	assert.Equal(s.T(), 1500.0, response["destination_balance"])

	// Test case: insufficient funds
	s.accountService.On("TransferBetweenAccounts", "low-balance-id", "dest-account-id", 500.0).Return(nil, services.ErrInsufficientFunds).Once()

	transferRequest["from_account_id"] = "low-balance-id"
	requestBody, _ = json.Marshal(transferRequest)

	req = httptest.NewRequest(http.MethodPost, "/accounts/transfer", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)

	// Test case: service error
	s.accountService.On("TransferBetweenAccounts", "error-id", "dest-account-id", 500.0).Return(nil, errors.New("database error")).Once()

	transferRequest["from_account_id"] = "error-id"
	requestBody, _ = json.Marshal(transferRequest)

	req = httptest.NewRequest(http.MethodPost, "/accounts/transfer", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusInternalServerError, resp.StatusCode)

	s.accountService.AssertExpectations(s.T())
}

// TestAccountControllerSuite runs the test suite
func TestAccountControllerSuite(t *testing.T) {
	suite.Run(t, new(AccountControllerTestSuite))
}
