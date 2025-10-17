package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"

	"github.com/sky0621/techcv/manager/backend/internal/domain"
	mysqlsqlc "github.com/sky0621/techcv/manager/backend/internal/infrastructure/mysql/sqlc"
)

// PublicURLRepository persists public URL entities in MySQL.
type PublicURLRepository struct {
	queries *mysqlsqlc.Queries
}

// NewPublicURLRepository constructs a new repository backed by sqlc queries.
func NewPublicURLRepository(db *sql.DB) *PublicURLRepository {
	return &PublicURLRepository{
		queries: mysqlsqlc.New(db),
	}
}

// Create inserts a new public URL record and returns the generated identifier.
func (r *PublicURLRepository) Create(ctx context.Context, urlKey string) (uint64, error) {
	result, err := r.queries.CreatePublicURL(ctx, urlKey)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	if id < 0 {
		return 0, fmt.Errorf("received negative last insert id: %d", id)
	}

	return uint64(id), nil
}

// GetActive fetches the most recently updated active public URL.
func (r *PublicURLRepository) GetActive(ctx context.Context) (*domain.PublicURL, error) {
	record, err := r.queries.GetActivePublicURL(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	entity, err := toDomainPublicURL(record)
	if err != nil {
		return nil, fmt.Errorf("convert record to domain model: %w", err)
	}

	return &entity, nil
}

// List returns all public URLs ordered by their update timestamp.
func (r *PublicURLRepository) List(ctx context.Context) ([]domain.PublicURL, error) {
	records, err := r.queries.ListPublicURLs(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]domain.PublicURL, 0, len(records))
	for _, record := range records {
		entity, err := toDomainPublicURL(record)
		if err != nil {
			return nil, fmt.Errorf("convert record to domain model: %w", err)
		}
		result = append(result, entity)
	}

	return result, nil
}

// Deactivate marks the specified public URL as inactive.
func (r *PublicURLRepository) Deactivate(ctx context.Context, id uint64) error {
	if id > math.MaxInt64 {
		return fmt.Errorf("public URL id %d exceeds max int64", id)
	}
	return r.queries.DeactivatePublicURL(ctx, int64(id))
}

func toDomainPublicURL(model mysqlsqlc.PublicUrl) (domain.PublicURL, error) {
	if model.ID < 0 {
		return domain.PublicURL{}, fmt.Errorf("public URL id must be non-negative: %d", model.ID)
	}

	return domain.PublicURL{
		ID:        uint64(model.ID),
		URLKey:    model.UrlKey,
		IsActive:  model.IsActive,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}, nil
}
