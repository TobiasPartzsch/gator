-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, last_fetched_at, name, url, user_id)
VALUES (@id, @created_at, @updated_at, @last_fetched_at, @name, @url, @user_id)
RETURNING *;

-- name: DeleteFeeds :exec
DELETE FROM feeds;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedsWithUsers :many
SELECT feeds.name, feeds.url, users.name as user_name 
FROM feeds 
JOIN users ON feeds.user_id = users.id;

-- name: GetFeedByURL :one
SELECT * FROM feeds
WHERE url = $1;

-- name: MarkFeedFetched :one
UPDATE feeds
SET last_fetched_at = NOW(),
updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST, id ASC
LIMIT 1;
