package auth

import (
	"context"

	"github.com/sky0621/techcv/manager/backend/internal/domain/user"
	"github.com/sky0621/techcv/manager/backend/internal/domain/uuidv7"
)

// UUIDTokenIssuer issues authentication tokens using UUID v7 values.
type UUIDTokenIssuer struct{}

// NewUUIDTokenIssuer constructs a new token issuer.
func NewUUIDTokenIssuer() UUIDTokenIssuer {
	return UUIDTokenIssuer{}
}

// Issue generates a new authentication token for the given user.
func (UUIDTokenIssuer) Issue(_ context.Context, _ user.User) (string, error) {
	return uuidv7.NewString()
}
