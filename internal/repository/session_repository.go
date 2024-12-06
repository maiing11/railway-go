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
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*entity.Session, error)
	UpdateSessionAccessToken(ctx context.Context, sessionID string, newAccessToken string, newExpiresAt time.Time) error
}

type sessionRepository struct {
	RedisClient *redis.Client
}

func NewSessionRepository(redisClient *redis.Client) SessionRepository {
	return &sessionRepository{RedisClient: redisClient}
}

func (r *sessionRepository) CreateSession(ctx context.Context, session *entity.Session) error {
	sessionJson, err := json.Marshal(session)
	if err != nil {
		return err
	}

	return r.RedisClient.Set(ctx, session.ID, sessionJson, time.Until(session.ExpiresAt)).Err()
}

func (r *sessionRepository) GetSessionByID(ctx context.Context, id string) (*entity.Session, error) {
	sessionData, err := r.RedisClient.Get(ctx, id).Result()
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
	return r.RedisClient.Del(ctx, id).Err()
}

func (r *sessionRepository) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*entity.Session, error) {
	sessionKey := "session:" + refreshToken
	val, err := r.RedisClient.Get(ctx, sessionKey).Bytes()
	if err == redis.Nil {
		return nil, entity.ErrSessionNotFound
	} else if err != nil {
		return nil, err
	}

	session := &entity.Session{}
	if err := json.Unmarshal(val, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (r *sessionRepository) UpdateSessionAccessToken(ctx context.Context, sessionID string, newAccessToken string, newExpiresAt time.Time) error {
	sessionKey := "session:" + sessionID
	updateData := map[string]interface{}{
		"access_token": newAccessToken,
		"expires_at":   newExpiresAt.Format(time.RFC3339),
	}
	return r.RedisClient.HMSet(ctx, sessionKey, updateData).Err()
}
