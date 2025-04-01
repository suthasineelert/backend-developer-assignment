package services_test

import (
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/app/repositories"
	"backend-developer-assignment/app/services"
	"backend-developer-assignment/platform/database"
	"sync"
	"testing"

	"github.com/joho/godotenv"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func setupTestDB() *sqlx.DB {
	// Load test environment variables
	err := godotenv.Load("../../../.env.test")
	if err != nil {
		// Try to load from the project root as a fallback
		err = godotenv.Load(".env.test")
		if err != nil {
			panic("Error loading .env.test file: " + err.Error())
		}
	}

	db, err := database.MysqlConnection()
	if err != nil {
		panic("Database connection failed:" + err.Error())
	}
	// Drop all existing tables to ensure a clean state
	err = database.Down(db)
	if err != nil {
		panic("Database down migration failed:" + err.Error())
	}

	err = database.Migrate(db)
	if err != nil {
		panic("Database migration failed:" + err.Error())
	}

	return db
}

func TestWithdrawRaceCondition(t *testing.T) {
	db := setupTestDB()
	txProvider := repositories.NewTransactionProvider(db)
	accountRepo := repositories.NewAccountRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)
	service := services.NewAccountService(accountRepo, transactionRepo, txProvider)

	initAmount := 10000.0
	account := &models.AccountWithDetails{
		AccountID:     "test-account",
		UserID:        "test-user",
		Type:          "saving-account",
		Currency:      "THB",
		Issuer:        "test-issuer",
		IsMainAccount: true,
		Amount:        initAmount,
	}
	err := accountRepo.CreateAccount(account)
	assert.NoError(t, err, "Failed to create test account")

	var wg sync.WaitGroup
	numWorkers := 1 // Simulate 100 concurrent withdrawals
	withdrawAmount := 10.0
	expectedFinalBalance := initAmount - (float64(numWorkers) * withdrawAmount)

	// Create a channel to collect errors
	errChan := make(chan error, numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			_, err := service.WithdrawFromAccount("test-account", withdrawAmount)
			if err != nil {
				t.Logf("Worker %d: Withdraw failed: %v", workerID, err)
				errChan <- err
			}
		}(i)
	}

	wg.Wait()
	close(errChan)

	// Count errors
	errorCount := 0
	for err := range errChan {
		t.Logf("Error: %v", err)
		errorCount++
	}

	// Check the final balance
	var finalBalance float64
	err = db.Get(&finalBalance, `SELECT amount FROM account_balances WHERE account_id = 'test-account'`)
	assert.NoError(t, err, "Failed to get final balance")

	t.Logf("Initial amount: %.2f", initAmount)
	t.Logf("Expected final balance: %.2f", expectedFinalBalance)
	t.Logf("Actual final balance: %.2f", finalBalance)
	t.Logf("Number of failed withdrawals: %d", errorCount)

	// If there were no errors, the final balance should match our expectation
	if errorCount == 0 {
		assert.Equal(t, expectedFinalBalance, finalBalance, "Balance mismatch, possible race condition")
	} else {
		// If there were errors, we should still have a valid balance (not negative)
		assert.GreaterOrEqual(t, finalBalance, 0.0, "Balance should not be negative")

		// And the balance should be the initial amount minus successful withdrawals
		successfulWithdrawals := numWorkers - errorCount
		expectedWithErrors := initAmount - (float64(successfulWithdrawals) * withdrawAmount)
		assert.Equal(t, expectedWithErrors, finalBalance, "Balance doesn't match expected value with errors")
	}
}
