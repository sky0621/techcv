package mysql

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestHealthRepositoryPing(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT 1 AS ok")).
		WillReturnRows(sqlmock.NewRows([]string{"ok"}).AddRow(int32(1)))
	mock.ExpectClose()
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Fatalf("failed to close db: %v", closeErr)
		}
		if expectationsErr := mock.ExpectationsWereMet(); expectationsErr != nil {
			t.Fatalf("unmet expectations: %v", expectationsErr)
		}
	}()

	repo := NewHealthRepository(db)
	if err := repo.Ping(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
