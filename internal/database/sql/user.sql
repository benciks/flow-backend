-- name: GetUsers :many
SELECT * FROM users;

-- name: FindUserByUsername :one
SELECT * FROM users WHERE username = @username;

-- name: FindUserById :one
SELECT * FROM users WHERE id = @id;

-- name: CreateUser :one
INSERT INTO users (password, username) VALUES (@password, @username) RETURNING *;