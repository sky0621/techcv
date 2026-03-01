-- name: CreateProfile :exec
INSERT INTO profiles (id, name, nickname)
VALUES (?, ?, ?);
