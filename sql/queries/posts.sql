-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostsByUser :one
SELECT * FROM posts
WHERE feed_id = $1
ORDER BY published_at
LIMIT $2;

-- name: GetPostByUrl :one
SELECT * FROM posts
WHERE url = $1;

-- name: UpdatePostByUrl :exec
UPDATE posts
SET updated_at = $2
WHERE url = $1;
