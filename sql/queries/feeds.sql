-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING *;

-- name: GetAllFeeds :many
SELECT * FROM feeds;

-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds
ORDER BY last_fetched_at
ASC NULLS FIRST
LIMIT $1;

-- name: MarkFeedFetched :one
UPDATE feeds
set last_fetched_at = $2 AND updated_at = $2
WHERE url = $1
    RETURNING *;

-- name: GetFeedByUrl :one
SELECT * FROM feeds
WHERE url = $1;
