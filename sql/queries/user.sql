-- name: CreateUser :one
INSERT INTO
  users (username, password_hash, email, full_name)
VALUES
  ($1, $2, $3, $4) RETURNING *;


-- name: FindUser :one 
SELECT
  *
FROM
  users
WHERE
  user_id = $1;


-- name: UpdateUser : one 
UPDATE
  users
SET
  username = $1,
  password_hash = $2,
  email = $3,
  full_name = $4
WHERE
  user_id = $5;


-- name: ListAllUser :many
SELECT
  *
FROM
  users;