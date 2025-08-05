# File: backend/internal/store/user_store.go
package store

import (
	"database/sql"
	"errors"

	"github.com/thisisjackii/quicklynks/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

// UserStore handles database operations for users.
type UserStore struct {
	db *sqlx.DB
}

// NewUserStore creates a new UserStore.
func NewUserStore(db *sqlx.DB) *UserStore {
	return &UserStore{db: db}
}

// Create inserts a new user into the database.
func (s *UserStore) Create(user *models.User) error {
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id, created_at`
	// Using QueryRowx because we need to scan the returned id and created_at back into the user struct.
	// The driver will handle the difference between $1 (Postgres) and ? (SQLite).
	return s.db.QueryRowx(s.db.Rebind(query), user.Username, user.Email, user.PasswordHash).StructScan(user)
}

// GetByEmail finds a user by their email address.
func (s *UserStore) GetByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE email = $1`
	err := s.db.Get(&user, s.db.Rebind(query), email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Return nil, nil if user not found
		}
		return nil, err
	}
	return &user, nil
}

// GetByUsername finds a user by their username.
func (s *UserStore) GetByUsername(username string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE username = $1`
	err := s.db.Get(&user, s.db.Rebind(query), username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found is not a fatal error
		}
		return nil, err
	}
	return &user, nil
}