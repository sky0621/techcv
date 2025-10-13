package mysql

import (
	"context"
	"database/sql"

	"github.com/sky0621/techcv/manager/backend/internal/infrastructure/mysql/sqlc"
)

// HealthRepository provides database-backed health checks.
type HealthRepository struct {
	queries *mysqlsqlc.Queries
}

// NewHealthRepository constructs a HealthRepository backed by sqlc queries.
func NewHealthRepository(db *sql.DB) *HealthRepository {
	return &HealthRepository{
		queries: mysqlsqlc.New(db),
	}
}

// Ping verifies the database connectivity.
func (r *HealthRepository) Ping(ctx context.Context) error {
	_, err := r.queries.Ping(ctx)
	return err
}
