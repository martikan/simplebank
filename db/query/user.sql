-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
    username,
    email,
    password,
    full_name
) VALUES (
    $1, $2, $3, $4
) RETURNING *;