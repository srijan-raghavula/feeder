-- +goose Up
CREATE TABLE (
    id TEXT PRIMARY KEY,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    name TEXT NOT NULL
);

-- +goose Down
DROP TABLE users;
