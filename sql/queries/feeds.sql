-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- -- name: GetUser :one
-- SELECT * FROM users WHERE name = $1;

-- name: DeleteFeeds :exec
DELETE FROM feeds;

-- -- name: GetUsers :many
-- SELECT * FROM users;
