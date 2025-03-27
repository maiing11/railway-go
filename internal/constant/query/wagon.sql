-- name: CreateWagon :one
INSERT INTO wagons (train_id, wagon_number, class_type, total_seats)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetWagon :one
SELECT * FROM wagons
WHERE id = $1 LIMIT 1;

-- name: ListWagons :many
SELECT * FROM wagons
ORDER BY id;

-- name: UpdateWagon :exec
UPDATE wagons
  set train_id = $2,
  wagon_number = $3,
  class_type = $4,
  total_seats = $5,
  updated_at = NOW()
WHERE id = $1;

-- name: DeleteWagon :exec
DELETE FROM wagons
WHERE id = $1;
