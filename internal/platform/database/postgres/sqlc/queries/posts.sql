-- name: CreatePost :exec
INSERT INTO posts (
  id, user_id, title, content, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6
);

-- name: GetPost :many
SELECT
  id, user_id, title, content, created_at, updated_at
FROM
  posts;

-- name: ListPosts :many
SELECT
  id, user_id, title, content, created_at, updated_at
FROM
  posts
WHERE
  (
    sqlc.arg(cursor_created_at)::timestamptz IS NULL
    OR (created_at, id) < (
      sqlc.arg(cursor_created_at)::timestamptz,
      sqlc.arg(cursor_id)::uuid
    )
  )
ORDER BY
  created_at DESC, id DESC
LIMIT
  sqlc.arg(pagination_limit);

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