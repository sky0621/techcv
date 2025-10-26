-- name: CreateOAuthState :exec
INSERT INTO oauth_states (
  state,
  code_verifier,
  nonce,
  redirect_uri,
  expires_at
) VALUES (
  sqlc.arg(state),
  sqlc.narg(code_verifier),
  sqlc.narg(nonce),
  sqlc.narg(redirect_uri),
  sqlc.arg(expires_at)
);

-- name: GetOAuthState :one
SELECT
  state,
  code_verifier,
  nonce,
  redirect_uri,
  expires_at,
  created_at,
  updated_at
FROM oauth_states
WHERE state = sqlc.arg(state)
LIMIT 1;

-- name: DeleteOAuthState :exec
DELETE FROM oauth_states
WHERE state = sqlc.arg(state);

-- name: DeleteExpiredOAuthStates :exec
DELETE FROM oauth_states
WHERE expires_at <= sqlc.arg(reference_time);
