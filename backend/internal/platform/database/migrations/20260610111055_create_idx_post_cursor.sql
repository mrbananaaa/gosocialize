-- +goose Up
CREATE INDEX IF NOT EXISTS idx_posts_cursor
ON posts (created_at DESC, id DESC);

-- +goose Down
DROP INDEX IF EXISTS idx_posts_cursor;
