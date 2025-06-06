package middleware

import (
	"context"
	"fmt"
	"railway-go/internal/constant/model"
	"railway-go/internal/usecase"
	"railway-go/internal/utils/token"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	SessionUsecase usecase.UserSessionUC
	TokenMaker     token.Maker
}

func NewAuthMiddleware(sessionUsecase usecase.UserSessionUC, tokenMaker token.Maker) *AuthMiddleware {
	return &AuthMiddleware{
		SessionUsecase: sessionUsecase,
		TokenMaker:     tokenMaker,
	}
}

// type contextKey string

// const fiberCtxKey contextKey = "fiberCtx"

func (m *AuthMiddleware) AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// check if session ID exists incookies
		sessionID := c.Cookies("session_id")
		if sessionID == "" {
			const fiberCtx string = "fiberCtx"
			session, err := m.SessionUsecase.CreateGuestSession(context.WithValue(c.UserContext(), fiberCtx, c))
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(model.BuildErrorResponse("invalid or expired guest session"))
			}

			// attach guest session to request context
			c.Locals("session", &session)
			return c.Next()
		}

		// check if session exists in redis
		session, err := m.SessionUsecase.GetSession(context.Background(), sessionID)
		if err != nil || session.IsBlocked || session.ExpiresAt.Before(time.Now()) {
			return c.Status(fiber.StatusUnauthorized).JSON(model.BuildErrorResponse("session not found or expired"))
		}

		if session.Role == "user" || session.Role == "admin" || session.Role == "general affairs" {
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

			// ctx := context.WithValue(c.Context(), fiberCtxKey, c)
			// Attach user/admin/ga details or session to request context
			c.Locals("user", payload)

		}

		c.Locals("session", &session)

		// fmt.Printf("session role :%s", session.Role)

		return c.Next()

	}
}

func (m *AuthMiddleware) GuestOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionID := c.Cookies("session_id")
		if sessionID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(model.BuildErrorResponse("guest session is required"))
		}

		session, err := m.SessionUsecase.CreateGuestSession(context.Background())
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(model.BuildErrorResponse("invalid or expired guest session"))
		}

		// attach session to request context
		c.Locals("session", session)

		return c.Next()
	}
}

func (m *AuthMiddleware) AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// sessionID := c.Cookies("session_id")
		session, ok := c.Locals("session").(*model.Session)
		if !ok || session.Role != "admin" {
			fmt.Println("Session missing or wrong type:", c.Locals("session"))
			return c.Status(fiber.StatusForbidden).JSON(model.BuildErrorResponse("admin access required"))
		}
		// fmt.Println("Session role:", session.Role)
		return c.Next()
	}
}

func (m *AuthMiddleware) GeneralAffairs() fiber.Handler {
	return func(c *fiber.Ctx) error {
		session, ok := c.Locals("session").(*model.Session)
		if !ok || session == nil {
			return c.Status(fiber.StatusForbidden).JSON(model.BuildErrorResponse("general affairs or admin access required"))
		}
		if session.Role != "general affairs" && session.Role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(model.BuildErrorResponse("general affairs or admin access required"))
		}
		return c.Next()
	}
}
