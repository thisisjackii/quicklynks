// File: /quicklynks/backend/internal/models/user.go
package models

import "gorm.io/gorm"

// User represents a user in the database.
type User struct {
	gorm.Model
	Email        string `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"not null" json:"-"` // Omit from JSON responses
	Links        []Link `json:"links,omitempty"`
}
