package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/sky0621/techcv/manager/backend/internal/domain"
	"github.com/sky0621/techcv/manager/backend/internal/domain/user"
	mysqlsqlc "github.com/sky0621/techcv/manager/backend/internal/infrastructure/mysql/sqlc"
)

// UserRepository persists user aggregates in MySQL.
type UserRepository struct {
	queries *mysqlsqlc.Queries
}

// NewUserRepository constructs a repository backed by sqlc queries.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		queries: mysqlsqlc.New(db),
	}
}

// Create inserts a new user record.
func (r *UserRepository) Create(ctx context.Context, entity user.User) error {
	idBytes, err := uuidToBytes(entity.ID())
	if err != nil {
		return err
	}

	params := mysqlsqlc.CreateUserParams{
		ID:                   idBytes,
		FirebaseUid:          entity.FirebaseUID().String(),
		Email:                entity.Email().String(),
		EmailVerified:        entity.EmailVerified(),
		DisplayName:          toNullString(entity.DisplayName()),
		PhotoUrl:             toNullString(entity.PhotoURL()),
		PhoneNumber:          toNullString(entity.PhoneNumber()),
		ProviderID:           entity.ProviderID(),
		FirebaseCreatedAt:    toNullTime(entity.FirebaseCreatedAt()),
		FirebaseLastSignInAt: toNullTime(entity.FirebaseLastSignInAt()),
		Bio:                  toNullString(entity.Bio()),
		IsActive:             entity.IsActive(),
		EmailVerifiedAt:      toNullTime(entity.EmailVerifiedAt()),
		LastLoginAt:          toNullTime(entity.LastLoginAt()),
		CreatedAt:            entity.CreatedAt(),
		UpdatedAt:            entity.UpdatedAt(),
	}

	if err := r.queries.CreateUser(ctx, params); err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

// FindByID retrieves a user by internal identifier.
func (r *UserRepository) FindByID(ctx context.Context, id string) (user.User, error) {
	idBytes, err := uuidToBytes(id)
	if err != nil {
		return user.User{}, err
	}

	record, err := r.queries.GetUserByID(ctx, idBytes)
	if errors.Is(err, sql.ErrNoRows) {
		return user.User{}, domain.NewNotFound(domain.ErrorCodeUserNotFound, "ユーザーが見つかりません")
	}
	if err != nil {
		return user.User{}, fmt.Errorf("get user by id: %w", err)
	}

	return toDomainUser(record)
}

// FindByEmail retrieves a user by email address.
func (r *UserRepository) FindByEmail(ctx context.Context, email user.Email) (user.User, error) {
	record, err := r.queries.GetUserByEmail(ctx, email.String())
	if errors.Is(err, sql.ErrNoRows) {
		return user.User{}, domain.NewNotFound(domain.ErrorCodeUserNotFound, "ユーザーが見つかりません")
	}
	if err != nil {
		return user.User{}, fmt.Errorf("get user by email: %w", err)
	}

	return toDomainUser(record)
}

// FindByFirebaseUID retrieves a user by Firebase UID.
func (r *UserRepository) FindByFirebaseUID(ctx context.Context, firebaseUID user.FirebaseUID) (user.User, error) {
	record, err := r.queries.GetUserByFirebaseUID(ctx, firebaseUID.String())
	if errors.Is(err, sql.ErrNoRows) {
		return user.User{}, domain.NewNotFound(domain.ErrorCodeUserNotFound, "ユーザーが見つかりません")
	}
	if err != nil {
		return user.User{}, fmt.Errorf("get user by firebase uid: %w", err)
	}

	return toDomainUser(record)
}

// UpdateFromFirebase refreshes persisted data based on Firebase profile details.
func (r *UserRepository) UpdateFromFirebase(ctx context.Context, entity user.User) error {
	params := mysqlsqlc.UpdateUserFromFirebaseParams{
		Email:                entity.Email().String(),
		EmailVerified:        entity.EmailVerified(),
		DisplayName:          toNullString(entity.DisplayName()),
		PhotoUrl:             toNullString(entity.PhotoURL()),
		PhoneNumber:          toNullString(entity.PhoneNumber()),
		ProviderID:           entity.ProviderID(),
		FirebaseCreatedAt:    toNullTime(entity.FirebaseCreatedAt()),
		FirebaseLastSignInAt: toNullTime(entity.FirebaseLastSignInAt()),
		LastLoginAt:          toNullTime(entity.LastLoginAt()),
		UpdatedAt:            entity.UpdatedAt(),
		FirebaseUid:          entity.FirebaseUID().String(),
	}

	result, err := r.queries.UpdateUserFromFirebase(ctx, params)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("retrieve update result: %w", err)
	}
	if affected == 0 {
		return domain.NewNotFound(domain.ErrorCodeUserNotFound, "ユーザーが見つかりません")
	}

	return nil
}

func toDomainUser(record mysqlsqlc.User) (user.User, error) {
	id, err := uuid.FromBytes(record.ID)
	if err != nil {
		return user.User{}, domain.NewInternal(domain.ErrorCodeInternalError, "ユーザーデータの変換に失敗しました", err)
	}

	snapshot := user.Snapshot{
		ID:                   id.String(),
		FirebaseUID:          record.FirebaseUid,
		Email:                record.Email,
		EmailVerified:        record.EmailVerified,
		DisplayName:          fromNullString(record.DisplayName),
		PhotoURL:             fromNullString(record.PhotoUrl),
		PhoneNumber:          fromNullString(record.PhoneNumber),
		ProviderID:           record.ProviderID,
		FirebaseCreatedAt:    fromNullTime(record.FirebaseCreatedAt),
		FirebaseLastSignInAt: fromNullTime(record.FirebaseLastSignInAt),
		Bio:                  fromNullString(record.Bio),
		IsActive:             record.IsActive,
		EmailVerifiedAt:      fromNullTime(record.EmailVerifiedAt),
		LastLoginAt:          fromNullTime(record.LastLoginAt),
		CreatedAt:            record.CreatedAt,
		UpdatedAt:            record.UpdatedAt,
	}

	return user.Restore(snapshot)
}

func uuidToBytes(id string) ([]byte, error) {
	parsed, err := uuid.Parse(strings.TrimSpace(id))
	if err != nil {
		return nil, domain.NewInternal(domain.ErrorCodeInternalError, "ユーザーIDが不正です", err)
	}
	return parsed.MarshalBinary()
}

func toNullString(value *string) sql.NullString {
	if value == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *value, Valid: true}
}

func fromNullString(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}

func toNullTime(value *time.Time) sql.NullTime {
	if value == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: value.UTC().Truncate(time.Microsecond), Valid: true}
}

func fromNullTime(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}
	ts := value.Time.UTC().Truncate(time.Microsecond)
	return &ts
}
