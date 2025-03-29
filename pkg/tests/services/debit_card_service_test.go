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
		mockBeginTxError error
		mockCreateErrors []error
		mockCommitError  error
		expectedError    error
		shouldGenerateID bool
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
				Status:      "active",
			},
			mockBeginTxError: nil,
			mockCreateErrors: []error{nil, nil, nil, nil},
			mockCommitError:  nil,
			expectedError:    nil,
			shouldGenerateID: false,
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
				Status:      "active",
			},
			mockBeginTxError: nil,
			mockCreateErrors: []error{nil, nil, nil, nil},
			mockCommitError:  nil,
			expectedError:    nil,
			shouldGenerateID: true,
		},
		{
			name: "Failure - Begin Transaction Error",
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
				Status:      "active",
			},
			mockBeginTxError: errors.New("failed to begin transaction"),
			mockCreateErrors: []error{},
			mockCommitError:  nil,
			expectedError:    errors.New("failed to begin transaction"),
			shouldGenerateID: false,
		},
		{
			name: "Failure - Create Card Error",
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
				Status:      "active",
			},
			mockBeginTxError: nil,
			mockCreateErrors: []error{errors.New("failed to create card"), nil, nil, nil},
			mockCommitError:  nil,
			expectedError:    errors.New("failed to create card"),
			shouldGenerateID: false,
		},
		{
			name: "Failure - Create Card Detail Error",
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
				Status:      "active",
			},
			mockBeginTxError: nil,
			mockCreateErrors: []error{nil, errors.New("failed to create card detail"), nil, nil},
			mockCommitError:  nil,
			expectedError:    errors.New("failed to create card detail"),
			shouldGenerateID: false,
		},
		{
			name: "Failure - Create Card Design Error",
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
				Status:      "active",
			},
			mockBeginTxError: nil,
			mockCreateErrors: []error{nil, nil, errors.New("failed to create card design"), nil},
			mockCommitError:  nil,
			expectedError:    errors.New("failed to create card design"),
			shouldGenerateID: false,
		},
		{
			name: "Failure - Create Card Status Error",
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
				Status:      "active",
			},
			mockBeginTxError: nil,
			mockCreateErrors: []error{nil, nil, nil, errors.New("failed to create card status")},
			mockCommitError:  nil,
			expectedError:    errors.New("failed to create card status"),
			shouldGenerateID: false,
		},
		{
			name: "Failure - Commit Error",
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
				Status:      "active",
			},
			mockBeginTxError: nil,
			mockCreateErrors: []error{nil, nil, nil, nil},
			mockCommitError:  errors.New("failed to commit transaction"),
			expectedError:    errors.New("failed to commit transaction"),
			shouldGenerateID: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Reset mocks
			s.debitCardRepository = new(mocks.DebitCardRepository)
			s.service = services.NewDebitCardService(s.debitCardRepository)
			mockTx := new(mocks.MockTx)

			// Setup mock transaction
			mockTx.On("Commit").Return(tc.mockCommitError)
			mockTx.On("Rollback").Return(nil)

			// Mock BeginTx
			s.debitCardRepository.On("BeginTx").Return(mockTx, tc.mockBeginTxError)

			if tc.mockBeginTxError == nil {
				// Mock CreateCardTx if we get past BeginTx
				if len(tc.mockCreateErrors) > 0 {
					s.debitCardRepository.On("CreateCardTx", mockTx, mock.AnythingOfType("*models.DebitCard")).Return(tc.mockCreateErrors[0])
				}

				// Mock CreateCardDetailTx if CreateCardTx succeeds
				if len(tc.mockCreateErrors) > 1 && tc.mockCreateErrors[0] == nil {
					s.debitCardRepository.On("CreateCardDetailTx", mockTx, mock.AnythingOfType("*models.DebitCardDetail")).Return(tc.mockCreateErrors[1])
				}

				// Mock CreateCardDesignTx if CreateCardDetailTx succeeds
				if len(tc.mockCreateErrors) > 2 && tc.mockCreateErrors[0] == nil && tc.mockCreateErrors[1] == nil {
					s.debitCardRepository.On("CreateCardDesignTx", mockTx, mock.AnythingOfType("*models.DebitCardDesign")).Return(tc.mockCreateErrors[2])
				}

				// Mock CreateCardStatusTx if CreateCardDesignTx succeeds
				if len(tc.mockCreateErrors) > 3 && tc.mockCreateErrors[0] == nil && tc.mockCreateErrors[1] == nil && tc.mockCreateErrors[2] == nil {
					s.debitCardRepository.On("CreateCardStatusTx", mockTx, mock.AnythingOfType("*models.DebitCardStatus")).Return(tc.mockCreateErrors[3])
				}
			}

			// Save the original CardID for later comparison
			originalCardID := tc.cardWithDetails.CardID

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
		mockBeginTxError   error
		mockUpdateError    error
		mockCommitError    error
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
			mockBeginTxError:   nil,
			mockUpdateError:    nil,
			mockCommitError:    nil,
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
			mockBeginTxError:   nil,
			mockUpdateError:    nil,
			mockCommitError:    nil,
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
			mockBeginTxError:   nil,
			mockUpdateError:    nil,
			mockCommitError:    nil,
			expectedError:      nil,
			shouldUpdateName:   true,
			shouldUpdateDesign: true,
		},
		{
			name: "Failure - Begin Transaction Error",
			card: &models.DebitCard{
				CardID: cardID,
				UserID: userID,
				Name:   "Old Name",
			},
			newName:            "New Name",
			newColor:           "#00FF00",
			newBorderColor:     "#0000FF",
			mockBeginTxError:   errors.New("failed to begin transaction"),
			mockUpdateError:    nil,
			mockCommitError:    nil,
			expectedError:      errors.New("failed to begin transaction"),
			shouldUpdateName:   true,
			shouldUpdateDesign: true,
		},
		{
			name: "Failure - Update Design Error",
			card: &models.DebitCard{
				CardID: cardID,
				UserID: userID,
				Name:   "Old Name",
			},
			newName:            "New Name",
			newColor:           "#00FF00",
			newBorderColor:     "#0000FF",
			mockBeginTxError:   nil,
			mockUpdateError:    errors.New("failed to update design"),
			mockCommitError:    nil,
			expectedError:      errors.New("failed to update design"),
			shouldUpdateName:   true,
			shouldUpdateDesign: true,
		},
		{
			name: "Failure - Commit Error",
			card: &models.DebitCard{
				CardID: cardID,
				UserID: userID,
				Name:   "Old Name",
			},
			newName:            "New Name",
			newColor:           "#00FF00",
			newBorderColor:     "#0000FF",
			mockBeginTxError:   nil,
			mockUpdateError:    nil,
			mockCommitError:    errors.New("failed to commit transaction"),
			expectedError:      errors.New("failed to commit transaction"),
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
			mockTx := new(mocks.MockTx)

			// Setup mock transaction
			mockTx.On("Commit").Return(tc.mockCommitError)
			mockTx.On("Rollback").Return(nil)

			// Mock BeginTx
			s.debitCardRepository.On("BeginTx").Return(mockTx, tc.mockBeginTxError)

			if tc.mockBeginTxError == nil && tc.shouldUpdateDesign {
				// Mock UpdateCardDesignTx if we need to update design
				s.debitCardRepository.On("UpdateCardDesignTx", mockTx, mock.MatchedBy(func(design *models.DebitCardDesign) bool {
					// Verify the design has the expected values
					return design.CardID == tc.card.CardID &&
						design.UserID == tc.card.UserID &&
						(tc.newColor == "" || design.Color == tc.newColor) &&
						(tc.newBorderColor == "" || design.BorderColor == tc.newBorderColor)
				})).Return(tc.mockUpdateError)
			}

			// Save original name for verification
			originalName := tc.card.Name

			// Call the service method
			err := s.service.UpdateCard(tc.card, tc.newName, tc.newColor, tc.newBorderColor)

			// Assert results
			if tc.expectedError != nil {
				assert.Error(s.T(), err)
				assert.Equal(s.T(), tc.expectedError.Error(), err.Error())
			} else {
				assert.NoError(s.T(), err)

				// Verify name was updated if expected
				if tc.shouldUpdateName {
					assert.Equal(s.T(), tc.newName, tc.card.Name)
				} else {
					assert.Equal(s.T(), originalName, tc.card.Name)
				}
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
