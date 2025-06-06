-- name: GetSchedule :one
SELECT * FROM  schedules
WHERE id = $1 LIMIT 1;

-- name: SearchSchedules :many
SELECT 
  s.id AS schedule_id,
  t.name AS train_name,
  r.source_station AS source_station,
  r.destination_station AS destination_station,
  s.departure_date,
  s.arrival_date,
  s.available_seats,
  s.price
FROM schedules s
JOIN routes r ON s.route_id = r.id
JOIN trains t ON s.train_id = t.id 
WHERE 
  source_station ILIKE '%' || $1 || '%' AND
  destination_station ILIKE '%' || $2 || '%' AND
  DATE(s.departure_date) = $3
ORDER BY s.departure_date;

-- name: ListSchedules :many
SELECT * FROM schedules
ORDER BY departure_date;

-- name: CreateSchedule :one
INSERT INTO schedules (
   train_id, departure_date, arrival_date, available_seats, price, route_id
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateSchedule :exec
UPDATE schedules
  set train_id = $2,
  route_id = $3,
  departure_date = $4,
  arrival_date = $5,
  price = $6,
  available_seats = $7
WHERE id = $1;

-- name: DeleteSchedule :exec
DELETE FROM schedules
WHERE id = $1;
