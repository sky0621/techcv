-- name: CreateUser :exec
INSERT INTO users (
    id,
    firebase_uid,
    email,
    email_verified,
    display_name,
    photo_url,
    phone_number,
    provider_id,
    firebase_created_at,
    firebase_last_sign_in_at,
    bio,
    is_active,
    email_verified_at,
    last_login_at,
    created_at,
    updated_at
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: GetUserByFirebaseUID :one
SELECT
    id,
    firebase_uid,
    email,
    email_verified,
    display_name,
    photo_url,
    phone_number,
    provider_id,
    firebase_created_at,
    firebase_last_sign_in_at,
    bio,
    is_active,
    email_verified_at,
    last_login_at,
    created_at,
    updated_at,
    deleted_at
FROM users
WHERE firebase_uid = ? AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserByEmail :one
SELECT
    id,
    firebase_uid,
    email,
    email_verified,
    display_name,
    photo_url,
    phone_number,
    provider_id,
    firebase_created_at,
    firebase_last_sign_in_at,
    bio,
    is_active,
    email_verified_at,
    last_login_at,
    created_at,
    updated_at,
    deleted_at
FROM users
WHERE email = ? AND deleted_at IS NULL
LIMIT 1;

-- name: GetUserByID :one
SELECT
    id,
    firebase_uid,
    email,
    email_verified,
    display_name,
    photo_url,
    phone_number,
    provider_id,
    firebase_created_at,
    firebase_last_sign_in_at,
    bio,
    is_active,
    email_verified_at,
    last_login_at,
    created_at,
    updated_at,
    deleted_at
FROM users
WHERE id = ? AND deleted_at IS NULL
LIMIT 1;

-- name: UpdateUserFromFirebase :execresult
UPDATE users
SET
    email = ?,
    email_verified = ?,
    display_name = ?,
    photo_url = ?,
    phone_number = ?,
    provider_id = ?,
    firebase_created_at = ?,
    firebase_last_sign_in_at = ?,
    last_login_at = ?,
    updated_at = ?
WHERE firebase_uid = ? AND deleted_at IS NULL;
