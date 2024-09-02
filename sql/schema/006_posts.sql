-- +goose Up
CREATE TABLE posts (
    id TEXT UNIQUE PRIMARY KEY,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    title TEXT NOT NULL,
    url TEXT UNIQUE,
    description TEXT,
    published_at TEXT,
    feed_id TEXT NOT NULL
);

-- +goose Down
DROP TABLE posts;
