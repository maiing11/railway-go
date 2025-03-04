package middleware

import (
	"context"
	"railway-go/internal/constant/model"
	"railway-go/internal/usecase"
	"railway-go/internal/utils/token"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	SessionUsecase usecase.UserSessionUsecase
	TokenMaker     token.Maker
}

func NewAuthMiddleware(sessionUsecase usecase.UserSessionUsecase, tokenMaker token.Maker) *AuthMiddleware {
	return &AuthMiddleware{
		SessionUsecase: sessionUsecase,
		TokenMaker:     tokenMaker,
	}
}

// type contextKey string

// const fiberCtxKey contextKey = "fiberCtx"

func (m *AuthMiddleware) AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(model.BuildErrorResponse("authorization header is missing"))
		}

		// Verify token
		payload, err := m.TokenMaker.VerifyToken(authHeader)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(model.BuildErrorResponse("invalid or expired token"))
		}

		// check if session exists in redis
		session, err := m.SessionUsecase.GetSession(context.Background(), payload.SessionID)
		if err != nil || session.IsBlocked || session.ExpiresAt.Before(time.Now()) {
			return c.Status(fiber.StatusUnauthorized).JSON(model.BuildErrorResponse("session not found or expired"))
		}
		// ctx := context.WithValue(c.Context(), fiberCtxKey, c)
		// Attach user details or session to request context
		c.Locals("user", payload)
		c.Locals("session", session)

		return c.Next()

	}
}

func (m *AuthMiddleware) GuestOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionID := c.Cookies("session_id")
		if sessionID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(model.BuildErrorResponse("guest session is required"))
		}

		session, err := m.SessionUsecase.GetGuestSession(context.Background(), sessionID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(model.BuildErrorResponse("invalid or expired guest session"))
		}

		// attach session to request context
		c.Locals("session", session)

		return c.Next()
	}
}
