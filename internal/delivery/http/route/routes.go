package route

import (
	"railway-go/internal/delivery/http"
	"railway-go/internal/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                   *fiber.App
	UserSessionController *http.UserSessionControllers
	AuthMiddleware        *middleware.AuthMiddleware
}

func (c *RouteConfig) Setup() {
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupAuthRoute() {

	c.App.Post("/register", c.UserSessionController.Register)
	c.App.Post("/login", c.UserSessionController.Login)

	c.App.Use("/guest", c.AuthMiddleware.GuestOnly())
	// c.App.Post("/guest/session", c.UserSessionController.)

	c.App.Use("/auth", c.AuthMiddleware.AuthRequired())
	c.App.Post("/logout", c.UserSessionController.Logout)
}
