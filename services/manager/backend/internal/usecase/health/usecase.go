// Package health provides use cases for service health reporting.
package health

import (
	"context"
	"time"

	"github.com/sky0621/techcv/manager/backend/internal/domain"
)

// Usecase exposes health check operations.
type Usecase struct{}

// New creates a new health usecase instance.
func New() *Usecase {
	return &Usecase{}
}

// Check reports the current health status.
func (u *Usecase) Check(ctx context.Context) (*domain.HealthStatus, error) {
	return &domain.HealthStatus{
		Status:    "ok",
		CheckedAt: time.Now().UTC(),
	}, nil
}
