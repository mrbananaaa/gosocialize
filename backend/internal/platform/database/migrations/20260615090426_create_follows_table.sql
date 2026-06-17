-- +goose Up
CREATE TABLE IF NOT EXISTS follows (
  follower_id UUID NOT NULL,
  followee_id UUID NOT NULL,

  PRIMARY KEY (follower_id, followee_id),

  FOREIGN KEY (follower_id)
    REFERENCES users (id),

  FOREIGN KEY (followee_id)
    REFERENCES users (id)
);

CREATE INDEX IF NOT EXISTS idx_follows_followee_id
ON follows(followee_id);

-- +goose Down
DROP INDEX IF EXISTS idx_follows_followee_id;

DROP TABLE IF EXISTS follows;
