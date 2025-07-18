// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: train.sql

package repository

import (
	"context"
)

const createTrain = `-- name: CreateTrain :one
INSERT INTO trains (
   name, capacity, created_at
) VALUES (
    $1, $2, now()
)
RETURNING id, name, capacity, created_at, updated_at
`

type CreateTrainParams struct {
	Name     string `db:"name" json:"name"`
	Capacity int32  `db:"capacity" json:"capacity"`
}

func (q *Queries) CreateTrain(ctx context.Context, arg CreateTrainParams) (Train, error) {
	row := q.db.QueryRow(ctx, createTrain, arg.Name, arg.Capacity)
	var i Train
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Capacity,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteTrain = `-- name: DeleteTrain :exec
DELETE FROM trains
WHERE id = $1
`

func (q *Queries) DeleteTrain(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteTrain, id)
	return err
}

const getTrain = `-- name: GetTrain :one
SELECT id, name, capacity, created_at, updated_at FROM  trains
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetTrain(ctx context.Context, id int64) (Train, error) {
	row := q.db.QueryRow(ctx, getTrain, id)
	var i Train
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Capacity,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listTrains = `-- name: ListTrains :many
SELECT id, name, capacity, created_at, updated_at FROM trains
ORDER BY name
`

func (q *Queries) ListTrains(ctx context.Context) ([]Train, error) {
	rows, err := q.db.Query(ctx, listTrains)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Train{}
	for rows.Next() {
		var i Train
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Capacity,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTrain = `-- name: UpdateTrain :exec
UPDATE trains
  set name = $2,
  capacity = $3,
  updated_at = NOW()
WHERE id = $1
`

type UpdateTrainParams struct {
	ID       int64  `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Capacity int32  `db:"capacity" json:"capacity"`
}

func (q *Queries) UpdateTrain(ctx context.Context, arg UpdateTrainParams) error {
	_, err := q.db.Exec(ctx, updateTrain, arg.ID, arg.Name, arg.Capacity)
	return err
}
