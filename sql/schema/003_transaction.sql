-- +goose Up
CREATE TABLE
  transactions (
    transaction_id SERIAL PRIMARY KEY,
    account_id INT REFERENCES accounts(account_id) ON DELETE CASCADE,
    recepient_id INT REFERENCES accounts(account_id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  );


-- +goose Down 
DROP TABLE
  transactions;