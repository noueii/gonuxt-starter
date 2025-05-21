-- name: CreateUser :one
INSERT INTO users(
	email, email_verified, name, hashed_password
) VALUES (
	$1, $2, $3, $4
) RETURNING *;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: UpdateUserBalance :one
UPDATE users
SET balance = $2
WHERE id = $1
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: UpdateUserById :one
UPDATE users
SET 
	name = COALESCE(sqlc.narg(name), name),
	hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
	balance = COALESCE(sqlc.narg(balance), balance),
	role = COALESCE(sqlc.narg(role), role),
	updated_at = NOW()
WHERE 
	id = sqlc.arg(id)
RETURNING *;
