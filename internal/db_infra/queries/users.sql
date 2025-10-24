-- name: GetUserDetailsById :one
SELECT users.id, users.email, users.hashed_password
  FROM users
 WHERE users.id = $1;

-- name: CreateUser :one
INSERT INTO users(email, hashed_password)
VALUES($1, $2)
RETURNING id;

-- name: DeleteUserById :exec
DELETE FROM users
WHERE id = $1;
