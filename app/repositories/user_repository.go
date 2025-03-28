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
	// Define User variable.
	user := &models.User{}

	// Define query string.
	query := `SELECT * FROM users WHERE user_id = ?`

	// Send query to database.
	err := r.DB.Get(user, query, id)
	if err != nil {
		// Return empty object and error.
		return user, err
	}

	// Return query result.
	return user, nil
}

// GetByName retrieves one User by given name.
func (r *UserRepositoryImpl) GetByName(name string) (*models.User, error) {
	// Define User variable.
	user := &models.User{}

	// Define query string.
	query := `SELECT * FROM users WHERE name = ?`

	// Send query to database.
	err := r.DB.Get(user, query, name)
	if err != nil {
		// Return empty object and error.
		return user, err
	}

	// Return query result.
	return user, nil
}

// Update performs an update on user information.
func (r *UserRepositoryImpl) Update(u *models.User) error {
	// Define query string.
	query := `UPDATE users SET name = ?, pin = ? WHERE user_id = ?`

	// Send query to database.
	_, err := r.DB.Exec(
		query,
		u.Name, u.PIN, u.UserID,
	)
	if err != nil {
		return err
	}

	return nil
}
