package http

import (
	"fmt"
	"railway-go/internal/constant/model"
	"railway-go/internal/usecase"
	"railway-go/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PassengerControllers interface {
	CreatePassenger(ctx *fiber.Ctx) error
	GetPassenger(ctx *fiber.Ctx) error
	UpdatePassenger(ctx *fiber.Ctx) error
	DeletePassenger(ctx *fiber.Ctx) error
}

type PassengerController struct {
	Log     *zap.Logger
	Usecase usecase.PassengerUC
	UserUC  usecase.UserSessionUC
}

func NewPassengerController(uc usecase.PassengerUC, userUC usecase.UserSessionUC, log *zap.Logger) PassengerControllers {
	return &PassengerController{
		Log:     log,
		Usecase: uc,
		UserUC:  userUC,
	}
}

func (c *PassengerController) CreatePassenger(ctx *fiber.Ctx) error {
	request := new(model.PassengerRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "failed to parse request body")
	}

	// validate required fields
	if request.Name == "" || request.IDNumber == "" {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, " fields are required")
	}

	sessionID := ctx.Cookies("session_id")
	var session model.Session
	if sessionID == "" {
		sess := ctx.Locals("session")
		fmt.Printf("session in locals: %v", sess)
		// if sess != nil {
		// 	return utils.HandleError(ctx, c.Log, nil, fiber.StatusUnauthorized, "session_id cookie missing")
		// }
		session, _ = sess.(model.Session)
	} else {
		session, err = c.UserUC.GetSession(ctx.UserContext(), sessionID)
		if err != nil {
			return utils.HandleError(ctx, c.Log, err, fiber.StatusNotFound, "failed to get session")
		}
	}

	if session.Role == "user" || session.Role == "admin" {
		userId, err := c.UserUC.GetUserIDFromSession(ctx.UserContext(), sessionID)
		if err != nil {
			return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to get user id")
		}
		request.UserID = *userId
	}
	c.Log.Info("params ", zap.Any("user_id", request.UserID))

	response, err := c.Usecase.CreatePassenger(ctx.UserContext(), *request)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to create passenger")
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.BuildSuccessResponse(response, nil))
}

func (c *PassengerController) GetPassenger(ctx *fiber.Ctx) error {
	request := ctx.Query("id")
	if request == "" {
		return utils.HandleError(ctx, c.Log, nil, fiber.StatusBadRequest, "passenger id is required")
	}

	response, err := c.Usecase.GetPassenger(ctx.UserContext(), uuid.MustParse(request))
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to get passenger")
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse(response, nil))
}

func (c *PassengerController) UpdatePassenger(ctx *fiber.Ctx) error {
	id := ctx.Query("id")
	if id == "" {
		return utils.HandleError(ctx, c.Log, nil, fiber.StatusBadRequest, "passenger id is required")
	}
	parsed, _ := uuid.Parse(id)
	passengerRequest := new(model.PassengerRequest)
	if err := ctx.BodyParser(passengerRequest); err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "failed to parse request body")
	}

	passengerRequest.ID = parsed

	if err := c.Usecase.UpdatePassenger(ctx.UserContext(), *passengerRequest); err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to update passenger")
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse("Passenger updated successfully", nil))
}

func (c *PassengerController) DeletePassenger(ctx *fiber.Ctx) error {
	request := ctx.Query("id")
	if request == "" {
		return utils.HandleError(ctx, c.Log, nil, fiber.StatusBadRequest, "passenger id is required")
	}

	err := c.Usecase.DeletePassenger(ctx.UserContext(), uuid.MustParse(request))
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to delete passenger")
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse("Passenger deleted successfully", nil))
}
