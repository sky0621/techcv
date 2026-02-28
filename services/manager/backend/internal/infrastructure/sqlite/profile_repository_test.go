package sqlite

import (
	"context"
	"testing"

	"github.com/sky0621/techcv/services/manager/backend/internal/domain"
)

func TestProfileRepository_Save(t *testing.T) {
	t.Parallel()

	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	if err := EnsureSchema(context.Background(), db); err != nil {
		t.Fatalf("EnsureSchema() error = %v", err)
	}

	repo := NewProfileRepository(db)
	profile := domain.Profile{
		ID:       "profile_1",
		Name:     "Alice",
		Nickname: "ali",
	}

	if err := repo.Save(context.Background(), profile); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	var got domain.Profile
	row := db.QueryRowContext(context.Background(), `SELECT id, name, nickname FROM profiles WHERE id = ?`, profile.ID)
	if err := row.Scan(&got.ID, &got.Name, &got.Nickname); err != nil {
		t.Fatalf("QueryRow().Scan() error = %v", err)
	}

	if got != profile {
		t.Fatalf("saved profile mismatch: got %+v, want %+v", got, profile)
	}
}
