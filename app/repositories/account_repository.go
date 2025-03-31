package repositories

import (
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/pkg/types"
	"sort"
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
	GetAccountWithDetailByID(accountID string) (*models.AccountWithDetails, error)
	GetAccountsWithDetailByUserID(userID string) ([]*models.AccountWithDetails, error)

	// Update operations
	UpdateAccountByID(accountID, userID string, updateFn func(account *models.AccountWithDetails) (bool, error)) error
	UnSetMainAccount(userID string) error
	SetMainAccount(accountID, userID string) error
	UpdateAccount(account *models.Account) error
	UpdateAccountDetail(detail *models.AccountDetail) error
	UpdateAccountBalance(accountID string, updateFn func(currentBalance float64) (float64, error)) error

	// Transfer operations
	TransferFunds(fromAccountID, toAccountID string, amount float64,
		updateFn func(sourceBalance, destBalance float64) (*types.TransferResult, error)) error

	// Create operations
	CreateAccount(account *models.AccountWithDetails) error

	// Delete operations
	DeleteAccount(accountID string) error
}

// AccountRepositoryImpl implements AccountRepository
type AccountRepositoryImpl struct {
	DB DB
}

// NewAccountRepository creates a new instance of AccountRepository
func NewAccountRepository(db DB) AccountRepository {
	return &AccountRepositoryImpl{
		DB: db,
	}
}

// GetAccountWithDetailByID retrieves a complete account with all related information by ID
func (r *AccountRepositoryImpl) GetAccountWithDetailByID(accountID string) (*models.AccountWithDetails, error) {
	account := &models.AccountWithDetails{}

	query := `
		SELECT 
			a.account_id, a.user_id, a.type, a.currency, a.account_number, a.issuer, a.created_at, a.updated_at, a.deleted_at,
			d.color, d.is_main_account, d.progress,
			b.amount
		FROM 
			accounts a
		LEFT JOIN 
			account_details d ON a.account_id = d.account_id
		LEFT JOIN 
			account_balances b ON a.account_id = b.account_id
		WHERE 
			a.account_id = ? AND a.deleted_at IS NULL
	`

	err := r.DB.Get(account, query, accountID)
	if err != nil {
		return nil, err
	}

	// Get flags separately as they are multiple
	flags, err := r.GetAccountFlagsByAccountID(accountID)
	if err == nil {
		account.Flags = flags
	}

	return account, nil
}

// GetAccountsWithDetailByUserID retrieves all complete accounts with related information for a user
func (r *AccountRepositoryImpl) GetAccountsWithDetailByUserID(userID string) ([]*models.AccountWithDetails, error) {
	// First, get all accounts with their basic details
	query := `
		SELECT 
			a.account_id, a.user_id, a.type, a.currency, a.account_number, a.issuer, a.created_at, a.updated_at,
			d.color, d.is_main_account, d.progress,
			b.amount,
			f.flag_id, f.flag_type, f.flag_value
		FROM 
			accounts a
		LEFT JOIN 
			account_details d ON a.account_id = d.account_id
		LEFT JOIN 
			account_balances b ON a.account_id = b.account_id
		LEFT JOIN
			account_flags f ON a.account_id = f.account_id AND f.deleted_at IS NULL
		WHERE 
			a.user_id = ? AND a.deleted_at IS NULL
	`

	// Execute the query
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map to store accounts by ID to avoid duplicates
	accountsMap := make(map[string]*models.AccountWithDetails)

	// Process the results
	for rows.Next() {
		var account models.AccountWithDetails
		var flag models.AccountFlag

		err := rows.Scan(
			&account.AccountID, &account.UserID, &account.Type, &account.Currency,
			&account.AccountNumber, &account.Issuer, &account.CreatedAt, &account.UpdatedAt,
			&account.Color, &account.IsMainAccount, &account.Progress,
			&account.Amount, &flag.FlagID, &flag.FlagType, &flag.FlagValue,
		)
		if err != nil {
			return nil, err
		}

		// Check if we already have this account in our map
		existingAccount, exists := accountsMap[account.AccountID]
		if !exists {
			// Initialize flags slice
			account.Flags = make([]*models.AccountFlag, 0)
			accountsMap[account.AccountID] = &account
			existingAccount = &account
		}

		// Add flag if it exists
		if flag.FlagID != 0 {
			flag.AccountID = account.AccountID
			flag.UserID = account.UserID

			// Add flag to the account
			existingAccount.Flags = append(existingAccount.Flags, &flag)
		}
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Convert map to slice
	accounts := make([]*models.AccountWithDetails, 0, len(accountsMap))
	for _, account := range accountsMap {
		accounts = append(accounts, account)
	}

	// Sort by created_at DESC to maintain the order
	sort.Slice(accounts, func(i, j int) bool {
		return accounts[i].CreatedAt.After(accounts[j].CreatedAt)
	})

	return accounts, nil
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

// UpdateAccountByID updates an account with all related details
// Will not update is_main_account in this function, will use another function to prevent confusion since is main account need to set all other main account to false
func (r *AccountRepositoryImpl) UpdateAccountByID(accountID, userID string, updateFn func(account *models.AccountWithDetails) (bool, error)) error {
	return runInTx(r.DB, func(tx *sqlx.Tx) error {
		// Get the existing account with details
		account := &models.AccountWithDetails{}
		query := `
			SELECT 
				a.account_id, a.user_id, a.type, a.currency, a.account_number, a.issuer, a.updated_at,
				d.color, d.progress,
			FROM 
				accounts a
			LEFT JOIN 
				account_details d ON a.account_id = d.account_id
			WHERE 
				a.account_id = ? AND a.deleted_at IS NULL
			FOR UPDATE
		`
		err := tx.Get(account, query, accountID)
		if err != nil {
			return err
		}

		// Get flags separately as they are multiple
		flags, err := r.GetAccountFlagsByAccountID(accountID)
		if err == nil {
			account.Flags = flags
		}

		// Apply the update function to modify the account
		updated, err := updateFn(account)
		if err != nil {
			return err
		}

		// If no changes were made, we can return early
		if !updated {
			return nil
		}

		now := time.Now()
		account.UpdatedAt = now

		// Use a single query with joins to update all related tables
		updateQuery := `
			UPDATE accounts a
			JOIN account_details d ON a.account_id = d.account_id
			SET 
				a.type = ?, 
				a.currency = ?,
				a.account_number = ?,
				a.issuer = ?,
				a.updated_at = ?,
				d.color = ?,
				d.progress = ?,
			WHERE 
				a.account_id = ? AND a.deleted_at IS NULL
		`
		_, err = tx.Exec(
			updateQuery,
			account.Type,
			account.Currency,
			account.AccountNumber,
			account.Issuer,
			now,
			account.Color,
			account.Progress,
			account.AccountID,
		)
		if err != nil {
			return err
		}

		// Handle flags updates if needed
		// delete and recreate all flags
		if len(account.Flags) > 0 {
			// First delete existing flags
			deleteQuery := `DELETE FROM account_flags WHERE account_id = ?`
			_, err = tx.Exec(deleteQuery, account.AccountID)
			if err != nil {
				return err
			}

			// Then insert new flags
			for _, flag := range account.Flags {
				flag.AccountID = account.AccountID
				flag.UserID = account.UserID
				flag.CreatedAt = now
				flag.UpdatedAt = now

				insertQuery := `INSERT INTO account_flags (account_id, user_id, flag_type, flag_value, created_at, updated_at) 
							   VALUES (?, ?, ?, ?, ?, ?)`
				_, err = tx.Exec(
					insertQuery,
					flag.AccountID,
					flag.UserID,
					flag.FlagType,
					flag.FlagValue,
					flag.CreatedAt,
					flag.UpdatedAt,
				)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// SetMainAccount sets a specific account as the main account
func (r *AccountRepositoryImpl) UnSetMainAccount(userID string) error {
	updatedAt := time.Now()

	query := `UPDATE account_details 
              SET is_main_account = ?, updated_at = ? 
              WHERE user_id = ? AND deleted_at IS NULL`
	_, err := r.DB.Exec(
		query,
		false,
		updatedAt,
		userID,
	)
	return err
}

// SetMainAccount sets a specific account as the main account
func (r *AccountRepositoryImpl) SetMainAccount(accountID, userID string) error {
	updatedAt := time.Now()

	query := `UPDATE account_details 
              SET is_main_account = ?, updated_at = ? 
              WHERE account_id = ? AND user_id = ? AND deleted_at IS NULL`
	_, err := r.DB.Exec(
		query,
		true,
		updatedAt,
		accountID,
		userID,
	)
	return err
}

// UpdateAccount updates an existing account
func (r *AccountRepositoryImpl) UpdateAccount(account *models.Account) error {
	account.UpdatedAt = time.Now()

	query := `UPDATE accounts 
              SET user_id = ?, type = ?, currency = ?, account_number = ?, issuer = ?, updated_at = ? 
              WHERE account_id = ? AND deleted_at IS NULL`
	_, err := r.DB.Exec(
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

// UpdateAccountDetail updates existing account details
func (r *AccountRepositoryImpl) UpdateAccountDetail(detail *models.AccountDetail) error {
	detail.UpdatedAt = time.Now()

	query := `UPDATE account_details 
              SET user_id = ?, color = ?, is_main_account = ?, progress = ?, updated_at = ? 
              WHERE account_id = ? AND deleted_at IS NULL`
	_, err := r.DB.Exec(
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

// UpdateAccountBalance updates an account balance with proper locking to prevent race conditions
func (r *AccountRepositoryImpl) UpdateAccountBalance(accountID string, updateFn func(currentBalance float64) (float64, error)) error {
	return runInTx(r.DB, func(tx *sqlx.Tx) error {
		// Get the current balance with a row lock
		var currentBalance float64
		query := `SELECT amount FROM account_balances WHERE account_id = ? FOR UPDATE`
		err := tx.Get(&currentBalance, query, accountID)
		if err != nil {
			return err
		}

		// Apply the update function
		newBalance, err := updateFn(currentBalance)
		if err != nil {
			return err
		}

		// Update the balance
		updateQuery := `UPDATE account_balances SET amount = ? WHERE account_id = ?`
		_, err = tx.Exec(updateQuery, newBalance, accountID)
		return err
	})
}

// CreateAccount adds a new account with all its details
func (r *AccountRepositoryImpl) CreateAccount(account *models.AccountWithDetails) error {
	return runInTx(r.DB, func(tx *sqlx.Tx) error {
		now := time.Now()
		account.CreatedAt = now
		account.UpdatedAt = now

		var err error

		// Create account
		query := `INSERT INTO accounts (account_id, user_id, type, currency, account_number, issuer, created_at, updated_at) 
				  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
		_, err = tx.Exec(
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
		if err != nil {
			return err
		}

		// Create account details
		query = `INSERT INTO account_details (account_id, user_id, color, is_main_account, progress, created_at, updated_at) 
				VALUES (?, ?, ?, ?, ?, ?, ?)`
		_, err = tx.Exec(
			query,
			account.AccountID,
			account.UserID,
			account.Color,
			account.IsMainAccount,
			account.Progress,
			account.CreatedAt,
			account.UpdatedAt,
		)
		if err != nil {
			return err
		}

		// Create account balance
		query = `INSERT INTO account_balances (account_id, user_id, amount, created_at, updated_at) 
				VALUES (?, ?, ?, ?, ?)`
		_, err = tx.Exec(
			query,
			account.AccountID,
			account.UserID,
			account.Amount,
			account.CreatedAt,
			account.UpdatedAt,
		)
		if err != nil {
			return err
		}

		// Create account flags if any
		if len(account.Flags) > 0 {
			for _, flag := range account.Flags {
				flag.AccountID = account.AccountID
				flag.UserID = account.UserID
				flag.CreatedAt = now
				flag.UpdatedAt = now

				query = `INSERT INTO account_flags (account_id, user_id, flag_type, flag_value, created_at, updated_at) 
						VALUES (?, ?, ?, ?, ?, ?)`
				_, err = tx.Exec(
					query,
					flag.AccountID,
					flag.UserID,
					flag.FlagType,
					flag.FlagValue,
					flag.CreatedAt,
					flag.UpdatedAt,
				)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// DeleteAccount marks an account as deleted without removing it
func (r *AccountRepositoryImpl) DeleteAccount(accountID string) error {
	now := time.Now()
	query := `UPDATE accounts SET deleted_at = ? WHERE account_id = ? AND deleted_at IS NULL`
	_, err := r.DB.Exec(query, now, accountID)
	return err
}

// TransferFunds transfers funds between accounts with proper locking to prevent race conditions
func (r *AccountRepositoryImpl) TransferFunds(fromAccountID, toAccountID string, amount float64,
	updateFn func(sourceBalance, destBalance float64) (*types.TransferResult, error)) error {

	return runInTx(r.DB, func(tx *sqlx.Tx) error {
		// Lock both accounts in a consistent order to prevent deadlocks
		// Always lock the account with the smaller ID first
		firstLockID, secondLockID := fromAccountID, toAccountID
		if fromAccountID > toAccountID {
			firstLockID, secondLockID = toAccountID, fromAccountID
		}

		// Get the balances with row locks
		var firstBalance, secondBalance float64
		var sourceBalance, destBalance float64

		// Lock first account
		query := `SELECT account_id, amount FROM account_balances WHERE account_id = ? FOR UPDATE`
		var firstAccountID string
		err := tx.QueryRow(query, firstLockID).Scan(&firstAccountID, &firstBalance)
		if err != nil {
			return err
		}

		// Lock second account
		var secondAccountID string
		err = tx.QueryRow(query, secondLockID).Scan(&secondAccountID, &secondBalance)
		if err != nil {
			return err
		}

		// Map the balances to source and destination
		if firstAccountID == fromAccountID {
			sourceBalance = firstBalance
			destBalance = secondBalance
		} else {
			sourceBalance = secondBalance
			destBalance = firstBalance
		}

		// Apply the update function
		result, err := updateFn(sourceBalance, destBalance)
		if err != nil {
			return err
		}

		// Update the source account balance
		updateQuery := `UPDATE account_balances SET amount = ? WHERE account_id = ?`
		_, err = tx.Exec(updateQuery, result.SourceBalance, fromAccountID)
		if err != nil {
			return err
		}

		// Update the destination account balance
		_, err = tx.Exec(updateQuery, result.DestinationBalance, toAccountID)
		return err
	})
}
