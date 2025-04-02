package services_test

import (
	"backend-developer-assignment/app/repositories"
	"backend-developer-assignment/app/services"
	mockCache "backend-developer-assignment/pkg/mocks/cache"
	mockRepo "backend-developer-assignment/pkg/mocks/repositories"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInitService(t *testing.T) {
	// Create mock repositories
	mockUserRepo := new(mockRepo.UserRepository)
	mockTransactionRepo := new(mockRepo.TransactionRepository)
	mockDebitCardRepo := new(mockRepo.DebitCardRepository)
	mockAccountRepo := new(mockRepo.AccountRepository)
	mockBannerRepo := new(mockRepo.BannerRepository)
	mockTxProvider := new(mockRepo.TxProvider)

	// Create mock redis client
	mockRedisClient := new(mockCache.RedisClient)

	// Create repository struct with mocks
	repo := &repositories.Repository{
		UserRepository:        mockUserRepo,
		TransactionRepository: mockTransactionRepo,
		DebitCardRepository:   mockDebitCardRepo,
		AccountRepository:     mockAccountRepo,
		BannerRepository:      mockBannerRepo,
	}
	// Initialize service
	service := services.InitService(repo, mockTxProvider, mockRedisClient)

	// Assert that all services are initialized
	assert.NotNil(t, service)
	assert.NotNil(t, service.UserService)
	assert.NotNil(t, service.TransactionService)
	assert.NotNil(t, service.DebitCardService)
	assert.NotNil(t, service.AccountService)
	assert.NotNil(t, service.BannerService)

	// Verify that the services are initialized with the correct dependencies
	// This is a bit tricky since we can't directly access the private fields
	// But we can test that the services are functional by calling methods on them

	// For example, we can set up expectations on the mocks and call methods on the services
	mockUserRepo.On("GetByID", mock.Anything).Return(nil, nil).Once()

	// Call a method on the user service
	_, err := service.UserService.GetUserByID("test-id")

	// Verify that the mock was called
	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}
