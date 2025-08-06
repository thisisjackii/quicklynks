// File: /quicklynks/backend/internal/models/link.go
package models

import "gorm.io/gorm"

// Link represents a shortened URL.
type Link struct {
	gorm.Model
	OriginalURL string  `gorm:"not null" json:"original_url"`
	ShortCode   string  `gorm:"uniqueIndex;not null" json:"short_code"`
	UserID      uint    `gorm:"not null" json:"user_id"`
	Clicks      []Click `json:"-"` // Omit from general link responses for performance
}
