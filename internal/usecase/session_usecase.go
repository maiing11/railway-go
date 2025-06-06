package usecase

import (
	"context"
	"errors"
	"railway-go/internal/constant/model"
	"railway-go/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (uc *UserSessionUsecase) RenewAccessToken(ctx context.Context, refreshToken string) (string, error) {
	// retrieve session by refresh token
	session, err := uc.Repo.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		uc.UseCase.Log.Warn("session not found", zap.String("refreshToken", refreshToken), zap.Error(err))
		return "", model.ErrSessionNotFound
	}

	// validate session
	if session.IsBlocked || time.Now().After(session.ExpiresAt) {
		uc.UseCase.Log.Warn("session is invalid or expired", zap.String("sessionID", session.ID))
		return "", model.ErrSessionInvalid
	}

	// generate new access token
	tokenDuration := uc.config.GetDuration("token.AccessTokenDuration")
	newAccessToken, _, err := uc.TokenMaker.CreateToken(session.UserID, session.ID, session.Role, tokenDuration)
	if err != nil {
		uc.UseCase.Log.Error("failed to create new access token", zap.Error(err))
		return "", err
	}

	// update session in redis
	err = uc.Repo.UpdateSessionAccessToken(ctx, session.ID, newAccessToken, time.Now().Add(tokenDuration))
	if err != nil {
		uc.UseCase.Log.Error("failed to update session with new access token", zap.Error(err))
		return "", err
	}

	uc.UseCase.Log.Info("Access token renewed successfully", zap.String("sessionID", session.ID))
	return newAccessToken, nil
}

func (u *UserSessionUsecase) CreateGuestSession(ctx context.Context) (model.Session, error) {

	expiry := time.Now().Add(5 * time.Hour)

	fiberCtx, ok := ctx.Value("fiberCtx").(*fiber.Ctx)
	if !ok || fiberCtx == nil {
		return model.Session{}, fiber.NewError(fiber.StatusBadRequest, "fiber context is nil")
	}

	session := &model.Session{
		ID:           uuid.NewString(),
		RefreshToken: "",
		Role:         "guest",
		UserAgent:    fiberCtx.Get("User-Agent"),
		ClientIP:     fiberCtx.IP(),
		IsBlocked:    false,
		ExpiresAt:    expiry,
	}

	// save session in redis
	if err := u.Repo.CreateSession(ctx, session); err != nil {
		u.UseCase.Log.Error("Failed to create guest session", zap.Error(err))
		return model.Session{}, errors.New("failed to create guest session")
	}

	cookie := &fiber.Cookie{
		Name:    "session_id",
		Value:   session.ID,
		Expires: session.ExpiresAt,
		Secure:  true,
		Path:    "/",
	}

	fiberCtx.Cookie(cookie)

	u.UseCase.Log.Info("Guest session created", zap.String("session_id", session.ID))
	return *session, nil
}

// GetSession retrieves a session from redis
func (u *UserSessionUsecase) GetSession(ctx context.Context, sessionID string) (model.Session, error) {

	session, err := u.Repo.GetSessionByID(ctx, sessionID)
	if err != nil {
		u.UseCase.Log.Warn("Failed to retrieve session", zap.String("session_id", session.ID), zap.Error(err))
		return model.Session{}, model.ErrSessionNotFound
	}

	// check if the sessionis expired or blocked
	if session.IsBlocked || session.ExpiresAt.Before(time.Now()) {
		u.UseCase.Log.Warn("session is expired or blocked", zap.String("session_id", session.ID))
		return model.Session{}, model.ErrSessionInvalid
	}

	return *session, nil
}

func (u *UserSessionUsecase) DeleteGuestSession(ctx context.Context, sessionID string) error {
	if err := u.Repo.DeleteSession(ctx, sessionID); err != nil {
		u.UseCase.Log.Error("Failed to delete guest session", zap.String("session_id", sessionID))
		return model.ErrSessionNotFound
	}

	u.UseCase.Log.Info("Guest session deleted", zap.String("session_id", sessionID))
	return nil
}

func (u *UserSessionUsecase) GetUserIDFromSession(ctx context.Context, sessionID string) (*uuid.UUID, error) {
	session, err := u.GetSession(ctx, sessionID)
	if err != nil {
		return nil, utils.WrapError(fiber.StatusUnauthorized, u.UseCase.Log, utils.Error, err, "failed to get session")
	}

	if session.UserID == nil {
		return nil, utils.WrapError(fiber.StatusUnauthorized, u.UseCase.Log, utils.Error, err, "user ID is nil")
	}

	return session.UserID, nil
}
