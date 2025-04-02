package services_test

import (
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/app/repositories"
	"backend-developer-assignment/app/services"
	mocks "backend-developer-assignment/pkg/mocks/repositories"
	"backend-developer-assignment/pkg/types"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// AccountServiceTestSuite defines the test suite
type AccountServiceTestSuite struct {
	suite.Suite
	accountRepository     *mocks.AccountRepository
	transactionRepository *mocks.TransactionRepository
	txProvider            *mocks.TxProvider
	service               services.AccountService
}

// SetupTest runs before each test
func (s *AccountServiceTestSuite) SetupTest() {
	s.accountRepository = new(mocks.AccountRepository)
	s.transactionRepository = new(mocks.TransactionRepository)
	s.txProvider = new(mocks.TxProvider)
	s.service = services.NewAccountService(s.accountRepository, s.transactionRepository, s.txProvider)
}

// TestGetAccountByID tests the GetAccountByID function
func (s *AccountServiceTestSuite) TestGetAccountByID() {
	now := time.Now()
	testCases := []struct {
		name            string
		accountID       string
		mockAccount     *models.Account
		mockError       error
		expectedAccount *models.Account
		expectedError   error
	}{
		{
			name:      "Success - Valid Account",
			accountID: "acc-123",
			mockAccount: &models.Account{
				AccountID: "acc-123",
				UserID:    "user-123",
				Type:      "saving-account",
				BaseModel: &models.BaseModel{
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			mockError: nil,
			expectedAccount: &models.Account{
				AccountID: "acc-123",
				UserID:    "user-123",
				Type:      "saving-account",
				BaseModel: &models.BaseModel{
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			expectedError: nil,
		},
		{
			name:            "Failure - Account Not Found",
			accountID:       "nonexistent-account",
			mockAccount:     nil,
			mockError:       errors.New("account not found"),
			expectedAccount: nil,
			expectedError:   errors.New("account not found"),
		},
		{
			name:            "Failure - Database Error",
			accountID:       "invalid-account-id",
			mockAccount:     nil,
			mockError:       errors.New("database connection failed"),
			expectedAccount: nil,
			expectedError:   errors.New("database connection failed"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Mock the repository method
			s.accountRepository.On("GetAccountByID", tc.accountID).Return(tc.mockAccount, tc.mockError).Once()

			// Call the service method
			account, err := s.service.GetAccountByID(tc.accountID)

			// Assert results
			if tc.expectedError != nil {
				assert.Error(s.T(), err)
				assert.Equal(s.T(), tc.expectedError.Error(), err.Error())
				assert.Nil(s.T(), account)
			} else {
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), tc.expectedAccount, account)
			}

			// Verify expected method calls
			s.accountRepository.AssertExpectations(s.T())
		})
	}
}

// TestGetAccountsByUserID tests the GetAccountsByUserID function
func (s *AccountServiceTestSuite) TestGetAccountsByUserID() {
	userID := "user-123"
	now := time.Now()

	testCases := []struct {
		name             string
		mockAccounts     []*models.Account
		mockError        error
		expectedAccounts []*models.Account
		expectedError    error
	}{
		{
			name: "Success - Multiple Accounts",
			mockAccounts: []*models.Account{
				{
					AccountID: "acc-123",
					UserID:    userID,
					Type:      "saving-account",
					BaseModel: &models.BaseModel{
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
				{
					AccountID: "acc-456",
					UserID:    userID,
					Type:      "credit-loan",
					BaseModel: &models.BaseModel{
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
			},
			mockError: nil,
			expectedAccounts: []*models.Account{
				{
					AccountID: "acc-123",
					UserID:    userID,
					Type:      "saving-account",
					BaseModel: &models.BaseModel{
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
				{
					AccountID: "acc-456",
					UserID:    userID,
					Type:      "credit-loan",
					BaseModel: &models.BaseModel{
						CreatedAt: now,
						UpdatedAt: now,
					},
				},
			},
			expectedError: nil,
		},
		{
			name:             "Success - No Accounts",
			mockAccounts:     []*models.Account{},
			mockError:        nil,
			expectedAccounts: []*models.Account{},
			expectedError:    nil,
		},
		{
			name:             "Failure - Database Error",
			mockAccounts:     nil,
			mockError:        errors.New("database connection failed"),
			expectedAccounts: nil,
			expectedError:    errors.New("database connection failed"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Mock the repository method
			s.accountRepository.On("GetAccountsByUserID", userID).Return(tc.mockAccounts, tc.mockError).Once()

			// Call the service method
			accounts, err := s.service.GetAccountsByUserID(userID)

			// Assert results
			if tc.expectedError != nil {
				assert.Error(s.T(), err)
				assert.Equal(s.T(), tc.expectedError.Error(), err.Error())
				assert.Nil(s.T(), accounts)
			} else {
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), tc.expectedAccounts, accounts)
			}

			// Verify expected method calls
			s.accountRepository.AssertExpectations(s.T())
		})
	}
}

// TestGetAccountWithDetailByID tests the GetAccountWithDetailByID function
func (s *AccountServiceTestSuite) TestGetAccountWithDetailByID() {
	now := time.Now()
	testCases := []struct {
		name            string
		accountID       string
		mockAccount     *models.AccountWithDetails
		mockError       error
		expectedAccount *models.AccountWithDetails
		expectedError   error
	}{
		{
			name:      "Success - Valid Account with Details",
			accountID: "acc-123",
			mockAccount: &models.AccountWithDetails{
				AccountID:     "acc-123",
				UserID:        "user-123",
				Type:          "saving-account",
				Currency:      "USD",
				AccountNumber: "123456789",
				Issuer:        "Bank",
				Color:         "#FF0000",
				Progress:      75,
				Amount:        1000.50,
				CreatedAt:     now,
				UpdatedAt:     now,
			},
			mockError: nil,
			expectedAccount: &models.AccountWithDetails{
				AccountID:     "acc-123",
				UserID:        "user-123",
				Type:          "saving-account",
				Currency:      "USD",
				AccountNumber: "123456789",
				Issuer:        "Bank",
				Color:         "#FF0000",
				Progress:      75,
				Amount:        1000.50,
				CreatedAt:     now,
				UpdatedAt:     now,
			},
			expectedError: nil,
		},
		{
			name:            "Failure - Account Not Found",
			accountID:       "nonexistent-account",
			mockAccount:     nil,
			mockError:       errors.New("account not found"),
			expectedAccount: nil,
			expectedError:   errors.New("account not found"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Mock the repository method
			s.accountRepository.On("GetAccountWithDetailByID", tc.accountID).Return(tc.mockAccount, tc.mockError).Once()

			// Call the service method
			account, err := s.service.GetAccountWithDetailByID(tc.accountID)

			// Assert results
			if tc.expectedError != nil {
				assert.Error(s.T(), err)
				assert.Equal(s.T(), tc.expectedError.Error(), err.Error())
				assert.Nil(s.T(), account)
			} else {
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), tc.expectedAccount, account)
			}

			// Verify expected method calls
			s.accountRepository.AssertExpectations(s.T())
		})
	}
}

// TestGetAccountsWithDetailByUserID tests the GetAccountsWithDetailByUserID function
func (s *AccountServiceTestSuite) TestGetAccountsWithDetailByUserID() {
	userID := "user-123"
	now := time.Now()

	testCases := []struct {
		name             string
		mockAccounts     []*models.AccountWithDetails
		mockError        error
		expectedAccounts []*models.AccountWithDetails
		expectedError    error
	}{
		{
			name: "Success - Multiple Accounts",
			mockAccounts: []*models.AccountWithDetails{
				{
					AccountID:     "acc-123",
					UserID:        userID,
					Type:          "saving-account",
					Currency:      "USD",
					AccountNumber: "123456789",
					Issuer:        "Bank A",
					Color:         "#FF0000",
					Progress:      75,
					Amount:        1000.50,
					CreatedAt:     now,
					UpdatedAt:     now,
				},
				{
					AccountID:     "acc-456",
					UserID:        userID,
					Type:          "credit-loan",
					Currency:      "EUR",
					AccountNumber: "987654321",
					Issuer:        "Bank B",
					Color:         "#00FF00",
					Progress:      50,
					Amount:        2500.75,
					CreatedAt:     now,
					UpdatedAt:     now,
				},
			},
			mockError: nil,
			expectedAccounts: []*models.AccountWithDetails{
				{
					AccountID:     "acc-123",
					UserID:        userID,
					Type:          "saving-account",
					Currency:      "USD",
					AccountNumber: "123456789",
					Issuer:        "Bank A",
					Color:         "#FF0000",
					Progress:      75,
					Amount:        1000.50,
					CreatedAt:     now,
					UpdatedAt:     now,
				},
				{
					AccountID:     "acc-456",
					UserID:        userID,
					Type:          "credit-loan",
					Currency:      "EUR",
					AccountNumber: "987654321",
					Issuer:        "Bank B",
					Color:         "#00FF00",
					Progress:      50,
					Amount:        2500.75,
					CreatedAt:     now,
					UpdatedAt:     now,
				},
			},
			expectedError: nil,
		},
		{
			name:             "Success - No Accounts",
			mockAccounts:     []*models.AccountWithDetails{},
			mockError:        nil,
			expectedAccounts: []*models.AccountWithDetails{},
			expectedError:    nil,
		},
		{
			name:             "Failure - Database Error",
			mockAccounts:     nil,
			mockError:        errors.New("database connection failed"),
			expectedAccounts: nil,
			expectedError:    errors.New("database connection failed"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Mock the repository method
			s.accountRepository.On("GetAccountsWithDetailByUserID", userID).Return(tc.mockAccounts, tc.mockError).Once()

			// Call the service method
			accounts, err := s.service.GetAccountsWithDetailByUserID(userID)

			// Assert results
			if tc.expectedError != nil {
				assert.Error(s.T(), err)
				assert.Equal(s.T(), tc.expectedError.Error(), err.Error())
				assert.Nil(s.T(), accounts)
			} else {
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), tc.expectedAccounts, accounts)
			}

			// Verify expected method calls
			s.accountRepository.AssertExpectations(s.T())
		})
	}
}

// TestCreateAccountWithDetailsWithProvidedID tests creating an account with a provided ID
func (s *AccountServiceTestSuite) TestCreateAccountWithDetailsWithProvidedID() {
	// Create test data
	accountID := "test-account-id"
	userID := "test-user-id"
	accountWithDetails := &models.AccountWithDetails{
		AccountID:     accountID,
		UserID:        userID,
		Type:          "savings",
		Currency:      "USD",
		AccountNumber: "123456789",
		Issuer:        "Test Bank",
		Amount:        1000.00,
		IsMainAccount: false,
		Color:         "#FF5733",
		Progress:      75,
	}

	// Mock repository behavior
	s.accountRepository.On("CreateAccount", accountWithDetails).Return(nil).Once()

	// Call the service method
	err := s.service.CreateAccountWithDetails(accountWithDetails)

	// Assert results
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), accountID, accountWithDetails.AccountID)
	s.accountRepository.AssertExpectations(s.T())
}

// TestCreateAccountWithDetailsWithGeneratedID tests creating an account with a generated ID
func (s *AccountServiceTestSuite) TestCreateAccountWithDetailsWithGeneratedID() {
	// Create test data
	userID := "test-user-id"
	accountWithDetails := &models.AccountWithDetails{
		AccountID:     "",
		UserID:        userID,
		Type:          "savings",
		Currency:      "USD",
		AccountNumber: "123456789",
		Issuer:        "Test Bank",
		Amount:        1000.00,
		IsMainAccount: false,
		Color:         "#FF5733",
		Progress:      75,
	}

	// Mock repository behavior with ID matcher
	s.accountRepository.On("CreateAccount", mock.MatchedBy(func(a *models.AccountWithDetails) bool {
		// Verify that an ID was generated (non-empty)
		return a.AccountID != "" && a.UserID == userID
	})).Return(nil).Once()

	// Call the service method
	err := s.service.CreateAccountWithDetails(accountWithDetails)

	// Assert results
	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), accountWithDetails.AccountID)
	_, err = uuid.Parse(accountWithDetails.AccountID)
	assert.NoError(s.T(), err) // ID should be a valid UUID
	s.accountRepository.AssertExpectations(s.T())
}

// TestCreateAccountWithDetailsError tests creating an account with a repository error
func (s *AccountServiceTestSuite) TestCreateAccountWithDetailsError() {
	// Create test data
	accountID := "test-account-id"
	userID := "test-user-id"
	accountWithDetails := &models.AccountWithDetails{
		AccountID:     accountID,
		UserID:        userID,
		Type:          "savings",
		Currency:      "USD",
		AccountNumber: "123456789",
		Issuer:        "Test Bank",
		Amount:        1000.00,
		IsMainAccount: false,
		Color:         "#FF5733",
		Progress:      75,
	}

	// Mock repository error
	expectedError := errors.New("database error")
	s.accountRepository.On("CreateAccount", accountWithDetails).Return(expectedError).Once()

	// Call the service method
	err := s.service.CreateAccountWithDetails(accountWithDetails)

	// Assert results
	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	s.accountRepository.AssertExpectations(s.T())
}

// TestUpdateAccountWithChanges tests updating an account with changes
func (s *AccountServiceTestSuite) TestUpdateAccountWithChanges() {
	// Create test data
	accountID := "test-account-id"
	userID := "test-user-id"
	account := &models.AccountWithDetails{
		AccountID:     accountID,
		UserID:        userID,
		Type:          "credit-loan",
		Currency:      "EUR",
		AccountNumber: "987654321",
		Issuer:        "New Bank",
		Amount:        1000.00,
		IsMainAccount: false,
		Color:         "#33FF57",
		Progress:      90,
	}

	// Mock repository behavior
	s.accountRepository.On("UpdateAccountByID", accountID, userID, mock.AnythingOfType("func(*models.AccountWithDetails) (bool, error)")).
		Run(func(args mock.Arguments) {
			// Extract the callback function
			callback := args.Get(2).(func(*models.AccountWithDetails) (bool, error))

			// Create a mock existing account
			existingAccount := &models.AccountWithDetails{
				AccountID:     accountID,
				UserID:        userID,
				Type:          "savings",
				Currency:      "USD",
				AccountNumber: "123456789",
				Issuer:        "Test Bank",
				Color:         "#FFF554",
				Progress:      75,
			}

			// Call the callback
			isUpdate, err := callback(existingAccount)

			// Assert callback results
			assert.True(s.T(), isUpdate)
			assert.NoError(s.T(), err)
			assert.Equal(s.T(), "credit-loan", existingAccount.Type)
			assert.Equal(s.T(), "EUR", existingAccount.Currency)
			assert.Equal(s.T(), "987654321", existingAccount.AccountNumber)
			assert.Equal(s.T(), "New Bank", existingAccount.Issuer)
			assert.Equal(s.T(), "#33FF57", existingAccount.Color)
			assert.Equal(s.T(), 90, existingAccount.Progress)
		}).
		Return(nil).Once()

	// Call the service method
	err := s.service.UpdateAccount(account)

	// Assert results
	assert.NoError(s.T(), err)
	s.accountRepository.AssertExpectations(s.T())
}

// TestUpdateAccountWithNoChanges tests updating an account with no changes
func (s *AccountServiceTestSuite) TestUpdateAccountWithNoChanges() {
	// Create test data
	accountID := "test-account-id"
	userID := "test-user-id"
	account := &models.AccountWithDetails{
		AccountID: accountID,
		UserID:    userID,
		// Empty fields should not trigger updates
	}

	// Mock repository behavior
	s.accountRepository.On("UpdateAccountByID", accountID, userID, mock.AnythingOfType("func(*models.AccountWithDetails) (bool, error)")).
		Run(func(args mock.Arguments) {
			// Extract the callback function
			callback := args.Get(2).(func(*models.AccountWithDetails) (bool, error))

			// Create a mock existing account
			existingAccount := &models.AccountWithDetails{
				AccountID:     accountID,
				UserID:        userID,
				Type:          "savings",
				Currency:      "USD",
				AccountNumber: "123456789",
				Issuer:        "Test Bank",
				Color:         "#FF5733",
				Progress:      75,
			}

			// Call the callback
			isUpdate, err := callback(existingAccount)

			// Assert callback results
			assert.False(s.T(), isUpdate)
			assert.NoError(s.T(), err)
			assert.Equal(s.T(), "savings", existingAccount.Type) // Unchanged
			assert.Equal(s.T(), "USD", existingAccount.Currency) // Unchanged
		}).
		Return(nil).Once()

	// Call the service method
	err := s.service.UpdateAccount(account)

	// Assert results
	assert.NoError(s.T(), err)
	s.accountRepository.AssertExpectations(s.T())
}

// TestUpdateAccountError tests updating an account with a repository error
func (s *AccountServiceTestSuite) TestUpdateAccountError() {
	// Create test data
	accountID := "test-account-id"
	userID := "test-user-id"
	account := &models.AccountWithDetails{
		AccountID: accountID,
		UserID:    userID,
		Type:      "checking",
	}

	// Mock repository error
	expectedError := errors.New("database error")
	s.accountRepository.On("UpdateAccountByID", accountID, userID, mock.AnythingOfType("func(*models.AccountWithDetails) (bool, error)")).
		Return(expectedError).Once()

	// Call the service method
	err := s.service.UpdateAccount(account)

	// Assert results
	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	s.accountRepository.AssertExpectations(s.T())
}

// TestSetMainAccountSuccess tests setting an account as the main account successfully
func (s *AccountServiceTestSuite) TestSetMainAccountSuccess() {
	// Create test data
	accountID := "test-account-id"
	userID := "test-user-id"
	account := &models.Account{
		AccountID: accountID,
		UserID:    userID,
	}

	// Mock repository behavior
	s.accountRepository.On("UnSetMainAccount", userID).Return(nil).Once()
	s.accountRepository.On("SetMainAccount", accountID, userID).Return(nil).Once()

	// Call the service method
	err := s.service.SetMainAccount(account)

	// Assert results
	assert.NoError(s.T(), err)
	s.accountRepository.AssertExpectations(s.T())
}

// TestSetMainAccountUnsetError tests setting an account as main with an error during unset
func (s *AccountServiceTestSuite) TestSetMainAccountUnsetError() {
	// Create test data
	accountID := "test-account-id"
	userID := "test-user-id"
	account := &models.Account{
		AccountID: accountID,
		UserID:    userID,
	}

	// Mock repository error
	expectedError := errors.New("database error")
	s.accountRepository.On("UnSetMainAccount", userID).Return(expectedError).Once()

	// Call the service method
	err := s.service.SetMainAccount(account)

	// Assert results
	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	s.accountRepository.AssertExpectations(s.T())
	s.accountRepository.AssertNotCalled(s.T(), "SetMainAccount")
}

// TestSetMainAccountSetError tests setting an account as main with an error during set
func (s *AccountServiceTestSuite) TestSetMainAccountSetError() {
	// Create test data
	accountID := "test-account-id"
	userID := "test-user-id"
	account := &models.Account{
		AccountID: accountID,
		UserID:    userID,
	}

	// Mock repository behavior
	s.accountRepository.On("UnSetMainAccount", userID).Return(nil).Once()

	// Mock repository error
	expectedError := errors.New("database error")
	s.accountRepository.On("SetMainAccount", accountID, userID).Return(expectedError).Once()

	// Call the service method
	err := s.service.SetMainAccount(account)

	// Assert results
	assert.Error(s.T(), err)
	assert.Equal(s.T(), expectedError, err)
	s.accountRepository.AssertExpectations(s.T())
}

func (s *AccountServiceTestSuite) TestTransferBetweenAccounts() {
	fromAccountID := "acc-123"
	toAccountID := "acc-456"
	amount := 100.0

	sourceAccount := &models.AccountWithDetails{
		AccountID:     fromAccountID,
		UserID:        "user-123",
		AccountNumber: "123456789",
	}

	destAccount := &models.AccountWithDetails{
		AccountID:     toAccountID,
		UserID:        "user-456",
		AccountNumber: "987654321",
	}

	expectedResult := &types.TransferResult{
		SourceBalance:      900.0,
		DestinationBalance: 200.0,
	}

	// Mock GetAccountWithDetailByID for source account
	s.accountRepository.On("GetAccountWithDetailByID", fromAccountID).Return(sourceAccount, nil)

	// Mock GetAccountWithDetailByID for destination account
	s.accountRepository.On("GetAccountWithDetailByID", toAccountID).Return(destAccount, nil)

	// Mock Transact method of txProvider
	s.txProvider.On("Transact", mock.AnythingOfType("func(repositories.Adapters) error")).
		Return(nil).
		Run(func(args mock.Arguments) {
			// Extract the transaction function
			txFunc := args.Get(0).(func(adapters repositories.Adapters) error)

			// Create mock adapters
			mockAdapters := repositories.Adapters{
				AccountRepository:     s.accountRepository,
				TransactionRepository: s.transactionRepository,
			}

			// Mock TransferFunds
			s.accountRepository.On("TransferFunds", fromAccountID, toAccountID, amount,
				mock.AnythingOfType("func(float64, float64) (*types.TransferResult, error)")).
				Return(nil).
				Run(func(args mock.Arguments) {
					// Extract and call the update function
					updateFn := args.Get(3).(func(float64, float64) (*types.TransferResult, error))
					result, _ := updateFn(1000.0, 100.0) // Source has 1000, dest has 100

					// Verify the result matches expected
					assert.Equal(s.T(), expectedResult.SourceBalance, result.SourceBalance)
					assert.Equal(s.T(), expectedResult.DestinationBalance, result.DestinationBalance)
				})

			// Mock Create for withdrawal transaction
			s.transactionRepository.On("Create", mock.MatchedBy(func(tx *models.Transaction) bool {
				return tx.AccountID == fromAccountID &&
					tx.Amount == amount &&
					tx.TransactionType == string(models.Transfer)
			})).Return(nil)

			// Mock Create for deposit transaction
			s.transactionRepository.On("Create", mock.MatchedBy(func(tx *models.Transaction) bool {
				return tx.AccountID == toAccountID &&
					tx.Amount == amount &&
					tx.TransactionType == string(models.Transfer)
			})).Return(nil)

			// Execute the transaction function
			err := txFunc(mockAdapters)
			assert.NoError(s.T(), err)
		})

	// Call the service method
	result, err := s.service.TransferBetweenAccounts(fromAccountID, toAccountID, amount)

	// Assert results
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), expectedResult.SourceBalance, result.SourceBalance)
	assert.Equal(s.T(), expectedResult.DestinationBalance, result.DestinationBalance)

	// Verify all mocks were called
	s.accountRepository.AssertExpectations(s.T())
	s.transactionRepository.AssertExpectations(s.T())
	s.txProvider.AssertExpectations(s.T())
}

// TestTransferBetweenAccountsWithInsufficientFunds tests the TransferBetweenAccounts function with insufficient funds
func (s *AccountServiceTestSuite) TestTransferBetweenAccountsWithInsufficientFunds() {
	fromAccountID := "acc-123"
	toAccountID := "acc-456"
	amount := 100.0

	sourceAccount := &models.AccountWithDetails{
		AccountID:     fromAccountID,
		UserID:        "user-123",
		AccountNumber: "123456789",
	}

	destAccount := &models.AccountWithDetails{
		AccountID:     toAccountID,
		UserID:        "user-456",
		AccountNumber: "987654321",
	}

	// Mock GetAccountWithDetailByID for source account
	s.accountRepository.On("GetAccountWithDetailByID", fromAccountID).Return(sourceAccount, nil)

	// Mock GetAccountWithDetailByID for destination account
	s.accountRepository.On("GetAccountWithDetailByID", toAccountID).Return(destAccount, nil)

	// Mock Transact method of txProvider with error
	s.txProvider.On("Transact", mock.AnythingOfType("func(repositories.Adapters) error")).
		Return(services.ErrInsufficientFunds).
		Run(func(args mock.Arguments) {
			// Extract the transaction function
			txFunc := args.Get(0).(func(adapters repositories.Adapters) error)

			// Create mock adapters
			mockAdapters := repositories.Adapters{
				AccountRepository:     s.accountRepository,
				TransactionRepository: s.transactionRepository,
			}

			// Mock TransferFunds with insufficient funds error
			s.accountRepository.On("TransferFunds", fromAccountID, toAccountID, amount,
				mock.AnythingOfType("func(float64, float64) (*types.TransferResult, error)")).
				Return(services.ErrInsufficientFunds).
				Run(func(args mock.Arguments) {
					// Extract and call the update function
					updateFn := args.Get(3).(func(float64, float64) (*types.TransferResult, error))
					_, err := updateFn(50.0, 100.0) // Source has only 50, not enough

					// Verify the error is insufficient funds
					assert.Equal(s.T(), services.ErrInsufficientFunds, err)
				})

			// Execute the transaction function
			err := txFunc(mockAdapters)
			assert.Equal(s.T(), services.ErrInsufficientFunds, err)
		})

	// Call the service method
	result, err := s.service.TransferBetweenAccounts(fromAccountID, toAccountID, amount)

	// Assert results
	assert.Error(s.T(), err)
	assert.Equal(s.T(), services.ErrInsufficientFunds, err)
	assert.Nil(s.T(), result)

	// Verify all mocks were called
	s.accountRepository.AssertExpectations(s.T())
	s.txProvider.AssertExpectations(s.T())
}

// TestTransferBetweenAccountsWithTransactionCreationError tests the TransferBetweenAccounts function with transaction creation error
func (s *AccountServiceTestSuite) TestTransferBetweenAccountsWithTransactionCreationError() {
	fromAccountID := "acc-123"
	toAccountID := "acc-456"
	amount := 100.0

	sourceAccount := &models.AccountWithDetails{
		AccountID:     fromAccountID,
		UserID:        "user-123",
		AccountNumber: "123456789",
	}

	destAccount := &models.AccountWithDetails{
		AccountID:     toAccountID,
		UserID:        "user-456",
		AccountNumber: "987654321",
	}

	txCreationError := errors.New("failed to create transaction record")

	// Mock GetAccountWithDetailByID for source account
	s.accountRepository.On("GetAccountWithDetailByID", fromAccountID).Return(sourceAccount, nil)

	// Mock GetAccountWithDetailByID for destination account
	s.accountRepository.On("GetAccountWithDetailByID", toAccountID).Return(destAccount, nil)

	// Mock Transact method of txProvider with error
	s.txProvider.On("Transact", mock.AnythingOfType("func(repositories.Adapters) error")).
		Return(txCreationError).
		Run(func(args mock.Arguments) {
			// Extract the transaction function
			txFunc := args.Get(0).(func(adapters repositories.Adapters) error)

			// Create mock adapters
			mockAdapters := repositories.Adapters{
				AccountRepository:     s.accountRepository,
				TransactionRepository: s.transactionRepository,
			}

			// Mock TransferFunds success
			s.accountRepository.On("TransferFunds", fromAccountID, toAccountID, amount,
				mock.AnythingOfType("func(float64, float64) (*types.TransferResult, error)")).
				Return(nil).
				Run(func(args mock.Arguments) {
					// Extract and call the update function
					updateFn := args.Get(3).(func(float64, float64) (*types.TransferResult, error))
					_, _ = updateFn(1000.0, 100.0) // Source has 1000, dest has 100
				})

			// Mock Create for withdrawal transaction with error
			s.transactionRepository.On("Create", mock.MatchedBy(func(tx *models.Transaction) bool {
				return tx.AccountID == fromAccountID
			})).Return(txCreationError)

			// Execute the transaction function
			err := txFunc(mockAdapters)
			assert.Equal(s.T(), txCreationError, err)
		})

	// Call the service method
	result, err := s.service.TransferBetweenAccounts(fromAccountID, toAccountID, amount)

	// Assert results
	assert.Error(s.T(), err)
	assert.Equal(s.T(), txCreationError, err)
	assert.Nil(s.T(), result)

	// Verify all mocks were called
	s.accountRepository.AssertExpectations(s.T())
	s.transactionRepository.AssertExpectations(s.T())
	s.txProvider.AssertExpectations(s.T())
}

// TestWithdrawFromAccount tests the WithdrawFromAccount function
func (s *AccountServiceTestSuite) TestWithdrawFromAccount() {
	accountID := "acc-123"
	amount := 50.0

	account := &models.AccountWithDetails{
		AccountID:     accountID,
		UserID:        "user-123",
		Type:          "saving-account",
		Currency:      "USD",
		AccountNumber: "123456789",
	}

	testCases := []struct {
		name            string
		currentBalance  float64
		amount          float64
		expectedError   error
		expectedBalance float64
		txCreateError   error
	}{
		{
			name:            "Success - Sufficient Funds",
			currentBalance:  100.0,
			amount:          50.0,
			expectedError:   nil,
			expectedBalance: 50.0,
			txCreateError:   nil,
		},
		{
			name:            "Failure - Insufficient Funds",
			currentBalance:  30.0,
			amount:          50.0,
			expectedError:   services.ErrInsufficientFunds,
			expectedBalance: 0.0,
			txCreateError:   nil,
		},
		{
			name:            "Failure - Transaction Creation Error",
			currentBalance:  100.0,
			amount:          50.0,
			expectedError:   errors.New("failed to create transaction record"),
			expectedBalance: 50.0,
			txCreateError:   errors.New("failed to create transaction record"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Reset mocks
			s.accountRepository = new(mocks.AccountRepository)
			s.transactionRepository = new(mocks.TransactionRepository)
			s.txProvider = new(mocks.TxProvider)
			s.service = services.NewAccountService(s.accountRepository, s.transactionRepository, s.txProvider)

			// Mock GetAccountWithDetailByID
			s.accountRepository.On("GetAccountWithDetailByID", accountID).Return(account, nil)

			// Mock Transact method of txProvider
			s.txProvider.On("Transact", mock.AnythingOfType("func(repositories.Adapters) error")).
				Return(tc.expectedError).
				Run(func(args mock.Arguments) {
					// Extract the transaction function
					txFunc := args.Get(0).(func(adapters repositories.Adapters) error)

					// Create mock adapters
					mockAdapters := repositories.Adapters{
						AccountRepository:     s.accountRepository,
						TransactionRepository: s.transactionRepository,
					}

					// Mock UpdateAccountBalance
					updateBalanceErr := tc.expectedError
					if tc.txCreateError != nil {
						// For transaction creation error, balance update succeeds
						updateBalanceErr = nil
					}

					s.accountRepository.On("UpdateAccountBalance", accountID,
						mock.AnythingOfType("func(float64) (float64, error)")).
						Return(updateBalanceErr).
						Run(func(args mock.Arguments) {
							// Extract and call the update function
							updateFn := args.Get(1).(func(float64) (float64, error))
							balance, err := updateFn(tc.currentBalance)

							if updateBalanceErr == nil {
								assert.NoError(s.T(), err)
								assert.Equal(s.T(), tc.expectedBalance, balance)
							} else {
								assert.Equal(s.T(), updateBalanceErr, err)
							}
						})

					if tc.expectedError == nil || tc.txCreateError != nil {
						// Mock Create for transaction record
						s.transactionRepository.On("Create", mock.MatchedBy(func(tx *models.Transaction) bool {
							return tx.AccountID == accountID &&
								tx.Amount == amount &&
								tx.TransactionType == string(models.Withdrawal)
						})).Return(tc.txCreateError)
					}

					// Execute the transaction function
					err := txFunc(mockAdapters)
					assert.Equal(s.T(), tc.expectedError, err)
				})

			// Call the service method
			balance, err := s.service.WithdrawFromAccount(accountID, amount)

			// Assert results
			if tc.expectedError != nil {
				assert.Error(s.T(), err)
				assert.Equal(s.T(), tc.expectedError.Error(), err.Error())
				assert.Equal(s.T(), 0.0, balance)
			} else {
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), tc.expectedBalance, balance)
			}

			// Verify all mocks were called
			s.accountRepository.AssertExpectations(s.T())
			s.transactionRepository.AssertExpectations(s.T())
			s.txProvider.AssertExpectations(s.T())
		})
	}
}

// TestDepositToAccount tests the DepositToAccount function
func (s *AccountServiceTestSuite) TestDepositToAccount() {
	accountID := "acc-123"
	amount := 50.0

	account := &models.AccountWithDetails{
		AccountID:     accountID,
		UserID:        "user-123",
		Type:          "saving-account",
		Currency:      "USD",
		AccountNumber: "123456789",
	}

	testCases := []struct {
		name            string
		currentBalance  float64
		amount          float64
		expectedError   error
		expectedBalance float64
	}{
		{
			name:            "Success - Deposit Funds",
			currentBalance:  100.0,
			amount:          50.0,
			expectedError:   nil,
			expectedBalance: 150.0,
		},
		{
			name:            "Failure - Transaction Creation Error",
			currentBalance:  100.0,
			amount:          50.0,
			expectedError:   errors.New("failed to create transaction record"),
			expectedBalance: 150.0,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Reset mocks
			s.accountRepository = new(mocks.AccountRepository)
			s.transactionRepository = new(mocks.TransactionRepository)
			s.txProvider = new(mocks.TxProvider)
			s.service = services.NewAccountService(s.accountRepository, s.transactionRepository, s.txProvider)

			// Mock GetAccountWithDetailByID
			s.accountRepository.On("GetAccountWithDetailByID", accountID).Return(account, nil)

			// Mock Transact method of txProvider
			s.txProvider.On("Transact", mock.AnythingOfType("func(repositories.Adapters) error")).
				Return(tc.expectedError).
				Run(func(args mock.Arguments) {
					// Extract the transaction function
					txFunc := args.Get(0).(func(adapters repositories.Adapters) error)

					// Create mock adapters
					mockAdapters := repositories.Adapters{
						AccountRepository:     s.accountRepository,
						TransactionRepository: s.transactionRepository,
					}

					// Mock UpdateAccountBalance
					s.accountRepository.On("UpdateAccountBalance", accountID,
						mock.AnythingOfType("func(float64) (float64, error)")).
						Return(nil).
						Run(func(args mock.Arguments) {
							// Extract and call the update function
							updateFn := args.Get(1).(func(float64) (float64, error))
							balance, err := updateFn(tc.currentBalance)

							assert.NoError(s.T(), err)
							assert.Equal(s.T(), tc.expectedBalance, balance)
						})

					// Mock Create for transaction record
					if tc.expectedError == nil {
						s.transactionRepository.On("Create", mock.MatchedBy(func(tx *models.Transaction) bool {
							return tx.AccountID == accountID &&
								tx.Amount == amount &&
								tx.TransactionType == string(models.Deposit)
						})).Return(nil)
					} else {
						s.transactionRepository.On("Create", mock.MatchedBy(func(tx *models.Transaction) bool {
							return tx.AccountID == accountID &&
								tx.Amount == amount &&
								tx.TransactionType == string(models.Deposit)
						})).Return(tc.expectedError)
					}

					// Execute the transaction function
					err := txFunc(mockAdapters)
					assert.Equal(s.T(), tc.expectedError, err)
				})

			// Call the service method
			balance, err := s.service.DepositToAccount(accountID, amount)

			// Assert results
			if tc.expectedError != nil {
				assert.Error(s.T(), err)
				assert.Equal(s.T(), tc.expectedError, err)
				assert.Equal(s.T(), 0.0, balance)
			} else {
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), tc.expectedBalance, balance)
			}

			// Verify all mocks were called
			s.accountRepository.AssertExpectations(s.T())
			s.transactionRepository.AssertExpectations(s.T())
			s.txProvider.AssertExpectations(s.T())
		})
	}
}

// TestAccountServiceSuite runs the test suite
func TestAccountServiceSuite(t *testing.T) {
	suite.Run(t, new(AccountServiceTestSuite))
}
