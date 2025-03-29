package repositories

import (
	"backend-developer-assignment/app/models"

	"github.com/jmoiron/sqlx"
)

// UserGreetingRepository is an interface for user repository
type UserGreetingRepository interface {
	GetByID(id string) (*models.UserGreeting, error)
	Update(u *models.UserGreeting) error
}

// UserGreetingsRepository will hold all the repository operations related to users.
type UserGreetingRepositoryImpl struct {
	DB *sqlx.DB
}

// NewUserGreetingsRepository creates a new instance of UserGreetingsRepository.
func NewUserGreetingsRepository(db *sqlx.DB) UserGreetingRepository {
	return &UserGreetingRepositoryImpl{
		DB: db,
	}
}

// GetByID get user greeting by ID.
func (r *UserGreetingRepositoryImpl) GetByID(id string) (*models.UserGreeting, error) {
	user := &models.UserGreeting{}

	query := `SELECT * FROM user_greetings WHERE user_id = ?`

	err := r.DB.Get(user, query, id)
	if err != nil {
		return user, err
	}

	return user, nil
}

// Update performs an update on user greeting information.
func (r *UserGreetingRepositoryImpl) Update(u *models.UserGreeting) error {
	query := `UPDATE user_greetings SET greeting = ? WHERE user_id = ?`

	_, err := r.DB.Exec(
		query,
		u.Greeting, u.UserID,
	)
	if err != nil {
		return err
	}

	return nil
}
