// File: /quicklynks/backend/internal/models/click.go
package models

import "gorm.io/gorm"

// Click represents a single click/redirect event on a Link.
type Click struct {
	gorm.Model
	LinkID    uint   `gorm:"not null" json:"link_id"`
	IPAddress string `json:"ip_address,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
}
