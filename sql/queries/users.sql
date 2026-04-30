-- name: CreateUser :one
INSERT INTO
  users (id, created_at, updated_at, email, name, password)
VALUES
  (?, ?, ?, ?, ?, ?)
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

-- name: UpdatePersonalInfo :one
UPDATE users
SET
  name = ?,
  email = ?,
  updated_at = ?
WHERE
  id = ?
RETURNING
  *;

-- name: UpdatePassword :one
UPDATE users
SET
  password = ?,
  updated_at = ?
WHERE
  id = ?
RETURNING
  *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE
  id = ?;

-- name: UpdateProfilePicture :one
UPDATE users
SET
  profile_picture_url = ?,
  updated_at = ?
WHERE
  id = ?
RETURNING
  *;
