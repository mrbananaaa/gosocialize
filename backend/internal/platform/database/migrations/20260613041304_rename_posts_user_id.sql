-- +goose Up
ALTER TABLE posts 
  DROP CONSTRAINT fk_user_post;

ALTER TABLE posts 
  RENAME COLUMN user_id TO author_id;

ALTER TABLE posts
  ADD CONSTRAINT fk_user_post
    FOREIGN KEY (author_id)
    REFERENCES users (id)
    ON DELETE CASCADE;

-- +goose Down
ALTER TABLE posts 
  DROP CONSTRAINT fk_user_post;

ALTER TABLE posts 
  RENAME COLUMN author_id TO user_id;

ALTER TABLE posts
  ADD CONSTRAINT fk_user_post
    FOREIGN KEY (user_id)
    REFERENCES users (id)
    ON DELETE CASCADE;
