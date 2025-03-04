package http

import (
	"railway-go/internal/usecase"

	"go.uber.org/zap"
)

type Controllers struct {
	Log     *zap.Logger
	Usecase *usecase.UserSessionUsecase
}

func NewController(usecase *usecase.UserSessionUsecase, log *zap.Logger) *Controllers {
	return &Controllers{
		Usecase: usecase,
		Log:     log,
	}
}
