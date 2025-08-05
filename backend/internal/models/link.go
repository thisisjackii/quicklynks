# File: backend/internal/models/link.go
package models

import "time"

// Link represents a link in the database.
type Link struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"-" db:"user_id"` // Hide user_id from public JSON
	Title     string    `json:"title" db:"title"`
	URL       string    `json:"url" db:"url"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}