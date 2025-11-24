// Package domain contains shared application domain types.
package domain

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
