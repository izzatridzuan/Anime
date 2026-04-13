-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC;

-- name: CreateUser :one
INSERT INTO users(email, name, password_hash, role)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $2, must_change_password = false
WHERE id = $1;

-- name: UpdateUser :one
UPDATE users
SET name = $2, email = $3, role = $4
WHERE id = $1
RETURNING *;

-- name: UpdateUserName :exec
UPDATE users SET name = $2 WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: ResetUserPassword :exec
UPDATE users 
SET password_hash = $2, must_change_password = true
WHERE id = $1;


-- name: SearchUsers :many
SELECT * FROM users
WHERE
    ($1::text = '' OR LOWER(email) LIKE LOWER('%' || $1 || '%'))
    OR ($1::text = '' OR LOWER(name) LIKE LOWER('%' || $1 || '%'))
ORDER BY created_at DESC;