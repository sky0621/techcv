package user

import (
	"errors"
	"testing"
	"time"

	"github.com/sky0621/techcv/manager/backend/internal/domain"
)

func TestNewUserFromFirebase(t *testing.T) {
	now := time.Date(2025, time.January, 1, 12, 0, 0, 0, time.UTC)
	display := "Example User"

	entity, err := NewUserFromFirebase(FirebaseUserData{
		UID:           "firebase-uid",
		Email:         "user@example.com",
		EmailVerified: true,
		DisplayName:   &display,
		ProviderID:    "google.com",
	}, now, func() (string, error) {
		return "018f48dc-1b1c-7738-8000-000000000001", nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if entity.ID() == "" {
		t.Fatalf("expected generated id")
	}
	if entity.FirebaseUID().String() != "firebase-uid" {
		t.Fatalf("unexpected firebase uid: %s", entity.FirebaseUID().String())
	}
	if entity.Email().String() != "user@example.com" {
		t.Fatalf("unexpected email: %s", entity.Email().String())
	}
	if !entity.EmailVerified() {
		t.Fatalf("expected email verified flag to be true")
	}
	if entity.DisplayName() == nil || *entity.DisplayName() != display {
		t.Fatalf("unexpected display name: %+v", entity.DisplayName())
	}
	if entity.ProviderID() != "google.com" {
		t.Fatalf("unexpected provider id: %s", entity.ProviderID())
	}
	if entity.LastLoginAt() == nil || !entity.LastLoginAt().Equal(now) {
		t.Fatalf("unexpected last login at: %+v", entity.LastLoginAt())
	}
	if entity.CreatedAt() != now || entity.UpdatedAt() != now {
		t.Fatalf("unexpected timestamps: created %v updated %v", entity.CreatedAt(), entity.UpdatedAt())
	}
}

func TestNewUserFromFirebaseInvalidProvider(t *testing.T) {
	_, err := NewUserFromFirebase(FirebaseUserData{
		UID:        "firebase-uid",
		Email:      "user@example.com",
		ProviderID: "",
	}, time.Now(), func() (string, error) {
		return "018f48dc-1b1c-7738-8000-000000000001", nil
	})
	if err == nil {
		t.Fatalf("expected validation error")
	}

	var appErr *domain.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected app error, got %v", err)
	}
	if appErr.Code != domain.ErrorCodeValidationError {
		t.Fatalf("unexpected error code: %s", appErr.Code)
	}
}
