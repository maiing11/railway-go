package repository

import (
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SQLStore provides all functions to execute SQL queries and transactions
type Store interface {
	Querier
	SessionRepository
	ReservationRepository
}

type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
	*redisRepository
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
