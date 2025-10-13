package user

import (
	"time"

	"github.com/sky0621/techcv/manager/backend/internal/domain"
	"github.com/sky0621/techcv/manager/backend/internal/domain/uuidv7"
)

// VerificationToken stores information required to verify user registration.
type VerificationToken struct {
	id           string
	email        Email
	token        string
	passwordHash string
	expiresAt    time.Time
	createdAt    time.Time
}

// NewVerificationToken creates a verification token with the provided TTL.
func NewVerificationToken(email Email, passwordHash string, now time.Time, ttl time.Duration) (VerificationToken, error) {
	if passwordHash == "" {
		return VerificationToken{}, domain.NewInternal(domain.ErrorCodeInvalidPasswordHash, "無効なパスワードハッシュです", nil)
	}

	id, err := uuidv7.NewString()
	if err != nil {
		return VerificationToken{}, domain.NewInternal(domain.ErrorCodeUUIDGenerationFailed, "トークンIDの生成に失敗しました", err)
	}

	tokenValue, err := uuidv7.NewString()
	if err != nil {
		return VerificationToken{}, domain.NewInternal(domain.ErrorCodeUUIDGenerationFailed, "確認トークンの生成に失敗しました", err)
	}

	createdAt := now.UTC().Truncate(time.Microsecond)
	expiresAt := createdAt.Add(ttl)

	return VerificationToken{
		id:           id,
		email:        email,
		token:        tokenValue,
		passwordHash: passwordHash,
		expiresAt:    expiresAt,
		createdAt:    createdAt,
	}, nil
}

// ID returns the internal identifier for the token.
func (t VerificationToken) ID() string {
	return t.id
}

// Email returns the associated email.
func (t VerificationToken) Email() Email {
	return t.email
}

// Token returns the externally visible token string.
func (t VerificationToken) Token() string {
	return t.token
}

// PasswordHash returns the hashed password tied to the token.
func (t VerificationToken) PasswordHash() string {
	return t.passwordHash
}

// ExpiresAt returns the expiration timestamp.
func (t VerificationToken) ExpiresAt() time.Time {
	return t.expiresAt
}

// CreatedAt returns the creation timestamp.
func (t VerificationToken) CreatedAt() time.Time {
	return t.createdAt
}

// IsExpired reports whether the token is expired relative to the supplied time.
func (t VerificationToken) IsExpired(reference time.Time) bool {
	return reference.UTC().After(t.expiresAt)
}
