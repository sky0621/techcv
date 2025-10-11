package domain

import "time"

// HealthStatus represents the status of the application.
type HealthStatus struct {
	Status    string    `json:"status"`
	CheckedAt time.Time `json:"checked_at"`
}
