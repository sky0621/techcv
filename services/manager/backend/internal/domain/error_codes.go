// Package domain contains shared application domain types.
package domain

const (
	ErrorCodeInvalidEmailFormat       = "INVALID_EMAIL_FORMAT"
	ErrorCodeInvalidPassword          = "INVALID_PASSWORD"
	ErrorCodeInvalidPasswordHash      = "INVALID_PASSWORD_HASH"
	ErrorCodePasswordHashFailed       = "PASSWORD_HASH_FAILED"
	ErrorCodeUUIDGenerationFailed     = "UUID_GENERATION_FAILED"
	ErrorCodeEmailAlreadyRegistered   = "EMAIL_ALREADY_REGISTERED"
	ErrorCodeUserNotFound             = "USER_NOT_FOUND"
	ErrorCodeTokenNotFound            = "TOKEN_NOT_FOUND"
	ErrorCodeInvalidJSON              = "INVALID_JSON"
	ErrorCodeInvalidRequest           = "INVALID_REQUEST"
	ErrorCodeInvalidVerificationToken = "INVALID_VERIFICATION_TOKEN" // #nosec G101 -- error code identifier, not a credential
	ErrorCodeTokenLookupFailed        = "TOKEN_LOOKUP_FAILED"        // #nosec G101 -- error code identifier, not a credential
	ErrorCodeVerificationTokenExpired = "VERIFICATION_TOKEN_EXPIRED"
	ErrorCodeUserLookupFailed         = "USER_LOOKUP_FAILED"
	ErrorCodeUserCreateFailed         = "USER_CREATE_FAILED"
	ErrorCodeTokenDeleteFailed        = "TOKEN_DELETE_FAILED"
	ErrorCodeAuthTokenIssueFailed     = "AUTH_TOKEN_ISSUE_FAILED" // #nosec G101 -- error code identifier, not a credential
	ErrorCodePasswordMismatch         = "PASSWORD_MISMATCH"
	ErrorCodeTokenCleanupFailed       = "TOKEN_CLEANUP_FAILED"
	ErrorCodeTokenSaveFailed          = "TOKEN_SAVE_FAILED" // #nosec G101 -- error code identifier, not a credential
	ErrorCodeVerificationURLError     = "VERIFICATION_URL_ERROR"
	ErrorCodeVerificationURLMissing   = "VERIFICATION_URL_MISSING"
	ErrorCodeEmailSendFailed          = "EMAIL_SEND_FAILED"
)
