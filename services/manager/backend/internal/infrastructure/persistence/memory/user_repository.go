package memory

import (
	"context"
	"sync"

	"github.com/sky0621/techcv/manager/backend/internal/domain"
	"github.com/sky0621/techcv/manager/backend/internal/domain/user"
)

// UserRepository provides an in-memory implementation of user persistence.
type UserRepository struct {
	mu    sync.RWMutex
	users map[string]user.User
}

// NewUserRepository constructs a new in-memory user repository.
func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]user.User),
	}
}

// ExistsByEmail reports whether a user with the provided email already exists.
func (r *UserRepository) ExistsByEmail(_ context.Context, email user.Email) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.users[email.String()]
	return ok, nil
}

// Create persists a new user aggregate.
func (r *UserRepository) Create(_ context.Context, u user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	email := u.Email().String()
	if _, exists := r.users[email]; exists {
		detail := domain.ErrorDetail{Field: "email", Code: domain.ErrorCodeEmailAlreadyRegistered, Message: "このメールアドレスは既に登録されています"}
		return domain.NewValidation(domain.ErrorCodeEmailAlreadyRegistered, "このメールアドレスは既に登録されています").WithDetails(detail)
	}

	r.users[email] = u
	return nil
}

// GetByEmail loads the user aggregate associated with the given email.
func (r *UserRepository) GetByEmail(_ context.Context, email user.Email) (user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if u, ok := r.users[email.String()]; ok {
		return u, nil
	}

	detail := domain.ErrorDetail{Field: "email", Code: domain.ErrorCodeUserNotFound, Message: "ユーザーが見つかりません"}
	return user.User{}, domain.NewNotFound(domain.ErrorCodeUserNotFound, "ユーザーが見つかりません").WithDetails(detail)
}
