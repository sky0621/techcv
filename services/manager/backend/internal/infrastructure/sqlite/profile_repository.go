package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sky0621/techcv/services/manager/backend/internal/domain/model"
	sqlcgen "github.com/sky0621/techcv/services/manager/backend/internal/infrastructure/sqlite/sqlc"
)

type ProfileRepository struct {
	queries *sqlcgen.Queries
}

func NewProfileRepository(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{queries: sqlcgen.New(db)}
}

func (r *ProfileRepository) Save(ctx context.Context, profile model.Profile) error {
	if err := r.queries.CreateProfile(ctx, sqlcgen.CreateProfileParams{
		ID:       profile.ID,
		Name:     profile.Name,
		Nickname: profile.Nickname,
	}); err != nil {
		return fmt.Errorf("create profile with sqlc: %w", err)
	}
	return nil
}
