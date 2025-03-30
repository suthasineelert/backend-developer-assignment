package services

import (
	"backend-developer-assignment/app/models"
	"backend-developer-assignment/app/repositories"

	"github.com/google/uuid"
)

// DebitCardService defines the interface for debit card operations
type DebitCardService interface {
	// Card operations
	GetCardByID(cardID string) (*models.DebitCard, error)
	GetCardsByUserID(userID string) ([]*models.DebitCard, error)
	GetCardWithDetailByID(cardID string) (*models.DebitCardWithDetails, error)
	GetCardWithDetailByUserID(userID string) ([]*models.DebitCardWithDetails, error)

	// Create operations
	CreateCardWithDetails(cardWithDetails *models.DebitCardWithDetails) error

	// Update operations
	UpdateCard(card *models.DebitCard, name, color, borderColor string) error

	// Delete operations
	DeleteCard(cardID string) error
}

// DebitCardServiceImpl implements DebitCardService
type DebitCardServiceImpl struct {
	debitCardRepository repositories.DebitCardRepository
}

// NewDebitCardService creates a new instance of DebitCardService
func NewDebitCardService(repo repositories.DebitCardRepository) DebitCardService {
	return &DebitCardServiceImpl{
		debitCardRepository: repo,
	}
}

// GetCardByID retrieves a debit card by ID
func (s *DebitCardServiceImpl) GetCardByID(cardID string) (*models.DebitCard, error) {
	return s.debitCardRepository.GetCardByID(cardID)
}

// GetCardsByUserID retrieves all debit cards for a user
func (s *DebitCardServiceImpl) GetCardsByUserID(userID string) ([]*models.DebitCard, error) {
	return s.debitCardRepository.GetCardsByUserID(userID)
}

// GetCardWithDetailByID retrieves a complete debit card with all related information by ID
func (s *DebitCardServiceImpl) GetCardWithDetailByID(cardID string) (*models.DebitCardWithDetails, error) {
	return s.debitCardRepository.GetCardWithDetailByID(cardID)
}

// GetCardWithDetailByUserID retrieves all complete debit cards with related information for a user
func (s *DebitCardServiceImpl) GetCardWithDetailByUserID(userID string) ([]*models.DebitCardWithDetails, error) {
	return s.debitCardRepository.GetCardWithDetailByUserID(userID)
}

// CreateCardWithDetails creates a new debit card with all related details
func (s *DebitCardServiceImpl) CreateCardWithDetails(cardWithDetails *models.DebitCardWithDetails) error {
	// Generate a new UUID if not provided
	if cardWithDetails.CardID == "" {
		cardWithDetails.CardID = uuid.New().String()
	}

	// Create card default status to in-progress
	cardWithDetails.Status = string(models.CardStatusInprogress)

	if err := s.debitCardRepository.CreateCard(cardWithDetails); err != nil {
		return err
	}

	return nil
}

// UpdateCard updates an existing debit card
func (s *DebitCardServiceImpl) UpdateCard(card *models.DebitCard, name, color, borderColor string) error {
	return s.debitCardRepository.UpdateCardByID(card.CardID, card.UserID, func(card *models.DebitCardWithDetails) (bool, error) {
		isUpdate := false

		if name != "" && card.Name != name {
			card.Name = name
			isUpdate = true
		}
		if color != "" && card.Color != color {
			card.Color = color
			isUpdate = true
		}
		if borderColor != "" && card.BorderColor != borderColor {
			card.BorderColor = borderColor
			isUpdate = true
		}

		return isUpdate, nil
	})
}

// DeleteCard marks a card as deleted without removing it
func (s *DebitCardServiceImpl) DeleteCard(cardID string) error {
	// update status card to inactive
	cardStatus := models.DebitCardStatus{
		CardID: cardID,
		Status: string(models.CardStatusInactive),
	}
	return s.debitCardRepository.UpdateCardStatus(&cardStatus)
}
