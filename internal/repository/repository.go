package repository

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SQLStore provides all functions to execute SQL queries and transactions
type Store interface {
	Querier
	ReservationSeatTX(ctx context.Context, arg CreateReservationParams) (ReservationTxResult, error)
	SessionRepository
}

type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
	*sessionRepository
}

// NewStore creates a new store
func NewStore(connPool *pgxpool.Pool, redis *redis.Client) Store {
	return &SQLStore{
		connPool:          connPool,
		Queries:           New(connPool),
		sessionRepository: &sessionRepository{RedisClient: redis},
	}
}
