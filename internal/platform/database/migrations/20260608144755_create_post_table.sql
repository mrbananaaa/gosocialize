-- +goose Up
CREATE TABLE IF NOT EXISTS posts (
  id UUID PRIMARY KEY NOT NULL,
  user_id UUID NOT NULL,
  title VARCHAR(255) NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE,
  updated_at TIMESTAMP WITH TIME ZONE,

  CONSTRAINT fk_user_post
    FOREIGN KEY (user_id)
    REFERENCES users (id)
    ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS 
  idx_posts_user_id 
  ON posts (user_id);

-- +goose Down
DROP CONSTRAINT IF EXISTS fk_user_post;
DROP INDEX IF EXISTS idx_posts_user_id;

DROP TABLE posts;
