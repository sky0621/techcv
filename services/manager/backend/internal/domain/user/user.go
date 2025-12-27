package user

import (
	"strings"
	"time"

	"github.com/sky0621/techcv/manager/backend/internal/domain"
)

// User represents an authenticated user backed by Firebase.
type User struct {
	id                   string
	firebaseUID          FirebaseUID
	email                Email
	emailVerified        bool
	displayName          *string
	photoURL             *string
	phoneNumber          *string
	providerID           string
	firebaseCreatedAt    *time.Time
	firebaseLastSignInAt *time.Time
	bio                  *string
	isActive             bool
	emailVerifiedAt      *time.Time
	lastLoginAt          *time.Time
	createdAt            time.Time
	updatedAt            time.Time
}

// FirebaseUserData holds Firebase-backed profile information required for persistence.
type FirebaseUserData struct {
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

// NewUserFromFirebase constructs a new user aggregate from Firebase profile data.
func NewUserFromFirebase(data FirebaseUserData, now time.Time, idGenerator func() (string, error)) (User, error) {
	firebaseUID, err := NewFirebaseUID(data.UID)
	if err != nil {
		return User{}, err
	}

	email, err := NewEmail(data.Email)
	if err != nil {
		return User{}, err
	}

	providerID, err := validateProviderID(data.ProviderID)
	if err != nil {
		return User{}, err
	}

	if idGenerator == nil {
		return User{}, domain.NewInternal(domain.ErrorCodeUUIDGenerationFailed, "ユーザーIDの生成に失敗しました", nil)
	}

	id, err := idGenerator()
	if err != nil {
		return User{}, domain.NewInternal(domain.ErrorCodeUUIDGenerationFailed, "ユーザーIDの生成に失敗しました", err)
	}

	ts := now.UTC().Truncate(time.Microsecond)

	return User{
		id:                   id,
		firebaseUID:          firebaseUID,
		email:                email,
		emailVerified:        data.EmailVerified,
		displayName:          normalizeOptionalString(data.DisplayName),
		photoURL:             normalizeOptionalString(data.PhotoURL),
		phoneNumber:          normalizeOptionalString(data.PhoneNumber),
		providerID:           providerID,
		firebaseCreatedAt:    normalizeOptionalTime(data.FirebaseCreatedAt),
		firebaseLastSignInAt: normalizeOptionalTime(data.FirebaseLastSignInAt),
		isActive:             true,
		lastLoginAt:          &ts,
		createdAt:            ts,
		updatedAt:            ts,
	}, nil
}

// UpdateFromFirebase refreshes Firebase-backed attributes and updates login timestamps.
func (u User) UpdateFromFirebase(data FirebaseUserData, now time.Time) (User, error) {
	if strings.TrimSpace(data.UID) == "" || data.UID != u.firebaseUID.String() {
		return User{}, domain.NewUnauthorized(domain.ErrorCodeInvalidToken, "Firebase UIDが一致しません", nil)
	}

	email, err := NewEmail(data.Email)
	if err != nil {
		return User{}, err
	}

	providerID, err := validateProviderID(data.ProviderID)
	if err != nil {
		return User{}, err
	}

	ts := now.UTC().Truncate(time.Microsecond)

	u.email = email
	u.emailVerified = data.EmailVerified
	u.displayName = normalizeOptionalString(data.DisplayName)
	u.photoURL = normalizeOptionalString(data.PhotoURL)
	u.phoneNumber = normalizeOptionalString(data.PhoneNumber)
	u.providerID = providerID
	u.firebaseCreatedAt = normalizeOptionalTime(data.FirebaseCreatedAt)
	u.firebaseLastSignInAt = normalizeOptionalTime(data.FirebaseLastSignInAt)
	u.lastLoginAt = &ts
	u.updatedAt = ts

	return u, nil
}

// Snapshot represents a persisted view of a user used for restoration from storage.
type Snapshot struct {
	ID                   string
	FirebaseUID          string
	Email                string
	EmailVerified        bool
	DisplayName          *string
	PhotoURL             *string
	PhoneNumber          *string
	ProviderID           string
	FirebaseCreatedAt    *time.Time
	FirebaseLastSignInAt *time.Time
	Bio                  *string
	IsActive             bool
	EmailVerifiedAt      *time.Time
	LastLoginAt          *time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// Restore reconstructs a user aggregate from its persisted snapshot.
func Restore(snapshot Snapshot) (User, error) {
	firebaseUID, err := NewFirebaseUID(snapshot.FirebaseUID)
	if err != nil {
		return User{}, err
	}

	email, err := NewEmail(snapshot.Email)
	if err != nil {
		return User{}, err
	}

	providerID, err := validateProviderID(snapshot.ProviderID)
	if err != nil {
		return User{}, err
	}

	id := strings.TrimSpace(snapshot.ID)
	if id == "" {
		return User{}, domain.NewInternal(domain.ErrorCodeInternalError, "ユーザーデータが不正です", nil)
	}

	return User{
		id:                   id,
		firebaseUID:          firebaseUID,
		email:                email,
		emailVerified:        snapshot.EmailVerified,
		displayName:          normalizeOptionalString(snapshot.DisplayName),
		photoURL:             normalizeOptionalString(snapshot.PhotoURL),
		phoneNumber:          normalizeOptionalString(snapshot.PhoneNumber),
		providerID:           providerID,
		firebaseCreatedAt:    normalizeOptionalTime(snapshot.FirebaseCreatedAt),
		firebaseLastSignInAt: normalizeOptionalTime(snapshot.FirebaseLastSignInAt),
		bio:                  normalizeOptionalString(snapshot.Bio),
		isActive:             snapshot.IsActive,
		emailVerifiedAt:      normalizeOptionalTime(snapshot.EmailVerifiedAt),
		lastLoginAt:          normalizeOptionalTime(snapshot.LastLoginAt),
		createdAt:            snapshot.CreatedAt.UTC().Truncate(time.Microsecond),
		updatedAt:            snapshot.UpdatedAt.UTC().Truncate(time.Microsecond),
	}, nil
}

// ID returns the user identifier.
func (u User) ID() string {
	return u.id
}

// FirebaseUID returns the associated Firebase UID.
func (u User) FirebaseUID() FirebaseUID {
	return u.firebaseUID
}

// Email returns the validated email value object.
func (u User) Email() Email {
	return u.email
}

// EmailVerified indicates whether Firebase marks the email as verified.
func (u User) EmailVerified() bool {
	return u.emailVerified
}

// DisplayName returns the optional profile name.
func (u User) DisplayName() *string {
	return u.displayName
}

// PhotoURL returns the optional profile image URL.
func (u User) PhotoURL() *string {
	return u.photoURL
}

// PhoneNumber returns the optional phone number.
func (u User) PhoneNumber() *string {
	return u.phoneNumber
}

// ProviderID returns the Firebase provider identifier.
func (u User) ProviderID() string {
	return u.providerID
}

// FirebaseCreatedAt returns the Firebase account creation timestamp when available.
func (u User) FirebaseCreatedAt() *time.Time {
	return u.firebaseCreatedAt
}

// FirebaseLastSignInAt returns the Firebase last sign-in timestamp when available.
func (u User) FirebaseLastSignInAt() *time.Time {
	return u.firebaseLastSignInAt
}

// Bio returns the optional bio text.
func (u User) Bio() *string {
	return u.bio
}

// IsActive indicates whether the user is marked active.
func (u User) IsActive() bool {
	return u.isActive
}

// EmailVerifiedAt returns the timestamp the application recorded for email verification, if any.
func (u User) EmailVerifiedAt() *time.Time {
	return u.emailVerifiedAt
}

// LastLoginAt returns the last login timestamp.
func (u User) LastLoginAt() *time.Time {
	return u.lastLoginAt
}

// CreatedAt returns when the user was created in the application.
func (u User) CreatedAt() time.Time {
	return u.createdAt
}

// UpdatedAt returns the last update timestamp.
func (u User) UpdatedAt() time.Time {
	return u.updatedAt
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func normalizeOptionalTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	ts := value.UTC().Truncate(time.Microsecond)
	return &ts
}

const maxProviderIDLength = 50

func validateProviderID(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		detail := domain.ErrorDetail{Field: "providerId", Code: domain.ErrorCodeValidationError, Message: "provider_idは必須です"}
		return "", domain.NewValidation(domain.ErrorCodeValidationError, "プロバイダー情報が不足しています").WithDetails(detail)
	}
	if len(trimmed) > maxProviderIDLength {
		detail := domain.ErrorDetail{Field: "providerId", Code: domain.ErrorCodeValidationError, Message: "provider_idは50文字以内で入力してください"}
		return "", domain.NewValidation(domain.ErrorCodeValidationError, "プロバイダー情報が不正です").WithDetails(detail)
	}
	return trimmed, nil
}
