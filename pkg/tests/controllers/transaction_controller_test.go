package controllers_test

import (
	"backend-developer-assignment/app/controllers"
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/pkg/middleware"
	mocks "backend-developer-assignment/pkg/mocks/services"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/suite"
)

type TransactionControllerTestSuite struct {
	suite.Suite
	app                    *fiber.App
	mockTransactionService *mocks.TransactionService
	testToken              string
	testUserID             string
}

func (s *TransactionControllerTestSuite) SetupTest() {
	s.app = fiber.New()
	s.mockTransactionService = new(mocks.TransactionService)
	s.testUserID = "000018b0e1a211ef95a30242ac180002"

	// Create a test JWT token
	claims := jwt.MapClaims{
		"sub": s.testUserID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("your-secret-key"))
	s.testToken = tokenString

	// Setup controller and routes
	transactionController := controllers.NewTransactionController(s.mockTransactionService)
	route := s.app.Group("/transactions", middleware.AuthProtected()...)
	route.Get("/", transactionController.ListTransactions)
}

func (s *TransactionControllerTestSuite) TestListTransactions_Success() {
	// Setup mock data
	mockTransactions := []*models.Transaction{
		{
			TransactionID:   "000018b0e1a211ef95a30242ac180003",
			UserID:          s.testUserID,
			Name:            "Test Transaction 1",
			Amount:          100.00,
			TransactionType: "deposit",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
		{
			TransactionID:   "000018b0e1a211ef95a30242ac180004",
			UserID:          s.testUserID,
			Name:            "Test Transaction 2",
			Amount:          200.00,
			TransactionType: "withdrawal",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}
	mockTotal := 2

	// Setup mock expectations
	s.mockTransactionService.On("GetTransactionsByUserID", s.testUserID, 1).Return(mockTransactions, mockTotal, nil)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
	req.Header.Set("Authorization", "Bearer "+s.testToken)

	// Test the endpoint
	resp, err := s.app.Test(req)
	s.NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)

	// Parse response
	var response struct {
		Transactions []*models.Transaction `json:"transactions"`
		Total        int                   `json:"total"`
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	s.NoError(err)

	// Verify response
	s.Equal(mockTotal, response.Total)
	s.Len(response.Transactions, 2)
	s.Equal(mockTransactions[0].TransactionID, response.Transactions[0].TransactionID)
	s.Equal(mockTransactions[1].TransactionID, response.Transactions[1].TransactionID)

	// Verify mock expectations
	s.mockTransactionService.AssertExpectations(s.T())
}

func (s *TransactionControllerTestSuite) TestListTransactions_WithPagination() {
	// Setup mock data
	mockTransactions := []*models.Transaction{
		{
			TransactionID:   "000018b0e1a211ef95a30242ac180005",
			UserID:          s.testUserID,
			Name:            "Test Transaction 3",
			Amount:          300.00,
			TransactionType: "deposit",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}
	mockTotal := 3 // Total of 3 transactions, but only returning 1 for page 2

	// Setup mock expectations for page 2
	s.mockTransactionService.On("GetTransactionsByUserID", s.testUserID, 2).Return(mockTransactions, mockTotal, nil)

	// Create request with page parameter
	req := httptest.NewRequest(http.MethodGet, "/transactions?page=2", nil)
	req.Header.Set("Authorization", "Bearer "+s.testToken)

	// Test the endpoint
	resp, err := s.app.Test(req)
	s.NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)

	// Parse response
	var response struct {
		Transactions []*models.Transaction `json:"transactions"`
		Total        int                   `json:"total"`
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	s.NoError(err)

	// Verify response
	s.Equal(mockTotal, response.Total)
	s.Len(response.Transactions, 1)
	s.Equal(mockTransactions[0].TransactionID, response.Transactions[0].TransactionID)

	// Verify mock expectations
	s.mockTransactionService.AssertExpectations(s.T())
}

func (s *TransactionControllerTestSuite) TestListTransactions_InvalidPage() {
	// Setup mock expectations for default page 1 when invalid page is provided
	mockTransactions := []*models.Transaction{}
	mockTotal := 0
	s.mockTransactionService.On("GetTransactionsByUserID", s.testUserID, 1).Return(mockTransactions, mockTotal, nil)

	// Create request with invalid page parameter
	req := httptest.NewRequest(http.MethodGet, "/transactions?page=invalid", nil)
	req.Header.Set("Authorization", "Bearer "+s.testToken)

	// Test the endpoint
	resp, err := s.app.Test(req)
	s.NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)

	// Verify mock expectations - should default to page 1
	s.mockTransactionService.AssertExpectations(s.T())
}

func (s *TransactionControllerTestSuite) TestListTransactions_ServiceError() {
	// Setup mock expectations
	s.mockTransactionService.On("GetTransactionsByUserID", s.testUserID, 1).
		Return(nil, 0, errors.New("database error"))

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
	req.Header.Set("Authorization", "Bearer "+s.testToken)

	// Test the endpoint
	resp, err := s.app.Test(req)
	s.NoError(err)
	s.Equal(http.StatusInternalServerError, resp.StatusCode)

	// Parse response
	var response struct {
		Message string `json:"message"`
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	s.NoError(err)

	// Verify response
	s.Equal("Failed to retrieve transactions", response.Message)

	// Verify mock expectations
	s.mockTransactionService.AssertExpectations(s.T())
}

func (s *TransactionControllerTestSuite) TestListTransactions_Unauthorized() {
	// Create request without auth token
	req := httptest.NewRequest(http.MethodGet, "/transactions", nil)

	// Test the endpoint
	resp, err := s.app.Test(req)
	s.NoError(err)
	s.Equal(http.StatusUnauthorized, resp.StatusCode)
}

// Run the test suite
func TestTransactionControllerTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionControllerTestSuite))
}
