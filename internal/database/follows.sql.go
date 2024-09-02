// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: follows.sql

package database

import (
	"context"
)

const deleteFollow = `-- name: DeleteFollow :exec
DELETE FROM follows
WHERE feed_id=$1 AND user_id=$2
`

type DeleteFollowParams struct {
	FeedID string
	UserID string
}

func (q *Queries) DeleteFollow(ctx context.Context, arg DeleteFollowParams) error {
	_, err := q.db.ExecContext(ctx, deleteFollow, arg.FeedID, arg.UserID)
	return err
}

const followFeed = `-- name: FollowFeed :one
INSERT INTO follows (id, created_at, updated_at, feed_id, user_id)
VALUES ($1, $2, $3, $4, $5)
    RETURNING id, created_at, updated_at, feed_id, user_id
`

type FollowFeedParams struct {
	ID        string
	CreatedAt string
	UpdatedAt string
	FeedID    string
	UserID    string
}

func (q *Queries) FollowFeed(ctx context.Context, arg FollowFeedParams) (Follow, error) {
	row := q.db.QueryRowContext(ctx, followFeed,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.FeedID,
		arg.UserID,
	)
	var i Follow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.FeedID,
		&i.UserID,
	)
	return i, err
}

const getFollows = `-- name: GetFollows :many
SELECT id, created_at, updated_at, feed_id, user_id FROM follows
WHERE user_id=$1
`

func (q *Queries) GetFollows(ctx context.Context, userID string) ([]Follow, error) {
	rows, err := q.db.QueryContext(ctx, getFollows, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Follow
	for rows.Next() {
		var i Follow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.FeedID,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFollowsByUser = `-- name: GetFollowsByUser :many
SELECT id, created_at, updated_at, feed_id, user_id FROM follows
WHERE user_id = $1
`

func (q *Queries) GetFollowsByUser(ctx context.Context, userID string) ([]Follow, error) {
	rows, err := q.db.QueryContext(ctx, getFollowsByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Follow
	for rows.Next() {
		var i Follow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.FeedID,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
