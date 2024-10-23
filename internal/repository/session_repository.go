package repository

import (
	"context"
	"encoding/json"
	"railway-go/internal/constant/entity"
	"time"

	"github.com/go-redis/redis/v8"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, session *entity.Session) error
	GetSessionByID(ctx context.Context, id string) (*entity.Session, error)
	DeleteSession(ctx context.Context, id string) error
}

type sessionRepository struct {
	redisClient *redis.Client
}

func NewSessionRepository(redisClient *redis.Client) SessionRepository {
	return &sessionRepository{redisClient: redisClient}
}

func (r *sessionRepository) CreateSession(ctx context.Context, session *entity.Session) error {
	sessionJson, err := json.Marshal(session)
	if err != nil {
		return err
	}

	return r.redisClient.Set(ctx, session.ID, sessionJson, time.Until(session.ExpiresAt)).Err()
}

func (r *sessionRepository) GetSessionByID(ctx context.Context, id string) (*entity.Session, error) {
	sessionData, err := r.redisClient.Get(ctx, id).Result()
	if err != nil {
		return nil, err
	}
	var session entity.Session
	err = json.Unmarshal([]byte(sessionData), &session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) DeleteSession(ctx context.Context, id string) error {
	return r.redisClient.Del(ctx, id).Err()
}
