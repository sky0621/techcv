package mysql

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

const (
	getActivePublicURLQuery = "-- name: GetActivePublicURL :one\n" +
		"SELECT\n" +
		"  id,\n" +
		"  url_key,\n" +
		"  is_active,\n" +
		"  created_at,\n" +
		"  updated_at\n" +
		"FROM public_urls\n" +
		"WHERE is_active = TRUE\n" +
		"ORDER BY updated_at DESC\n" +
		"LIMIT 1\n"
	createPublicURLQuery = "-- name: CreatePublicURL :execresult\n" +
		"INSERT INTO public_urls (url_key)\n" +
		"VALUES (?)\n"
	listPublicURLsQuery = "-- name: ListPublicURLs :many\n" +
		"SELECT\n" +
		"  id,\n" +
		"  url_key,\n" +
		"  is_active,\n" +
		"  created_at,\n" +
		"  updated_at\n" +
		"FROM public_urls\n" +
		"ORDER BY updated_at DESC\n"
	deactivatePublicURLQuery = "-- name: DeactivatePublicURL :exec\n" +
		"UPDATE public_urls\n" +
		"SET is_active = FALSE,\n" +
		"    updated_at = CURRENT_TIMESTAMP(6)\n" +
		"WHERE id = ?\n"
)

func TestPublicURLRepositoryGetActive(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Fatalf("failed to close db: %v", closeErr)
		}
	}()

	now := time.Now()
	rows := sqlmock.
		NewRows([]string{"id", "url_key", "is_active", "created_at", "updated_at"}).
		AddRow(int64(1), "active-key", true, now, now)

	mock.ExpectQuery(regexp.QuoteMeta(getActivePublicURLQuery)).WillReturnRows(rows)

	repo := NewPublicURLRepository(db)
	result, err := repo.GetActive(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil || result.URLKey != "active-key" {
		t.Fatalf("unexpected result: %+v", result)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestPublicURLRepositoryCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Fatalf("failed to close db: %v", closeErr)
		}
	}()

	mock.ExpectExec(regexp.QuoteMeta(createPublicURLQuery)).
		WithArgs("new-key").
		WillReturnResult(sqlmock.NewResult(10, 1))

	repo := NewPublicURLRepository(db)
	id, err := repo.Create(context.Background(), "new-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if id != 10 {
		t.Fatalf("unexpected id: got %d, want %d", id, 10)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestPublicURLRepositoryList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Fatalf("failed to close db: %v", closeErr)
		}
	}()

	now := time.Now()
	rows := sqlmock.
		NewRows([]string{"id", "url_key", "is_active", "created_at", "updated_at"}).
		AddRow(int64(1), "first", true, now, now).
		AddRow(int64(2), "second", false, now, now)

	mock.ExpectQuery(regexp.QuoteMeta(listPublicURLsQuery)).WillReturnRows(rows)

	repo := NewPublicURLRepository(db)
	results, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("unexpected result length: got %d, want %d", len(results), 2)
	}

	if results[0].URLKey != "first" || !results[0].IsActive {
		t.Fatalf("unexpected first result: %+v", results[0])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestPublicURLRepositoryDeactivate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Fatalf("failed to close db: %v", closeErr)
		}
	}()

	mock.ExpectExec(regexp.QuoteMeta(deactivatePublicURLQuery)).
		WithArgs(int64(5)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	repo := NewPublicURLRepository(db)
	if err := repo.Deactivate(context.Background(), 5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
