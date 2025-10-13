package publicurl

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sky0621/techcv/manager/backend/internal/domain"
)

type mockRepository struct {
	listResult         []domain.PublicURL
	listErr            error
	getActiveResponses []*domain.PublicURL
	getActiveErr       error
	getActiveCalls     int
	createdKeys        []string
	createErr          error
	deactivatedIDs     []uint64
	deactivateErr      error
}

func (m *mockRepository) Create(ctx context.Context, urlKey string) (uint64, error) {
	if m.createErr != nil {
		return 0, m.createErr
	}
	m.createdKeys = append(m.createdKeys, urlKey)
	return uint64(len(m.createdKeys)), nil
}

func (m *mockRepository) GetActive(ctx context.Context) (*domain.PublicURL, error) {
	if m.getActiveErr != nil {
		return nil, m.getActiveErr
	}
	if m.getActiveCalls >= len(m.getActiveResponses) {
		return nil, nil
 	}
	result := m.getActiveResponses[m.getActiveCalls]
	m.getActiveCalls++
	if result == nil {
		return nil, nil
	}
	return result, nil
}

func (m *mockRepository) List(ctx context.Context) ([]domain.PublicURL, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.listResult, nil
}

func (m *mockRepository) Deactivate(ctx context.Context, id uint64) error {
	if m.deactivateErr != nil {
		return m.deactivateErr
	}
	m.deactivatedIDs = append(m.deactivatedIDs, id)
	return nil
}

func TestList(t *testing.T) {
	expected := []domain.PublicURL{
		{ID: 1, URLKey: "first"},
		{ID: 2, URLKey: "second"},
	}

	repo := &mockRepository{
		listResult: expected,
	}

	usecase := New(repo)

	results, err := usecase.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != len(expected) {
		t.Fatalf("unexpected result length: got %d, want %d", len(results), len(expected))
	}

	for i := range expected {
		if results[i].ID != expected[i].ID || results[i].URLKey != expected[i].URLKey {
			t.Fatalf("unexpected element at %d: got %+v, want %+v", i, results[i], expected[i])
		}
	}
}

func TestGenerateWithExistingActive(t *testing.T) {
	now := time.Now()
	existing := &domain.PublicURL{
		ID:        1,
		URLKey:    "existing",
		IsActive:  true,
		CreatedAt: now.Add(-time.Hour),
		UpdatedAt: now.Add(-time.Minute),
	}
	created := &domain.PublicURL{
		ID:        2,
		URLKey:    "generated-key",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	repo := &mockRepository{
		getActiveResponses: []*domain.PublicURL{existing, created},
	}

	usecase := New(repo)
	usecase.keygen = func() (string, error) {
		return "generated-key", nil
	}

	result, err := usecase.Generate(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatalf("expected result, got nil")
	}

	if result.URLKey != created.URLKey {
		t.Fatalf("unexpected URL key: got %q, want %q", result.URLKey, created.URLKey)
	}

	if len(repo.deactivatedIDs) != 1 || repo.deactivatedIDs[0] != existing.ID {
		t.Fatalf("expected deactivate to be called with %d, got %+v", existing.ID, repo.deactivatedIDs)
	}

	if len(repo.createdKeys) != 1 || repo.createdKeys[0] != "generated-key" {
		t.Fatalf("expected create to be called with generated-key, got %+v", repo.createdKeys)
	}
}

func TestGenerateCreateError(t *testing.T) {
	repo := &mockRepository{
		getActiveResponses: []*domain.PublicURL{nil},
		createErr:          errors.New("insert failed"),
	}

	usecase := New(repo)
	usecase.keygen = func() (string, error) {
		return "key", nil
	}

	_, err := usecase.Generate(context.Background())
	if err == nil {
		t.Fatalf("expected error but got nil")
	}

	appErr, ok := err.(*domain.AppError)
	if !ok {
		t.Fatalf("expected AppError, got %T", err)
	}

	if appErr.Code != "public_url.create_failed" {
		t.Fatalf("unexpected error code: got %q, want %q", appErr.Code, "public_url.create_failed")
	}
}
