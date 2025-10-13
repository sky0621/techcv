package health

import (
	"context"
	"time"

	"github.com/sky0621/techcv/manager/backend/internal/domain"
)

// Repository exposes the persistence operations required by the health use case.
type Repository interface {
	Ping(ctx context.Context) error
}

type noopRepository struct{}

func (noopRepository) Ping(context.Context) error {
	return nil
}

// Usecase exposes health check operations.
type Usecase struct {
	repo Repository
	now  func() time.Time
}

// New creates a new health usecase instance.
func New(repo Repository) *Usecase {
	if repo == nil {
		repo = noopRepository{}
	}

	return &Usecase{
		repo: repo,
		now: func() time.Time {
			return time.Now().UTC()
		},
	}
}

// Check reports the current health status.
func (u *Usecase) Check(ctx context.Context) (*domain.HealthStatus, error) {
	if err := u.repo.Ping(ctx); err != nil {
		return nil, domain.NewInternal("health.database_unavailable", "database connectivity failed", err)
	}

	return &domain.HealthStatus{
		Status:    "ok",
		CheckedAt: u.now(),
	}, nil
}
