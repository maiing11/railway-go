package usecase

import (
	"railway-go/internal/repository"
	"railway-go/internal/utils/token"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Usecase struct {
	Repo       repository.Store
	Logger     *zap.Logger
	Validate   *validator.Validate
	TokenMaker token.Maker
	config     *viper.Viper
}
