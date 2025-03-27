-- name: CreateSeat :one
INSERT INTO seats (wagon_id, seat_number, seat_row, is_available)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetSeat :one
SELECT * FROM seats
WHERE id = $1 LIMIT 1;

-- name: ListSeats :many
SELECT * FROM seats
ORDER BY id;

-- name: UpdateSeat :exec
UPDATE seats
  set wagon_id = $2,
  seat_number = $3,
  is_available = $4,
  seat_row = $5,
  updated_at = NOW()
WHERE id = $1;

-- name: DeleteSeat :exec
DELETE FROM seats
WHERE id = $1;
