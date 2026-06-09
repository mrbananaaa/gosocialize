-- name: CreatePost :exec
INSERT INTO posts (
  id, user_id, title, content, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6
);