package user

import (
	"strings"

	"github.com/sky0621/techcv/manager/backend/internal/domain"
)

// FirebaseUID represents a validated Firebase user identifier.
type FirebaseUID struct {
	value string
}

// NewFirebaseUID validates and constructs a FirebaseUID value object.
func NewFirebaseUID(raw string) (FirebaseUID, error) {
	trimmed := strings.TrimSpace(raw)
	detail := domain.ErrorDetail{Field: "firebaseUid", Code: domain.ErrorCodeValidationError, Message: "Firebase UIDが指定されていません"}
	if trimmed == "" {
		return FirebaseUID{}, domain.NewValidation(domain.ErrorCodeValidationError, "Firebase UIDが無効です").WithDetails(detail)
	}

	if len(trimmed) > 128 {
		detail.Message = "Firebase UIDは128文字以内で入力してください"
		return FirebaseUID{}, domain.NewValidation(domain.ErrorCodeValidationError, "Firebase UIDが無効です").WithDetails(detail)
	}

	return FirebaseUID{value: trimmed}, nil
}

// String returns the string representation of the Firebase UID.
func (u FirebaseUID) String() string {
	return u.value
}
