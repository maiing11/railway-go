package http

import (
	"context"
	"railway-go/internal/constant/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// Register handles the user registration process.
// It parses the request body into a RegisterUserRequest struct,
// validates the required fields, and calls the Usecase Register method.
// If any error occurs during parsing, validation, or registration,
// it returns an appropriate error response.
// On successful registration, it returns a success response.
//
// @Summary Register a new user
// @Tags users
// @Accept json
// @Produce json
// @Param request body model.RegisterUserRequest true "User registration request"
// @Success 201 {object} model.SuccessResponse "User registered successfully"
// @Failure 400 {object} model.ErrorResponse "Invalid request body or missing required fields"
// @Router /register [post]
func (c *Controllers) Register(ctx *fiber.Ctx) error {
	request := new(model.RegisterUserRequest)
	err := ctx.BodyParser(request)

	if err != nil {
		c.Log.Warn("Failed to parse request body", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(model.BuildErrorResponse("invalid request body"))
	}

	// validate required fields
	if request.Name == "" || request.Email == "" || request.PhoneNumber == "" || request.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.BuildErrorResponse("All fields are required"))
	}

	err = c.Usecase.Register(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warn("Failed to register user", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(model.BuildErrorResponse(err.Error()))
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.BuildSuccessResponse("User registered successfully", nil))
}

type contextKey = string

const fiberCtx contextKey = "fiberCtx"

func (c *Controllers) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warn("Failed to parse request body: %+v", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(model.BuildErrorResponse("Invalid request body"))
	}

	ctxWithFiber := context.WithValue(ctx.UserContext(), fiberCtx, ctx)

	token, err := c.Usecase.Login(ctxWithFiber, request)
	if err != nil {
		c.Log.Warn("Failed to login", zap.Error(err))
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.BuildErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse(token, nil))
}

func (c *Controllers) Logout(ctx *fiber.Ctx) error {
	// extract session_id from cookies

	sessionID := ctx.Cookies("session_id")
	if sessionID == "" {
		c.Log.Warn("Failed to retrieve session id")
		return ctx.Status(fiber.StatusBadRequest).JSON(model.BuildErrorResponse("Session ID not found"))
	}

	if err := c.Usecase.Logout(ctx.Context(), sessionID); err != nil {
		c.Log.Warn("failed to logout", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.BuildErrorResponse(err.Error()))
	}

	// clear the session cookie
	ctx.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		SameSite: "Lax",
		Secure:   true,
	})

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse("User logged out successfully", nil))
}
