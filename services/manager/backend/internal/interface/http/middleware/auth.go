package middleware

import (
	"context"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/sky0621/techcv/manager/backend/internal/domain"
)

const firebaseUIDContextKey = "firebase_uid"

// FirebaseTokenVerifier matches the use case interface while avoiding echo dependency in the inner layer.
type FirebaseTokenVerifier interface {
	VerifyIDToken(ctx context.Context, idToken string) (string, error)
}

// FirebaseAuth middleware verifies Firebase ID tokens and stores the UID in the context.
func FirebaseAuth(verifier FirebaseTokenVerifier) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := bearerToken(c.Request().Header.Get(echo.HeaderAuthorization))
			if token == "" {
				return domain.NewUnauthorized(domain.ErrorCodeInvalidToken, "Authorizationヘッダーが不足しています", nil)
			}

			uid, err := verifier.VerifyIDToken(c.Request().Context(), token)
			if err != nil {
				return err
			}

			c.Set(firebaseUIDContextKey, uid)
			return next(c)
		}
	}
}

// FirebaseUIDFromContext fetches the Firebase UID assigned by the authentication middleware.
func FirebaseUIDFromContext(c echo.Context) (string, bool) {
	if c == nil {
		return "", false
	}
	val := c.Get(firebaseUIDContextKey)
	uid, ok := val.(string)
	if !ok || strings.TrimSpace(uid) == "" {
		return "", false
	}
	return uid, true
}

func bearerToken(header string) string {
	parts := strings.Fields(header)
	if len(parts) != 2 {
		return ""
	}
	if !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
