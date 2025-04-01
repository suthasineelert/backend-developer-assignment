package services_test

import (
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/app/services"
	"backend-developer-assignment/pkg/configs"
	mockCache "backend-developer-assignment/pkg/mocks/cache"
	mocks "backend-developer-assignment/pkg/mocks/repositories"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
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
	ctx                   context.Context
}

// SetupTest runs before each test
func (s *TransactionServiceTestSuite) SetupTest() {
	s.transactionRepository = new(mocks.TransactionRepository)
	s.redisClient = new(mockCache.RedisClient)
	s.service = services.NewTransactionService(s.transactionRepository, s.redisClient)
	s.ctx = context.Background()
}

// TestGetTransactionByIDCacheHit tests the GetTransactionByID function with a cache hit
func (s *TransactionServiceTestSuite) TestGetTransactionByIDCacheHit() {
	transactionID := "000018b0e1a211ef95a30242ac180002"
	now := time.Now().Truncate(time.Second)

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

	cacheKey := fmt.Sprintf("transaction:%s", transactionID)

	// Mock Redis cache hit
	s.redisClient.On("Get", s.ctx, cacheKey).Return(string(cachedData), nil).Once()

	// Call the service method
	transaction, err := s.service.GetTransactionByID(transactionID)

	// Assert results
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedTransaction.TransactionID, transaction.TransactionID)
	assert.Equal(s.T(), expectedTransaction.UserID, transaction.UserID)
	assert.Equal(s.T(), expectedTransaction.Name, transaction.Name)
	assert.Equal(s.T(), expectedTransaction.Image, transaction.Image)
	assert.Equal(s.T(), expectedTransaction.IsBank, transaction.IsBank)
	assert.Equal(s.T(), expectedTransaction.Amount, transaction.Amount)
	assert.Equal(s.T(), expectedTransaction.TransactionType, transaction.TransactionType)

	// Verify that repository was not called (cache hit)
	s.redisClient.AssertExpectations(s.T())
	s.transactionRepository.AssertNotCalled(s.T(), "GetByID")
}

// TestGetTransactionByIDCacheMiss tests the GetTransactionByID function with a cache miss
func (s *TransactionServiceTestSuite) TestGetTransactionByIDCacheMiss() {
	transactionID := "000018b0e1a211ef95a30242ac180002"
	now := time.Now().Truncate(time.Second)

	mockTransaction := &models.Transaction{
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

	cacheKey := fmt.Sprintf("transaction:%s", transactionID)

	// Mock Redis cache miss
	s.redisClient.On("Get", s.ctx, cacheKey).Return("", errors.New("cache miss")).Once()

	// Mock repository call
	s.transactionRepository.On("GetByID", transactionID).Return(mockTransaction, nil).Once()

	// Mock cache set
	transactionData, _ := json.Marshal(mockTransaction)
	s.redisClient.On("Set", s.ctx, cacheKey, transactionData, mock.Anything).Return(nil).Once()

	// Call the service method
	transaction, err := s.service.GetTransactionByID(transactionID)

	// Assert results
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), mockTransaction.TransactionID, transaction.TransactionID)
	assert.Equal(s.T(), mockTransaction.UserID, transaction.UserID)
	assert.Equal(s.T(), mockTransaction.Name, transaction.Name)
	assert.Equal(s.T(), mockTransaction.Image, transaction.Image)
	assert.Equal(s.T(), mockTransaction.IsBank, transaction.IsBank)
	assert.Equal(s.T(), mockTransaction.Amount, transaction.Amount)
	assert.Equal(s.T(), mockTransaction.TransactionType, transaction.TransactionType)

	// Verify expected method calls
	s.redisClient.AssertExpectations(s.T())
	s.transactionRepository.AssertExpectations(s.T())
}

// TestGetTransactionByIDDatabaseError tests the GetTransactionByID function with a database error
func (s *TransactionServiceTestSuite) TestGetTransactionByIDDatabaseError() {
	transactionID := "000018b0e1a211ef95a30242ac180002"
	expectedError := errors.New("database error")

	cacheKey := fmt.Sprintf("transaction:%s", transactionID)

	// Mock Redis cache miss
	s.redisClient.On("Get", s.ctx, cacheKey).Return("", errors.New("cache miss")).Once()

	// Mock repository error
	s.transactionRepository.On("GetByID", transactionID).Return(nil, expectedError).Once()

	// Call the service method
	transaction, err := s.service.GetTransactionByID(transactionID)

	// Assert results
	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError.Error(), err.Error())
	assert.Nil(s.T(), transaction)

	// Verify expected method calls
	s.redisClient.AssertExpectations(s.T())
	s.transactionRepository.AssertExpectations(s.T())
}

// TestGetTransactionByIDCacheUnmarshalError tests the GetTransactionByID function with a cache unmarshal error
func (s *TransactionServiceTestSuite) TestGetTransactionByIDCacheUnmarshalError() {
	transactionID := "000018b0e1a211ef95a30242ac180002"
	now := time.Now().Truncate(time.Second)

	mockTransaction := &models.Transaction{
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

	cacheKey := fmt.Sprintf("transaction:%s", transactionID)

	// Mock Redis cache hit but with invalid JSON
	s.redisClient.On("Get", s.ctx, cacheKey).Return("invalid json", nil).Once()

	// Mock repository call as fallback
	s.transactionRepository.On("GetByID", transactionID).Return(mockTransaction, nil).Once()

	// Mock cache set
	transactionData, _ := json.Marshal(mockTransaction)
	s.redisClient.On("Set", s.ctx, cacheKey, transactionData, mock.Anything).Return(nil).Once()

	// Call the service method
	transaction, err := s.service.GetTransactionByID(transactionID)

	// Assert results
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), mockTransaction.TransactionID, transaction.TransactionID)
	assert.Equal(s.T(), mockTransaction.UserID, transaction.UserID)

	// Verify expected method calls
	s.redisClient.AssertExpectations(s.T())
	s.transactionRepository.AssertExpectations(s.T())
}

// TestGetTransactionsByUserIDCacheHit tests the GetTransactionsByUserID function with a cache hit
func (s *TransactionServiceTestSuite) TestGetTransactionsByUserIDCacheHit() {
	userID := "000018b0e1a211ef95a30242ac180003"
	page := 1
	now := time.Now().Truncate(time.Second)
	expectedCount := 2

	expectedTransactions := []*models.Transaction{
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
	}

	// Serialize the transactions for the cache
	cachedData, _ := json.Marshal(expectedTransactions)
	cachedCount := fmt.Sprintf("%d", expectedCount)

	cacheKey := fmt.Sprintf("transactions:user:%s:page:%d", userID, page)
	countCacheKey := fmt.Sprintf("transactions:user:%s:count", userID)

	// Mock Redis cache hit
	s.redisClient.On("Get", s.ctx, cacheKey).Return(string(cachedData), nil).Once()
	s.redisClient.On("Get", s.ctx, countCacheKey).Return(cachedCount, nil).Once()

	// Call the service method
	transactions, count, err := s.service.GetTransactionsByUserID(userID, page)

	// Assert results
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedCount, count)
	assert.Len(s.T(), transactions, 2)
	assert.Equal(s.T(), expectedTransactions[0].TransactionID, transactions[0].TransactionID)
	assert.Equal(s.T(), expectedTransactions[1].TransactionID, transactions[1].TransactionID)

	// Verify that repository was not called (cache hit)
	s.redisClient.AssertExpectations(s.T())
	s.transactionRepository.AssertNotCalled(s.T(), "GetByUserIDWithPagination")
}

// TestGetTransactionsByUserIDCacheMiss tests the GetTransactionsByUserID function with a cache miss
func (s *TransactionServiceTestSuite) TestGetTransactionsByUserIDCacheMiss() {
	userID := "000018b0e1a211ef95a30242ac180003"
	page := 1
	now := time.Now().Truncate(time.Second)
	expectedCount := 2
	perPage := configs.DEFAULT_PAGE_SIZE
	offset := (page - 1) * perPage
	orderBy := "created_at desc"

	mockTransactions := []*models.Transaction{
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
	}

	cacheKey := fmt.Sprintf("transactions:user:%s:page:%d", userID, page)
	countCacheKey := fmt.Sprintf("transactions:user:%s:count", userID)

	// Mock Redis cache miss
	s.redisClient.On("Get", s.ctx, cacheKey).Return("", errors.New("cache miss")).Once()
	s.redisClient.On("Get", s.ctx, countCacheKey).Return("", errors.New("cache miss")).Once()

	// Mock repository call
	s.transactionRepository.On("GetByUserIDWithPagination", userID, orderBy, perPage, offset).
		Return(mockTransactions, expectedCount, nil).Once()

	// Mock cache set
	transactionsData, _ := json.Marshal(mockTransactions)
	s.redisClient.On("Set", s.ctx, cacheKey, transactionsData, mock.Anything).Return(nil).Once()
	s.redisClient.On("Set", s.ctx, countCacheKey, fmt.Sprintf("%d", expectedCount), mock.Anything).Return(nil).Once()

	// Call the service method
	transactions, count, err := s.service.GetTransactionsByUserID(userID, page)

	// Assert results
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedCount, count)
	assert.Len(s.T(), transactions, 2)
	assert.Equal(s.T(), mockTransactions[0].TransactionID, transactions[0].TransactionID)
	assert.Equal(s.T(), mockTransactions[1].TransactionID, transactions[1].TransactionID)

	// Verify expected method calls
	s.redisClient.AssertExpectations(s.T())
	s.transactionRepository.AssertExpectations(s.T())
}

// TestGetTransactionsByUserIDDatabaseError tests the GetTransactionsByUserID function with a database error
func (s *TransactionServiceTestSuite) TestGetTransactionsByUserIDDatabaseError() {
	userID := "000018b0e1a211ef95a30242ac180003"
	page := 1
	perPage := configs.DEFAULT_PAGE_SIZE
	offset := (page - 1) * perPage
	orderBy := "created_at desc"
	expectedError := errors.New("database connection failed")

	cacheKey := fmt.Sprintf("transactions:user:%s:page:%d", userID, page)
	countCacheKey := fmt.Sprintf("transactions:user:%s:count", userID)

	// Mock Redis cache miss
	s.redisClient.On("Get", s.ctx, cacheKey).Return("", errors.New("cache miss")).Once()
	s.redisClient.On("Get", s.ctx, countCacheKey).Return("", errors.New("cache miss")).Once()

	// Mock repository error
	s.transactionRepository.On("GetByUserIDWithPagination", userID, orderBy, perPage, offset).
		Return(nil, 0, expectedError).Once()

	// Call the service method
	transactions, count, err := s.service.GetTransactionsByUserID(userID, page)

	// Assert results
	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError.Error(), err.Error())
	assert.Nil(s.T(), transactions)
	assert.Equal(s.T(), 0, count)

	// Verify expected method calls
	s.redisClient.AssertExpectations(s.T())
	s.transactionRepository.AssertExpectations(s.T())
}

// TestCreateTransactionWithExistingID tests the CreateTransaction function with an existing ID
func (s *TransactionServiceTestSuite) TestCreateTransactionWithExistingID() {
	now := time.Now().Truncate(time.Second)
	transactionID := "000018b0e1a211ef95a30242ac180002"
	userID := "000018b0e1a211ef95a30242ac180003"

	transaction := &models.Transaction{
		TransactionID:   transactionID,
		UserID:          userID,
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

	// Mock repository call
	s.transactionRepository.On("Create", transaction).Return(nil).Once()

	// Mock cache operations
	userTransactionsPattern := fmt.Sprintf("transactions:user:%s:*", userID)
	s.redisClient.On("Delete", s.ctx, strings.TrimSuffix(userTransactionsPattern, "*")).Return(nil).Once()

	// Mock cache set for the new transaction
	cacheKey := fmt.Sprintf("transaction:%s", transactionID)
	transactionData, _ := json.Marshal(transaction)
	s.redisClient.On("Set", s.ctx, cacheKey, transactionData, mock.Anything).Return(nil).Once()

	// Call the service method
	err := s.service.CreateTransaction(transaction)

	// Assert results
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), transactionID, transaction.TransactionID) // ID should remain unchanged

	// Verify expected method calls
	s.transactionRepository.AssertExpectations(s.T())
	s.redisClient.AssertExpectations(s.T())
}

// TestCreateTransactionWithGeneratedID tests the CreateTransaction function with a generated ID
func (s *TransactionServiceTestSuite) TestCreateTransactionWithGeneratedID() {
	now := time.Now().Truncate(time.Second)
	userID := "000018b0e1a211ef95a30242ac180003"

	transaction := &models.Transaction{
		TransactionID:   "", // Empty ID should trigger generation
		UserID:          userID,
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

	// Mock repository call with ID matcher
	s.transactionRepository.On("Create", mock.MatchedBy(func(t *models.Transaction) bool {
		// Verify that an ID was generated (non-empty)
		return t.TransactionID != "" &&
			t.UserID == userID &&
			t.Name == "Test Transaction" &&
			t.Amount == 100.50
	})).Return(nil).Once()

	// Mock cache operations
	userTransactionsPattern := fmt.Sprintf("transactions:user:%s:*", userID)
	s.redisClient.On("Delete", s.ctx, strings.TrimSuffix(userTransactionsPattern, "*")).Return(nil).Once()

	// Mock cache set for the new transaction (with any ID)
	s.redisClient.On("Set", s.ctx, mock.MatchedBy(func(key string) bool {
		return strings.HasPrefix(key, "transaction:")
	}), mock.Anything, mock.Anything).Return(nil).Once()

	// Call the service method
	err := s.service.CreateTransaction(transaction)

	// Assert results
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), transaction.TransactionID) // ID should be generated
	_, err = uuid.Parse(transaction.TransactionID)
	assert.NoError(s.T(), err) // ID should be a valid UUID

	// Verify expected method calls
	s.transactionRepository.AssertExpectations(s.T())
	s.redisClient.AssertExpectations(s.T())
}

// TestCreateTransactionDatabaseError tests the CreateTransaction function with a database error
func (s *TransactionServiceTestSuite) TestCreateTransactionDatabaseError() {
	now := time.Now().Truncate(time.Second)
	transactionID := "000018b0e1a211ef95a30242ac180002"
	userID := "000018b0e1a211ef95a30242ac180003"
	expectedError := errors.New("database error")

	transaction := &models.Transaction{
		TransactionID:   transactionID,
		UserID:          userID,
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

	// Mock repository error
	s.transactionRepository.On("Create", transaction).Return(expectedError).Once()

	// Call the service method
	err := s.service.CreateTransaction(transaction)

	// Assert results
	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError.Error(), err.Error())

	// Verify expected method calls
	s.transactionRepository.AssertExpectations(s.T())
	// Cache operations should not be called on error
	s.redisClient.AssertNotCalled(s.T(), "Delete")
	s.redisClient.AssertNotCalled(s.T(), "Set")
}

// Run the test suite
func TestTransactionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionServiceTestSuite))
}
