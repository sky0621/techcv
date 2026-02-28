package register

import (
	"context"
	"errors"
	"testing"

	"github.com/sky0621/techcv/services/manager/backend/internal/domain"
)

type fixedIDGen struct {
	id string
}

func (f fixedIDGen) NewID() string {
	return f.id
}

type stubProfileRepository struct {
	saveErr error
	saved   []domain.Profile
}

func (s *stubProfileRepository) Save(_ context.Context, profile domain.Profile) error {
	if s.saveErr != nil {
		return s.saveErr
	}
	s.saved = append(s.saved, profile)
	return nil
}

func TestService_Register(t *testing.T) {
	t.Parallel()

	repo := &stubProfileRepository{}
	svc := NewService(repo, fixedIDGen{id: "profile_1"})

	got, err := svc.Register(context.Background(), "Alice", "ali")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	if got.ID != "profile_1" || got.Name != "Alice" || got.Nickname != "ali" {
		t.Fatalf("unexpected profile: %+v", got)
	}
	if len(repo.saved) != 1 {
		t.Fatalf("unexpected save count: %d", len(repo.saved))
	}
}

func TestService_Register_InvalidInput(t *testing.T) {
	t.Parallel()

	repo := &stubProfileRepository{}
	svc := NewService(repo, fixedIDGen{id: "profile_1"})

	_, err := svc.Register(context.Background(), " ", "nick")
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got: %v", err)
	}
	if len(repo.saved) != 0 {
		t.Fatalf("profile should not be saved on invalid input")
	}
}

func TestService_Register_SaveError(t *testing.T) {
	t.Parallel()

	repo := &stubProfileRepository{saveErr: errors.New("db error")}
	svc := NewService(repo, fixedIDGen{id: "profile_1"})

	_, err := svc.Register(context.Background(), "Alice", "ali")
	if err == nil {
		t.Fatalf("expected error")
	}
}
