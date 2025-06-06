package http

import (
	"railway-go/internal/constant/model"
	"railway-go/internal/usecase"
	"railway-go/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type SeatControllers interface {
	CreateSeat(ctx *fiber.Ctx) error
	GetSeat(ctx *fiber.Ctx) error
	GetSeats(ctx *fiber.Ctx) error
	UpdateSeat(ctx *fiber.Ctx) error
	DeleteSeat(ctx *fiber.Ctx) error
}

type SeatController struct {
	Log     *zap.Logger
	Usecase usecase.SeatUC
}

func NewSeatController(log *zap.Logger, uc usecase.SeatUC) SeatControllers {
	return &SeatController{
		Log:     log,
		Usecase: uc,
	}
}

func (c *SeatController) CreateSeat(ctx *fiber.Ctx) error {
	// Implement the logic to create a seat
	request := new(model.SeatRequest)

	if err := ctx.BodyParser(request); err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "failed to parse request body")
	}

	seat, err := c.Usecase.CreateSeat(ctx.UserContext(), request)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to create seat")
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.BuildSuccessResponse(seat, nil))
}

func (c *SeatController) GetSeat(ctx *fiber.Ctx) error {
	// Implement the logic to get a seat by ID
	req := ctx.Query("id")
	if req == "" {
		return utils.HandleError(ctx, c.Log, nil, fiber.StatusBadRequest, "seat id is required")
	}

	// convert string to int64
	seatID, _ := strconv.ParseInt(req, 10, 64)

	seat, err := c.Usecase.GetSeat(ctx.UserContext(), seatID)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to get seat")
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse(seat, nil))
}

func (c *SeatController) GetSeats(ctx *fiber.Ctx) error {
	// Implement the logic to get all seats
	req := ctx.QueryInt("wagon_id")

	seats, err := c.Usecase.GetAllSeat(ctx.UserContext(), int64(req))
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to get seats")
	}

	paging := &model.PageMetaData{
		Page: 1,
		Size: 10,
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse(seats, paging))
}

func (c *SeatController) UpdateSeat(ctx *fiber.Ctx) error {
	// Implement the logic to update a seat
	request := new(model.SeatRequest)

	if err := ctx.BodyParser(request); err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "failed to parse request body")
	}

	seatID := ctx.Params("id")
	if seatID == "" {
		return utils.HandleError(ctx, c.Log, nil, fiber.StatusBadRequest, "seat id is required")
	}

	// convert string to int64
	id, _ := strconv.ParseInt(seatID, 10, 64)

	if err := c.Usecase.UpdateSeat(ctx.UserContext(), id, request); err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to update seat")
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse("successfully updating seat", nil))
}

func (c *SeatController) DeleteSeat(ctx *fiber.Ctx) error {
	// Implement the logic to delete a seat
	request := ctx.Params("id")
	if request == "" {
		return utils.HandleError(ctx, c.Log, nil, fiber.StatusBadRequest, "seat id is required")
	}

	// convert string to int64
	seatID, _ := strconv.ParseInt(request, 10, 64)

	if err := c.Usecase.DeleteSeat(ctx.UserContext(), seatID); err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to delete seat")
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse("successfully deleting seat", nil))
}
