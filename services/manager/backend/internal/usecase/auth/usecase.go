// Package auth provides Firebase-backed authentication use cases.
package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/sky0621/techcv/manager/backend/internal/domain"
	"github.com/sky0621/techcv/manager/backend/internal/domain/user"
)

// FirebaseAuthClient abstracts Firebase authentication operations used by the use case layer.
type FirebaseAuthClient interface {
	VerifyIDToken(ctx context.Context, idToken string) (string, error)
	GetUser(ctx context.Context, uid string) (FirebaseUser, error)
}

// FirebaseUser represents the subset of Firebase user data required by the application.
type FirebaseUser struct {
	UID                  string
	Email                string
	EmailVerified        bool
	DisplayName          *string
	PhotoURL             *string
	PhoneNumber          *string
	ProviderID           string
	FirebaseCreatedAt    *time.Time
	FirebaseLastSignInAt *time.Time
}

// UserRepository defines the persistence contract needed for authentication flows.
type UserRepository interface {
	Create(ctx context.Context, user user.User) error
	FindByID(ctx context.Context, id string) (user.User, error)
	FindByEmail(ctx context.Context, email user.Email) (user.User, error)
	FindByFirebaseUID(ctx context.Context, firebaseUID user.FirebaseUID) (user.User, error)
	UpdateFromFirebase(ctx context.Context, user user.User) error
}

// Clock abstracts time retrieval to ease testing.
type Clock interface {
	Now() time.Time
}

// IDGenerator defines a UUID generator.
type IDGenerator func() (string, error)

// AuthResult represents the API response payload for authentication operations.
type AuthResult struct {
	UserID      string  `json:"userId"`
	FirebaseUID string  `json:"firebaseUid"`
	Email       string  `json:"email"`
	DisplayName *string `json:"displayName,omitempty"`
}

// Service coordinates Firebase authentication flows.
type Service struct {
	auth       FirebaseAuthClient
	repo       UserRepository
	clock      Clock
	generateID IDGenerator
}

// New constructs a Service instance.
func New(auth FirebaseAuthClient, repo UserRepository, clock Clock, generator IDGenerator) *Service {
	return &Service{
		auth:       auth,
		repo:       repo,
		clock:      clock,
		generateID: generator,
	}
}

// Register creates a new application user based on Firebase authentication.
func (s *Service) Register(ctx context.Context, firebaseUID string) (*AuthResult, error) {
	uid, err := user.NewFirebaseUID(firebaseUID)
	if err != nil {
		return nil, err
	}

	firebaseUser, err := s.auth.GetUser(ctx, uid.String())
	if err != nil {
		return nil, wrapRegistrationError(err)
	}

	_, err = s.repo.FindByFirebaseUID(ctx, uid)
	if err == nil {
		return nil, domain.NewConflict(domain.ErrorCodeUserAlreadyExists, "ユーザーが既に登録されています")
	}
	if !isNotFound(err) {
		return nil, wrapRegistrationError(err)
	}

	now := s.clock.Now()
	entity, err := user.NewUserFromFirebase(toDomainFirebaseUser(firebaseUser), now, s.generateID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, entity); err != nil {
		return nil, wrapRegistrationError(err)
	}

	return toAuthResult(entity), nil
}

// Login refreshes user information based on Firebase authentication.
func (s *Service) Login(ctx context.Context, firebaseUID string) (*AuthResult, error) {
	uid, err := user.NewFirebaseUID(firebaseUID)
	if err != nil {
		return nil, err
	}

	firebaseUser, err := s.auth.GetUser(ctx, uid.String())
	if err != nil {
		return nil, wrapLoginError(err)
	}

	existing, err := s.repo.FindByFirebaseUID(ctx, uid)
	if err != nil {
		if isNotFound(err) {
			now := s.clock.Now()
			entity, createErr := user.NewUserFromFirebase(toDomainFirebaseUser(firebaseUser), now, s.generateID)
			if createErr != nil {
				return nil, wrapLoginError(createErr)
			}
			if createErr := s.repo.Create(ctx, entity); createErr != nil {
				return nil, wrapLoginError(createErr)
			}
			return toAuthResult(entity), nil
		}
		return nil, wrapLoginError(err)
	}

	updated, err := existing.UpdateFromFirebase(toDomainFirebaseUser(firebaseUser), s.clock.Now())
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdateFromFirebase(ctx, updated); err != nil {
		if isNotFound(err) {
			return nil, err
		}
		return nil, wrapLoginError(err)
	}

	return toAuthResult(updated), nil
}

func toDomainFirebaseUser(firebaseUser FirebaseUser) user.FirebaseUserData {
	return user.FirebaseUserData{
		UID:                  strings.TrimSpace(firebaseUser.UID),
		Email:                firebaseUser.Email,
		EmailVerified:        firebaseUser.EmailVerified,
		DisplayName:          firebaseUser.DisplayName,
		PhotoURL:             firebaseUser.PhotoURL,
		PhoneNumber:          firebaseUser.PhoneNumber,
		ProviderID:           firebaseUser.ProviderID,
		FirebaseCreatedAt:    firebaseUser.FirebaseCreatedAt,
		FirebaseLastSignInAt: firebaseUser.FirebaseLastSignInAt,
	}
}

func toAuthResult(entity user.User) *AuthResult {
	return &AuthResult{
		UserID:      entity.ID(),
		FirebaseUID: entity.FirebaseUID().String(),
		Email:       entity.Email().String(),
		DisplayName: entity.DisplayName(),
	}
}

func isNotFound(err error) bool {
	var appErr *domain.AppError
	if !errors.As(err, &appErr) {
		return false
	}
	return appErr.StatusCode == http.StatusNotFound
}

func wrapRegistrationError(err error) error {
	if domain.IsAppError(err) {
		return err
	}
	return domain.NewInternal(domain.ErrorCodeRegistrationFailed, "ユーザー登録に失敗しました", err)
}

func wrapLoginError(err error) error {
	if domain.IsAppError(err) {
		return err
	}
	return domain.NewInternal(domain.ErrorCodeLoginFailed, "ログイン処理に失敗しました", err)
}
