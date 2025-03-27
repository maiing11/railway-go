package repository

import (
	"context"
	"encoding/json"
	"railway-go/internal/constant/model"
	"time"

	"github.com/go-redis/redis/v8"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, session *model.Session) error
	GetSessionByID(ctx context.Context, id string) (*model.Session, error)
	DeleteSession(ctx context.Context, id string) error
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*model.Session, error)
	UpdateSessionAccessToken(ctx context.Context, sessionID string, newAccessToken string, newExpiresAt time.Time) error
}

func (r *redisRepository) CreateSession(ctx context.Context, session *model.Session) error {
	sessionJson, err := json.Marshal(session)
	if err != nil {
		return err
	}

	return r.RedisClient.Set(ctx, session.ID, sessionJson, time.Until(session.ExpiresAt)).Err()
}

func (r *redisRepository) GetSessionByID(ctx context.Context, id string) (*model.Session, error) {
	sessionData, err := r.RedisClient.Get(ctx, id).Result()
	if err != nil {
		return nil, err
	}
	var session model.Session
	err = json.Unmarshal([]byte(sessionData), &session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *redisRepository) DeleteSession(ctx context.Context, id string) error {
	return r.RedisClient.Del(ctx, id).Err()
}

func (r *redisRepository) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*model.Session, error) {
	sessionKey := "session:" + refreshToken
	val, err := r.RedisClient.Get(ctx, sessionKey).Bytes()
	if err == redis.Nil {
		return nil, model.ErrSessionNotFound
	} else if err != nil {
		return nil, err
	}

	session := &model.Session{}
	if err := json.Unmarshal(val, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (r *redisRepository) UpdateSessionAccessToken(ctx context.Context, sessionID string, newAccessToken string, newExpiresAt time.Time) error {
	sessionKey := "session:" + sessionID
	updateData := map[string]interface{}{
		"access_token": newAccessToken,
		"expires_at":   newExpiresAt.Format(time.RFC3339),
	}
	return r.RedisClient.HMSet(ctx, sessionKey, updateData).Err()
}
