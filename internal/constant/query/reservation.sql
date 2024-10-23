-- name: GetReservation :one
SELECT * FROM  reservations
WHERE id = $1 LIMIT 1;

-- name: ListReservations :many
SELECT * FROM reservations
WHERE
 passenger_id = $1 OR
 schedule_id = $2
ORDER BY booking_date
LIMIT $3
OFFSET $4;

-- name: CreateReservation :one
INSERT INTO reservations (
   passenger_id, schedule_id, seat_number, booking_date, payment_id, status
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateReservation :exec
UPDATE reservations
  set  passenger_id = $2 , schedule_id = $3, seat_number = $4, booking_date = $5, payment_id = $6, status = $7
WHERE id = $1;


-- name: DeleteReservation :exec
DELETE FROM reservations
WHERE id = $1;
