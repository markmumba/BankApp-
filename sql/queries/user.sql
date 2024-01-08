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


-- name: ListAllUser :many
SELECT
  *
FROM
  users;
