-- +goose Up
ALTER TABLE users
ADD COLUMN balance INT NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE users
DROP COLUMN balance;
