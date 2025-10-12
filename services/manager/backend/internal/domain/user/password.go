package user

import (
	"unicode"

	"golang.org/x/crypto/bcrypt"

	"github.com/sky0621/techcv/manager/backend/internal/domain"
)

const invalidPasswordMessage = "パスワードは8文字以上で、英字と数字を含む必要があります"

// Password represents a validated password prior to hashing.
type Password struct {
	value string
}

// NewPassword validates the given raw password against domain rules.
func NewPassword(raw string) (Password, error) {
	if len(raw) < 8 {
		return Password{}, domain.NewValidation("INVALID_PASSWORD", invalidPasswordMessage)
	}

	var hasLetter, hasDigit bool
	for _, r := range raw {
		if unicode.IsLetter(r) {
			hasLetter = true
		}
		if unicode.IsDigit(r) {
			hasDigit = true
		}
	}

	if !hasLetter || !hasDigit {
		return Password{}, domain.NewValidation("INVALID_PASSWORD", invalidPasswordMessage)
	}

	return Password{value: raw}, nil
}

// Hash returns the bcrypt hash of the password.
func (p Password) Hash() (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(p.value), bcrypt.DefaultCost)
	if err != nil {
		return "", domain.NewInternal("PASSWORD_HASH_FAILED", "パスワードのハッシュ化に失敗しました", err)
	}
	return string(hashed), nil
}
