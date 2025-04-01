package services_test

import (
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/app/repositories"
	"backend-developer-assignment/app/services"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransferRaceCondition(t *testing.T) {
	db := setupTestDB()
	txProvider := repositories.NewTransactionProvider(db)
	accountRepo := repositories.NewAccountRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)
	service := services.NewAccountService(accountRepo, transactionRepo, txProvider)

	// Create source account with initial balance
	sourceInitAmount := 10000.0
	sourceAccount := &models.AccountWithDetails{
		AccountID:     "source-account",
		UserID:        "test-user",
		Type:          "saving-account",
		Currency:      "THB",
		Issuer:        "test-issuer",
		IsMainAccount: true,
		Amount:        sourceInitAmount,
		AccountNumber: "SRC-123456",
	}
	err := accountRepo.CreateAccount(sourceAccount)
	assert.NoError(t, err, "Failed to create source account")

	// Create destination account
	destInitAmount := 0.0
	destAccount := &models.AccountWithDetails{
		AccountID:     "dest-account",
		UserID:        "test-user",
		Type:          "saving-account",
		Currency:      "THB",
		Issuer:        "test-issuer",
		IsMainAccount: false,
		Amount:        destInitAmount,
		AccountNumber: "DST-654321",
	}
	err = accountRepo.CreateAccount(destAccount)
	assert.NoError(t, err, "Failed to create destination account")

	var wg sync.WaitGroup
	numWorkers := 10 // Simulate 10 concurrent transfers
	transferAmount := 100.0
	expectedSourceBalance := sourceInitAmount - (float64(numWorkers) * transferAmount)
	expectedDestBalance := destInitAmount + (float64(numWorkers) * transferAmount)

	// Create a channel to collect errors
	errChan := make(chan error, numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			_, err := service.TransferBetweenAccounts("source-account", "dest-account", transferAmount)
			if err != nil {
				t.Logf("Worker %d: Transfer failed: %v", workerID, err)
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

	// Check the final balances
	var finalSourceBalance float64
	err = db.Get(&finalSourceBalance, `SELECT amount FROM account_balances WHERE account_id = 'source-account'`)
	assert.NoError(t, err, "Failed to get source account balance")

	var finalDestBalance float64
	err = db.Get(&finalDestBalance, `SELECT amount FROM account_balances WHERE account_id = 'dest-account'`)
	assert.NoError(t, err, "Failed to get destination account balance")

	t.Logf("Initial source amount: %.2f", sourceInitAmount)
	t.Logf("Initial destination amount: %.2f", destInitAmount)
	t.Logf("Expected final source balance: %.2f", expectedSourceBalance)
	t.Logf("Expected final destination balance: %.2f", expectedDestBalance)
	t.Logf("Actual final source balance: %.2f", finalSourceBalance)
	t.Logf("Actual final destination balance: %.2f", finalDestBalance)
	t.Logf("Number of failed transfers: %d", errorCount)

	// If there were no errors, the final balances should match our expectations
	if errorCount == 0 {
		assert.Equal(t, expectedSourceBalance, finalSourceBalance, "Source balance mismatch, possible race condition")
		assert.Equal(t, expectedDestBalance, finalDestBalance, "Destination balance mismatch, possible race condition")
	} else {
		// If there were errors, we should still have valid balances
		assert.GreaterOrEqual(t, finalSourceBalance, 0.0, "Source balance should not be negative")

		// Calculate expected balances with errors
		successfulTransfers := numWorkers - errorCount
		expectedSourceWithErrors := sourceInitAmount - (float64(successfulTransfers) * transferAmount)
		expectedDestWithErrors := destInitAmount + (float64(successfulTransfers) * transferAmount)

		assert.Equal(t, expectedSourceWithErrors, finalSourceBalance, "Source balance doesn't match expected value with errors")
		assert.Equal(t, expectedDestWithErrors, finalDestBalance, "Destination balance doesn't match expected value with errors")
	}

	// Verify that the sum of both accounts remains constant (conservation of money)
	totalBefore := sourceInitAmount + destInitAmount
	totalAfter := finalSourceBalance + finalDestBalance
	assert.Equal(t, totalBefore, totalAfter, "Total money in the system should remain constant")
}
