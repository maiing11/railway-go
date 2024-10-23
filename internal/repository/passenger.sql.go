// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: passenger.sql

package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createPassenger = `-- name: CreatePassenger :one
INSERT INTO passengers (
  id, name, id_number, user_id
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, name, id_number, user_id, created_at, updated_at
`

type CreatePassengerParams struct {
	ID       int64       `db:"id" json:"id"`
	Name     string      `db:"name" json:"name"`
	IDNumber string      `db:"id_number" json:"id_number"`
	UserID   pgtype.UUID `db:"user_id" json:"user_id"`
}

func (q *Queries) CreatePassenger(ctx context.Context, arg CreatePassengerParams) (Passenger, error) {
	row := q.db.QueryRow(ctx, createPassenger,
		arg.ID,
		arg.Name,
		arg.IDNumber,
		arg.UserID,
	)
	var i Passenger
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.IDNumber,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deletePassenger = `-- name: DeletePassenger :exec
DELETE FROM passengers
WHERE id = $1
`

func (q *Queries) DeletePassenger(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deletePassenger, id)
	return err
}

const getPassenger = `-- name: GetPassenger :one
SELECT id, name, id_number, user_id, created_at, updated_at FROM  passengers
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetPassenger(ctx context.Context, id int64) (Passenger, error) {
	row := q.db.QueryRow(ctx, getPassenger, id)
	var i Passenger
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.IDNumber,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listPassengers = `-- name: ListPassengers :many
SELECT id, name, id_number, user_id, created_at, updated_at FROM passengers
ORDER BY name
`

func (q *Queries) ListPassengers(ctx context.Context) ([]Passenger, error) {
	rows, err := q.db.Query(ctx, listPassengers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Passenger{}
	for rows.Next() {
		var i Passenger
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.IDNumber,
			&i.UserID,
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

const updatePassenger = `-- name: UpdatePassenger :exec
UPDATE passengers
  set name = $2,
  id_number = $3,
  user_id = $4
WHERE id = $1
`

type UpdatePassengerParams struct {
	ID       int64       `db:"id" json:"id"`
	Name     string      `db:"name" json:"name"`
	IDNumber string      `db:"id_number" json:"id_number"`
	UserID   pgtype.UUID `db:"user_id" json:"user_id"`
}

func (q *Queries) UpdatePassenger(ctx context.Context, arg UpdatePassengerParams) error {
	_, err := q.db.Exec(ctx, updatePassenger,
		arg.ID,
		arg.Name,
		arg.IDNumber,
		arg.UserID,
	)
	return err
}
