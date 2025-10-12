package auth

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/sky0621/techcv/manager/backend/internal/domain"
	"github.com/sky0621/techcv/manager/backend/internal/domain/user"
)

func TestRegisterUsecase_Success(t *testing.T) {
	userRepo := newFakeUserRepo()
	tokenRepo := newFakeTokenRepo()
	mailer := &fakeMailer{}
	clock := fixedClock{now: time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)}

	uc := NewRegisterUsecase(
		userRepo,
		tokenRepo,
		mailer,
		clock,
		RegisterConfig{VerificationURLBase: "https://example.com/verify", VerificationTTL: time.Hour},
	)

	out, err := uc.Execute(context.Background(), RegisterInput{
		Email:                "guest@example.com",
		Password:             "Passw0rd",
		PasswordConfirmation: "Passw0rd",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out.Message == "" {
		t.Fatalf("expected message to be set")
	}

	expectedExpiry := clock.now.Add(time.Hour)
	if !out.ExpiresAt.Equal(expectedExpiry) {
		t.Fatalf("expiresAt mismatch: got %v want %v", out.ExpiresAt, expectedExpiry)
	}

	if len(tokenRepo.tokens) != 1 {
		t.Fatalf("expected token saved")
	}

	saved := tokenRepo.tokens[0]
	if saved.Email().String() != "guest@example.com" {
		t.Fatalf("unexpected email stored: %s", saved.Email().String())
	}

	if saved.PasswordHash() == "Passw0rd" || saved.PasswordHash() == "" {
		t.Fatalf("password hash not stored correctly")
	}

	if !strings.Contains(mailer.lastURL, saved.Token()) {
		t.Fatalf("verification url should contain token")
	}

	if mailer.sentTo.String() != "guest@example.com" {
		t.Fatalf("unexpected recipient: %s", mailer.sentTo.String())
	}
}

func TestRegisterUsecase_EmailAlreadyExists(t *testing.T) {
	userRepo := newFakeUserRepo()
	email, _ := user.NewEmail("guest@example.com")
	userRepo.existing[email.String()] = true

	uc := NewRegisterUsecase(
		userRepo,
		newFakeTokenRepo(),
		&fakeMailer{},
		fixedClock{now: time.Now()},
		RegisterConfig{VerificationURLBase: "https://example.com/verify"},
	)

	_, err := uc.Execute(context.Background(), RegisterInput{
		Email:                "guest@example.com",
		Password:             "Passw0rd",
		PasswordConfirmation: "Passw0rd",
	})
	if err == nil {
		t.Fatalf("expected error")
	}

	var appErr *domain.AppError
	if !errors.As(err, &appErr) || appErr.Code != "EMAIL_ALREADY_REGISTERED" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRegisterUsecase_PasswordMismatch(t *testing.T) {
	uc := NewRegisterUsecase(
		newFakeUserRepo(),
		newFakeTokenRepo(),
		&fakeMailer{},
		fixedClock{now: time.Now()},
		RegisterConfig{VerificationURLBase: "https://example.com/verify"},
	)

	_, err := uc.Execute(context.Background(), RegisterInput{
		Email:                "guest@example.com",
		Password:             "Passw0rd",
		PasswordConfirmation: "wrong",
	})

	var appErr *domain.AppError
	if err == nil || !errors.As(err, &appErr) || appErr.Code != "PASSWORD_MISMATCH" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRegisterUsecase_InvalidEmail(t *testing.T) {
	uc := NewRegisterUsecase(
		newFakeUserRepo(),
		newFakeTokenRepo(),
		&fakeMailer{},
		fixedClock{now: time.Now()},
		RegisterConfig{VerificationURLBase: "https://example.com/verify"},
	)

	_, err := uc.Execute(context.Background(), RegisterInput{
		Email:                "not-email",
		Password:             "Passw0rd",
		PasswordConfirmation: "Passw0rd",
	})

	var appErr *domain.AppError
	if err == nil || !errors.As(err, &appErr) || appErr.Code != "INVALID_EMAIL_FORMAT" {
		t.Fatalf("unexpected error: %v", err)
	}
}

type fixedClock struct {
	now time.Time
}

func (c fixedClock) Now() time.Time {
	return c.now
}

type fakeMailer struct {
	sentTo   user.Email
	lastURL  string
	calls    int
	lastExpr time.Time
	fail     bool
}

func (m *fakeMailer) SendVerificationEmail(_ context.Context, email user.Email, verificationURL string, expiresAt time.Time) error {
	if m.fail {
		return errors.New("send failed")
	}
	m.sentTo = email
	m.lastURL = verificationURL
	m.lastExpr = expiresAt
	m.calls++
	return nil
}

type fakeUserRepo struct {
	existing map[string]bool
	users    map[string]user.User
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{
		existing: make(map[string]bool),
		users:    make(map[string]user.User),
	}
}

func (r *fakeUserRepo) ExistsByEmail(_ context.Context, email user.Email) (bool, error) {
	return r.existing[email.String()], nil
}

func (r *fakeUserRepo) Create(_ context.Context, u user.User) error {
	r.existing[u.Email().String()] = true
	r.users[u.Email().String()] = u
	return nil
}

func (r *fakeUserRepo) GetByEmail(_ context.Context, email user.Email) (user.User, error) {
	if u, ok := r.users[email.String()]; ok {
		return u, nil
	}
	return user.User{}, domain.NewNotFound("USER_NOT_FOUND", "ユーザーが見つかりません")
}

type fakeTokenRepo struct {
	tokens []user.VerificationToken
}

func newFakeTokenRepo() *fakeTokenRepo {
	return &fakeTokenRepo{tokens: []user.VerificationToken{}}
}

func (r *fakeTokenRepo) Save(_ context.Context, token user.VerificationToken) error {
	r.tokens = append(r.tokens, token)
	return nil
}

func (r *fakeTokenRepo) FindByToken(_ context.Context, token string) (user.VerificationToken, error) {
	for _, t := range r.tokens {
		if t.Token() == token {
			return t, nil
		}
	}
	return user.VerificationToken{}, domain.NewNotFound("TOKEN_NOT_FOUND", "確認トークンが見つかりません")
}

func (r *fakeTokenRepo) DeleteByToken(_ context.Context, token string) error {
	filtered := r.tokens[:0]
	for _, t := range r.tokens {
		if t.Token() != token {
			filtered = append(filtered, t)
		}
	}
	r.tokens = filtered
	return nil
}

func (r *fakeTokenRepo) DeleteByEmail(_ context.Context, email user.Email) error {
	filtered := r.tokens[:0]
	for _, t := range r.tokens {
		if t.Email().String() != email.String() {
			filtered = append(filtered, t)
		}
	}
	r.tokens = filtered
	return nil
}
