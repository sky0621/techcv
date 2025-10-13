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
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT 1 AS ok")).
		WillReturnRows(sqlmock.NewRows([]string{"ok"}).AddRow(int32(1)))

	repo := NewHealthRepository(db)
	if err := repo.Ping(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
