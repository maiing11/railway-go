package usecase

import (
	"context"
	"railway-go/internal/constant/model"
	"railway-go/internal/utils"
	"railway-go/internal/utils/token"
	"time"

	"railway-go/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Usecase struct {
	Repo     repository.Store
	Logger   *zap.Logger
	Validate *validator.Validate
	TokenMaker token.Maker
	Session  repository.SessionRepository
}

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

func (u *Usecase) Login(ctx context.Context, request model.LoginUserRequest) error {
	if err := u.Validate.Struct(request); err != nil {
		u.Logger.Warn("invalid request body", zap.Any("request", request))
		return fiber.ErrBadRequest
	}

	user, err := u.Repo.GetUserByEmail(ctx, request.Email)
	if err != nil{
		u.Logger.Warn("failed find user", zap.String("email", request.Email))
		return fiber.ErrUnauthorized
	}

	err = utils.CheckPassword(request.Password, user.Password)
	if err != nil {
		u.Logger.Warn("failed to compare user password with bycrype hash", zap.Error(err))
		return fiber.ErrUnauthorized
	}
	
	role := "user"
	
	token, _ ,err := u.TokenMaker.CreateToken(&user.Email, role, 24*time.Hour)
	if err != nil{
		u.Logger.Error("failed to create token", zap.Error(err))
	}
	

}
