-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUsers :many
SELECT * FROM users;

-- name: GetCurrentUser :one
SELECT * FROM users WHERE name = $1;

-- name: GetUserNames :many
SELECT name FROM users;

-- name: ResetUsers :exec
DELETE FROM users;

