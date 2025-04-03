-- name: CreateUser :one
INSERT INTO users(
	name, hashed_password
) VALUES (
	$1, $2
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

-- name: GetUserByName :one
SELECT * FROM users WHERE name = $1 LIMIT 1;


