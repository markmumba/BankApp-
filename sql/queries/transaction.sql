-- name: SaveTransaction :one
INSERT INTO
  transactions (account_id, recepient_id, type)
VALUES
  ($1, $2, $3) RETURNING *;
