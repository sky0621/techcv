package user

import (
	"errors"
	"testing"

	"github.com/sky0621/techcv/manager/backend/internal/domain"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{
			name:  "valid email",
			input: "User@example.com",
		},
		{
			name:      "empty",
			input:     "",
			wantError: true,
		},
		{
			name:      "invalid format",
			input:     "not-an-email",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := NewEmail(tt.input)
			if tt.wantError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}

				var appErr *domain.AppError
				if !domain.IsAppError(err) || !errors.As(err, &appErr) {
					t.Fatalf("expected app error, got %v", err)
				}

				if appErr.Code != domain.ErrorCodeInvalidEmailFormat {
					t.Fatalf("unexpected code: %s", appErr.Code)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if email.String() != "user@example.com" {
				t.Fatalf("expected normalized email, got %s", email.String())
			}
		})
	}
}
