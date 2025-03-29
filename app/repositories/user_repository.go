package repositories

import (
	"backend-developer-assignment/app/models"

	"github.com/jmoiron/sqlx"
)

// UserRepository is an interface for user repository
type UserRepository interface {
	GetByID(id string) (*models.User, error)
	GetByName(name string) (*models.User, error)
	Update(u *models.User) error
}

// UserRepository will hold all the repository operations related to users.
type UserRepositoryImpl struct {
	DB *sqlx.DB
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *sqlx.DB) UserRepository {
	return &UserRepositoryImpl{
		DB: db,
	}
}

// GetByID retrieves one User by given ID.
func (r *UserRepositoryImpl) GetByID(id string) (*models.User, error) {
	user := &models.User{}

	query := `SELECT * FROM users WHERE user_id = ? and deleted_at IS NULL`

	// Send query to database.
	err := r.DB.Get(user, query, id)
	if err != nil {
		return user, err
	}

	return user, nil
}

// GetByName retrieves one User by given name.
func (r *UserRepositoryImpl) GetByName(name string) (*models.User, error) {
	user := &models.User{}

	query := `SELECT user_id, name FROM users WHERE name = ? and deleted_at IS NULL`

	err := r.DB.Get(user, query, name)
	if err != nil {
		return user, err
	}

	return user, nil
}

// Update performs an update on user information.
func (r *UserRepositoryImpl) Update(u *models.User) error {
	var query string
	var err error

	// Check if PIN is provided
	if u.PIN != "" {
		// Update both name and PIN
		query = `UPDATE users SET name = ?, pin = ? WHERE user_id = ?`
		_, err = r.DB.Exec(query, u.Name, u.PIN, u.UserID)
	} else {
		// Update only name
		query = `UPDATE users SET name = ? WHERE user_id = ?`
		_, err = r.DB.Exec(query, u.Name, u.UserID)
	}

	if err != nil {
		return err
	}

	return nil
}
