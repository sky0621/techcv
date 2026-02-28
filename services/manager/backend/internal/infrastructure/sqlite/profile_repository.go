package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sky0621/techcv/services/manager/backend/internal/domain"
)

type ProfileRepository struct {
	db *sql.DB
}

func NewProfileRepository(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (r *ProfileRepository) Save(ctx context.Context, profile domain.Profile) error {
	const q = `
INSERT INTO profiles (id, name, nickname)
VALUES (?, ?, ?)
`

	if _, err := r.db.ExecContext(ctx, q, profile.ID, profile.Name, profile.Nickname); err != nil {
		return fmt.Errorf("insert profile: %w", err)
	}
	return nil
}
