// Package firebase provides Firebase Admin SDK-backed authentication utilities.
package firebase

import (
	"context"
	"fmt"
	"strings"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"

	"github.com/sky0621/techcv/manager/backend/internal/domain"
	appconfig "github.com/sky0621/techcv/manager/backend/internal/infrastructure/config"
	usecaseauth "github.com/sky0621/techcv/manager/backend/internal/usecase/auth"
)

// AuthService bridges Firebase Admin SDK operations for the application.
type AuthService struct {
	client *auth.Client
}

// NewAuthService initializes a Firebase Admin SDK client.
func NewAuthService(ctx context.Context, cfg appconfig.FirebaseConfig) (*AuthService, error) {
	opts := []option.ClientOption{}
	if strings.TrimSpace(cfg.CredentialsPath) != "" {
		opts = append(opts, option.WithCredentialsFile(cfg.CredentialsPath))
	}

	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: cfg.ProjectID}, opts...)
	if err != nil {
		return nil, fmt.Errorf("initialize firebase app: %w", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("initialize firebase auth client: %w", err)
	}

	return &AuthService{client: client}, nil
}

// VerifyIDToken validates the given Firebase ID token and returns its UID.
func (s *AuthService) VerifyIDToken(ctx context.Context, idToken string) (string, error) {
	token := strings.TrimSpace(idToken)
	if token == "" {
		return "", domain.NewUnauthorized(domain.ErrorCodeInvalidToken, "Firebase IDトークンが指定されていません", nil)
	}

	parsed, err := s.client.VerifyIDToken(ctx, token)
	if err != nil {
		return "", domain.NewUnauthorized(domain.ErrorCodeInvalidToken, "Firebase IDトークンの検証に失敗しました", err)
	}

	return parsed.UID, nil
}

// GetUser fetches Firebase user information for the given UID.
func (s *AuthService) GetUser(ctx context.Context, uid string) (usecaseauth.FirebaseUser, error) {
	record, err := s.client.GetUser(ctx, strings.TrimSpace(uid))
	if err != nil {
		return usecaseauth.FirebaseUser{}, domain.NewUnauthorized(domain.ErrorCodeInvalidToken, "Firebaseユーザーの取得に失敗しました", err)
	}

	providerID := record.ProviderID
	if len(record.ProviderUserInfo) > 0 && record.ProviderUserInfo[0] != nil && record.ProviderUserInfo[0].ProviderID != "" {
		providerID = record.ProviderUserInfo[0].ProviderID
	}
	if providerID == "" {
		providerID = "firebase"
	}

	var firebaseCreatedAt *time.Time
	if record.UserMetadata != nil && record.UserMetadata.CreationTimestamp != 0 {
		ts := time.UnixMilli(record.UserMetadata.CreationTimestamp).UTC()
		firebaseCreatedAt = &ts
	}

	var firebaseLastSignInAt *time.Time
	if record.UserMetadata != nil && record.UserMetadata.LastLogInTimestamp != 0 {
		ts := time.UnixMilli(record.UserMetadata.LastLogInTimestamp).UTC()
		firebaseLastSignInAt = &ts
	}

	return usecaseauth.FirebaseUser{
		UID:                  record.UID,
		Email:                record.Email,
		EmailVerified:        record.EmailVerified,
		DisplayName:          nilIfEmpty(record.DisplayName),
		PhotoURL:             nilIfEmpty(record.PhotoURL),
		PhoneNumber:          nilIfEmpty(record.PhoneNumber),
		ProviderID:           providerID,
		FirebaseCreatedAt:    firebaseCreatedAt,
		FirebaseLastSignInAt: firebaseLastSignInAt,
	}, nil
}

func nilIfEmpty(value string) *string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	v := strings.TrimSpace(value)
	return &v
}
