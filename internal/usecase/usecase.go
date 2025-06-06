package usecase

import (
	"railway-go/internal/repository"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type UseCase struct {
	Repo     repository.Store
	Log      *zap.Logger
	Validate *validator.Validate
}

func NewUseCase(
	repo repository.Store,
	log *zap.Logger,
	validate *validator.Validate,
) *UseCase {
	return &UseCase{
		Repo:     repo,
		Log:      log,
		Validate: validate,
	}
}
