-- +goose Up
CREATE TABLE
  accounts (
    account_id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES Users(user_id) ON DELETE CASCADE,
    account_number VARCHAR (60) UNIQUE NOT NULL,
    account_type VARCHAR(20) NOT NULL,
    balance DECIMAL(15, 2) NOT NULL DEFAULT 0.0,
    date_opened TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, account_type)
  );

CREATE INDEX  idx_account_number ON accounts(account_number);

-- +goose Down 
DROP TABLE
  accounts;
