package config

import (
	"railway-go/internal/delivery/http"
	"railway-go/internal/delivery/http/middleware"
	"railway-go/internal/delivery/http/route"
	"railway-go/internal/repository"
	"railway-go/internal/usecase"
	"railway-go/internal/utils/token"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type BootstrapConfig struct {
	DB          *pgxpool.Pool
	App         *fiber.App
	Log         *zap.Logger
	Validate    *validator.Validate
	Config      *viper.Viper
	TokenMaker  token.Maker
	RedisClient *redis.Client
}

// do bootstrap here
func Boostrap(config BootstrapConfig) {
	// setup repositories
	repo := repository.NewStore(config.DB, config.RedisClient)

	// setup usecases
	userSessionUC := usecase.NewUserSessionUsecase(repo, config.Log, config.Validate, config.TokenMaker, config.Config)

	// setup controlers
	userSesioncontroller := http.NewUserSessionController(userSessionUC, config.Log)

	// setup middlewares
	userSessionMiddlewares := middleware.NewAuthMiddleware(*userSessionUC, config.TokenMaker)

	// setup routes
	routeConfig := route.RouteConfig{
		App:                   config.App,
		UserSessionController: userSesioncontroller,
		AuthMiddleware:        userSessionMiddlewares,
	}

	routeConfig.Setup()
}
