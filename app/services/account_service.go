package services

import (
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/app/repositories"
	"time"

	"github.com/google/uuid"
)

// AccountService defines the interface for account operations
type AccountService interface {
	// Get account operations
	GetAccountByID(accountID string) (*models.Account, error)
	GetAccountsByUserID(userID string) ([]*models.Account, error)
	GetAccountWithDetails(accountID string) (*models.AccountWithDetails, error)
	GetAccountsWithDetailsByUserID(userID string) ([]*models.AccountWithDetails, error)

	// Create operations
	CreateAccountWithDetails(accountDetails *models.AccountWithDetails) error

	// Update operations
	UpdateAccount(account *models.Account) error
	UpdateAccountBalance(accountID string, amount float64) error

	// Delete operations
	DeleteAccount(accountID string) error
}

// AccountServiceImpl implements AccountService
type AccountServiceImpl struct {
	accountRepository repositories.AccountRepository
}

// NewAccountService creates a new instance of AccountService
func NewAccountService(repo repositories.AccountRepository) AccountService {
	return &AccountServiceImpl{
		accountRepository: repo,
	}
}

// GetAccountByID retrieves an account by ID
func (s *AccountServiceImpl) GetAccountByID(accountID string) (*models.Account, error) {
	return s.accountRepository.GetAccountByID(accountID)
}

// GetAccountsByUserID retrieves all accounts for a user
func (s *AccountServiceImpl) GetAccountsByUserID(userID string) ([]*models.Account, error) {
	return s.accountRepository.GetAccountsByUserID(userID)
}

// GetAccountWithDetails retrieves an account with all its related details
func (s *AccountServiceImpl) GetAccountWithDetails(accountID string) (*models.AccountWithDetails, error) {
	// Get the base account
	account, err := s.accountRepository.GetAccountByID(accountID)
	if err != nil {
		return nil, err
	}

	// Get account details
	detail, err := s.accountRepository.GetAccountDetailByID(accountID)
	if err != nil {
		return nil, err
	}

	// Get account balance
	balance, err := s.accountRepository.GetAccountBalanceByID(accountID)
	if err != nil {
		return nil, err
	}

	// Get account flags
	flags, err := s.accountRepository.GetAccountFlagsByAccountID(accountID)
	if err != nil {
		return nil, err
	}

	// Combine all data into a single flattened structure
	result := &models.AccountWithDetails{
		AccountID:     account.AccountID,
		UserID:        account.UserID,
		Type:          account.Type,
		Currency:      account.Currency,
		AccountNumber: account.AccountNumber,
		Issuer:        account.Issuer,
		CreatedAt:     account.CreatedAt,
		UpdatedAt:     account.UpdatedAt,

		Color:         detail.Color,
		IsMainAccount: detail.IsMainAccount,
		Progress:      detail.Progress,

		Amount: balance.Amount,

		Flags: flags,
	}

	return result, nil
}

// GetAccountsWithDetailsByUserID retrieves all accounts with details for a user
func (s *AccountServiceImpl) GetAccountsWithDetailsByUserID(userID string) ([]*models.AccountWithDetails, error) {
	// Get all accounts for the user
	accounts, err := s.accountRepository.GetAccountsByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Create a slice to hold the results
	results := make([]*models.AccountWithDetails, 0, len(accounts))

	// For each account, get the related details
	for _, account := range accounts {
		// Get account details
		detail, err := s.accountRepository.GetAccountDetailByID(account.AccountID)
		if err != nil {
			continue // Skip accounts with missing details
		}

		// Get account balance
		balance, err := s.accountRepository.GetAccountBalanceByID(account.AccountID)
		if err != nil {
			continue // Skip accounts with missing balance
		}

		// Get account flags
		flags, err := s.accountRepository.GetAccountFlagsByAccountID(account.AccountID)
		if err != nil {
			flags = []*models.AccountFlag{} // Use empty slice if no flags
		}

		// Combine all data into a single flattened structure
		result := &models.AccountWithDetails{
			AccountID:     account.AccountID,
			UserID:        account.UserID,
			Type:          account.Type,
			Currency:      account.Currency,
			AccountNumber: account.AccountNumber,
			Issuer:        account.Issuer,
			CreatedAt:     account.CreatedAt,
			UpdatedAt:     account.UpdatedAt,

			Color:         detail.Color,
			IsMainAccount: detail.IsMainAccount,
			Progress:      detail.Progress,

			Amount: balance.Amount,

			Flags: flags,
		}

		results = append(results, result)
	}

	return results, nil
}

// CreateAccountWithDetails creates a new account with all related details
func (s *AccountServiceImpl) CreateAccountWithDetails(accountDetails *models.AccountWithDetails) error {
	// Generate a new UUID if not provided
	if accountDetails.AccountID == "" {
		accountDetails.AccountID = uuid.New().String()
	}

	// Set current time for timestamps if not provided
	now := time.Now()
	if accountDetails.CreatedAt.IsZero() {
		accountDetails.CreatedAt = now
	}
	if accountDetails.UpdatedAt.IsZero() {
		accountDetails.UpdatedAt = now
	}

	// Create individual model objects from the flattened structure
	account := &models.Account{
		AccountID:     accountDetails.AccountID,
		UserID:        accountDetails.UserID,
		Type:          accountDetails.Type,
		Currency:      accountDetails.Currency,
		AccountNumber: accountDetails.AccountNumber,
		Issuer:        accountDetails.Issuer,
	}

	detail := &models.AccountDetail{
		AccountID:     accountDetails.AccountID,
		UserID:        accountDetails.UserID,
		Color:         accountDetails.Color,
		IsMainAccount: accountDetails.IsMainAccount,
		Progress:      accountDetails.Progress,
	}

	balance := &models.AccountBalance{
		AccountID: accountDetails.AccountID,
		UserID:    accountDetails.UserID,
		Amount:    accountDetails.Amount,
	}

	// Use a transaction to ensure all operations succeed or fail together
	tx, err := s.accountRepository.BeginTx()
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				// Log the rollback error but return the original error
				// Consider using a logger here instead of just ignoring the error
				// logger.Error("Failed to rollback transaction", zap.Error(rbErr))
			}
		}
	}()

	// Create the main account
	if err := s.accountRepository.CreateAccountTx(tx, account); err != nil {
		return err
	}

	// Create account details
	if err := s.accountRepository.CreateAccountDetailTx(tx, detail); err != nil {
		return err
	}

	// Create account balance
	if err := s.accountRepository.CreateAccountBalanceTx(tx, balance); err != nil {
		return err
	}

	// Create account flags if any
	for _, flag := range accountDetails.Flags {
		flag.AccountID = accountDetails.AccountID
		flag.UserID = accountDetails.UserID
		if err := s.accountRepository.CreateAccountFlagTx(tx, flag); err != nil {
			return err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// UpdateAccount updates an existing account
func (s *AccountServiceImpl) UpdateAccount(account *models.Account) error {
	// Use a transaction to ensure all operations succeed or fail together
	tx, err := s.accountRepository.BeginTx()
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Update the account
	if err := s.accountRepository.UpdateAccountTx(tx, account); err != nil {
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// UpdateAccountBalance updates the balance of an account
func (s *AccountServiceImpl) UpdateAccountBalance(accountID string, amount float64) error {
	// Get the current balance
	balance, err := s.accountRepository.GetAccountBalanceByID(accountID)
	if err != nil {
		return err
	}

	// Update the balance amount
	balance.Amount = amount

	// Use a transaction to update the balance
	tx, err := s.accountRepository.BeginTx()
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Update the balance
	if err := s.accountRepository.UpdateAccountBalanceTx(tx, balance); err != nil {
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// DeleteAccount marks an account as deleted
func (s *AccountServiceImpl) DeleteAccount(accountID string) error {
	return s.accountRepository.DeleteAccount(accountID)
}
