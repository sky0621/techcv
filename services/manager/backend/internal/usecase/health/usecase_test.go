package health

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sky0621/techcv/manager/backend/internal/domain"
)

type stubRepository struct {
	err error
}

func (s stubRepository) Ping(context.Context) error {
	return s.err
}

func TestUsecaseCheckSuccess(t *testing.T) {
	repo := stubRepository{}
	usecase := New(repo)

	expectedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	usecase.now = func() time.Time {
		return expectedTime
	}

	status, err := usecase.Check(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if status.Status != "ok" {
		t.Fatalf("unexpected status: got %q, want %q", status.Status, "ok")
	}

	if !status.CheckedAt.Equal(expectedTime) {
		t.Fatalf("unexpected checkedAt: got %v, want %v", status.CheckedAt, expectedTime)
	}
}

func TestUsecaseCheckFailure(t *testing.T) {
	wantErr := errors.New("ping failed")
	repo := stubRepository{err: wantErr}
	usecase := New(repo)

	status, err := usecase.Check(context.Background())
	if err == nil {
		t.Fatalf("expected error but got nil")
	}

	if status != nil {
		t.Fatalf("expected nil status on failure")
	}

	appErr, ok := err.(*domain.AppError)
	if !ok {
		t.Fatalf("expected AppError, got %T", err)
	}

	if !errors.Is(appErr, wantErr) {
		t.Fatalf("expected wrapped error %v, got %v", wantErr, appErr.Err)
	}

	if appErr.Code != "health.database_unavailable" {
		t.Fatalf("unexpected error code: got %q, want %q", appErr.Code, "health.database_unavailable")
	}
}
