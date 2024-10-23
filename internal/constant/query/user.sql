-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY name;

-- name: CreateUser :exec
INSERT INTO users (
  name, email, password, phoneNumber
) VALUES (
    $1, $2, $3, $4
); 

-- name: UpdateUser :exec
UPDATE users
  set name = $2,
  email = $3,
  phoneNumber = $4
WHERE id = $1;

-- name: UpdateUserPassword :exec
UPDATE users
  set password = $2
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: CountUserByEmail :one
SELECT COUNT(*) as count FROM users
WHERE email = $1;
