-- name: GetUsers :many
SELECT * FROM users;

-- name: CreateUser :exec
INSERT INTO users (
  id, email, username, password, name, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
);