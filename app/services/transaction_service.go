package services

import (
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/app/repositories"
	"backend-developer-assignment/pkg/configs"

	"github.com/google/uuid"
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
}

// NewTransactionService creates a new TransactionService.
func NewTransactionService(transactionRepository repositories.TransactionRepository) TransactionService {
	return &TransactionServiceImpl{
		TransactionRepository: transactionRepository,
	}
}

// GetTransactionByID retrieves a transaction by ID.
func (s *TransactionServiceImpl) GetTransactionByID(id string) (*models.Transaction, error) {
	return s.TransactionRepository.GetByID(id)
}

// GetTransactionsByUserID retrieves all transactions for a user.
func (s *TransactionServiceImpl) GetTransactionsByUserID(userID string, page int) ([]*models.Transaction, int, error) {
	orderBy := "created_at desc"
	perPage := configs.DEFAULT_PAGE_SIZE
	offset := (page - 1) * perPage
	return s.TransactionRepository.GetByUserIDWithPagination(userID, orderBy, perPage, offset)
}

// CreateTransaction creates a new transaction.
func (s *TransactionServiceImpl) CreateTransaction(transaction *models.Transaction) error {
	// Generate a new UUID if not provided
	if transaction.TransactionID == "" {
		transaction.TransactionID = uuid.New().String()
	}

	return s.TransactionRepository.Create(transaction)
}
