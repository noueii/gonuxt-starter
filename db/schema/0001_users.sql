-- +goose Up
CREATE TABLE users(
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	role TEXT NOT NULL DEFAULT 'user',
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	name TEXT UNIQUE NOT NULL,
	hashed_password TEXT NOT NULL
);

-- +goose Down
DROP TABLE users;
