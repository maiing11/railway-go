package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func NewFiber(v *viper.Viper) *fiber.App {
	var app = fiber.New(fiber.Config{
		AppName:      v.GetString("app.name"),
		ErrorHandler: NewErrorHandler(),
		Prefork:      v.GetBool("web.prefork"),
	})
	return app
}

func NewErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		return c.Status(code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
}
