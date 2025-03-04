package usecase

import (
	"context"
	"railway-go/internal/constant/model"
	"railway-go/internal/repository"
	"railway-go/internal/utils/token"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type UserSessionUC interface {
	Register(ctx context.Context, request *model.RegisterUserRequest) error
	Login(ctx context.Context, request model.LoginUserRequest) (map[string]any, error)
	Logout(ctx context.Context, sessionID string) error
	RenewAccessToken(ctx context.Context, refreshToken string) (string, error)
	CreateGuestSession(ctx context.Context, userAgent, clientIp string) (*model.Session, error)
	GetGuestSession(ctx context.Context, sessionID string) (*model.Session, error)
}

type UserSessionUsecase struct {
	Repo       repository.Store
	Logger     *zap.Logger
	Validate   *validator.Validate
	TokenMaker token.Maker
	config     *viper.Viper
}

func NewUsecase(
	repo repository.Store,
	logger *zap.Logger,
	validate *validator.Validate,
	tokenMaker token.Maker,
	config *viper.Viper,
) *UserSessionUsecase {
	return &UserSessionUsecase{
		Repo:       repo,
		Logger:     logger,
		Validate:   validate,
		TokenMaker: tokenMaker,
		config:     config,
	}
}
