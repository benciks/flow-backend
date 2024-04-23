-- name: GetUsers :many
SELECT * FROM users;

-- name: FindUserByUsername :one
SELECT * FROM users WHERE username = @username;

-- name: FindUserById :one
SELECT * FROM users WHERE id = @id;

-- name: CreateUser :one
INSERT INTO users (password, username) VALUES (@password, @username) RETURNING *;

-- name: SaveUserUUID :one
UPDATE users SET taskd_uuid = @uuid WHERE id = @id RETURNING *;

-- name: SaveTimewID :one
UPDATE users SET timew_id = @timew_id WHERE id = @id RETURNING *;