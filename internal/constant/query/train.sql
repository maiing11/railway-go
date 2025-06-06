-- name: GetTrain :one
SELECT * FROM  trains
WHERE id = $1 LIMIT 1;

-- name: ListTrains :many
SELECT * FROM trains
ORDER BY name;

-- name: CreateTrain :one
INSERT INTO trains (
   name, capacity, created_at
) VALUES (
    $1, $2, now()
)
RETURNING *;

-- name: UpdateTrain :exec
UPDATE trains
  set name = $2,
  capacity = $3,
  updated_at = NOW()
WHERE id = $1;


-- name: DeleteTrain :exec
DELETE FROM trains
WHERE id = $1;
