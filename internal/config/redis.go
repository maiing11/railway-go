package config

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewRedisClient(cfg *viper.Viper, log *zap.Logger) *redis.Client {
	addr := fmt.Sprintf("%s:%d", cfg.GetString("redis.host"), cfg.GetInt("redis.port"))

	options := &redis.Options{
		Addr:     addr,
		Password: cfg.GetString("password"),
		DB:       0,
	}

	client := redis.NewClient(options)

	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		log.Sugar().Fatalf("failed to connect to Redis: %v", err)
	}

	return client
}
