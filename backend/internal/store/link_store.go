# File: backend/internal/store/link_store.go
package store

import (
	"github.com/thisisjackii/quicklynks/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

// LinkStore handles database operations for links.
type LinkStore struct {
	db *sqlx.DB
}

// NewLinkStore creates a new LinkStore.
func NewLinkStore(db *sqlx.DB) *LinkStore {
	return &LinkStore{db: db}
}

// Create inserts a new link into the database.
func (s *LinkStore) Create(link *models.Link) error {
	query := `INSERT INTO links (user_id, title, url) VALUES ($1, $2, $3) RETURNING id, created_at`
	return s.db.QueryRowx(s.db.Rebind(query), link.UserID, link.Title, link.URL).StructScan(link)
}

// GetByUserID retrieves all links for a given user ID.
func (s *LinkStore) GetByUserID(userID int) ([]models.Link, error) {
	var links []models.Link
	query := `SELECT id, title, url, created_at FROM links WHERE user_id = $1 ORDER BY created_at DESC`
	err := s.db.Select(&links, s.db.Rebind(query), userID)
	return links, err
}

// Delete removes a link from the database.
func (s *LinkStore) Delete(id, userID int) error {
	query := `DELETE FROM links WHERE id = $1 AND user_id = $2`
	result, err := s.db.Exec(s.db.Rebind(query), id, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sqlx.ErrNoRows // Or a custom "not found or not authorized" error
	}
	return nil
}