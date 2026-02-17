-- +goose Up
CREATE TABLE IF NOT EXISTS chirps (
    id UUID PRIMARY KEY,
    body VARCHAR(140) NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS chirps;