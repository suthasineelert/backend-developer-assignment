package controllers_test

import (
	"backend-developer-assignment/app/controllers"
	"backend-developer-assignment/app/services"
	mockServices "backend-developer-assignment/pkg/mocks/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitController(t *testing.T) {
	// Create mock services
	mockUserService := new(mockServices.UserService)
	mockTransactionService := new(mockServices.TransactionService)
	mockDebitCardService := new(mockServices.DebitCardService)
	mockAccountService := new(mockServices.AccountService)
	mockBannerService := new(mockServices.BannerService)

	// Create service struct with mocks
	service := &services.Service{
		UserService:        mockUserService,
		TransactionService: mockTransactionService,
		DebitCardService:   mockDebitCardService,
		AccountService:     mockAccountService,
		BannerService:      mockBannerService,
	}

	// Initialize controller
	controller := controllers.InitController(service)

	// Assert that all controllers are initialized
	assert.NotNil(t, controller)
	assert.NotNil(t, controller.AuthController)
	assert.NotNil(t, controller.UserController)
	assert.NotNil(t, controller.TransactionController)
	assert.NotNil(t, controller.DebitCardController)
	assert.NotNil(t, controller.AccountController)
	assert.NotNil(t, controller.BannerController)

	// Verify that the controllers are initialized with the correct services
	// This is a bit tricky since we can't directly access the private fields
	// But we can test that the controllers are of the correct type

	// We can also verify that the controllers are initialized by checking their type
	assert.IsType(t, controllers.AuthController{}, controller.AuthController)
	assert.IsType(t, controllers.UserController{}, controller.UserController)
	assert.IsType(t, controllers.TransactionController{}, controller.TransactionController)
	assert.IsType(t, controllers.DebitCardController{}, controller.DebitCardController)
	assert.IsType(t, controllers.AccountController{}, controller.AccountController)
	assert.IsType(t, controllers.BannerController{}, controller.BannerController)
}
