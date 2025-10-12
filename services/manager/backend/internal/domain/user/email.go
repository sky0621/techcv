package user

import (
	"net/mail"
	"strings"

	"github.com/sky0621/techcv/manager/backend/internal/domain"
)

const invalidEmailMessage = "メールアドレスの形式が正しくありません"

// Email models an email address with validation.
type Email struct {
	value string
}

// NewEmail validates and constructs an Email value object.
func NewEmail(raw string) (Email, error) {
	trimmed := strings.TrimSpace(raw)
	detail := domain.ErrorDetail{Field: "email", Code: "INVALID_EMAIL_FORMAT", Message: invalidEmailMessage}
	if trimmed == "" {
		return Email{}, domain.NewValidation("INVALID_EMAIL_FORMAT", invalidEmailMessage).WithDetails(detail)
	}

	if _, err := mail.ParseAddress(trimmed); err != nil {
		return Email{}, domain.NewValidation("INVALID_EMAIL_FORMAT", invalidEmailMessage).WithDetails(detail)
	}

	normalized := strings.ToLower(trimmed)
	return Email{value: normalized}, nil
}

// String returns the canonical email string.
func (e Email) String() string {
	return e.value
}
