package auth

import (
	"context"
	"net/url"
	"time"

	"github.com/sky0621/techcv/manager/backend/internal/domain"
	"github.com/sky0621/techcv/manager/backend/internal/domain/user"
)

// RegisterInput captures the data required to initiate registration.
type RegisterInput struct {
	Email                string
	Password             string
	PasswordConfirmation string
}

// RegisterOutput represents the response of a successful registration initiation.
type RegisterOutput struct {
	Message   string
	ExpiresAt time.Time
}

// RegisterUsecase coordinates user registration via email verification.
type RegisterUsecase struct {
	users  user.UserRepository
	tokens user.VerificationTokenRepository
	mailer Mailer
	clock  Clock
	config RegisterConfig
}

// NewRegisterUsecase builds a RegisterUsecase with the given dependencies.
func NewRegisterUsecase(
	users user.UserRepository,
	tokens user.VerificationTokenRepository,
	mailer Mailer,
	clock Clock,
	config RegisterConfig,
) *RegisterUsecase {
	if config.VerificationTTL == 0 {
		config.VerificationTTL = DefaultVerificationTTL
	}
	return &RegisterUsecase{
		users:  users,
		tokens: tokens,
		mailer: mailer,
		clock:  clock,
		config: config,
	}
}

// Execute performs the registration flow.
func (uc *RegisterUsecase) Execute(ctx context.Context, in RegisterInput) (RegisterOutput, error) {
	email, err := user.NewEmail(in.Email)
	if err != nil {
		return RegisterOutput{}, err
	}

	password, err := user.NewPassword(in.Password)
	if err != nil {
		return RegisterOutput{}, err
	}

	if in.Password != in.PasswordConfirmation {
		detail := domain.ErrorDetail{Field: "password_confirmation", Code: domain.ErrorCodePasswordMismatch, Message: "確認用パスワードが一致しません"}
		return RegisterOutput{}, domain.NewValidation(domain.ErrorCodePasswordMismatch, "パスワードが一致しません").WithDetails(detail)
	}

	exists, err := uc.users.ExistsByEmail(ctx, email)
	if err != nil {
		return RegisterOutput{}, domain.NewInternal(domain.ErrorCodeUserLookupFailed, "ユーザー情報の取得に失敗しました", err)
	}

	if exists {
		detail := domain.ErrorDetail{Field: "email", Code: domain.ErrorCodeEmailAlreadyRegistered, Message: "このメールアドレスは既に登録されています"}
		return RegisterOutput{}, domain.NewValidation(domain.ErrorCodeEmailAlreadyRegistered, "このメールアドレスは既に登録されています").WithDetails(detail)
	}

	if err := uc.tokens.DeleteByEmail(ctx, email); err != nil {
		return RegisterOutput{}, domain.NewInternal(domain.ErrorCodeTokenCleanupFailed, "確認トークンの初期化に失敗しました", err)
	}

	hashed, err := password.Hash()
	if err != nil {
		return RegisterOutput{}, err
	}

	now := uc.clock.Now()
	token, err := user.NewVerificationToken(email, hashed, now, uc.config.VerificationTTL)
	if err != nil {
		return RegisterOutput{}, err
	}

	if err := uc.tokens.Save(ctx, token); err != nil {
		return RegisterOutput{}, domain.NewInternal(domain.ErrorCodeTokenSaveFailed, "確認トークンの保存に失敗しました", err)
	}

	verificationURL, err := buildVerificationURL(uc.config.VerificationURLBase, token.Token())
	if err != nil {
		return RegisterOutput{}, domain.NewInternal(domain.ErrorCodeVerificationURLError, "確認メールのURL生成に失敗しました", err)
	}

	if err := uc.mailer.SendVerificationEmail(ctx, email, verificationURL, token.ExpiresAt()); err != nil {
		return RegisterOutput{}, domain.NewInternal(domain.ErrorCodeEmailSendFailed, "確認メールの送信に失敗しました", err)
	}

	return RegisterOutput{
		Message:   "確認メールを送信しました。メールに記載されたリンクをクリックして登録を完了してください",
		ExpiresAt: token.ExpiresAt(),
	}, nil
}

func buildVerificationURL(base string, token string) (string, error) {
	if base == "" {
		return "", domain.NewInternal(domain.ErrorCodeVerificationURLMissing, "確認用URLのベースが設定されていません", nil)
	}

	parsed, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	q := parsed.Query()
	q.Set("token", token)
	parsed.RawQuery = q.Encode()
	return parsed.String(), nil
}
