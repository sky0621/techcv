package email

import (
	"context"
	"log/slog"
	"time"

	"github.com/sky0621/techcv/manager/backend/internal/domain/user"
)

// LogMailer sends verification emails by logging their content.
type LogMailer struct {
	logger *slog.Logger
}

// NewLogMailer constructs a logging mailer.
func NewLogMailer(logger *slog.Logger) LogMailer {
	return LogMailer{logger: logger}
}

// SendVerificationEmail records the verification email details in the log.
func (m LogMailer) SendVerificationEmail(_ context.Context, email user.Email, verificationURL string, expiresAt time.Time) error {
	m.logger.Info("verification email dispatched",
		slog.String("email", email.String()),
		slog.String("verification_url", verificationURL),
		slog.Time("expires_at", expiresAt),
	)
	return nil
}
