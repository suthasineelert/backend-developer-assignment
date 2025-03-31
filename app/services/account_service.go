package services

import (
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/app/repositories"
	"backend-developer-assignment/pkg/types"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Custom errors for account operations
var (
	ErrInsufficientFunds = errors.New("insufficient funds")
)

// AccountService defines the interface for account operations
type AccountService interface {
	// Account operations
	GetAccountByID(accountID string) (*models.Account, error)
	GetAccountsByUserID(userID string) ([]*models.Account, error)
	GetAccountWithDetailByID(accountID string) (*models.AccountWithDetails, error)
	GetAccountsWithDetailByUserID(userID string) ([]*models.AccountWithDetails, error)

	// Create operations
	CreateAccountWithDetails(accountWithDetails *models.AccountWithDetails) error

	// Update operations
	UpdateAccount(account *models.AccountWithDetails) error
	SetMainAccount(account *models.Account) error

	// Transaction operations
	WithdrawFromAccount(accountID string, amount float64) (float64, error)
	TransferBetweenAccounts(fromAccountID, toAccountID string, amount float64) (*types.TransferResult, error)
	DepositToAccount(accountID string, amount float64) (float64, error)

	// Delete operations
	DeleteAccount(accountID string) error
}

// AccountServiceImpl implements AccountService
type AccountServiceImpl struct {
	accountRepository     repositories.AccountRepository
	transactionRepository repositories.TransactionRepository
	txProvider            repositories.TxProvider
}

// NewAccountService creates a new instance of AccountService
func NewAccountService(accountRepo repositories.AccountRepository, transactionRepo repositories.TransactionRepository, txProvider repositories.TxProvider) AccountService {
	return &AccountServiceImpl{
		accountRepository:     accountRepo,
		transactionRepository: transactionRepo,
		txProvider:            txProvider,
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

// GetAccountWithDetailByID retrieves a complete account with all related information by ID
func (s *AccountServiceImpl) GetAccountWithDetailByID(accountID string) (*models.AccountWithDetails, error) {
	return s.accountRepository.GetAccountWithDetailByID(accountID)
}

// GetAccountsWithDetailByUserID retrieves all complete accounts with related information for a user
func (s *AccountServiceImpl) GetAccountsWithDetailByUserID(userID string) ([]*models.AccountWithDetails, error) {
	return s.accountRepository.GetAccountsWithDetailByUserID(userID)
}

// CreateAccountWithDetails creates a new account with all related details
func (s *AccountServiceImpl) CreateAccountWithDetails(accountWithDetails *models.AccountWithDetails) error {
	// Generate a new UUID if not provided
	if accountWithDetails.AccountID == "" {
		accountWithDetails.AccountID = uuid.New().String()
	}

	if err := s.accountRepository.CreateAccount(accountWithDetails); err != nil {
		return err
	}

	return nil
}

// UpdateAccount updates an existing account
func (s *AccountServiceImpl) UpdateAccount(account *models.AccountWithDetails) error {
	return s.accountRepository.UpdateAccountByID(account.AccountID, account.UserID, func(accountWithDetails *models.AccountWithDetails) (bool, error) {
		isUpdate := false

		// Update account fields
		if account.Type != "" && accountWithDetails.Type != account.Type {
			accountWithDetails.Type = account.Type
			isUpdate = true
		}
		if account.Currency != "" && accountWithDetails.Currency != account.Currency {
			accountWithDetails.Currency = account.Currency
			isUpdate = true
		}
		if account.AccountNumber != "" && accountWithDetails.AccountNumber != account.AccountNumber {
			accountWithDetails.AccountNumber = account.AccountNumber
			isUpdate = true
		}
		if account.Issuer != "" && accountWithDetails.Issuer != account.Issuer {
			accountWithDetails.Issuer = account.Issuer
			isUpdate = true
		}

		// Update account detail fields
		if account.Color != "" && accountWithDetails.Color != account.Color {
			accountWithDetails.Color = account.Color
			isUpdate = true
		}
		if account.Progress > 0 && accountWithDetails.Progress != account.Progress {
			accountWithDetails.Progress = account.Progress
			isUpdate = true
		}

		return isUpdate, nil
	})
}

// UpdateAccount updates an existing account
func (s *AccountServiceImpl) SetMainAccount(account *models.Account) error {
	if err := s.accountRepository.UnSetMainAccount(account.UserID); err != nil {
		logger.Error("Unable to unset main account", zap.String("account_id", account.AccountID), zap.String("user_id", account.UserID), zap.Error(err))
		return err
	}

	if err := s.accountRepository.SetMainAccount(account.AccountID, account.UserID); err != nil {
		logger.Error("Unable to set main account", zap.String("account_id", account.AccountID), zap.String("user_id", account.UserID), zap.Error(err))
		return err
	}

	return nil
}

// WithdrawFromAccount withdraws money from an account with proper locking to prevent race conditions
func (s *AccountServiceImpl) WithdrawFromAccount(accountID string, amount float64) (float64, error) {
	var updatedBalance float64

	// Get account details for transaction record
	account, err := s.GetAccountWithDetailByID(accountID)
	if err != nil {
		logger.Error("Failed to get account details", zap.String("account_id", accountID), zap.Error(err))
		return 0, err
	}

	// Use transaction provider to handle the transaction
	err = s.txProvider.Transact(func(adapters repositories.Adapters) error {
		// Update account balance within transaction
		balanceErr := adapters.AccountRepository.UpdateAccountBalance(accountID, func(currentBalance float64) (float64, error) {
			// Check if there are sufficient funds
			if currentBalance < amount {
				return 0, ErrInsufficientFunds
			}

			// Calculate the new balance
			updatedBalance = currentBalance - amount
			return updatedBalance, nil
		})

		if balanceErr != nil {
			return balanceErr
		}

		// Create withdrawal transaction record
		withdrawalTx := &models.Transaction{
			BaseModel:       &models.BaseModel{},
			TransactionID:   uuid.New().String(),
			UserID:          account.UserID,
			Name:            "Withdrawal",
			IsBank:          true,
			Amount:          amount,
			TransactionType: string(models.Withdrawal),
			AccountID:       accountID,
		}

		// Save the transaction record
		if err := adapters.TransactionRepository.Create(withdrawalTx); err != nil {
			logger.Error("Failed to create withdrawal transaction record",
				zap.String("account_id", accountID),
				zap.Error(err))
			return err
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return updatedBalance, nil
}

// DepositToAccount deposits money to an account with proper locking to prevent race conditions
func (s *AccountServiceImpl) DepositToAccount(accountID string, amount float64) (float64, error) {
	var updatedBalance float64

	// Get account details for transaction record
	account, err := s.GetAccountWithDetailByID(accountID)
	if err != nil {
		logger.Error("Failed to get account details", zap.String("account_id", accountID), zap.Error(err))
		return 0, err
	}

	// Use transaction provider to handle the transaction
	err = s.txProvider.Transact(func(adapters repositories.Adapters) error {
		// Update account balance within transaction
		balanceErr := adapters.AccountRepository.UpdateAccountBalance(accountID, func(currentBalance float64) (float64, error) {
			// Calculate the new balance
			updatedBalance = currentBalance + amount
			return updatedBalance, nil
		})

		if balanceErr != nil {
			return balanceErr
		}

		// Create deposit transaction record
		depositTx := &models.Transaction{
			BaseModel:       &models.BaseModel{},
			TransactionID:   uuid.New().String(),
			UserID:          account.UserID,
			Name:            "Deposit",
			IsBank:          true,
			Amount:          amount,
			TransactionType: string(models.Deposit),
			AccountID:       accountID,
		}

		// Save the transaction record
		if err := adapters.TransactionRepository.Create(depositTx); err != nil {
			logger.Error("Failed to create deposit transaction record",
				zap.String("account_id", accountID),
				zap.Error(err))
			return err
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return updatedBalance, nil
}

// TransferBetweenAccounts transfers money between accounts with proper locking to prevent race conditions
func (s *AccountServiceImpl) TransferBetweenAccounts(fromAccountID, toAccountID string, amount float64) (*types.TransferResult, error) {
	// Use a transaction with row locking to prevent race conditions
	result := &types.TransferResult{}

	// Get source and destination account details for transaction records
	sourceAccount, err := s.GetAccountWithDetailByID(fromAccountID)
	if err != nil {
		logger.Error("Failed to get source account details", zap.String("account_id", fromAccountID), zap.Error(err))
		return nil, err
	}

	destAccount, err := s.GetAccountWithDetailByID(toAccountID)
	if err != nil {
		logger.Error("Failed to get destination account details", zap.String("account_id", toAccountID), zap.Error(err))
		return nil, err
	}

	// Begin a database transaction that encompasses both the fund transfer and transaction record creation
	err = s.txProvider.Transact(func(adapters repositories.Adapters) error {
		// Transfer funds within the transaction
		transferErr := adapters.AccountRepository.TransferFunds(fromAccountID, toAccountID, amount, func(sourceBalance, destBalance float64) (*types.TransferResult, error) {
			// Check if source account has sufficient funds
			if sourceBalance < amount {
				return nil, ErrInsufficientFunds
			}

			// Calculate the new balances
			result.SourceBalance = sourceBalance - amount
			result.DestinationBalance = destBalance + amount

			return result, nil
		})

		if transferErr != nil {
			// If transfer fails, the transaction will be rolled back
			return transferErr
		}

		// Create withdrawal transaction record for source account
		withdrawalTx := &models.Transaction{
			BaseModel:       &models.BaseModel{},
			TransactionID:   uuid.New().String(),
			UserID:          sourceAccount.UserID,
			Name:            "Transfer to " + destAccount.AccountNumber,
			IsBank:          true,
			Amount:          amount,
			TransactionType: string(models.Transfer),
			AccountID:       fromAccountID,
		}

		// Create deposit transaction record for destination account
		depositTx := &models.Transaction{
			BaseModel:       &models.BaseModel{},
			TransactionID:   uuid.New().String(),
			UserID:          destAccount.UserID,
			Name:            "Transfer from " + sourceAccount.AccountNumber,
			IsBank:          true,
			Amount:          amount,
			TransactionType: string(models.Transfer),
			AccountID:       toAccountID,
		}

		// Save the transaction records within the same database transaction
		if err := adapters.TransactionRepository.Create(withdrawalTx); err != nil {
			logger.Error("Failed to create withdrawal transaction record",
				zap.String("from_account", fromAccountID),
				zap.String("to_account", toAccountID),
				zap.Error(err))
			// Return error to trigger rollback
			return err
		}

		if err := adapters.TransactionRepository.Create(depositTx); err != nil {
			logger.Error("Failed to create deposit transaction record",
				zap.String("from_account", fromAccountID),
				zap.String("to_account", toAccountID),
				zap.Error(err))
			// Return error to trigger rollback
			return err
		}

		return nil
	})

	if err != nil {
		// If any part of the transaction failed, return the error
		return nil, err
	}

	return result, nil
}

// DeleteAccount marks an account as deleted without removing it
func (s *AccountServiceImpl) DeleteAccount(accountID string) error {
	return s.accountRepository.DeleteAccount(accountID)
}
