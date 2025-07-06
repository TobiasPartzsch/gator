-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (@id, @created_at, @updated_at, @name
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE name = $1;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUsers :many
SELECT * FROM users;
