package user

import (
	"testing"
	"time"
)

func TestNewVerificationToken(t *testing.T) {
	email, err := NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("unexpected email error: %v", err)
	}

	now := time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC)
	token, err := NewVerificationToken(email, "hashed", now, 24*time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if token.ID() == "" {
		t.Fatalf("expected id to be generated")
	}

	if token.Token() == "" {
		t.Fatalf("expected token to be generated")
	}

	if token.ExpiresAt() != now.Add(24*time.Hour) {
		t.Fatalf("unexpected expiresAt: %v", token.ExpiresAt())
	}

	if token.CreatedAt() != now {
		t.Fatalf("unexpected createdAt: %v", token.CreatedAt())
	}

	if token.Email().String() != email.String() {
		t.Fatalf("mismatched email")
	}

	if token.PasswordHash() != "hashed" {
		t.Fatalf("unexpected password hash")
	}

	if token.IsExpired(now.Add(23 * time.Hour)) {
		t.Fatalf("should not be expired yet")
	}

	if !token.IsExpired(now.Add(25 * time.Hour)) {
		t.Fatalf("expected token to be expired")
	}
}
