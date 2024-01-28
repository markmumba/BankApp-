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
  a.balance,
  a.date_opened
FROM
  users u
  JOIN accounts a ON u.user_id = a.user_id
WHERE
  u.user_id = $1;


-- name: FindAccountByAccNo :one
SELECT
  *
FROM
  accounts
WHERE
  account_number = $1;


-- name: FindAccountById :one
SELECT
  *
FROM
  accounts
WHERE
  account_id = $1;


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
  t.transaction_id,
  t.account_id,
  t.recepient_id,
  t.amount,
  t.type,
  t.timestamp,
  a.account_type
FROM
  transactions t
  JOIN accounts a ON t.account_id = a.account_id
WHERE
  t.account_id = $1
  OR t.recepient_id = $1
ORDER BY
  timestamp DESC;


-- name: CheckToSave :exec
UPDATE
  accounts
SET
  balance = $1
WHERE
  account_id = $2
  AND account_type = $3;


-- name: CloseAccount :exec
DELETE FROM
  accounts
WHERE
  account_id = $1;


-- name: GetAllAccounts :many 
SELECT
  *
FROM
  accounts
ORDER BY
  user_id;