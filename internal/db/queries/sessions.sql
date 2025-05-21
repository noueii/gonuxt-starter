-- name: CreateSession :one
INSERT INTO sessions(
	id, user_id, refresh_token, expires_at, user_agent, client_ip
) VALUES ( $1, $2, $3, $4, $5, $6 )
RETURNING *;
	
-- name: GetSessionById :one
SELECT sessions.*, users.email AS email, users.name AS username, users.role as role
FROM sessions
INNER JOIN users ON users.id = sessions.user_id
WHERE sessions.id = $1 
LIMIT 1;
