-- +goose Up
CREATE TABLE sessions(
	id UUID PRIMARY KEY,
	user_id UUID NOT NULL,
	refresh_token TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	expires_at TIMESTAMP NOT NULL,
	user_agent TEXT NOT NULL,
	client_ip TEXT NOT NULL,
	is_revoked BOOL DEFAULT FALSE,

	CONSTRAINT fk_user
	FOREIGN KEY (user_id)
	REFERENCES users(id)
	ON DELETE CASCADE
);

-- +goose Down
DROP TABLE sessions;
