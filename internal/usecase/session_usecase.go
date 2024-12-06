package usecase

import (
	"context"
	"railway-go/internal/constant/entity"
	"time"

	"go.uber.org/zap"
)

func (uc *Usecase) RenewAccessToken(ctx context.Context, refreshToken string) (string, error) {
	// retrieve session by refresh token
	session, err := uc.Repo.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		uc.Logger.Warn("session not found", zap.String("refreshToken", refreshToken), zap.Error(err))
		return "", entity.ErrSessionNotFound
	}

	// validate session
	if session.IsBlocked || time.Now().After(session.ExpiresAt) {
		uc.Logger.Warn("session is invalid or expired", zap.String("sessionID", session.ID))
		return "", entity.ErrSessionInvalid
	}

	// generate new access token
	tokenDuration := uc.config.GetDuration("token.AccessTokenDuration")
	newAccessToken, _, err := uc.TokenMaker.CreateToken(session.UserID, session.Role, tokenDuration)
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
