package usecase

import (
	"context"
	"errors"
	"railway-go/internal/constant/model"
	"railway-go/internal/repository"
	"railway-go/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (uc *UserSessionUsecase) Register(ctx context.Context, request *model.RegisterUserRequest) error {

	err := uc.Validate.Struct(request)
	if err != nil {
		uc.Logger.Warn("invalid request body: ", zap.Error(err))
		return fiber.ErrBadRequest
	}

	count, err := uc.Repo.CountUserByEmail(ctx, request.Email)
	if err != nil {
		uc.Logger.Warn("failed to count user from database", zap.Error(err))
		return fiber.ErrInternalServerError
	}

	if count > 0 {
		uc.Logger.Warn("email already exists", zap.String("email", request.Email))
		return fiber.ErrConflict
	}

	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		uc.Logger.Warn("failed to hash password", zap.Error(err))
		return fiber.ErrInternalServerError
	}
	id := uuid.New()
	user := repository.CreateUserParams{
		ID:          id,
		Name:        request.Name,
		Email:       request.Email,
		Password:    hashedPassword,
		Phonenumber: request.PhoneNumber,
	}

	if err = uc.Repo.CreateUser(ctx, user); err != nil {
		uc.Logger.Warn("failed to create user", zap.Error(err))
		return fiber.ErrInternalServerError
	}

	uc.Logger.Info("user registered successfully", zap.String("email", user.Email))
	return nil
}

// user login
func (uc *UserSessionUsecase) Login(ctx context.Context, request *model.LoginUserRequest) (map[string]any, error) {
	if ctx == nil {
		return nil, errors.New("context is nil")
	}

	fiberCtx, ok := ctx.Value("fiberCtx").(*fiber.Ctx)
	if !ok || fiberCtx == nil {
		return nil, errors.New("fiber.Ctx is missing in context")
	}

	if err := uc.Validate.Struct(request); err != nil {
		uc.Logger.Warn("invalid request body", zap.Any("request", request))
		return nil, fiber.ErrBadRequest
	}

	user, err := uc.Repo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		uc.Logger.Warn("failed find user", zap.String("email", request.Email))
		return nil, fiber.ErrUnauthorized
	}

	err = utils.CheckPassword(request.Password, user.Password)
	if err != nil {
		uc.Logger.Warn("failed to compare user password with bycrype hash", zap.Error(err))
		return nil, fiber.ErrUnauthorized
	}

	// generate session
	session := &model.Session{
		ID:           uuid.NewString(),
		RefreshToken: "",
		UserID:       &user.ID,
		Role:         "user",
		IsBlocked:    false,
		ExpiresAt:    time.Now().Add(uc.config.GetDuration("Token.RefreshTokenDuration")),
	}

	if err := uc.Repo.CreateSession(ctx, session); err != nil {
		uc.Logger.Error("failed to save session", zap.Error(err))
		return nil, err
	}

	// generate tokens with sesssion ID
	accessDuration := uc.config.GetDuration("Token.AccessTokenDuration")
	accessToken, accessPayload, err := uc.TokenMaker.CreateToken(&user.ID, session.ID, session.Role, accessDuration)
	if err != nil {
		uc.Logger.Error("failed to create token", zap.Error(err))
		return nil, fiber.ErrInternalServerError
	}
	// uc.Logger.Info("generate access token", zap.Any("payload", accessToken))

	refreshDuration := uc.config.GetDuration("Token.RefreshTokenDuration")
	refreshToken, refreshPayload, err := uc.TokenMaker.CreateToken(&user.ID, session.ID, session.Role, refreshDuration)
	if err != nil {
		uc.Logger.Error("failed to create refresh token", zap.Error(err))
		return nil, fiber.ErrInternalServerError
	}

	// update session with refresh token
	session.RefreshToken = refreshToken
	if err := uc.Repo.UpdateSessionAccessToken(ctx, session.ID, accessToken, session.ExpiresAt); err != nil {
		uc.Logger.Error("failed to update session", zap.Error(err))
		return nil, fiber.ErrInternalServerError
	}

	// set session ID in cookies
	cookie := &fiber.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Expires:  session.ExpiresAt,
		HTTPOnly: true,
		SameSite: "Lax",
		Secure:   true,
	}

	fiberCtx.Cookie(cookie)

	response := map[string]any{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": map[string]any{
			"id":    user.ID,
			"email": user.Email,
			"role":  session.Role,
		},
		"access_payload":  accessPayload,
		"refresh_payload": refreshPayload,
	}
	uc.Logger.Info("Login successful", zap.Any("userID", user.ID))
	return response, nil

}

// func (u *UserSessionUsecase) ForgotPassword(ctx context.Context, email string) error {
// 	user, err := u.Repo.GetUserByEmail(ctx, email)
// }

func (u *UserSessionUsecase) Logout(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		u.Logger.Error("Session ID is empty")
		return errors.New("invalid session")

	}
	err := u.Repo.DeleteSession(ctx, sessionID)
	if err != nil {
		u.Logger.Error("Failed to log out", zap.String("sessionID", sessionID), zap.Error(err))
		return err
	}
	u.Logger.Info("User logged out successfully", zap.String("sessionID", sessionID))
	return nil
}

func (u *UserSessionUsecase) GetSession(ctx context.Context, sessionID string) (*model.Session, error) {
	session, err := u.Repo.GetSessionByID(ctx, sessionID)
	if err != nil {
		u.Logger.Error("Failed to get session", zap.String("sessionID", sessionID), zap.Error(err))
		return nil, err
	}
	return session, nil
}
