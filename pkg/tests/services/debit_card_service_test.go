package services_test

import (
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/app/services"
	mocks "backend-developer-assignment/pkg/mocks/repositories"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// DebitCardServiceTestSuite defines the test suite
type DebitCardServiceTestSuite struct {
	suite.Suite
	debitCardRepository *mocks.DebitCardRepository
	service             services.DebitCardService
}

// SetupTest runs before each test
func (s *DebitCardServiceTestSuite) SetupTest() {
	s.debitCardRepository = new(mocks.DebitCardRepository)
	s.service = services.NewDebitCardService(s.debitCardRepository)
}

// TestGetCardByID tests the GetCardByID function
func (s *DebitCardServiceTestSuite) TestGetCardByID() {
	now := time.Now()
	testCases := []struct {
		name          string
		cardID        string
		mockCard      *models.DebitCard
		mockError     error
		expectedCard  *models.DebitCard
		expectedError error
	}{
		{
			name:   "Success - Valid Card",
			cardID: "card-123",
			mockCard: &models.DebitCard{
				CardID: "card-123",
				UserID: "user-123",
				Name:   "Test Card",
				BaseModel: &models.BaseModel{
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			mockError: nil,
			expectedCard: &models.DebitCard{
				CardID: "card-123",
				UserID: "user-123",
				Name:   "Test Card",
				BaseModel: &models.BaseModel{
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			expectedError: nil,
		},
		{
			name:          "Failure - Card Not Found",
			cardID:        "nonexistent-card",
			mockCard:      nil,
			mockError:     errors.New("card not found"),
			expectedCard:  nil,
			expectedError: errors.New("card not found"),
		},
		{
			name:          "Failure - Database Error",
			cardID:        "invalid-card-id",
			mockCard:      nil,
			mockError:     errors.New("database connection failed"),
			expectedCard:  nil,
			expectedError: errors.New("database connection failed"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Mock the repository method
			s.debitCardRepository.On("GetCardByID", tc.cardID).Return(tc.mockCard, tc.mockError).Once()

			// Call the service method
			card, err := s.service.GetCardByID(tc.cardID)

			// Assert results
			if tc.expectedError != nil {
				assert.Error(s.T(), err)
				assert.Equal(s.T(), tc.expectedError.Error(), err.Error())
				assert.Nil(s.T(), card)
			} else {
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), tc.expectedCard, card)
			}

			// Verify expected method calls
			s.debitCardRepository.AssertExpectations(s.T())
		})
	}
}

// TestGetCardWithDetailByID tests the GetCardWithDetailByID function
func (s *DebitCardServiceTestSuite) TestGetCardWithDetailByID() {
	now := time.Now()
	testCases := []struct {
		name          string
		cardID        string
		mockCard      *models.DebitCardWithDetails
		mockError     error
		expectedCard  *models.DebitCardWithDetails
		expectedError error
	}{
		{
			name:   "Success - Valid Card with Details",
			cardID: "card-123",
			mockCard: &models.DebitCardWithDetails{
				CardID:      "card-123",
				UserID:      "user-123",
				Name:        "Test Card",
				Issuer:      "Visa",
				Number:      "4111111111111111",
				Color:       "#FF0000",
				BorderColor: "#000000",
				Status:      "active",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			mockError: nil,
			expectedCard: &models.DebitCardWithDetails{
				CardID:      "card-123",
				UserID:      "user-123",
				Name:        "Test Card",
				Issuer:      "Visa",
				Number:      "4111111111111111",
				Color:       "#FF0000",
				BorderColor: "#000000",
				Status:      "active",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			expectedError: nil,
		},
		{
			name:          "Failure - Card Not Found",
			cardID:        "nonexistent-card",
			mockCard:      nil,
			mockError:     errors.New("card not found"),
			expectedCard:  nil,
			expectedError: errors.New("card not found"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Mock the repository method
			s.debitCardRepository.On("GetCardWithDetailByID", tc.cardID).Return(tc.mockCard, tc.mockError).Once()

			// Call the service method
			card, err := s.service.GetCardWithDetailByID(tc.cardID)

			// Assert results
			if tc.expectedError != nil {
				assert.Error(s.T(), err)
				assert.Equal(s.T(), tc.expectedError.Error(), err.Error())
				assert.Nil(s.T(), card)
			} else {
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), tc.expectedCard, card)
			}

			// Verify expected method calls
			s.debitCardRepository.AssertExpectations(s.T())
		})
	}
}

// TestGetCardWithDetailByUserID tests the GetCardWithDetailByUserID function
func (s *DebitCardServiceTestSuite) TestGetCardWithDetailByUserID() {
	userID := "user-123"
	now := time.Now()

	testCases := []struct {
		name          string
		mockCards     []*models.DebitCardWithDetails
		mockError     error
		expectedCards []*models.DebitCardWithDetails
		expectedError error
	}{
		{
			name: "Success - Multiple Cards",
			mockCards: []*models.DebitCardWithDetails{
				{
					CardID:      "card-123",
					UserID:      userID,
					Name:        "Card 1",
					Issuer:      "Visa",
					Number:      "4111111111111111",
					Color:       "#FF0000",
					BorderColor: "#000000",
					Status:      "active",
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					CardID:      "card-456",
					UserID:      userID,
					Name:        "Card 2",
					Issuer:      "Mastercard",
					Number:      "5555555555554444",
					Color:       "#00FF00",
					BorderColor: "#FFFFFF",
					Status:      "inactive",
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			mockError: nil,
			expectedCards: []*models.DebitCardWithDetails{
				{
					CardID:      "card-123",
					UserID:      userID,
					Name:        "Card 1",
					Issuer:      "Visa",
					Number:      "4111111111111111",
					Color:       "#FF0000",
					BorderColor: "#000000",
					Status:      "active",
					CreatedAt:   now,
					UpdatedAt:   now,
				},
				{
					CardID:      "card-456",
					UserID:      userID,
					Name:        "Card 2",
					Issuer:      "Mastercard",
					Number:      "5555555555554444",
					Color:       "#00FF00",
					BorderColor: "#FFFFFF",
					Status:      "inactive",
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			expectedError: nil,
		},
		{
			name:          "Success - No Cards",
			mockCards:     []*models.DebitCardWithDetails{},
			mockError:     nil,
			expectedCards: []*models.DebitCardWithDetails{},
			expectedError: nil,
		},
		{
			name:          "Failure - Database Error",
			mockCards:     nil,
			mockError:     errors.New("database connection failed"),
			expectedCards: nil,
			expectedError: errors.New("database connection failed"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Mock the repository method
			s.debitCardRepository.On("GetCardWithDetailByUserID", userID).Return(tc.mockCards, tc.mockError).Once()

			// Call the service method
			cards, err := s.service.GetCardWithDetailByUserID(userID)

			// Assert results
			if tc.expectedError != nil {
				assert.Error(s.T(), err)
				assert.Equal(s.T(), tc.expectedError.Error(), err.Error())
				assert.Nil(s.T(), cards)
			} else {
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), tc.expectedCards, cards)
			}

			// Verify expected method calls
			s.debitCardRepository.AssertExpectations(s.T())
		})
	}
}

// TestCreateCardWithDetails tests the CreateCardWithDetails function
func (s *DebitCardServiceTestSuite) TestCreateCardWithDetails() {
	now := time.Now()
	userID := "user-123"

	testCases := []struct {
		name             string
		cardWithDetails  *models.DebitCardWithDetails
		mockError        error
		expectedError    error
		shouldGenerateID bool
		expectedStatus   string
	}{
		{
			name: "Success - With Existing ID",
			cardWithDetails: &models.DebitCardWithDetails{
				CardID:      "card-123",
				UserID:      userID,
				Name:        "Test Card",
				CreatedAt:   now,
				UpdatedAt:   now,
				Issuer:      "Visa",
				Number:      "4111111111111111",
				Color:       "#FF0000",
				BorderColor: "#000000",
			},
			mockError:        nil,
			expectedError:    nil,
			shouldGenerateID: false,
			expectedStatus:   string(models.CardStatusInprogress),
		},
		{
			name: "Success - Generate New ID",
			cardWithDetails: &models.DebitCardWithDetails{
				CardID:      "",
				UserID:      userID,
				Name:        "Test Card",
				CreatedAt:   now,
				UpdatedAt:   now,
				Issuer:      "Visa",
				Number:      "4111111111111111",
				Color:       "#FF0000",
				BorderColor: "#000000",
			},
			mockError:        nil,
			expectedError:    nil,
			shouldGenerateID: true,
			expectedStatus:   string(models.CardStatusInprogress),
		},
		{
			name: "Failure - Repository Error",
			cardWithDetails: &models.DebitCardWithDetails{
				CardID:      "card-123",
				UserID:      userID,
				Name:        "Test Card",
				CreatedAt:   now,
				UpdatedAt:   now,
				Issuer:      "Visa",
				Number:      "4111111111111111",
				Color:       "#FF0000",
				BorderColor: "#000000",
			},
			mockError:        errors.New("failed to create card"),
			expectedError:    errors.New("failed to create card"),
			shouldGenerateID: false,
			expectedStatus:   string(models.CardStatusInprogress),
		},
	}

	for i := range testCases {
		tc := &testCases[i] // Use pointer to avoid copying the struct
		s.Run(tc.name, func() {
			// Reset mocks
			s.debitCardRepository = new(mocks.DebitCardRepository)
			s.service = services.NewDebitCardService(s.debitCardRepository)

			// Save the original CardID for later comparison
			originalCardID := tc.cardWithDetails.CardID

			// Mock the CreateCard method
			s.debitCardRepository.On("CreateCard", mock.MatchedBy(func(card *models.DebitCardWithDetails) bool {
				// Verify the card has the expected values and status is set to in-progress
				if tc.shouldGenerateID {
					return card.UserID == tc.cardWithDetails.UserID &&
						card.Name == tc.cardWithDetails.Name &&
						card.Status == tc.expectedStatus
				}
				return card.CardID == tc.cardWithDetails.CardID &&
					card.UserID == tc.cardWithDetails.UserID &&
					card.Name == tc.cardWithDetails.Name &&
					card.Status == tc.expectedStatus
			})).Return(tc.mockError)

			// Call the service method
			err := s.service.CreateCardWithDetails(tc.cardWithDetails)

			// Assert results
			if tc.expectedError != nil {
				assert.Error(s.T(), err)
				assert.Equal(s.T(), tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(s.T(), err)
				if tc.shouldGenerateID {
					assert.NotEmpty(s.T(), tc.cardWithDetails.CardID)
					_, err := uuid.Parse(tc.cardWithDetails.CardID)
					assert.NoError(s.T(), err, "Generated ID should be a valid UUID")
				} else {
					assert.Equal(s.T(), originalCardID, tc.cardWithDetails.CardID)
				}
				// Verify status was set correctly
				assert.Equal(s.T(), tc.expectedStatus, tc.cardWithDetails.Status)
			}

			// Verify expected method calls
			s.debitCardRepository.AssertExpectations(s.T())
		})
	}
}

// TestUpdateCard tests the UpdateCard function
func (s *DebitCardServiceTestSuite) TestUpdateCard() {
	userID := "user-123"
	cardID := "card-123"

	testCases := []struct {
		name               string
		card               *models.DebitCard
		newName            string
		newColor           string
		newBorderColor     string
		mockError          error
		expectedError      error
		shouldUpdateName   bool
		shouldUpdateDesign bool
	}{
		{
			name: "Success - Update Name Only",
			card: &models.DebitCard{
				CardID: cardID,
				UserID: userID,
				Name:   "Old Name",
			},
			newName:            "New Name",
			newColor:           "",
			newBorderColor:     "",
			mockError:          nil,
			expectedError:      nil,
			shouldUpdateName:   true,
			shouldUpdateDesign: false,
		},
		{
			name: "Success - Update Design Only",
			card: &models.DebitCard{
				CardID: cardID,
				UserID: userID,
				Name:   "Card Name",
			},
			newName:            "",
			newColor:           "#00FF00",
			newBorderColor:     "#0000FF",
			mockError:          nil,
			expectedError:      nil,
			shouldUpdateName:   false,
			shouldUpdateDesign: true,
		},
		{
			name: "Success - Update Both Name and Design",
			card: &models.DebitCard{
				CardID: cardID,
				UserID: userID,
				Name:   "Old Name",
			},
			newName:            "New Name",
			newColor:           "#00FF00",
			newBorderColor:     "#0000FF",
			mockError:          nil,
			expectedError:      nil,
			shouldUpdateName:   true,
			shouldUpdateDesign: true,
		},
		{
			name: "Failure - Repository Error",
			card: &models.DebitCard{
				CardID: cardID,
				UserID: userID,
				Name:   "Old Name",
			},
			newName:            "New Name",
			newColor:           "#00FF00",
			newBorderColor:     "#0000FF",
			mockError:          errors.New("failed to update card"),
			expectedError:      errors.New("failed to update card"),
			shouldUpdateName:   true,
			shouldUpdateDesign: true,
		},
	}

	for i := range testCases {
		tc := &testCases[i]
		s.Run(tc.name, func() {
			// Reset mocks
			s.debitCardRepository = new(mocks.DebitCardRepository)
			s.service = services.NewDebitCardService(s.debitCardRepository)

			// Mock the UpdateCardByID method
			s.debitCardRepository.On("UpdateCardByID", tc.card.CardID, tc.card.UserID, mock.AnythingOfType("func(*models.DebitCardWithDetails) (bool, error)")).
				Return(tc.mockError).
				Run(func(args mock.Arguments) {
					// Extract and call the callback function to verify it works correctly
					updateFn := args.Get(2).(func(*models.DebitCardWithDetails) (bool, error))
					
					// Create a test card with details to pass to the callback
					cardWithDetails := &models.DebitCardWithDetails{
						CardID:      tc.card.CardID,
						UserID:      tc.card.UserID,
						Name:        tc.card.Name,
						Color:       "original-color",
						BorderColor: "original-border",
					}
					
					// Call the update function
					updated, _ := updateFn(cardWithDetails)
					
					// Verify the card was updated as expected
					if tc.shouldUpdateName {
						assert.Equal(s.T(), tc.newName, cardWithDetails.Name)
					}
					if tc.shouldUpdateDesign {
						if tc.newColor != "" {
							assert.Equal(s.T(), tc.newColor, cardWithDetails.Color)
						}
						if tc.newBorderColor != "" {
							assert.Equal(s.T(), tc.newBorderColor, cardWithDetails.BorderColor)
						}
					}
					
					// Verify the update flag is set correctly
					expectedUpdate := tc.shouldUpdateName || tc.shouldUpdateDesign
					assert.Equal(s.T(), expectedUpdate, updated)
				})

			// Call the service method
			err := s.service.UpdateCard(tc.card, tc.newName, tc.newColor, tc.newBorderColor)

			// Assert results
			if tc.expectedError != nil {
				assert.Error(s.T(), err)
				assert.Equal(s.T(), tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(s.T(), err)
			}

			// Verify expected method calls
			s.debitCardRepository.AssertExpectations(s.T())
		})
	}
}

// TestDeleteCard tests the DeleteCard function
func (s *DebitCardServiceTestSuite) TestDeleteCard() {
	testCases := []struct {
		name          string
		cardID        string
		mockError     error
		expectedError error
	}{
		{
			name:          "Success - Card Deleted",
			cardID:        "card-123",
			mockError:     nil,
			expectedError: nil,
		},
		{
			name:          "Failure - Card Not Found",
			cardID:        "nonexistent-card",
			mockError:     errors.New("card not found"),
			expectedError: errors.New("card not found"),
		},
		{
			name:          "Failure - Database Error",
			cardID:        "card-123",
			mockError:     errors.New("database connection failed"),
			expectedError: errors.New("database connection failed"),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Mock the repository method with the expected DebitCardStatus object
			s.debitCardRepository.On("UpdateCardStatus", mock.MatchedBy(func(status *models.DebitCardStatus) bool {
				return status.CardID == tc.cardID && status.Status == string(models.CardStatusInactive)
			})).Return(tc.mockError).Once()

			// Call the service method
			err := s.service.DeleteCard(tc.cardID)

			// Assert results
			if tc.expectedError != nil {
				assert.Error(s.T(), err)
				assert.Equal(s.T(), tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(s.T(), err)
			}

			// Verify expected method calls
			s.debitCardRepository.AssertExpectations(s.T())
		})
	}
}

// Run the test suite
func TestDebitCardServiceTestSuite(t *testing.T) {
	suite.Run(t, new(DebitCardServiceTestSuite))
}
