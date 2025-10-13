// Package auth defines use cases and supporting interfaces for authentication flows.
package auth

import (
	"context"
	"time"

	"github.com/sky0621/techcv/manager/backend/internal/domain/user"
)

// Clock abstracts the source of current time for easier testing.
type Clock interface {
	Now() time.Time
}

// Mailer sends verification emails to guests.
type Mailer interface {
	SendVerificationEmail(ctx context.Context, email user.Email, verificationURL string, expiresAt time.Time) error
}

// AuthTokenIssuer creates authentication tokens after verification.
type AuthTokenIssuer interface {
	Issue(ctx context.Context, user user.User) (string, error)
}

// TransactionManager executes operations within a transaction boundary.
type TransactionManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// RegisterConfig holds configuration for the registration process.
type RegisterConfig struct {
	VerificationURLBase string
	VerificationTTL     time.Duration
}

// DefaultVerificationTTL represents the default lifetime for verification tokens.
const DefaultVerificationTTL = 24 * time.Hour
