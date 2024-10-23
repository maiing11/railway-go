-- name: GetTrain :one
SELECT * FROM  trains
WHERE id = $1 LIMIT 1;

-- name: ListTrains :many
SELECT * FROM trains
ORDER BY name;

-- name: CreateTrain :one
INSERT INTO trains (
   name, capacity
) VALUES (
    $1, $2
)
RETURNING *;

-- name: UpdateTrain :exec
UPDATE trains
  set name = $2,
  capacity = $3
WHERE id = $1;


-- name: DeleteTrain :exec
DELETE FROM trains
WHERE id = $1;
