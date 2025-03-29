package repositories

import (
	"backend-developer-assignment/app/models"
	"time"

	"github.com/jmoiron/sqlx"
)

// DebitCardRepository is an interface for debit card repository operations
type DebitCardRepository interface {
	// Get card operations
	GetCardByID(cardID string) (*models.DebitCard, error)
	GetCardsByUserID(userID string) ([]*models.DebitCard, error)
	GetCardDetailByID(cardID string) (*models.DebitCardDetail, error)
	GetCardDesignByID(cardID string) (*models.DebitCardDesign, error)
	GetCardStatusByID(cardID string) (*models.DebitCardStatus, error)
	GetCardWithDetailByID(cardID string) (*models.DebitCardWithDetail, error)
	GetCardWithDetailByUserID(userID string) ([]*models.DebitCardWithDetail, error)

	// Update Card operations
	UpdateCardTx(tx DBTransaction, card *models.DebitCard) error
	UpdateCardDetail(detail *models.DebitCardDetail) error
	UpdateCardDesignTx(tx DBTransaction, design *models.DebitCardDesign) error
	UpdateCardStatus(status *models.DebitCardStatus) error

	// Create Card operations
	BeginTx() (DBTransaction, error)
	CreateCardTx(tx DBTransaction, card *models.DebitCard) error
	CreateCardDetailTx(tx DBTransaction, detail *models.DebitCardDetail) error
	CreateCardDesignTx(tx DBTransaction, design *models.DebitCardDesign) error
	CreateCardStatusTx(tx DBTransaction, status *models.DebitCardStatus) error

	// Transaction operations
	CommitTx(tx DBTransaction) error
	RollbackTx(tx DBTransaction) error

	// Delete Card operations
	DeleteCard(cardID string) error
}

// DebitCardRepositoryImpl implements DebitCardRepository
type DebitCardRepositoryImpl struct {
	DB *sqlx.DB
}

// NewDebitCardRepository creates a new instance of DebitCardRepository
func NewDebitCardRepository(db *sqlx.DB) DebitCardRepository {
	return &DebitCardRepositoryImpl{
		DB: db,
	}
}

// GetCardWithDetailByID retrieves a complete debit card with all related information by ID
func (r *DebitCardRepositoryImpl) GetCardWithDetailByID(cardID string) (*models.DebitCardWithDetail, error) {
	card := &models.DebitCardWithDetail{}

	query := `
		SELECT 
			c.card_id, c.user_id, c.name, c.created_at, c.updated_at, c.deleted_at,
			d.issuer, d.number,
			ds.color, ds.border_color,
			s.status
		FROM 
			debit_cards c
		LEFT JOIN 
			debit_card_details d ON c.card_id = d.card_id
		LEFT JOIN 
			debit_card_design ds ON c.card_id = ds.card_id
		LEFT JOIN 
			debit_card_status s ON c.card_id = s.card_id
		WHERE 
			c.card_id = ? AND c.deleted_at IS NULL
	`

	err := r.DB.Get(card, query, cardID)
	if err != nil {
		return nil, err
	}

	return card, nil
}

// GetCardWithDetailByUserID retrieves all complete debit cards with related information for a user
func (r *DebitCardRepositoryImpl) GetCardWithDetailByUserID(userID string) ([]*models.DebitCardWithDetail, error) {
	cards := []*models.DebitCardWithDetail{}

	query := `
		SELECT 
			c.card_id, c.user_id, c.name, c.created_at, c.updated_at, c.deleted_at,
			d.issuer, d.number,
			ds.color, ds.border_color,
			s.status
		FROM 
			debit_cards c
		LEFT JOIN 
			debit_card_details d ON c.card_id = d.card_id
		LEFT JOIN 
			debit_card_design ds ON c.card_id = ds.card_id
		LEFT JOIN 
			debit_card_status s ON c.card_id = s.card_id
		WHERE 
			c.user_id = ? AND c.deleted_at IS NULL
		ORDER BY
			c.created_at DESC
	`

	err := r.DB.Select(&cards, query, userID)
	if err != nil {
		return nil, err
	}

	return cards, nil
}

// GetCardByID retrieves a debit card by ID
func (r *DebitCardRepositoryImpl) GetCardByID(cardID string) (*models.DebitCard, error) {
	card := &models.DebitCard{}
	query := `SELECT * FROM debit_cards WHERE card_id = ? AND deleted_at IS NULL`
	err := r.DB.Get(card, query, cardID)
	if err != nil {
		return nil, err
	}
	return card, nil
}

// GetCardsByUserID retrieves all debit cards for a user
func (r *DebitCardRepositoryImpl) GetCardsByUserID(userID string) ([]*models.DebitCard, error) {
	cards := []*models.DebitCard{}
	query := `SELECT * FROM debit_cards WHERE user_id = ? AND deleted_at IS NULL`
	err := r.DB.Select(&cards, query, userID)
	if err != nil {
		return nil, err
	}
	return cards, nil
}

// UpdateCard updates an existing debit card
func (r *DebitCardRepositoryImpl) UpdateCardTx(tx DBTransaction, card *models.DebitCard) error {
	card.UpdatedAt = time.Now()

	query := `UPDATE debit_cards 
              SET user_id = ?, name = ?, updated_at = ? 
              WHERE card_id = ? AND deleted_at IS NULL`
	_, err := tx.Exec(
		query,
		card.UserID,
		card.Name,
		card.UpdatedAt,
		card.CardID,
	)
	return err
}

// DeleteCard marks a card as deleted without removing it
func (r *DebitCardRepositoryImpl) DeleteCard(cardID string) error {
	now := time.Now()
	query := `UPDATE debit_cards SET deleted_at = ? WHERE card_id = ? AND deleted_at IS NULL`
	_, err := r.DB.Exec(query, now, cardID)
	return err
}

// GetCardDetailByID retrieves card details by card ID
func (r *DebitCardRepositoryImpl) GetCardDetailByID(cardID string) (*models.DebitCardDetail, error) {
	detail := &models.DebitCardDetail{}
	query := `SELECT * FROM debit_card_details WHERE card_id = ?`
	err := r.DB.Get(detail, query, cardID)
	if err != nil {
		return nil, err
	}
	return detail, nil
}

// UpdateCardDetail updates existing card details
func (r *DebitCardRepositoryImpl) UpdateCardDetail(detail *models.DebitCardDetail) error {
	query := `UPDATE debit_card_details 
              SET user_id = ?, issuer = ?, number = ? 
              WHERE card_id = ?`
	_, err := r.DB.Exec(
		query,
		detail.UserID,
		detail.Issuer,
		detail.Number,
		detail.CardID,
	)
	return err
}

// GetCardDesignByID retrieves card design by card ID
func (r *DebitCardRepositoryImpl) GetCardDesignByID(cardID string) (*models.DebitCardDesign, error) {
	design := &models.DebitCardDesign{}
	query := `SELECT * FROM debit_card_design WHERE card_id = ?`
	err := r.DB.Get(design, query, cardID)
	if err != nil {
		return nil, err
	}
	return design, nil
}

// UpdateCardDesign updates existing card design
func (r *DebitCardRepositoryImpl) UpdateCardDesignTx(tx DBTransaction, design *models.DebitCardDesign) error {
	query := `UPDATE debit_card_design 
              SET user_id = ?, color = ?, border_color = ? 
              WHERE card_id = ?`
	_, err := tx.Exec(
		query,
		design.UserID,
		design.Color,
		design.BorderColor,
		design.CardID,
	)
	return err
}

// GetCardStatusByID retrieves card status by card ID
func (r *DebitCardRepositoryImpl) GetCardStatusByID(cardID string) (*models.DebitCardStatus, error) {
	status := &models.DebitCardStatus{}
	query := `SELECT * FROM debit_card_status WHERE card_id = ?`
	err := r.DB.Get(status, query, cardID)
	if err != nil {
		return nil, err
	}
	return status, nil
}

// UpdateCardStatus updates existing card status
func (r *DebitCardRepositoryImpl) UpdateCardStatus(status *models.DebitCardStatus) error {
	query := `UPDATE debit_card_status 
              SET user_id = ?, status = ? 
              WHERE card_id = ?`
	_, err := r.DB.Exec(
		query,
		status.UserID,
		status.Status,
		status.CardID,
	)
	return err
}

// BeginTx starts a new transaction
func (r *DebitCardRepositoryImpl) BeginTx() (DBTransaction, error) {
	tx, err := r.DB.Beginx()
	if err != nil {
		return nil, err
	}
	return &sqlxTransaction{tx: tx}, nil
}

// CommitTx commits a transaction
func (r *DebitCardRepositoryImpl) CommitTx(tx DBTransaction) error {
	if sqlxTx, ok := tx.(*sqlxTransaction); ok {
		return sqlxTx.tx.Commit()
	}
	return tx.Commit()
}

// RollbackTx rolls back a transaction
func (r *DebitCardRepositoryImpl) RollbackTx(tx DBTransaction) error {
	if sqlxTx, ok := tx.(*sqlxTransaction); ok {
		return sqlxTx.tx.Rollback()
	}
	return tx.Rollback()
}

// CreateCardTx adds a new debit card within a transaction
func (r *DebitCardRepositoryImpl) CreateCardTx(tx DBTransaction, card *models.DebitCard) error {
	now := time.Now()
	card.CreatedAt = now
	card.UpdatedAt = now

	query := `INSERT INTO debit_cards (card_id, user_id, name, created_at, updated_at) 
              VALUES (?, ?, ?, ?, ?)`
	_, err := tx.Exec(
		query,
		card.CardID,
		card.UserID,
		card.Name,
		card.CreatedAt,
		card.UpdatedAt,
	)
	return err
}

// CreateCardDetailTx adds new card details within a transaction
func (r *DebitCardRepositoryImpl) CreateCardDetailTx(tx DBTransaction, detail *models.DebitCardDetail) error {
	query := `INSERT INTO debit_card_details (card_id, user_id, issuer, number) 
              VALUES (?, ?, ?, ?)`
	_, err := tx.Exec(
		query,
		detail.CardID,
		detail.UserID,
		detail.Issuer,
		detail.Number,
	)
	return err
}

// CreateCardDesignTx adds new card design within a transaction
func (r *DebitCardRepositoryImpl) CreateCardDesignTx(tx DBTransaction, design *models.DebitCardDesign) error {
	query := `INSERT INTO debit_card_design (card_id, user_id, color, border_color) 
              VALUES (?, ?, ?, ?)`
	_, err := tx.Exec(
		query,
		design.CardID,
		design.UserID,
		design.Color,
		design.BorderColor,
	)
	return err
}

// CreateCardStatusTx adds new card status within a transaction
func (r *DebitCardRepositoryImpl) CreateCardStatusTx(tx DBTransaction, status *models.DebitCardStatus) error {
	query := `INSERT INTO debit_card_status (card_id, user_id, status) 
              VALUES (?, ?, ?)`
	_, err := tx.Exec(
		query,
		status.CardID,
		status.UserID,
		status.Status,
	)
	return err
}
