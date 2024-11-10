-- name: CreateUser :one
INSERT INTO users ( created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE id  = $1;
-- name: GetUserWithName :one
SELECT * FROM users WHERE name  = $1;

-- name: UserExists :one
SELECT EXISTS (
    SELECT 1
    FROM users
    WHERE name = $1
) AS exists;


-- name: Reset :exec
DELETE FROM users;

-- name: GetUsers :many
SELECT * FROM users;
