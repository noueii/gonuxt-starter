-- +goose Up
CREATE TABLE users(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	role TEXT NOT NULL DEFAULT 'user',
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	name TEXT NOT NULL,
	email TEXT UNIQUE NOT NULL,
	email_verified BOOL NOT NULL DEFAULT FALSE,
	hashed_password TEXT
);

-- +goose Down
DROP TABLE users;
