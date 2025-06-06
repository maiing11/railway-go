-- name: GetPassenger :one
SELECT * FROM  passengers
WHERE id = $1 LIMIT 1;

-- name: ListPassengers :many
SELECT * FROM passengers
ORDER BY name;

-- name: CreatePassenger :one
INSERT INTO passengers (id, name, id_number, user_id) 
VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdatePassenger :exec
UPDATE passengers
  set name = $2,
  id_number = $3,
  user_id = $4
WHERE id = $1;


-- name: DeletePassenger :exec
DELETE FROM passengers
WHERE id = $1;

-- name: GetPassengerByUser :one
SELECT * FROM passengers
WHERE user_id = $1;
