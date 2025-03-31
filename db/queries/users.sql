-- name: CreateUser :one
INSERT INTO users(
	name
) VALUES (
	$1
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




