-- name: CreateFeed :one
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
-- name: GetFeeds :many
SELECT * from feeds_follows WHERE user_id=$1;

-- name : GetUserByApiKey :one
-- SELECT * FROM feeds WHERE api_key * $1;

-- name : DeletFeedsFollow :exec
DELETE  from feed_follows WHERE id = $1 AND user_id = $1;