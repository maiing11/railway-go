-- name: CreateStation :one
INSERT INTO stations (code, station_name)
VALUES ($1, $2)
RETURNING *;

-- name: GetStation :one
SELECT * FROM stations
WHERE id = $1 LIMIT 1;

-- name: GetStationByCode :one
SELECT * FROM stations
WHERE code = $1 LIMIT 1;

-- name: GetStationByName :one
SELECT * FROM stations
WHERE station_name = $1 LIMIT 1;

-- name: ListStations :many
SELECT * FROM stations
ORDER BY id;

-- name: UpdateStation :exec
UPDATE stations
  set code = $2,
  station_name = $3,
  updated_at = NOW()
WHERE id = $1;

-- name: DeleteStation :exec
DELETE FROM stations
WHERE id = $1;

