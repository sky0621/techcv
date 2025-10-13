package user

import "context"

// UserRepository defines persistence operations for user aggregates.
type UserRepository interface {
	ExistsByEmail(ctx context.Context, email Email) (bool, error)
	Create(ctx context.Context, user User) error
	GetByEmail(ctx context.Context, email Email) (User, error)
}

// VerificationTokenRepository defines persistence operations for email verification tokens.
type VerificationTokenRepository interface {
	Save(ctx context.Context, token VerificationToken) error
	FindByToken(ctx context.Context, token string) (VerificationToken, error)
	DeleteByToken(ctx context.Context, token string) error
	DeleteByEmail(ctx context.Context, email Email) error
}
