-- name: GetRoute :one
SELECT * FROM  routes
WHERE id = $1 LIMIT 1;

-- name: ListRoute :many
SELECT * FROM routes
ORDER BY id;

-- name: CreateRoute :one
INSERT INTO routes (
   source_station, destination_station, travel_time
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: UpdateRoute :exec
UPDATE routes
  set source_station = $2,
  destination_station = $3
WHERE id = $1;


-- name: DeleteRoute :exec
DELETE FROM routes
WHERE id = $1;
