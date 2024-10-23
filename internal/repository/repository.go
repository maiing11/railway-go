package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// SQLStore provides all functions to execute SQL queries and transactions
type Store interface {
	Querier
	ReservationSeatTX(ctx context.Context, arg CreateReservationParams) (ReservationTxResult, error)
}

type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
}

// NewStore creates a new store
func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
