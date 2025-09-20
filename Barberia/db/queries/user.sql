-- name: CreateUser :exec
INSERT INTO users (nombre, email, telefono)
VALUES (?, ?, ?);

-- name: ListUsers :many
SELECT id, nombre, telefono, email
FROM users;
