package utils

import (
	"fmt"
	"railway-go/internal/constant/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func ToUUID(pg pgtype.UUID) (uuid.UUID, error) {
	if !pg.Valid {
		return uuid.Nil, nil
	}

	u := uuid.UUID(pg.Bytes)
	return u, nil
}

func ToPgUUID(u uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: u,
		Valid: true,
	}
}

type LogLevel string

const (
	Info  LogLevel = "info"
	Debug LogLevel = "debug"
	Error LogLevel = "error"
	Warn  LogLevel = "warn"
	Fatal LogLevel = "fatal"
	Panic LogLevel = "panic"
	Trace LogLevel = "trace"
)

func WrapError(code int, log *zap.Logger, level LogLevel, err error, msg string) *fiber.Error {
	if err == nil {
		return nil
	}

	wrapped := fiber.NewError(code, fmt.Sprintf("%s: %s", msg, err.Error()))

	switch level {
	case Info:
		log.Info(msg, zap.Error(err))
	case Debug:
		log.Debug(msg, zap.Error(err))
	case Error:
		log.Error(msg, zap.Error(err))
	case Warn:
		log.Warn(msg, zap.Error(err))
	case Fatal:
		log.Fatal(msg, zap.Error(err))
	case Panic:
		log.Panic(msg, zap.Error(err))
	case Trace:
		log.Debug(msg, zap.Error(err))
	default:
		log.Error("unknown log level", zap.String("level", string(level)))
		log.Error(msg, zap.Error(err))
	}

	return wrapped
}

func HandleError(ctx *fiber.Ctx, log *zap.Logger, err error, statusCode int, msg string) error {
	log.Warn(msg, zap.Error(err))
	return ctx.Status(statusCode).JSON(model.BuildErrorResponse(msg))
}
