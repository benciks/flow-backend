-- name: CreateSession :one
INSERT INTO sessions (user_id, token)
VALUES (@user_id, @token)
RETURNING *;

-- name: DeleteSessionByUserID :one
DELETE FROM sessions
WHERE user_id = @user_id
RETURNING *;

-- name: GetSessionByToken :one
SELECT * FROM sessions
WHERE token = @token;