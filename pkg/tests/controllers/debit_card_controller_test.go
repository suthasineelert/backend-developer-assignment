package controllers_test

import (
	"backend-developer-assignment/app/controllers"
	"backend-developer-assignment/app/models"
	mocks "backend-developer-assignment/pkg/mocks/services"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// DebitCardControllerTestSuite defines the test suite
type DebitCardControllerTestSuite struct {
	suite.Suite
	app              *fiber.App
	debitCardService *mocks.DebitCardService
	controller       *controllers.DebitCardController
	testUserID       string
	testCardID       string
	testCardData     *models.DebitCardWithDetails
}

// SetupTest runs before each test
func (s *DebitCardControllerTestSuite) SetupTest() {
	s.app = fiber.New()
	s.debitCardService = new(mocks.DebitCardService)
	s.controller = controllers.NewDebitCardController(s.debitCardService)
	s.testUserID = "test-user-id"
	s.testCardID = "test-card-id"

	now := time.Now()
	s.testCardData = &models.DebitCardWithDetails{
		CardID:      s.testCardID,
		UserID:      s.testUserID,
		Name:        "Test Card",
		Issuer:      "Test Bank",
		Color:       "#FF0000",
		BorderColor: "#000000",
		Number:      "1234567890123456",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Setup routes
	s.app.Get("/cards", func(c *fiber.Ctx) error {
		c.Locals("userID", s.testUserID)
		return s.controller.ListDebitCards(c)
	})

	s.app.Get("/cards/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", s.testUserID)
		return s.controller.GetDebitCard(c)
	})

	s.app.Post("/cards", func(c *fiber.Ctx) error {
		c.Locals("userID", s.testUserID)
		return s.controller.CreateDebitCard(c)
	})

	s.app.Put("/cards/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", s.testUserID)
		return s.controller.UpdateDebitCard(c)
	})

	s.app.Delete("/cards/:id", func(c *fiber.Ctx) error {
		c.Locals("userID", s.testUserID)
		return s.controller.DeleteDebitCard(c)
	})
}

// TestListDebitCards tests the ListDebitCards controller method
func (s *DebitCardControllerTestSuite) TestListDebitCards() {
	// Test case: successful retrieval of cards
	s.debitCardService.On("GetCardWithDetailByUserID", s.testUserID).Return([]*models.DebitCardWithDetails{s.testCardData}, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/cards", http.NoBody)
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

	var cards []*models.DebitCardWithDetails
	err = json.NewDecoder(resp.Body).Decode(&cards)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), cards, 1)
	assert.Equal(s.T(), s.testCardID, cards[0].CardID)

	// Test case: error retrieving cards
	s.debitCardService.On("GetCardWithDetailByUserID", s.testUserID).Return(nil, errors.New("database error")).Once()

	req = httptest.NewRequest(http.MethodGet, "/cards", http.NoBody)
	resp, err = s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusInternalServerError, resp.StatusCode)

	s.debitCardService.AssertExpectations(s.T())
}

// TestGetDebitCard tests the GetDebitCard controller method
func (s *DebitCardControllerTestSuite) TestGetDebitCard() {
	// Test case: successful retrieval of card
	s.debitCardService.On("GetCardWithDetailByID", s.testCardID).Return(s.testCardData, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/cards/"+s.testCardID, http.NoBody)
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

	var card models.DebitCardWithDetails
	err = json.NewDecoder(resp.Body).Decode(&card)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), s.testCardID, card.CardID)

	// Test case: card not found
	s.debitCardService.On("GetCardWithDetailByID", "nonexistent-id").Return(nil, errors.New("Debit card not found")).Once()

	req = httptest.NewRequest(http.MethodGet, "/cards/nonexistent-id", http.NoBody)
	resp, err = s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusNotFound, resp.StatusCode)

	s.debitCardService.AssertExpectations(s.T())
}

// TestCreateDebitCard tests the CreateDebitCard controller method
func (s *DebitCardControllerTestSuite) TestCreateDebitCard() {
	// Test case: successful card creation
	createRequest := map[string]interface{}{
		"name":         "New Card",
		"issuer":       "New Bank",
		"color":        "#00FF00",
		"border_color": "#FFFFFF",
	}

	requestBody, _ := json.Marshal(createRequest)

	newCard := &models.DebitCardWithDetails{
		CardID:      s.testCardID,
		UserID:      s.testUserID,
		Name:        "New Card",
		Issuer:      "New Bank",
		Color:       "#00FF00",
		BorderColor: "#FFFFFF",
	}

	s.debitCardService.On("CreateCardWithDetails", mock.AnythingOfType("*models.DebitCardWithDetails")).Return(nil).Once()
	s.debitCardService.On("GetCardWithDetailByID", mock.AnythingOfType("string")).Return(newCard, nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/cards", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusCreated, resp.StatusCode)

	var card models.DebitCardWithDetails
	err = json.NewDecoder(resp.Body).Decode(&card)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "New Card", card.Name)
	assert.Equal(s.T(), "New Bank", card.Issuer)
	assert.Equal(s.T(), "#00FF00", card.Color)
	assert.Equal(s.T(), "#FFFFFF", card.BorderColor)

	// Test case: validation error
	invalidRequest := map[string]interface{}{
		"name":   "+-?!$", // Non-alpha num space name
		"issuer": "New Bank",
	}

	requestBody, _ = json.Marshal(invalidRequest)

	req = httptest.NewRequest(http.MethodPost, "/cards", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)

	// Test case: service error
	s.debitCardService.On("CreateCardWithDetails", mock.AnythingOfType("*models.DebitCardWithDetails")).Return(errors.New("database error")).Once()

	requestBody, _ = json.Marshal(createRequest)

	req = httptest.NewRequest(http.MethodPost, "/cards", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusInternalServerError, resp.StatusCode)

	s.debitCardService.AssertExpectations(s.T())
}

// TestUpdateDebitCard tests the UpdateDebitCard controller method
func (s *DebitCardControllerTestSuite) TestUpdateDebitCard() {
	// Test case: successful card update
	updateRequest := map[string]interface{}{
		"name":         "Updated Card",
		"color":        "#0000FF",
		"border_color": "#EEEEEE",
	}

	requestBody, _ := json.Marshal(updateRequest)

	s.debitCardService.On("GetCardByID", s.testCardID).Return(&models.DebitCard{
		CardID: s.testCardID,
		UserID: s.testUserID,
	}, nil).Once()

	updatedCard := *s.testCardData
	updatedCard.Name = "Updated Card"
	updatedCard.Color = "#0000FF"
	updatedCard.BorderColor = "#EEEEEE"

	s.debitCardService.On("UpdateCard", mock.AnythingOfType("*models.DebitCard"), "Updated Card", "#0000FF", "#EEEEEE").Return(nil).Once()
	s.debitCardService.On("GetCardWithDetailByID", s.testCardID).Return(&updatedCard, nil).Once()

	req := httptest.NewRequest(http.MethodPut, "/cards/"+s.testCardID, bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

	var card models.DebitCardWithDetails
	err = json.NewDecoder(resp.Body).Decode(&card)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "Updated Card", card.Name)
	assert.Equal(s.T(), "#0000FF", card.Color)
	assert.Equal(s.T(), "#EEEEEE", card.BorderColor)

	// Test case: card not found
	s.debitCardService.On("GetCardByID", "nonexistent-id").Return(nil, errors.New("Debit card not found")).Once()

	req = httptest.NewRequest(http.MethodPut, "/cards/nonexistent-id", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusNotFound, resp.StatusCode)

	// Test case: service error
	s.debitCardService.On("GetCardByID", "error-id").Return(&models.DebitCard{
		CardID: "error-id",
		UserID: s.testUserID,
	}, nil).Once()
	s.debitCardService.On("UpdateCard", mock.AnythingOfType("*models.DebitCard"), "Updated Card", "#0000FF", "#EEEEEE").Return(errors.New("database error")).Once()

	req = httptest.NewRequest(http.MethodPut, "/cards/error-id", bytes.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusInternalServerError, resp.StatusCode)

	s.debitCardService.AssertExpectations(s.T())
}

// TestDeleteDebitCard tests the DeleteDebitCard controller method
func (s *DebitCardControllerTestSuite) TestDeleteDebitCard() {
	// Test case: successful card deletion
	s.debitCardService.On("GetCardByID", s.testCardID).Return(&models.DebitCard{
		CardID: s.testCardID,
		UserID: s.testUserID,
	}, nil).Once()
	s.debitCardService.On("DeleteCard", s.testCardID).Return(nil).Once()

	req := httptest.NewRequest(http.MethodDelete, "/cards/"+s.testCardID, http.NoBody)
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusNoContent, resp.StatusCode)

	// Test case: card not found
	s.debitCardService.On("GetCardByID", "nonexistent-id").Return(nil, errors.New("Debit card not found")).Once()

	req = httptest.NewRequest(http.MethodDelete, "/cards/nonexistent-id", http.NoBody)
	resp, err = s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusNotFound, resp.StatusCode)

	// Test case: service error
	s.debitCardService.On("GetCardByID", "error-id").Return(&models.DebitCard{
		CardID: "error-id",
		UserID: s.testUserID,
	}, nil).Once()
	s.debitCardService.On("DeleteCard", "error-id").Return(errors.New("database error")).Once()

	req = httptest.NewRequest(http.MethodDelete, "/cards/error-id", http.NoBody)
	resp, err = s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusInternalServerError, resp.StatusCode)

	s.debitCardService.AssertExpectations(s.T())
}

// TestDebitCardControllerSuite runs the test suite
func TestDebitCardControllerSuite(t *testing.T) {
	suite.Run(t, new(DebitCardControllerTestSuite))
}
