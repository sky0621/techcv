// Package middleware はHTTPミドルウェアを提供します。
package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
)

const firebaseUIDContextKey = "firebase_uid"

// FirebaseUIDFromContext は認証ミドルウェアで設定したFirebase UIDを取得します。
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
