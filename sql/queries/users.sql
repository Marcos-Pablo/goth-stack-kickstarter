-- name: CreateUser :one
INSERT INTO
  users (id, created_at, updated_at, email, password)
VALUES
  (?, ?, ?, ?, ?)
RETURNING
  *;

-- name: GetUserByEmail :one
SELECT
  *
FROM
  users
WHERE
  email = ?
LIMIT
  1;

-- name: GetUserById :one
SELECT
  *
FROM
  users
WHERE
  id = ?
LIMIT
  1;
