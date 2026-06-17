-- +goose Up
CREATE INDEX IF NOT EXISTS idx_posts_author_id
  ON posts (author_id);

-- +goose Down
DROP INDEX IF EXISTS idx_posts_author_id;
