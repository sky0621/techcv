package domain

import "time"

// PublicURL represents a sharable URL entry managed by the system.
type PublicURL struct {
	ID        uint64    `json:"id"`
	URLKey    string    `json:"url_key"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
