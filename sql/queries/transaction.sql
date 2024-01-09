-- name: SaveTransaction :exec
INSERT INTO
  transactions (account_id, recepient_id,amount, type)
VALUES
  ($1, $2, $3 ,$4) RETURNING *;
