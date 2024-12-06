package usecase

import (
	"context"
	"railway-go/internal/constant/entity"
	"railway-go/internal/constant/model"
	"railway-go/internal/utils"

	"railway-go/internal/repository"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func (u *Usecase) Register(ctx context.Context, request *model.RegisterUserRequest) error {
	err := u.Validate.Struct(request)
	if err != nil {
		u.Logger.Warn("invalid request body: ", zap.Error(err))
		return fiber.ErrBadRequest
	}

	count, err := u.Repo.CountUserByEmail(ctx, request.Email)
	if err != nil {
		u.Logger.Warn("failed to count user from database", zap.Error(err))
		return fiber.ErrInternalServerError
	}

	if count > 0 {
		u.Logger.Warn("email already exists", zap.String("email", request.Email))
		return fiber.ErrConflict
	}

	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		u.Logger.Warn("failed to hash password", zap.Error(err))
		return fiber.ErrInternalServerError
	}

	user := repository.CreateUserParams{
		Name:        request.Name,
		Email:       request.Email,
		Password:    hashedPassword,
		Phonenumber: request.PhoneNumber,
	}

	if err = u.Repo.CreateUser(ctx, user); err != nil {
		u.Logger.Warn("failed to create user", zap.Error(err))
		return fiber.ErrInternalServerError
	}

	u.Logger.Info("user registered successfully", zap.String("email", user.Email))
	return nil
}

// user login
func (u *Usecase) Login(ctx context.Context, request model.LoginUserRequest) (map[string]interface{}, error) {
	if err := u.Validate.Struct(request); err != nil {
		u.Logger.Warn("invalid request body", zap.Any("request", request))
		return nil, fiber.ErrBadRequest
	}

	user, err := u.Repo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		u.Logger.Warn("failed find user", zap.String("email", request.Email))
		return nil, fiber.ErrUnauthorized
	}

	err = utils.CheckPassword(request.Password, user.Password)
	if err != nil {
		u.Logger.Warn("failed to compare user password with bycrype hash", zap.Error(err))
		return nil, fiber.ErrUnauthorized
	}

	role := "user"
	accessDuration := u.config.GetDuration("Token.AccessTokenDuration")
	accessToken, accessPayload, err := u.TokenMaker.CreateToken(&user.ID, role, accessDuration)
	if err != nil {
		u.Logger.Error("failed to create token", zap.Error(err))
	}

	refreshDuration := u.config.GetDuration("Token.RefreshTokenDuration")
	refreshToken, refreshPayload, err := u.TokenMaker.CreateToken(&user.ID, role, refreshDuration)
	if err != nil {
		u.Logger.Error("failed to create refresh token", zap.Error(err))
	}

	session := &entity.Session{
		ID:           refreshPayload.ID.String(),
		RefreshToken: refreshToken,
		UserID:       &user.ID,
		Role:         role,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	}

	if err := u.Repo.CreateSession(ctx, session); err != nil {
		u.Logger.Error("failed to save error", zap.Error(err))
		return nil, err
	}

	response := map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
			"role":  session.Role,
		},
		"access_payload":  accessPayload,
		"refresh_payload": refreshPayload,
	}
	u.Logger.Info("Login successful", zap.Any("userID", user.ID))
	return response, nil

}
