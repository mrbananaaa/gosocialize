-- name: UserFeedsFirst :many
SELECT
  p.id,
  p.author_id,
  p.title,
  p.content,
  p.created_at,
  p.updated_at
FROM posts p
WHERE p.author_id IN (
  SELECT f.followee_id
  FROM follows f
  WHERE f.follower_id = sqlc.arg(user_id)::uuid

  UNION

  SELECT sqlc.arg(user_id)::uuid
)
ORDER BY p.created_at DESC, p.id DESC
LIMIT sqlc.arg(pagination_limit);

-- name: UserFeedsAfter :many
SELECT
  p.id,
  p.author_id,
  p.title,
  p.content,
  p.created_at,
  p.updated_at
FROM posts p
WHERE p.author_id IN (
  SELECT f.followee_id
  FROM follows f
  WHERE f.follower_id = sqlc.arg(user_id)::uuid

  UNION

  SELECT sqlc.arg(user_id)::uuid
)
AND (
  p.created_at,
  p.id
) < (
  sqlc.arg(cursor_created_at)::timestamptz,
  sqlc.arg(cursor_id)::uuid
)
ORDER BY p.created_at DESC, p.id DESC
LIMIT sqlc.arg(pagination_limit);
