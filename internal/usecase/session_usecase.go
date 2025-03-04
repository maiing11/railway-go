package usecase

import (
	"context"
	"errors"
	"railway-go/internal/constant/model"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (uc *UserSessionUsecase) RenewAccessToken(ctx context.Context, refreshToken string) (string, error) {
	// retrieve session by refresh token
	session, err := uc.Repo.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		uc.Logger.Warn("session not found", zap.String("refreshToken", refreshToken), zap.Error(err))
		return "", model.ErrSessionNotFound
	}

	// validate session
	if session.IsBlocked || time.Now().After(session.ExpiresAt) {
		uc.Logger.Warn("session is invalid or expired", zap.String("sessionID", session.ID))
		return "", model.ErrSessionInvalid
	}

	// generate new access token
	tokenDuration := uc.config.GetDuration("token.AccessTokenDuration")
	newAccessToken, _, err := uc.TokenMaker.CreateToken(session.UserID, session.ID, session.Role, tokenDuration)
	if err != nil {
		uc.Logger.Error("failed to create new access token", zap.Error(err))
		return "", err
	}

	// update session in redis
	err = uc.Repo.UpdateSessionAccessToken(ctx, session.ID, newAccessToken, time.Now().Add(tokenDuration))
	if err != nil {
		uc.Logger.Error("failed to update session with new access token", zap.Error(err))
		return "", err
	}

	uc.Logger.Info("Access token renewed successfully", zap.String("sessionID", session.ID))
	return newAccessToken, nil
}

func (u *UserSessionUsecase) CreateGuestSession(ctx context.Context, userAgent, clientIp string) (*model.Session, error) {
	sessionID, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	expiry := time.Now().Add(24 * time.Hour)

	session := &model.Session{
		ID:        sessionID.String(),
		Role:      "guest",
		UserAgent: userAgent,
		ClientIP:  clientIp,
		IsBlocked: false,
		ExpiresAt: expiry,
	}

	// save session in redis
	if err := u.Repo.CreateSession(ctx, session); err != nil {
		u.Logger.Error("Failed to create guest session", zap.Error(err))
		return nil, errors.New("failed to create guest session")
	}
	u.Logger.Info("Guest session created", zap.String("session_id", session.ID))
	return session, nil
}

// GetGuestSession retrieves a guest session from redis
func (u *UserSessionUsecase) GetGuestSession(ctx context.Context, sessionID string) (*model.Session, error) {
	session, err := u.Repo.GetSessionByID(ctx, sessionID)
	if err != nil {
		u.Logger.Warn("Failed to retrieve guest session", zap.String("session_id", session.ID), zap.Error(err))
		return nil, model.ErrSessionNotFound
	}

	// check if the sessionis expired or blocked
	if session.IsBlocked || session.ExpiresAt.Before(time.Now()) {
		u.Logger.Warn("Guest session is expired or blocked", zap.String("session_id", session.ID))
		return nil, model.ErrSessionInvalid
	}

	return session, nil
}

func (u *UserSessionUsecase) DeleteGuestSession(ctx context.Context, sessionID string) error {
	if err := u.Repo.DeleteSession(ctx, sessionID); err != nil {
		u.Logger.Error("Failed to delete guest session", zap.String("session_id", sessionID))
		return model.ErrSessionNotFound
	}

	u.Logger.Info("Guest session deleted", zap.String("session_id", sessionID))
	return nil
}
