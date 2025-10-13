-- name: CreatePublicURL :execresult
INSERT INTO public_urls (url_key)
VALUES (?);

-- name: GetActivePublicURL :one
SELECT
  id,
  url_key,
  is_active,
  created_at,
  updated_at
FROM public_urls
WHERE is_active = TRUE
ORDER BY updated_at DESC
LIMIT 1;

-- name: ListPublicURLs :many
SELECT
  id,
  url_key,
  is_active,
  created_at,
  updated_at
FROM public_urls
ORDER BY updated_at DESC;

-- name: DeactivatePublicURL :exec
UPDATE public_urls
SET is_active = FALSE,
    updated_at = CURRENT_TIMESTAMP(6)
WHERE id = ?;
