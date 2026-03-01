package repository

import (
	"context"

	"github.com/sky0621/techcv/services/manager/backend/internal/domain/model"
)

type ProfileRepository interface {
	Save(ctx context.Context, profile model.Profile) error
}
