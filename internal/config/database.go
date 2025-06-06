package config

import (
	"context"
	"fmt"
	"time"

	// "time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewDatabase(v *viper.Viper, log *zap.Logger) (*pgxpool.Pool, error) {
	v.BindEnv("database.port", "DB_PORT")
	v.BindEnv("database.name", "DB_NAME")
	database := v.GetString("database.name")
	username := v.GetString("database.username")
	port := v.GetInt("database.port")
	password := v.GetString("database.password")
	host := v.GetString("database.host")
	// idleConnection := v.GetInt("database.pool.idle")
	maxConnection := v.GetInt("database.pool.max")
	// maxLifeTimeConnection := v.GetInt("database.pool.lifetime")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", username, password, host, port, database)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Sugar().Fatalf("Unable to parse database config: %v", err)
		return nil, err
	}

	// config.MaxConnIdleTime = time.Duration(idleConnection)
	// if maxConnection < 1 {
	// 	maxConnection = 1
	// }
	config.MaxConns = int32(maxConnection)
	// config.MaxConnLifetime = time.Duration(maxLifeTimeConnection)

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Sugar().Fatalf("unable to create connection pool: %v", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.Ping(ctx)
	if err != nil {
		log.Sugar().Fatalf("database ping failed: %v", err)
		return nil, err
	}

	log.Info("connected to postgresql successfully")
	return db, nil
}
