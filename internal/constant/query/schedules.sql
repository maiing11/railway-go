-- name: GetSchedule :one
SELECT * FROM  schedules
WHERE id = $1 LIMIT 1;

-- name: ListSchedules :many
SELECT * FROM trains
ORDER BY name;

-- name: CreateSchedule :one
INSERT INTO schedules (
   train_id, departure_time, arrival_time, available_seats, price, route_id
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateSchedule :exec
UPDATE schedules
  set train_id = $2,
  route_id = $3,
  departure_time = $4,
  arrival_time = $5,
  price = $6,
  available_seats = $7
WHERE id = $1;


-- name: DeleteSchedule :exec
DELETE FROM schedules
WHERE id = $1;
