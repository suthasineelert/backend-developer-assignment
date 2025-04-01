package services_test

import (
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/app/services"
	mockCache "backend-developer-assignment/pkg/mocks/cache"
	mocks "backend-developer-assignment/pkg/mocks/repositories"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// TransactionServiceTestSuite defines the test suite
type TransactionServiceTestSuite struct {
	suite.Suite
	transactionRepository *mocks.TransactionRepository
	redisClient           *mockCache.RedisClient
	service               services.TransactionService
}

// SetupTest runs before each test
func (s *TransactionServiceTestSuite) SetupTest() {
	s.transactionRepository = new(mocks.TransactionRepository)
	s.redisClient = new(mockCache.RedisClient)
	s.service = services.NewTransactionService(s.transactionRepository, s.redisClient)
}

// TestGetTransactionByID tests the GetTransactionByID function
func (s *TransactionServiceTestSuite) TestGetTransactionByID() {
	now := time.Now()
	testCases := []struct {
		name                string
		transactionID       string
		mockTransaction     *models.Transaction
		mockError           error
		expectedTransaction *models.Transaction
		expectedError       error
	}{
		{
			name:          "Success - Valid Transaction",
			transactionID: "000018b0e1a211ef95a30242ac180002",
			mockTransaction: &models.Transaction{
				TransactionID:   "000018b0e1a211ef95a30242ac180002",
				UserID:          "000018b0e1a211ef95a30242ac180003",
				Name:            "Test Transaction",
				Image:           "transaction.jpg",
				IsBank:          true,
				Amount:          100.50,
				TransactionType: "deposit",
				BaseModel: &models.BaseModel{
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			mockError: nil,
			expectedTransaction: &models.Transaction{
				TransactionID:   "000018b0e1a211ef95a30242ac180002",
				UserID:          "000018b0e1a211ef95a30242ac180003",
				Name:            "Test Transaction",
				Image:           "transaction.jpg",
				IsBank:          true,
				Amount:          100.50,
				TransactionType: "deposit",
				BaseModel: &models.BaseModel{
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			expectedError: nil,
		},
		{
			name:                "Failure - Transaction Not Found",
			transactionID:       "000018b0e1a211ef95a30242ac180001",
			mockTransaction:     nil,
			mockError:           errors.New("transaction not found"),
			expectedTransaction: nil,
			expectedError:       errors.New("transaction not found"),
		},
		{
			name:                "Failure - Database Error",
			transactionID:       "invalid-transaction-id",
			mockTransaction:     nil,
			mockError:           errors.New("database connection failed"),
			expectedTransaction: nil,
			expectedError:       errors.New("database connection failed"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			ctx := context.Background()
			cacheKey := fmt.Sprintf("transaction:%s", tc.transactionID)

			// Mock Redis cache miss
			s.redisClient.On("Get", ctx, cacheKey).Return("", errors.New("cache miss")).Once()

			// Mock the repository method
			s.transactionRepository.On("GetByID", tc.transactionID).Return(tc.mockTransaction, tc.mockError).Once()

			if tc.mockError == nil && tc.mockTransaction != nil {
				// Mock Redis cache set (only for successful DB retrieval)
				s.redisClient.On("Set", ctx, cacheKey, mock.Anything, mock.Anything).Return(nil).Once()
			}

			// Call the service method
			transaction, err := s.service.GetTransactionByID(tc.transactionID)

			// Assert results
			if tc.expectedError != nil {
				assert.Error(s.T(), err)
				assert.Equal(s.T(), tc.expectedError.Error(), err.Error())
				assert.Nil(s.T(), transaction)
			} else {
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), tc.expectedTransaction, transaction)
			}

			// Verify expected method calls
			s.transactionRepository.AssertExpectations(s.T())
			s.redisClient.AssertExpectations(s.T())
		})
	}
}

// TestGetTransactionsByUserID tests the GetTransactionsByUserID function
func (s *TransactionServiceTestSuite) TestGetTransactionsByUserID() {
	userID := "000018b0e1a211ef95a30242ac180003"
	now := time.Now()
	page := 1
	mockTotal := 2

	testCases := []struct {
		name                 string
		mockTransactions     []*models.Transaction
		mockTotal            int
		mockError            error
		expectedTransactions []*models.Transaction
		expectedTotal        int
		expectedError        error
	}{
		{
			name: "Success - Multiple Transactions",
			mockTransactions: []*models.Transaction{
				{
					TransactionID:   "000018b0e1a211ef95a30242ac180002",
					UserID:          userID,
					Name:            "Transaction 1",
					Image:           "transaction1.jpg",
					IsBank:          true,
					Amount:          100.50,
					TransactionType: "deposit",
					BaseModel: &models.BaseModel{
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
				{
					TransactionID:   "000018b0e1a211ef95a30242ac180004",
					UserID:          userID,
					Name:            "Transaction 2",
					Image:           "transaction2.jpg",
					IsBank:          false,
					Amount:          200.75,
					TransactionType: "withdrawal",
					BaseModel: &models.BaseModel{
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
			},
			mockTotal: mockTotal,
			mockError: nil,
			expectedTransactions: []*models.Transaction{
				{
					TransactionID:   "000018b0e1a211ef95a30242ac180002",
					UserID:          userID,
					Name:            "Transaction 1",
					Image:           "transaction1.jpg",
					IsBank:          true,
					Amount:          100.50,
					TransactionType: "deposit",
					BaseModel: &models.BaseModel{
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
				{
					TransactionID:   "000018b0e1a211ef95a30242ac180004",
					UserID:          userID,
					Name:            "Transaction 2",
					Image:           "transaction2.jpg",
					IsBank:          false,
					Amount:          200.75,
					TransactionType: "withdrawal",
					BaseModel: &models.BaseModel{
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
			},
			expectedTotal: mockTotal,
			expectedError: nil,
		},
		{
			name:                 "Success - No Transactions",
			mockTransactions:     []*models.Transaction{},
			mockTotal:            0,
			mockError:            nil,
			expectedTransactions: []*models.Transaction{},
			expectedTotal:        0,
			expectedError:        nil,
		},
		{
			name:                 "Failure - Database Error",
			mockTransactions:     nil,
			mockTotal:            0,
			mockError:            errors.New("database connection failed"),
			expectedTransactions: nil,
			expectedTotal:        0,
			expectedError:        errors.New("database connection failed"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Mock the repository method
			s.transactionRepository.On("GetByUserIDWithPagination", userID, mock.Anything, mock.Anything, mock.Anything).Return(tc.mockTransactions, tc.mockTotal, tc.mockError).Once()

			// Call the service method
			transactions, total, err := s.service.GetTransactionsByUserID(userID, page)

			// Assert results
			if tc.expectedError != nil {
				assert.Error(s.T(), err)
				assert.Equal(s.T(), tc.expectedError.Error(), err.Error())
				assert.Nil(s.T(), transactions)
				assert.Equal(s.T(), 0, total)
			} else {
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), tc.expectedTransactions, transactions)
				assert.Equal(s.T(), tc.expectedTotal, total)
			}

			// Verify expected method calls
			s.transactionRepository.AssertExpectations(s.T())
		})
	}
}

// TestCreateTransaction tests the CreateTransaction function
func (s *TransactionServiceTestSuite) TestCreateTransaction() {
	now := time.Now()
	testCases := []struct {
		name             string
		transaction      *models.Transaction
		mockError        error
		expectedError    error
		shouldGenerateID bool
	}{
		{
			name: "Success - With Existing ID",
			transaction: &models.Transaction{
				TransactionID:   "000018b0e1a211ef95a30242ac180002",
				UserID:          "000018b0e1a211ef95a30242ac180003",
				Name:            "Test Transaction",
				Image:           "transaction.jpg",
				IsBank:          true,
				Amount:          100.50,
				TransactionType: "deposit",
				BaseModel: &models.BaseModel{
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			mockError:        nil,
			expectedError:    nil,
			shouldGenerateID: false,
		},
		{
			name: "Success - Generate New ID",
			transaction: &models.Transaction{
				TransactionID:   "",
				UserID:          "000018b0e1a211ef95a30242ac180003",
				Name:            "Test Transaction",
				Image:           "transaction.jpg",
				IsBank:          true,
				Amount:          100.50,
				TransactionType: "deposit",
				BaseModel: &models.BaseModel{
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			mockError:        nil,
			expectedError:    nil,
			shouldGenerateID: true,
		},
		{
			name: "Failure - Database Error",
			transaction: &models.Transaction{
				TransactionID:   "000018b0e1a211ef95a30242ac180002",
				UserID:          "000018b0e1a211ef95a30242ac180003",
				Name:            "Test Transaction",
				Image:           "transaction.jpg",
				IsBank:          true,
				Amount:          100.50,
				TransactionType: "deposit",
				BaseModel: &models.BaseModel{
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			mockError:        errors.New("database connection failed"),
			expectedError:    errors.New("database connection failed"),
			shouldGenerateID: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// For the case where we need to generate a new ID, we can't know the exact ID in advance
			// So we use a matcher that checks any non-empty string
			if tc.shouldGenerateID {
				s.transactionRepository.On("Create", mock.MatchedBy(func(t *models.Transaction) bool {
					return t.TransactionID != "" && t.UserID == tc.transaction.UserID &&
						t.Name == tc.transaction.Name && t.Amount == tc.transaction.Amount &&
						t.TransactionType == tc.transaction.TransactionType
				})).Return(tc.mockError).Once()
			} else {
				s.transactionRepository.On("Create", tc.transaction).Return(tc.mockError).Once()
			}

			// Call the service method
			err := s.service.CreateTransaction(tc.transaction)

			// Assert results
			if tc.expectedError != nil {
				assert.Error(s.T(), err)
				assert.Equal(s.T(), tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(s.T(), err)
				if tc.shouldGenerateID {
					assert.NotEmpty(s.T(), tc.transaction.TransactionID)
				}
			}

			// Verify expected method calls
			s.transactionRepository.AssertExpectations(s.T())
		})
	}
}

// TestGetTransactionByIDCacheHit tests the GetTransactionByID function with a cache hit
func (s *TransactionServiceTestSuite) TestGetTransactionByIDCacheHit() {
	transactionID := "000018b0e1a211ef95a30242ac180002"
	now := time.Now()

	expectedTransaction := &models.Transaction{
		TransactionID:   transactionID,
		UserID:          "000018b0e1a211ef95a30242ac180003",
		Name:            "Test Transaction",
		Image:           "transaction.jpg",
		IsBank:          true,
		Amount:          100.50,
		TransactionType: "deposit",
		BaseModel: &models.BaseModel{
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	// Serialize the transaction for the cache
	cachedData, _ := json.Marshal(expectedTransaction)

	ctx := context.Background()
	cacheKey := fmt.Sprintf("transaction:%s", transactionID)

	// Mock Redis cache hit
	s.redisClient.On("Get", ctx, cacheKey).Return(string(cachedData), nil).Once()

	// Call the service method
	transaction, err := s.service.GetTransactionByID(transactionID)

	// Assert results
	assert.NoError(s.T(), err)
	
	// Compare fields individually, ignoring exact time comparison
	assert.Equal(s.T(), expectedTransaction.TransactionID, transaction.TransactionID)
	assert.Equal(s.T(), expectedTransaction.UserID, transaction.UserID)
	assert.Equal(s.T(), expectedTransaction.Name, transaction.Name)
	assert.Equal(s.T(), expectedTransaction.Image, transaction.Image)
	assert.Equal(s.T(), expectedTransaction.IsBank, transaction.IsBank)
	assert.Equal(s.T(), expectedTransaction.Amount, transaction.Amount)
	assert.Equal(s.T(), expectedTransaction.TransactionType, transaction.TransactionType)
	
	// For time fields, just check if they're not zero
	assert.False(s.T(), transaction.CreatedAt.IsZero())
	assert.False(s.T(), transaction.UpdatedAt.IsZero())

	// Verify that repository was not called (cache hit)
	s.redisClient.AssertExpectations(s.T())
	s.transactionRepository.AssertNotCalled(s.T(), "GetByID")
}

// Run the test suite
func TestTransactionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionServiceTestSuite))
}
