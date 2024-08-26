-- +goose Up
CREATE TABLE follows (
    id TEXT PRIMARY KEY,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    feed_id TEXT NOT NULL REFERENCES feeds(id) ON DElETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DElETE CASCADE
);

-- +goose Down
DROP TABLE follows;
