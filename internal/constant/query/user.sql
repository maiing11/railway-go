-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY name;

-- name: CreateUser :exec
INSERT INTO users (
 id, name, email, password, phone_number, role, created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, now()
); 

-- name: UpdateUser :exec
UPDATE users
  set name = $2,
  email = $3,
  phone_number = $4,
  role = $5,
  updated_at = now()
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
