// Package domain はアプリケーション全体で共有するドメイン型とエラー定義を提供します。
package domain

// ErrorCodeInvalidEmailFormat などのエラーコードを定義します。
const (
	ErrorCodeInvalidEmailFormat   = "INVALID_EMAIL_FORMAT"
	ErrorCodeInvalidPassword      = "INVALID_PASSWORD"
	ErrorCodeInvalidPasswordHash  = "INVALID_PASSWORD_HASH"
	ErrorCodePasswordHashFailed   = "PASSWORD_HASH_FAILED"
	ErrorCodeUUIDGenerationFailed = "UUID_GENERATION_FAILED"
	ErrorCodeInvalidToken         = "INVALID_TOKEN"
	ErrorCodeUserAlreadyExists    = "USER_ALREADY_EXISTS"
	ErrorCodeUserNotFound         = "USER_NOT_FOUND"
	ErrorCodeRegistrationFailed   = "REGISTRATION_FAILED"
	ErrorCodeLoginFailed          = "LOGIN_FAILED"
	ErrorCodeInternalError        = "INTERNAL_ERROR"
	ErrorCodeValidationError      = "VALIDATION_ERROR"
)
