package publicurl

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/sky0621/techcv/manager/backend/internal/domain"
)

// Repository defines the persistence operations required by the public URL use case.
type Repository interface {
	Create(ctx context.Context, urlKey string) (uint64, error)
	GetActive(ctx context.Context) (*domain.PublicURL, error)
	List(ctx context.Context) ([]domain.PublicURL, error)
	Deactivate(ctx context.Context, id uint64) error
}

// Usecase orchestrates public URL management.
type Usecase struct {
	repo   Repository
	keygen func() (string, error)
}

// New constructs a new Usecase instance.
func New(repo Repository) *Usecase {
	return &Usecase{
		repo:   repo,
		keygen: generateKey,
	}
}

// List returns the stored public URLs.
func (u *Usecase) List(ctx context.Context) ([]domain.PublicURL, error) {
	urls, err := u.repo.List(ctx)
	if err != nil {
		return nil, domain.NewInternal("public_url.list_failed", "failed to list public URLs", err)
	}
	return urls, nil
}

// GetActive returns the currently active public URL, if one exists.
func (u *Usecase) GetActive(ctx context.Context) (*domain.PublicURL, error) {
	url, err := u.repo.GetActive(ctx)
	if err != nil {
		return nil, domain.NewInternal("public_url.fetch_failed", "failed to fetch active public URL", err)
	}
	return url, nil
}

// Generate deactivates the current URL (if any) and issues a new random key.
func (u *Usecase) Generate(ctx context.Context) (*domain.PublicURL, error) {
	key, err := u.keygen()
	if err != nil {
		return nil, domain.NewInternal("public_url.key_generation_failed", "failed to generate public URL key", err)
	}

	active, err := u.repo.GetActive(ctx)
	if err != nil {
		return nil, domain.NewInternal("public_url.fetch_failed", "failed to fetch active public URL", err)
	}

	if active != nil {
		if err := u.repo.Deactivate(ctx, active.ID); err != nil {
			return nil, domain.NewInternal("public_url.deactivate_failed", "failed to deactivate existing public URL", err)
		}
	}

	if _, err := u.repo.Create(ctx, key); err != nil {
		return nil, domain.NewInternal("public_url.create_failed", "failed to create public URL", err)
	}

	created, err := u.repo.GetActive(ctx)
	if err != nil {
		return nil, domain.NewInternal("public_url.fetch_failed", "failed to fetch active public URL", err)
	}

	return created, nil
}

func generateKey() (string, error) {
	const keyLength = 16
	buf := make([]byte, keyLength)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate random bytes: %w", err)
	}
	return hex.EncodeToString(buf), nil
}
