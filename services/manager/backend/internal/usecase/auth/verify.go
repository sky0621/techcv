package auth

import (
	"context"
	"strings"
	"time"

	"github.com/sky0621/techcv/manager/backend/internal/domain"
	"github.com/sky0621/techcv/manager/backend/internal/domain/user"
)

// VerifyInput captures the token supplied by the guest.
type VerifyInput struct {
	Token string
}

// VerifiedUser represents the user data returned after successful verification.
type VerifiedUser struct {
	ID              string
	Email           string
	Name            *string
	Bio             *string
	IsActive        bool
	EmailVerifiedAt time.Time
	LastLoginAt     *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// VerifyOutput bundles the results of a verification attempt.
type VerifyOutput struct {
	Message   string
	AuthToken string
	User      VerifiedUser
}

// VerifyUsecase finalises registration by validating and consuming verification tokens.
type VerifyUsecase struct {
	users  user.UserRepository
	tokens user.VerificationTokenRepository
	tx     TransactionManager
	clock  Clock
	issuer AuthTokenIssuer
}

// NewVerifyUsecase constructs a VerifyUsecase instance.
func NewVerifyUsecase(
	users user.UserRepository,
	tokens user.VerificationTokenRepository,
	tx TransactionManager,
	clock Clock,
	issuer AuthTokenIssuer,
) *VerifyUsecase {
	return &VerifyUsecase{
		users:  users,
		tokens: tokens,
		tx:     tx,
		clock:  clock,
		issuer: issuer,
	}
}

// Execute validates the token and creates a fully verified user account.
func (uc *VerifyUsecase) Execute(ctx context.Context, in VerifyInput) (VerifyOutput, error) {
	tokenValue := strings.TrimSpace(in.Token)
	if tokenValue == "" {
		detail := domain.ErrorDetail{Field: "token", Code: "INVALID_VERIFICATION_TOKEN", Message: "確認トークンを指定してください"}
		return VerifyOutput{}, domain.NewValidation("INVALID_VERIFICATION_TOKEN", "確認トークンを指定してください").WithDetails(detail)
	}

	record, err := uc.tokens.FindByToken(ctx, tokenValue)
	if err != nil {
		if domain.IsAppError(err) {
			return VerifyOutput{}, err
		}
		return VerifyOutput{}, domain.NewInternal("TOKEN_LOOKUP_FAILED", "確認トークンの取得に失敗しました", err)
	}

	now := uc.clock.Now()
	if record.IsExpired(now) {
		_ = uc.tokens.DeleteByToken(ctx, tokenValue)
		detail := domain.ErrorDetail{Field: "token", Code: "VERIFICATION_TOKEN_EXPIRED", Message: "確認リンクが無効または期限切れです"}
		return VerifyOutput{}, domain.NewValidation("VERIFICATION_TOKEN_EXPIRED", "確認リンクが無効または期限切れです。再度登録をお試しください").WithDetails(detail)
	}

	exists, err := uc.users.ExistsByEmail(ctx, record.Email())
	if err != nil {
		return VerifyOutput{}, domain.NewInternal("USER_LOOKUP_FAILED", "ユーザー情報の取得に失敗しました", err)
	}

	if exists {
		detail := domain.ErrorDetail{Field: "email", Code: "EMAIL_ALREADY_REGISTERED", Message: "このメールアドレスは既に登録されています"}
		return VerifyOutput{}, domain.NewValidation("EMAIL_ALREADY_REGISTERED", "このメールアドレスは既に登録されています").WithDetails(detail)
	}

	newUser, err := user.NewUser(record.Email(), record.PasswordHash(), now)
	if err != nil {
		return VerifyOutput{}, err
	}

	if err := uc.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.users.Create(txCtx, newUser); err != nil {
			return domain.NewInternal("USER_CREATE_FAILED", "ユーザーの作成に失敗しました", err)
		}

		if err := uc.tokens.DeleteByToken(txCtx, tokenValue); err != nil {
			return domain.NewInternal("TOKEN_DELETE_FAILED", "確認トークンの削除に失敗しました", err)
		}
		return nil
	}); err != nil {
		return VerifyOutput{}, err
	}

	authToken, err := uc.issuer.Issue(ctx, newUser)
	if err != nil {
		return VerifyOutput{}, domain.NewInternal("AUTH_TOKEN_ISSUE_FAILED", "認証トークンの発行に失敗しました", err)
	}

	return VerifyOutput{
		Message:   "登録が完了しました",
		AuthToken: authToken,
		User:      toVerifiedUser(newUser),
	}, nil
}

func toVerifiedUser(u user.User) VerifiedUser {
	return VerifiedUser{
		ID:              u.ID(),
		Email:           u.Email().String(),
		Name:            u.Name(),
		Bio:             u.Bio(),
		IsActive:        u.IsActive(),
		EmailVerifiedAt: u.EmailVerifiedAt(),
		LastLoginAt:     u.LastLoginAt(),
		CreatedAt:       u.CreatedAt(),
		UpdatedAt:       u.UpdatedAt(),
	}
}
