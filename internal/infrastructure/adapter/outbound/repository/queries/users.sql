-- name: CreateUser :exec
INSERT INTO users (id, email, name, nick, password_hash, is_admin, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetUserByEmail :one
SELECT id, email, name, nick, password_hash, is_admin, created_at, updated_at, deleted_at
FROM users
WHERE email = $1 AND deleted_at IS NULL;

-- name: GetUserByID :one
SELECT id, email, name, nick, is_admin, created_at, updated_at, deleted_at
FROM users
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateUser :one
UPDATE users
SET name = $2, nick = $3, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING id, email, name, nick, is_admin, created_at, updated_at, deleted_at;

-- name: DeleteUser :execrows
UPDATE users
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;
