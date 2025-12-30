package mysql

import (
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/sky0621/techcv/manager/backend/internal/domain"
)

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
