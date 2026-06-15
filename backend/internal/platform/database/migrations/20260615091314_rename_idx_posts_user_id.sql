-- +goose Up
ALTER INDEX IF EXISTS
  idx_posts_user_id RENAME TO idx_posts_author_id;

-- +goose Down
ALTER INDEX IF EXISTS
  idx_posts_author_id RENAME TO idx_posts_user_id;
