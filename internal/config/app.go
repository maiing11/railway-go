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
	usecase := usecase.NewUsecase(repo, config.Log, config.Validate, config.TokenMaker, config.Config)

	// setup controlers
	controller := http.NewController(usecase, config.Log)

	// setup middlewares
	middlewares := middleware.NewAuthMiddleware(*usecase, config.TokenMaker)

	// setup routes
	routeConfig := route.RouteConfig{
		App:            config.App,
		UserController: controller,
		AuthMiddleware: middlewares,
	}

	routeConfig.Setup()
}
