package repository

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SQLStore provides all functions to execute SQL queries and transactions
type Store interface {
	Querier
	SessionRepository
	ReservationRepository
	BeginTransaction(ctx context.Context) (Transaction, error)
}

type Transaction interface {
	Querier
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
	*redisRepository
}

type SQLTransaction struct {
	TX pgx.Tx
	*Queries
}

type redisRepository struct {
	RedisClient *redis.Client
}

func NewRedis(redis *redis.Client) *redisRepository {
	return &redisRepository{RedisClient: redis}
}

// NewStore creates a new store
func NewStore(connPool *pgxpool.Pool, redis *redis.Client) Store {
	return &SQLStore{
		connPool:        connPool,
		Queries:         New(connPool),
		redisRepository: NewRedis(redis),
	}
}

// BeginTransaction exec transaction
func (store *SQLStore) BeginTransaction(ctx context.Context) (Transaction, error) {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return &SQLTransaction{
		TX:      tx,
		Queries: New(tx),
	}, nil
}

func (t *SQLTransaction) Commit(ctx context.Context) error {
	return t.TX.Commit(ctx)
}

func (t *SQLTransaction) Rollback(ctx context.Context) error {
	return t.TX.Rollback(ctx)
}
