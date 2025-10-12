package user

import (
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/sky0621/techcv/manager/backend/internal/domain"
)

func TestNewPassword(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{
			name:  "valid password",
			input: "Passw0rd",
		},
		{
			name:      "too short",
			input:     "Pw1",
			wantError: true,
		},
		{
			name:      "missing digit",
			input:     "Password",
			wantError: true,
		},
		{
			name:      "missing letter",
			input:     "12345678",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pwd, err := NewPassword(tt.input)
			if tt.wantError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}

				var appErr *domain.AppError
				if !domain.IsAppError(err) || !errors.As(err, &appErr) {
					t.Fatalf("expected app error, got %v", err)
				}

				if appErr.Code != "INVALID_PASSWORD" {
					t.Fatalf("unexpected code: %s", appErr.Code)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			hash, err := pwd.Hash()
			if err != nil {
				t.Fatalf("hash error: %v", err)
			}

			if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(tt.input)); err != nil {
				t.Fatalf("bcrypt compare failed: %v", err)
			}
		})
	}
}
