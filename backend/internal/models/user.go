# File: backend/internal/models/user.go
package models

import "time"

// User represents a user in the database.
type User struct {
	ID           int       `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"` // Never expose password hash
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
}