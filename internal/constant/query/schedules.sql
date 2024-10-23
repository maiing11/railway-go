-- name: GetSchedule :one
SELECT * FROM  schedules
WHERE id = $1 LIMIT 1;

-- name: ListSchedules :many
SELECT * FROM trains
ORDER BY name;

-- name: CreateSchedule :one
INSERT INTO schedules (
   train_id, class_type, departure_date, available_seats, price, route_id
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateSchedule :exec
UPDATE schedules
  set train_id = $2,
  class_type = $3,
  departure_date = $4,
  available_seats= $5,
  price = $6,
  route_id = $7
WHERE id = $1;


-- name: DeleteSchedule :exec
DELETE FROM schedules
WHERE id = $1;
