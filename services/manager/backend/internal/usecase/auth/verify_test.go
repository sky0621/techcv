package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sky0621/techcv/manager/backend/internal/domain"
	"github.com/sky0621/techcv/manager/backend/internal/domain/user"
)

const guestEmailAddress = "guest@example.com"

func TestVerifyUsecase_Success(t *testing.T) {
	userRepo := newFakeUserRepo()
	tokenRepo := newFakeTokenRepo()
	tx := &fakeTxManager{}
	clock := fixedClock{now: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)}
	issuer := &fakeTokenIssuer{}

	email, _ := user.NewEmail(guestEmailAddress)
	token, _ := user.NewVerificationToken(email, "hashed", clock.now.Add(-time.Hour), 24*time.Hour)
	tokenRepo.tokens = append(tokenRepo.tokens, token)

	uc := NewVerifyUsecase(userRepo, tokenRepo, tx, clock, issuer)

	out, err := uc.Execute(context.Background(), VerifyInput{Token: token.Token()})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out.AuthToken != "issued-token" {
		t.Fatalf("unexpected auth token: %s", out.AuthToken)
	}

	if out.User.Email != guestEmailAddress {
		t.Fatalf("unexpected email: %s", out.User.Email)
	}

	if len(tokenRepo.tokens) != 0 {
		t.Fatalf("expected token to be deleted")
	}

	if !tx.called {
		t.Fatalf("expected transaction to be used")
	}
}

func TestVerifyUsecase_TokenExpired(t *testing.T) {
	userRepo := newFakeUserRepo()
	tokenRepo := newFakeTokenRepo()
	tx := &fakeTxManager{}
	clock := fixedClock{now: time.Date(2024, 1, 2, 10, 0, 0, 0, time.UTC)}
	issuer := &fakeTokenIssuer{}

	email, _ := user.NewEmail(guestEmailAddress)
	token, _ := user.NewVerificationToken(email, "hashed", clock.now.Add(-48*time.Hour), 24*time.Hour)
	tokenRepo.tokens = append(tokenRepo.tokens, token)

	uc := NewVerifyUsecase(userRepo, tokenRepo, tx, clock, issuer)

	_, err := uc.Execute(context.Background(), VerifyInput{Token: token.Token()})
	if err == nil {
		t.Fatalf("expected error")
	}

	var appErr *domain.AppError
	if !errors.As(err, &appErr) || appErr.Code != domain.ErrorCodeVerificationTokenExpired {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVerifyUsecase_TokenNotFound(t *testing.T) {
	userRepo := newFakeUserRepo()
	tokenRepo := newFakeTokenRepo()
	tx := &fakeTxManager{}
	clock := fixedClock{now: time.Now()}
	issuer := &fakeTokenIssuer{}

	uc := NewVerifyUsecase(userRepo, tokenRepo, tx, clock, issuer)

	_, err := uc.Execute(context.Background(), VerifyInput{Token: "unknown"})

	var appErr *domain.AppError
	if err == nil || !errors.As(err, &appErr) || appErr.Code != domain.ErrorCodeTokenNotFound {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVerifyUsecase_EmailAlreadyRegistered(t *testing.T) {
	userRepo := newFakeUserRepo()
	tokenRepo := newFakeTokenRepo()
	tx := &fakeTxManager{}
	clock := fixedClock{now: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)}
	issuer := &fakeTokenIssuer{}

	email, _ := user.NewEmail(guestEmailAddress)
	token, _ := user.NewVerificationToken(email, "hashed", clock.now.Add(-time.Hour), 24*time.Hour)
	tokenRepo.tokens = append(tokenRepo.tokens, token)

	userRepo.existing[email.String()] = true

	uc := NewVerifyUsecase(userRepo, tokenRepo, tx, clock, issuer)

	_, err := uc.Execute(context.Background(), VerifyInput{Token: token.Token()})
	if err == nil {
		t.Fatalf("expected error")
	}

	var appErr *domain.AppError
	if !errors.As(err, &appErr) || appErr.Code != domain.ErrorCodeEmailAlreadyRegistered {
		t.Fatalf("unexpected error: %v", err)
	}
}

// fakeTxManager implements TransactionManager for tests.
type fakeTxManager struct {
	called bool
	fail   bool
}

func (m *fakeTxManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	m.called = true
	if m.fail {
		return errors.New("tx failed")
	}
	return fn(ctx)
}

// fakeTokenIssuer issues fake auth tokens in tests.
type fakeTokenIssuer struct {
	fail bool
}

func (i *fakeTokenIssuer) Issue(_ context.Context, _ user.User) (string, error) {
	if i.fail {
		return "", errors.New("issue failed")
	}
	return "issued-token", nil
}
