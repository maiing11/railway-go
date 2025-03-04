package main

import (
	"fmt"
	"railway-go/internal/config"
	"railway-go/internal/utils/token"

	"go.uber.org/zap"
)

func main() {
	// Load Configuration
	viperConfig := config.NewViper()

	// Initialize Logger
	log, err := config.NewLogger()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}

	// Initialize Database
	db, err := config.NewDatabase(viperConfig, log)
	if err != nil {
		log.Fatal("failed to initialize database", zap.Error(err))
	}

	// Initialize Validator
	validate := config.NewValidator(viperConfig)

	// Initialize Fiber App
	app := config.NewFiber(viperConfig)

	// Initialize Token Manager (PASETO)
	tokenSecret := viperConfig.GetString("token.secret")
	fmt.Println("Token Secret Length:", len(tokenSecret))

	tokenMaker, err := token.NewTokenManager(tokenSecret)
	if err != nil {
		log.Fatal("failed to initialize token maker", zap.Error(err))
	}

	// Initialize Redis Client
	redisClient := config.NewRedisClient(viperConfig, log)

	// Bootstrap the application
	config.Boostrap(config.BootstrapConfig{
		DB:          db,
		App:         app,
		Log:         log,
		Validate:    validate,
		Config:      viperConfig,
		TokenMaker:  tokenMaker,
		RedisClient: redisClient,
	})

	// Start the Web Server
	webPort := viperConfig.GetInt("web.port")
	err = app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatal("failed to start web server", zap.Error(err))
	}
}
