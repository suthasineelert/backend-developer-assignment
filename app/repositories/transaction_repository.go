package repositories

import (
	"backend-developer-assignment/app/models"
	"time"
)

// TransactionRepository is an interface for transaction repository
type TransactionRepository interface {
	GetByID(id string) (*models.Transaction, error)
	GetByUserIDWithPagination(userID, orderBy string, limit, offset int) ([]*models.Transaction, int, error)
	Create(transaction *models.Transaction) error
	Update(transaction *models.Transaction) error
}

// TransactionRepositoryImpl will hold all the repository operations related to transactions.
type TransactionRepositoryImpl struct {
	DB DB
}

// NewTransactionRepository creates a new instance of TransactionRepository.
func NewTransactionRepository(db DB) TransactionRepository {
	return &TransactionRepositoryImpl{
		DB: db,
	}
}

// GetByID retrieves one Transaction by given ID.
func (r *TransactionRepositoryImpl) GetByID(id string) (*models.Transaction, error) {
	transaction := &models.Transaction{}

	query := `SELECT transaction_id, account_id, user_id, name, image, isBank, amount, transaction_type, created_at, updated_at
	 FROM transactions WHERE transaction_id = ? and deleted_at IS NULL`

	err := r.DB.Get(transaction, query, id)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// GetByUserID retrieves all transactions for a given user ID.
func (r *TransactionRepositoryImpl) GetByUserID(userID string) ([]*models.Transaction, error) {
	transactions := []*models.Transaction{}

	query := `SELECT transaction_id, account_id, user_id, name, image, isBank, amount, transaction_type, created_at, updated_at
	 FROM transactions WHERE user_id = ? and deleted_at IS NULL ORDER BY created_at DESC`

	err := r.DB.Select(&transactions, query, userID)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

// GetByUserIDWithPagination retrieves paginated transactions for a given user ID.
func (r *TransactionRepositoryImpl) GetByUserIDWithPagination(userID, orderBy string, limit, offset int) ([]*models.Transaction, int, error) {
	transactions := []*models.Transaction{}

	// Query for paginated results
	query := `SELECT transaction_id, account_id, user_id, name, image, isBank, amount, transaction_type, created_at, updated_at
	FROM transactions WHERE user_id = ? and deleted_at IS NULL ORDER BY ? DESC LIMIT ? OFFSET ?`

	err := r.DB.Select(&transactions, query, userID, orderBy, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Get total count for pagination metadata
	var total int
	countQuery := `SELECT COUNT(*) FROM transactions WHERE user_id = ? and deleted_at IS NULL`
	err = r.DB.Get(&total, countQuery, userID)
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// Create adds a new transaction to the database.
func (r *TransactionRepositoryImpl) Create(transaction *models.Transaction) error {
	now := time.Now()
	transaction.CreatedAt = now
	transaction.UpdatedAt = now

	query := `INSERT INTO transactions (
		transaction_id, user_id, account_id, name, image, isBank, amount, transaction_type, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.DB.Exec(
		query,
		transaction.TransactionID,
		transaction.UserID,
		transaction.AccountID,
		transaction.Name,
		transaction.Image,
		transaction.IsBank,
		transaction.Amount,
		transaction.TransactionType,
		transaction.CreatedAt,
		transaction.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

// Update modifies an existing transaction.
func (r *TransactionRepositoryImpl) Update(transaction *models.Transaction) error {
	transaction.UpdatedAt = time.Now()

	query := `UPDATE transactions 
              SET user_id = ?, name = ?, image = ?, isBank = ?, amount = ?, 
                  transaction_type = ?, updated_at = ? 
              WHERE transaction_id = ? and deleted_at IS NULL`

	_, err := r.DB.Exec(
		query,
		transaction.UserID,
		transaction.Name,
		transaction.Image,
		transaction.IsBank,
		transaction.Amount,
		transaction.TransactionType,
		transaction.UpdatedAt,
		transaction.TransactionID,
	)
	if err != nil {
		return err
	}

	return nil
}
