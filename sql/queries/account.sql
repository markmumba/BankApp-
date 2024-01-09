-- name: CreateAccount :one
INSERT INTO
  accounts (user_id, account_number, account_type, balance)
VALUES
  (
    $1,
    lpad(floor(random() * 10 ^ 15):: text, 15, '0'),
    $2,
    0.00
  ) RETURNING *;


-- name: FindAccount :many
SELECT
  u.user_id,
  u.username,
  a.account_id,
  a.account_number,
  a.account_type,
  a.balance
FROM
  users u
  JOIN accounts a ON u.user_id = a.user_id
WHERE
  u.user_id = $1;


-- name: Deposit :one
UPDATE
  accounts
SET
  balance = $1
WHERE
  account_id = $2 RETURNING *;


-- name: Withdraw :one 
UPDATE
  accounts
SET
  balance = $1
WHERE
  account_id = $2 RETURNING *;


-- name: ViewTransactions :many 
SELECT
  *
FROM
  transactions
WHERE
  account_id = $1
  OR recepient_id = $1
ORDER BY
  timestamp DESC;


-- name: CloseAccount :exec
DELETE FROM
  accounts
WHERE
  account_id = $1;