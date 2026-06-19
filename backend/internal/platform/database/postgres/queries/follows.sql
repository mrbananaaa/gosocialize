-- name: FollowUser :exec
INSERT INTO follows (
  followee_id, follower_id
) VALUES (
  $1, $2
);

-- name: UnfollowUser :exec
DELETE
FROM follows
WHERE
  followee_id = $1 AND follower_id = $2;

-- name: FollowersCount :one
SELECT
  COUNT(*)
FROM follows
WHERE followee_id = $1;
