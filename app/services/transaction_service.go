package services

import (
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/app/repositories"
	"backend-developer-assignment/pkg/configs"
	"backend-developer-assignment/pkg/types"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Cache expiration time constants
const (
	TransactionCacheDuration     = 10 * time.Minute
	TransactionListCacheDuration = 5 * time.Minute
)

// TransactionService interface defines the methods for transaction business logic
type TransactionService interface {
	GetTransactionByID(id string) (*models.Transaction, error)
	GetTransactionsByUserID(userID string, page int) ([]*models.Transaction, int, error)
	CreateTransaction(transaction *models.Transaction) error
}

// TransactionServiceImpl contains business logic related to transactions.
type TransactionServiceImpl struct {
	TransactionRepository repositories.TransactionRepository
	redisClient           types.CacheClient
}

// NewTransactionService creates a new TransactionService.
func NewTransactionService(transactionRepository repositories.TransactionRepository, redisClient types.CacheClient) TransactionService {
	return &TransactionServiceImpl{
		TransactionRepository: transactionRepository,
		redisClient:           redisClient,
	}
}

// GetTransactionByID retrieves a transaction by ID.
func (s *TransactionServiceImpl) GetTransactionByID(id string) (*models.Transaction, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("transaction:%s", id)

	// Try to get from cache first
	cachedData, err := s.redisClient.Get(ctx, cacheKey)
	if err == nil {
		// Cache hit - deserialize and return
		var transaction models.Transaction
		if err := json.Unmarshal([]byte(cachedData), &transaction); err == nil {
			return &transaction, nil
		}
	}

	// Cache miss or error - get from database
	transaction, err := s.TransactionRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Store in cache for future requests
	if transactionData, err := json.Marshal(transaction); err == nil {
		if err := s.redisClient.Set(ctx, cacheKey, transactionData, TransactionCacheDuration); err != nil {
			logger.Warn("Failed to set cache for transaction", zap.String("transaction_id", id), zap.Error(err))
		}
	}

	return transaction, nil
}

// GetTransactionsByUserID retrieves all transactions for a user.
func (s *TransactionServiceImpl) GetTransactionsByUserID(userID string, page int) ([]*models.Transaction, int, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("transactions:user:%s:page:%d", userID, page)
	countCacheKey := fmt.Sprintf("transactions:user:%s:count", userID)

	// Try to get from cache first
	cachedData, err := s.redisClient.Get(ctx, cacheKey)
	cachedCount, countErr := s.redisClient.Get(ctx, countCacheKey)

	if err == nil && countErr == nil {
		// Cache hit - deserialize and return
		var transactions []*models.Transaction
		var count int

		if err := json.Unmarshal([]byte(cachedData), &transactions); err == nil {
			if _, err := fmt.Sscanf(cachedCount, "%d", &count); err == nil {
				return transactions, count, nil
			}
		}
	}

	// Cache miss or error - get from database
	orderBy := "created_at desc"
	perPage := configs.DEFAULT_PAGE_SIZE
	offset := (page - 1) * perPage
	transactions, count, err := s.TransactionRepository.GetByUserIDWithPagination(userID, orderBy, perPage, offset)
	if err != nil {
		return nil, 0, err
	}

	// Store in cache for future requests
	if transactionsData, err := json.Marshal(transactions); err == nil {
		if err := s.redisClient.Set(ctx, cacheKey, transactionsData, TransactionListCacheDuration); err != nil {
			logger.Warn("Failed to set cache for transactions", zap.String("user_id", userID), zap.Error(err))
		}
		if err := s.redisClient.Set(ctx, countCacheKey, fmt.Sprintf("%d", count), TransactionListCacheDuration); err != nil {
			logger.Warn("Failed to set cache for transactions count", zap.String("user_id", userID), zap.Error(err))
		}
	}

	return transactions, count, nil
}

// CreateTransaction creates a new transaction.
func (s *TransactionServiceImpl) CreateTransaction(transaction *models.Transaction) error {
	// Generate a new UUID if not provided
	if transaction.TransactionID == "" {
		transaction.TransactionID = uuid.New().String()
	}

	// Create transaction in database
	err := s.TransactionRepository.Create(transaction)
	if err != nil {
		return err
	}

	// Invalidate user transactions cache
	ctx := context.Background()
	userTransactionsPattern := fmt.Sprintf("transactions:user:%s:*", transaction.UserID)
	s.invalidateCache(ctx, userTransactionsPattern)

	// Cache the new transaction
	cacheKey := fmt.Sprintf("transaction:%s", transaction.TransactionID)
	if transactionData, err := json.Marshal(transaction); err == nil {
		if err := s.redisClient.Set(ctx, cacheKey, transactionData, TransactionCacheDuration); err != nil {
			logger.Warn("Failed to set cache for new transaction", zap.String("transaction_id", transaction.TransactionID), zap.Error(err))
		}
	}

	return nil
}

// invalidateCache invalidates cache entries matching the given pattern
func (s *TransactionServiceImpl) invalidateCache(ctx context.Context, pattern string) {
	if err := s.redisClient.Delete(ctx, strings.TrimSuffix(pattern, "*")); err != nil {
		logger.Warn("Failed to invalidate cache", zap.String("pattern", pattern), zap.Error(err))
	}
}
