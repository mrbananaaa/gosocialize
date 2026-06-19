-- name: FollowUser :exec
INSERT INTO follows (follower_id, followee_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: UnfollowUser :exec
DELETE
FROM follows
WHERE follower_id = $1 
  AND followee_id = $2;

-- name: IsFollowing :one
SELECT EXISTS (
  SELECT 1
  FROM follows
  WHERE follower_id = $1
    AND followee_id = $2
);

-- name: CountFollowers :one
SELECT COUNT(*)
FROM follows
WHERE followee_id = $1;

-- name: CountFollowing :one
SELECT COUNT(*)
FROM follows
WHERE follower_id = $1;