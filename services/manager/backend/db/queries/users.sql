-- name: CreateUser :exec
INSERT INTO users (
  id,
  email,
  password_hash,
  google_id,
  name,
  profile_image_url,
  bio,
  is_active,
  email_verified_at,
  last_login_at,
  created_at,
  updated_at
) VALUES (
  UUID_TO_BIN(sqlc.arg(id), TRUE),
  sqlc.arg(email),
  sqlc.narg(password_hash),
  sqlc.narg(google_id),
  sqlc.narg(name),
  sqlc.narg(profile_image_url),
  sqlc.narg(bio),
  sqlc.arg(is_active),
  sqlc.narg(email_verified_at),
  sqlc.narg(last_login_at),
  sqlc.arg(created_at),
  sqlc.arg(updated_at)
);

-- name: GetUserByID :one
SELECT
  BIN_TO_UUID(id, TRUE) AS id,
  email,
  password_hash,
  google_id,
  name,
  profile_image_url,
  bio,
  is_active,
  email_verified_at,
  last_login_at,
  created_at,
  updated_at
FROM users
WHERE id = UUID_TO_BIN(sqlc.arg(id), TRUE);

-- name: GetUserByEmail :one
SELECT
  BIN_TO_UUID(id, TRUE) AS id,
  email,
  password_hash,
  google_id,
  name,
  profile_image_url,
  bio,
  is_active,
  email_verified_at,
  last_login_at,
  created_at,
  updated_at
FROM users
WHERE email = sqlc.arg(email)
LIMIT 1;

-- name: GetUserByGoogleID :one
SELECT
  BIN_TO_UUID(id, TRUE) AS id,
  email,
  password_hash,
  google_id,
  name,
  profile_image_url,
  bio,
  is_active,
  email_verified_at,
  last_login_at,
  created_at,
  updated_at
FROM users
WHERE google_id = sqlc.arg(google_id)
LIMIT 1;

-- name: UpdateUserGoogleID :exec
UPDATE users
SET
  google_id = sqlc.arg(google_id),
  updated_at = sqlc.arg(updated_at)
WHERE id = UUID_TO_BIN(sqlc.arg(id), TRUE)
  AND google_id IS NULL;

-- name: UpdateUserLastLogin :exec
UPDATE users
SET
  last_login_at = sqlc.arg(last_login_at),
  updated_at = sqlc.arg(updated_at)
WHERE id = UUID_TO_BIN(sqlc.arg(id), TRUE);
