-- name: CreatePost :exec
INSERT INTO posts (
  id, user_id, title, content, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6
);

-- name: ListPostsFirst :many
SELECT
  id, user_id, title, content, created_at, updated_at
FROM posts
ORDER BY created_at DESC, id DESC
LIMIT sqlc.arg(pagination_limit);

-- name: ListPostsAfter :many
SELECT
  id, user_id, title, content, created_at, updated_at
FROM posts
WHERE
  (created_at, id) < (
    sqlc.arg(cursor_created_at)::timestamptz,
    sqlc.arg(cursor_id)::uuid
  )
ORDER BY created_at DESC, id DESC
LIMIT sqlc.arg(pagination_limit);

-- name: FindPostByID :one
SELECT
  id, user_id, title, content, created_at, updated_at
FROM posts
WHERE
  id = sqlc.arg(post_id);

-- name: UpdatePost :exec
UPDATE posts
SET
  title = $2,
  content = $3,
  updated_at = NOW()
WHERE id = $1;

-- name: DeletePost :exec
DELETE 
FROM posts
WHERE id = $1;