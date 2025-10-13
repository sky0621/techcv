package user

import (
	"fmt"
	"time"

	"github.com/sky0621/techcv/manager/backend/internal/domain"
	"github.com/sky0621/techcv/manager/backend/internal/domain/uuidv7"
)

// User represents the aggregate root for an authenticated user.
type User struct {
	id              string
	email           Email
	passwordHash    string
	name            *string
	bio             *string
	isActive        bool
	emailVerifiedAt time.Time
	lastLoginAt     *time.Time
	createdAt       time.Time
	updatedAt       time.Time
}

// NewUser constructs a fully verified user aggregate.
func NewUser(email Email, passwordHash string, now time.Time) (User, error) {
	if passwordHash == "" {
		return User{}, domain.NewInternal(domain.ErrorCodeInvalidPasswordHash, "無効なパスワードハッシュです", fmt.Errorf("empty password hash"))
	}

	id, err := uuidv7.NewString()
	if err != nil {
		return User{}, domain.NewInternal(domain.ErrorCodeUUIDGenerationFailed, "ユーザーIDの生成に失敗しました", err)
	}

	ts := now.UTC().Truncate(time.Microsecond)

	return User{
		id:              id,
		email:           email,
		passwordHash:    passwordHash,
		isActive:        true,
		emailVerifiedAt: ts,
		lastLoginAt:     &ts,
		createdAt:       ts,
		updatedAt:       ts,
	}, nil
}

// ID returns the user's identifier.
func (u User) ID() string {
	return u.id
}

// Email returns the user's email.
func (u User) Email() Email {
	return u.email
}

// PasswordHash returns the hashed password.
func (u User) PasswordHash() string {
	return u.passwordHash
}

// Name returns the optional profile name.
func (u User) Name() *string {
	return u.name
}

// Bio returns the optional bio.
func (u User) Bio() *string {
	return u.bio
}

// IsActive indicates whether the user is active.
func (u User) IsActive() bool {
	return u.isActive
}

// EmailVerifiedAt returns the timestamp when the email was verified.
func (u User) EmailVerifiedAt() time.Time {
	return u.emailVerifiedAt
}

// CreatedAt returns the creation timestamp.
func (u User) CreatedAt() time.Time {
	return u.createdAt
}

// UpdatedAt returns the last update timestamp.
func (u User) UpdatedAt() time.Time {
	return u.updatedAt
}

// LastLoginAt returns the timestamp of the previous login, if recorded.
func (u User) LastLoginAt() *time.Time {
	return u.lastLoginAt
}

// WithLastLogin updates the last login timestamp and returns a copy.
func (u User) WithLastLogin(t time.Time) User {
	ts := t.UTC().Truncate(time.Microsecond)
	u.lastLoginAt = &ts
	u.updatedAt = ts
	return u
}
