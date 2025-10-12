package memory

import (
	"context"
	"sync"

	"github.com/sky0621/techcv/manager/backend/internal/domain"
	"github.com/sky0621/techcv/manager/backend/internal/domain/user"
)

// VerificationTokenRepository provides in-memory storage for verification tokens.
type VerificationTokenRepository struct {
	mu            sync.RWMutex
	tokensByValue map[string]user.VerificationToken
	tokensByEmail map[string]map[string]struct{}
}

// NewVerificationTokenRepository constructs a new repository instance.
func NewVerificationTokenRepository() *VerificationTokenRepository {
	return &VerificationTokenRepository{
		tokensByValue: make(map[string]user.VerificationToken),
		tokensByEmail: make(map[string]map[string]struct{}),
	}
}

// Save persists or replaces a verification token.
func (r *VerificationTokenRepository) Save(_ context.Context, token user.VerificationToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tokensByValue[token.Token()] = token

	emailKey := token.Email().String()
	bucket, ok := r.tokensByEmail[emailKey]
	if !ok {
		bucket = make(map[string]struct{})
		r.tokensByEmail[emailKey] = bucket
	}
	bucket[token.Token()] = struct{}{}
	return nil
}

// FindByToken retrieves a token by its value.
func (r *VerificationTokenRepository) FindByToken(_ context.Context, token string) (user.VerificationToken, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if t, ok := r.tokensByValue[token]; ok {
		return t, nil
	}

	detail := domain.ErrorDetail{Field: "token", Code: "TOKEN_NOT_FOUND", Message: "確認トークンが見つかりません"}
	return user.VerificationToken{}, domain.NewNotFound("TOKEN_NOT_FOUND", "確認トークンが見つかりません").WithDetails(detail)
}

// DeleteByToken removes a token using its value.
func (r *VerificationTokenRepository) DeleteByToken(_ context.Context, token string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	record, ok := r.tokensByValue[token]
	if !ok {
		return nil
	}

	delete(r.tokensByValue, token)
	emailKey := record.Email().String()
	if bucket, exists := r.tokensByEmail[emailKey]; exists {
		delete(bucket, token)
		if len(bucket) == 0 {
			delete(r.tokensByEmail, emailKey)
		}
	}
	return nil
}

// DeleteByEmail removes all tokens associated with the given email.
func (r *VerificationTokenRepository) DeleteByEmail(_ context.Context, email user.Email) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	emailKey := email.String()
	bucket, exists := r.tokensByEmail[emailKey]
	if !exists {
		return nil
	}

	for token := range bucket {
		delete(r.tokensByValue, token)
	}
	delete(r.tokensByEmail, emailKey)
	return nil
}
