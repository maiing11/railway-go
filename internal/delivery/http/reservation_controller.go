package http

import (
	"math"
	"railway-go/internal/constant/model"
	"railway-go/internal/usecase"
	"railway-go/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ReservationControllers interface {
	CreateReservation(ctx *fiber.Ctx) error
	GetDetailReservation(ctx *fiber.Ctx) error
	CancelReservation(ctx *fiber.Ctx) error
	DeleteReservation(ctx *fiber.Ctx) error
	GetAllReservations(ctx *fiber.Ctx) error
}

type ReservationController struct {
	Log     *zap.Logger
	Usecase usecase.ReservationUC
	UserUC  usecase.UserSessionUC
}

func NewReservationController(usecase usecase.ReservationUC, log *zap.Logger, userUC usecase.UserSessionUC) ReservationControllers {
	return &ReservationController{
		Log:     log,
		Usecase: usecase,
		UserUC:  userUC,
	}
}

func (c *ReservationController) CreateReservation(ctx *fiber.Ctx) error {
	request := new(model.ReservationRequest)

	err := ctx.BodyParser(request)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "failed to parse request body")
	}

	// validate required fields
	if request.ScheduleID == 0 || request.WagonID == 0 || request.Seat_id == 0 {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "All fields are required")
	}

	// get user id from session
	sessionID := ctx.Cookies("session_id")
	time.Sleep(100 * time.Millisecond)
	if sessionID == "" {
		return utils.HandleError(ctx, c.Log, nil, fiber.StatusNotFound, "session id is required")
	}

	session, err := c.UserUC.GetSession(ctx.UserContext(), sessionID)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusNotFound, "failed to get session")
	}

	if session.Role == "user" || session.Role == "admin" {
		userID, err := c.UserUC.GetUserIDFromSession(ctx.UserContext(), session.ID)
		if err != nil {
			return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to get user id from session")
		}
		request.UserId = *userID
	}

	// c.Log.Info("request ", zap.Any("request at controller", request))

	response, err := c.Usecase.CreateReservation(ctx.UserContext(), *request)

	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.BuildSuccessResponse(response, nil))
}

func (c *ReservationController) GetDetailReservation(ctx *fiber.Ctx) error {
	id := ctx.Query("id")

	if id == "" {
		return utils.HandleError(ctx, c.Log, nil, fiber.StatusBadRequest, "id is required")
	}

	reservationID, err := uuid.Parse(id)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "invalid reservation id")
	}

	response, err := c.Usecase.GetDetailReservation(ctx.UserContext(), reservationID)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to get reservation detail")
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse(response, nil))
}

func (c *ReservationController) CancelReservation(ctx *fiber.Ctx) error {
	id := ctx.Query("id")

	if id == "" {
		return utils.HandleError(ctx, c.Log, nil, fiber.StatusBadRequest, "id is required")
	}

	reservationID, err := uuid.Parse(id)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "invalid reservation id")
	}

	err = c.Usecase.CancelReservation(ctx.UserContext(), reservationID)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to cancel reservation")
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (c *ReservationController) DeleteReservation(ctx *fiber.Ctx) error {
	id := ctx.Query("id")

	if id == "" {
		return utils.HandleError(ctx, c.Log, nil, fiber.StatusBadRequest, "id is required")
	}

	reservationID, err := uuid.Parse(id)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "invalid reservation id")
	}

	err = c.Usecase.DeleteReservation(ctx.UserContext(), reservationID)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse("successful delete reservation", nil))
}

func (c *ReservationController) GetAllReservations(ctx *fiber.Ctx) error {

	response, total, err := c.Usecase.GetAllReservations(ctx.UserContext())
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to get all reservations")
	}

	paging := &model.PageMetaData{
		Page:      1,
		Size:      10,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(10))),
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse(response, paging))
}
