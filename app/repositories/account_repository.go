package repositories

import (
	"backend-developer-assignment/app/models"
	"time"

	"github.com/jmoiron/sqlx"
)

// AccountRepository is an interface for account repository operations
type AccountRepository interface {
	// Get account operations
	GetAccountByID(accountID string) (*models.Account, error)
	GetAccountsByUserID(userID string) ([]*models.Account, error)
	GetAccountDetailByID(accountID string) (*models.AccountDetail, error)
	GetAccountBalanceByID(accountID string) (*models.AccountBalance, error)
	GetAccountFlagsByAccountID(accountID string) ([]*models.AccountFlag, error)

	// Create operations
	BeginTx() (DBTransaction, error)
	CreateAccountTx(tx DBTransaction, account *models.Account) error
	CreateAccountDetailTx(tx DBTransaction, detail *models.AccountDetail) error
	CreateAccountBalanceTx(tx DBTransaction, balance *models.AccountBalance) error
	CreateAccountFlagTx(tx DBTransaction, flag *models.AccountFlag) error

	// Update operations
	UpdateAccountTx(tx DBTransaction, account *models.Account) error
	UpdateAccountDetailTx(tx DBTransaction, detail *models.AccountDetail) error
	UpdateAccountBalanceTx(tx DBTransaction, balance *models.AccountBalance) error

	// Delete operations
	DeleteAccount(accountID string) error
}

// AccountRepositoryImpl implements AccountRepository
type AccountRepositoryImpl struct {
	DB *sqlx.DB
}

// NewAccountRepository creates a new instance of AccountRepository
func NewAccountRepository(db *sqlx.DB) AccountRepository {
	return &AccountRepositoryImpl{
		DB: db,
	}
}

// GetAccountByID retrieves an account by ID
func (r *AccountRepositoryImpl) GetAccountByID(accountID string) (*models.Account, error) {
	account := &models.Account{}
	query := `SELECT * FROM accounts WHERE account_id = ? AND deleted_at IS NULL`
	err := r.DB.Get(account, query, accountID)
	if err != nil {
		return nil, err
	}
	return account, nil
}

// GetAccountsByUserID retrieves all accounts for a user
func (r *AccountRepositoryImpl) GetAccountsByUserID(userID string) ([]*models.Account, error) {
	accounts := []*models.Account{}
	query := `SELECT * FROM accounts WHERE user_id = ? AND deleted_at IS NULL`
	err := r.DB.Select(&accounts, query, userID)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

// GetAccountDetailByID retrieves account details by account ID
func (r *AccountRepositoryImpl) GetAccountDetailByID(accountID string) (*models.AccountDetail, error) {
	detail := &models.AccountDetail{}
	query := `SELECT * FROM account_details WHERE account_id = ? AND deleted_at IS NULL`
	err := r.DB.Get(detail, query, accountID)
	if err != nil {
		return nil, err
	}
	return detail, nil
}

// GetAccountBalanceByID retrieves account balance by account ID
func (r *AccountRepositoryImpl) GetAccountBalanceByID(accountID string) (*models.AccountBalance, error) {
	balance := &models.AccountBalance{}
	query := `SELECT * FROM account_balances WHERE account_id = ? AND deleted_at IS NULL`
	err := r.DB.Get(balance, query, accountID)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

// GetAccountFlagsByAccountID retrieves all flags for an account
func (r *AccountRepositoryImpl) GetAccountFlagsByAccountID(accountID string) ([]*models.AccountFlag, error) {
	flags := []*models.AccountFlag{}
	query := `SELECT * FROM account_flags WHERE account_id = ? AND deleted_at IS NULL`
	err := r.DB.Select(&flags, query, accountID)
	if err != nil {
		return nil, err
	}
	return flags, nil
}

// BeginTx starts a new transaction
func (r *AccountRepositoryImpl) BeginTx() (DBTransaction, error) {
	tx, err := r.DB.Beginx()
	if err != nil {
		return nil, err
	}
	return &sqlxTransaction{tx: tx}, nil
}

// CreateAccountTx creates a new account within a transaction
func (r *AccountRepositoryImpl) CreateAccountTx(tx DBTransaction, account *models.Account) error {
	now := time.Now()
	account.CreatedAt = now
	account.UpdatedAt = now

	query := `INSERT INTO accounts (account_id, user_id, type, currency, account_number, issuer, created_at, updated_at) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := tx.Exec(
		query,
		account.AccountID,
		account.UserID,
		account.Type,
		account.Currency,
		account.AccountNumber,
		account.Issuer,
		account.CreatedAt,
		account.UpdatedAt,
	)
	return err
}

// CreateAccountDetailTx creates new account details within a transaction
func (r *AccountRepositoryImpl) CreateAccountDetailTx(tx DBTransaction, detail *models.AccountDetail) error {
	now := time.Now()
	detail.CreatedAt = now
	detail.UpdatedAt = now

	query := `INSERT INTO account_details (account_id, user_id, color, is_main_account, progress, created_at, updated_at) 
              VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := tx.Exec(
		query,
		detail.AccountID,
		detail.UserID,
		detail.Color,
		detail.IsMainAccount,
		detail.Progress,
		detail.CreatedAt,
		detail.UpdatedAt,
	)
	return err
}

// CreateAccountBalanceTx creates new account balance within a transaction
func (r *AccountRepositoryImpl) CreateAccountBalanceTx(tx DBTransaction, balance *models.AccountBalance) error {
	now := time.Now()
	balance.CreatedAt = now
	balance.UpdatedAt = now

	query := `INSERT INTO account_balances (account_id, user_id, amount, created_at, updated_at) 
              VALUES (?, ?, ?, ?, ?)`
	_, err := tx.Exec(
		query,
		balance.AccountID,
		balance.UserID,
		balance.Amount,
		balance.CreatedAt,
		balance.UpdatedAt,
	)
	return err
}

// CreateAccountFlagTx creates a new account flag within a transaction
func (r *AccountRepositoryImpl) CreateAccountFlagTx(tx DBTransaction, flag *models.AccountFlag) error {
	now := time.Now()
	flag.CreatedAt = now
	flag.UpdatedAt = now

	query := `INSERT INTO account_flags (account_id, user_id, flag_type, flag_value, created_at, updated_at) 
              VALUES (?, ?, ?, ?, ?, ?)`
	_, err := tx.Exec(
		query,
		flag.AccountID,
		flag.UserID,
		flag.FlagType,
		flag.FlagValue,
		flag.CreatedAt,
		flag.UpdatedAt,
	)
	return err
}

// UpdateAccountTx updates an existing account within a transaction
func (r *AccountRepositoryImpl) UpdateAccountTx(tx DBTransaction, account *models.Account) error {
	account.UpdatedAt = time.Now()

	query := `UPDATE accounts 
              SET user_id = ?, type = ?, currency = ?, account_number = ?, issuer = ?, updated_at = ? 
              WHERE account_id = ? AND deleted_at IS NULL`
	_, err := tx.Exec(
		query,
		account.UserID,
		account.Type,
		account.Currency,
		account.AccountNumber,
		account.Issuer,
		account.UpdatedAt,
		account.AccountID,
	)
	return err
}

// UpdateAccountDetailTx updates existing account details within a transaction
func (r *AccountRepositoryImpl) UpdateAccountDetailTx(tx DBTransaction, detail *models.AccountDetail) error {
	detail.UpdatedAt = time.Now()

	query := `UPDATE account_details 
              SET user_id = ?, color = ?, is_main_account = ?, progress = ?, updated_at = ? 
              WHERE account_id = ? AND deleted_at IS NULL`
	_, err := tx.Exec(
		query,
		detail.UserID,
		detail.Color,
		detail.IsMainAccount,
		detail.Progress,
		detail.UpdatedAt,
		detail.AccountID,
	)
	return err
}

// UpdateAccountBalanceTx updates existing account balance within a transaction
func (r *AccountRepositoryImpl) UpdateAccountBalanceTx(tx DBTransaction, balance *models.AccountBalance) error {
	balance.UpdatedAt = time.Now()

	query := `UPDATE account_balances 
              SET user_id = ?, amount = ?, updated_at = ? 
              WHERE account_id = ? AND deleted_at IS NULL`
	_, err := tx.Exec(
		query,
		balance.UserID,
		balance.Amount,
		balance.UpdatedAt,
		balance.AccountID,
	)
	return err
}

// DeleteAccount marks an account as deleted without removing it
func (r *AccountRepositoryImpl) DeleteAccount(accountID string) error {
	now := time.Now()
	query := `UPDATE accounts SET deleted_at = ? WHERE account_id = ? AND deleted_at IS NULL`
	_, err := r.DB.Exec(query, now, accountID)
	return err
}
