-- +goose Up
CREATE TABLE
  users (
    user_id UUID PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash CHAR(60) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    date_joined TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  );


CREATE INDEX idx_users_on_email ON users(email);


CREATE INDEX idx_users_on_username ON users(username);


-- +goose Down 
DROP TABLE
  users;