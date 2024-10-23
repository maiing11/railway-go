package config

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func newDatabase(v *viper.Viper, log *zap.Logger) (*pgxpool.Pool, error) {
	username := v.GetString("database.username")
	password := v.GetString("database.password")
	host := v.GetString("database.host")
	port := v.GetInt("database.port")
	database := v.GetString("database.name")
	idleConnection := v.GetInt("database.pool.idle")
	maxConnection := v.GetInt("database.pool.max")
	maxLifeTimeConnection := v.GetInt("database.pool.lifetime")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Sugar().Fatalf("Unable to parse database config: %v", err)
		return nil, err
	}

	config.MaxConnIdleTime = time.Duration(idleConnection)
	config.MaxConns = int32(maxConnection)
	config.MaxConnLifetime = time.Duration(maxLifeTimeConnection)

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Sugar().Fatalf("unable to create connection pool: %v", err)
		return nil, err
	}

	return db, nil
}
