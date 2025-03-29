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
	GetCardWithDetailByID(cardID string) (*models.DebitCardWithDetail, error)
	GetCardWithDetailByUserID(userID string) ([]*models.DebitCardWithDetail, error)

	// Create operations
	CreateCardWithDetails(card *models.DebitCard, detail *models.DebitCardDetail,
		design *models.DebitCardDesign, status *models.DebitCardStatus) error

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
func (s *DebitCardServiceImpl) GetCardWithDetailByID(cardID string) (*models.DebitCardWithDetail, error) {
	return s.debitCardRepository.GetCardWithDetailByID(cardID)
}

// GetCardWithDetailByUserID retrieves all complete debit cards with related information for a user
func (s *DebitCardServiceImpl) GetCardWithDetailByUserID(userID string) ([]*models.DebitCardWithDetail, error) {
	return s.debitCardRepository.GetCardWithDetailByUserID(userID)
}

// CreateCardWithDetails creates a new debit card with all related details
func (s *DebitCardServiceImpl) CreateCardWithDetails(card *models.DebitCard, detail *models.DebitCardDetail,
	design *models.DebitCardDesign, status *models.DebitCardStatus) error {

	// Generate a new UUID if not provided
	if card.CardID == "" {
		card.CardID = uuid.New().String()
	}

	// Set the same card ID and user ID for all related entities
	detail.CardID = card.CardID
	detail.UserID = card.UserID

	design.CardID = card.CardID
	design.UserID = card.UserID

	status.CardID = card.CardID
	status.UserID = card.UserID

	// Use a transaction to ensure all operations succeed or fail together
	tx, err := s.debitCardRepository.BeginTx()
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Create the main card
	if err = s.debitCardRepository.CreateCardTx(tx, card); err != nil {
		return err
	}

	// Create card details
	if err = s.debitCardRepository.CreateCardDetailTx(tx, detail); err != nil {
		return err
	}

	// Create card design
	if err = s.debitCardRepository.CreateCardDesignTx(tx, design); err != nil {
		return err
	}

	// Create card status
	if err = s.debitCardRepository.CreateCardStatusTx(tx, status); err != nil {
		return err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// UpdateCard updates an existing debit card
func (s *DebitCardServiceImpl) UpdateCard(card *models.DebitCard, name, color, borderColor string) error {
	// Use a transaction to ensure all operations succeed or fail together
	tx, err := s.debitCardRepository.BeginTx()
	if err != nil {
		return err
	}

	// Defer a rollback in case anything fails
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// update name
	if name != "" {
		card.Name = name
	}

	var cardDesign models.DebitCardDesign
	cardDesign.CardID = card.CardID
	cardDesign.UserID = card.UserID

	// update color
	if color != "" {
		cardDesign.Color = color
	}
	// update border color
	if borderColor != "" {
		cardDesign.BorderColor = borderColor
	}

	if borderColor != "" || color != "" {
		// Update card design
		if err = s.debitCardRepository.UpdateCardDesignTx(tx, &cardDesign); err != nil {
			return err
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// DeleteCard marks a card as deleted without removing it
func (s *DebitCardServiceImpl) DeleteCard(cardID string) error {
	return s.debitCardRepository.DeleteCard(cardID)
}
