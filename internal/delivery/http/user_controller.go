package http

import (
	"context"
	"railway-go/internal/constant/model"
	"railway-go/internal/usecase"
	"railway-go/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserControllers interface {
	Register(ctx *fiber.Ctx) error
	RegisterAdmin(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	Logout(ctx *fiber.Ctx) error
	RenewAccessToken(ctx *fiber.Ctx) error
	GetUserIDFromSession(ctx *fiber.Ctx) (*uuid.UUID, error)
}
type UserController struct {
	Log     *zap.Logger
	Usecase usecase.UserSessionUC
}

func NewUserSessionController(usecase usecase.UserSessionUC, log *zap.Logger) UserControllers {
	return &UserController{
		Usecase: usecase,
		Log:     log,
	}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	request := new(model.RegisterUserRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "failed to parse request body")
	}

	// validate required fields
	if request.Name == "" || request.Email == "" || request.Password == "" {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "All fields are required")
	}

	err = c.Usecase.Register(ctx.UserContext(), request)
	if err != nil {

		return ctx.Status(fiber.StatusBadRequest).JSON(model.BuildErrorResponse(err.Error()))
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.BuildSuccessResponse("User registered successfully", nil))
}

// @Summary      Register a new admin
// @Description  Registers a new admin with the provided name, email, and password.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body  body      model.RegisterAdminRequest  true  "Admin registration data"
// @Success      201   {object}  model.SuccessResponse       "Admin registered successfully"
// @Failure      400   {object}  model.ErrorResponse         "Invalid input or registration error"
// @Router       /register-admin [post]
func (c *UserController) RegisterAdmin(ctx *fiber.Ctx) error {
	request := new(model.RegisterAdminRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warn("Failed to parse request body", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(model.BuildErrorResponse("invalid request body"))
	}

	// validate required fields
	if request.Name == "" || request.Email == "" || request.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.BuildErrorResponse("All fields are required"))
	}

	err = c.Usecase.RegisterAdmin(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warn("Failed to register", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(model.BuildErrorResponse(err.Error()))
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.BuildSuccessResponse("registered successfully", nil))
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.Warn("Failed to parse request body: %+v", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(model.BuildErrorResponse("Invalid request body"))
	}

	const fiberCtx string = "fiberCtx"
	ctxWithFiber := context.WithValue(ctx.UserContext(), fiberCtx, ctx)

	// Pass the fiber.Ctx directly to the Usecase without using context.WithValue
	token, err := c.Usecase.Login(ctxWithFiber, request)
	if err != nil {
		c.Log.Warn("Failed to login", zap.Error(err))
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.BuildErrorResponse(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse(token, nil))
}

//

func (c *UserController) Logout(ctx *fiber.Ctx) error {
	// extract session_id from cookies

	sessionID := ctx.Cookies("session_id")
	if sessionID == "" {
		c.Log.Warn("Failed to retrieve session id")
		return ctx.Status(fiber.StatusBadRequest).JSON(model.BuildErrorResponse("Session ID not found"))
	}

	if err := c.Usecase.Logout(ctx.UserContext(), sessionID); err != nil {
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

func (c *UserController) RenewAccessToken(ctx *fiber.Ctx) error {
	// extract session_id from cookies
	sessionID := ctx.Cookies("session_id")
	if sessionID == "" {
		return utils.HandleError(ctx, c.Log, nil, fiber.StatusBadRequest, "session id is required")
	}

	// renew access token
	token, err := c.Usecase.RenewAccessToken(ctx.UserContext(), sessionID)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusUnauthorized, "failed to renew access token")
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse(token, nil))
}

func (c *UserController) GetUserIDFromSession(ctx *fiber.Ctx) (*uuid.UUID, error) {
	sessionID := ctx.Cookies("session_id")
	if sessionID == "" {
		return nil, utils.HandleError(ctx, c.Log, nil, fiber.StatusBadRequest, "session id is required")
	}

	userID, err := c.Usecase.GetUserIDFromSession(ctx.UserContext(), sessionID)
	if err != nil {
		return nil, utils.HandleError(ctx, c.Log, err, fiber.StatusUnauthorized, "failed to get user id from session")
	}

	return userID, nil
}
