package usecase

import (
	"context"
	"railway-go/internal/constant/model"
)

type UserSessionUC interface {
	Register(ctx context.Context, request *model.RegisterUserRequest) error
	Login(ctx context.Context, request model.LoginUserRequest) (map[string]any, error)
	Logout(ctx context.Context, sessionID string) error
	RenewAccessToken(ctx context.Context, refreshToken string) (string, error)
	CreateGuestSession(ctx context.Context, userAgent, clientIp string) (*model.Session, error)
	GetGuestSession(ctx context.Context, sessionID string) (*model.Session, error)
}
