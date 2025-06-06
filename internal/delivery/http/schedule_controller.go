package http

import (
	"railway-go/internal/constant/model"
	"railway-go/internal/usecase"
	"railway-go/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"go.uber.org/zap"
)

type ScheduleControllers interface {
	CreateSchedule(ctx *fiber.Ctx) error
	UpdateSchedule(ctx *fiber.Ctx) error
	GetSchedule(ctx *fiber.Ctx) error
	DeleteSchedule(ctx *fiber.Ctx) error
	SearchSchedules(ctx *fiber.Ctx) error
}

type ScheduleController struct {
	Log     *zap.Logger
	Usecase usecase.ScheduleUC
}

func NewScheduleController(usecase usecase.ScheduleUC, log *zap.Logger) ScheduleControllers {
	return &ScheduleController{
		Log:     log,
		Usecase: usecase,
	}
}

func (c *ScheduleController) CreateSchedule(ctx *fiber.Ctx) error {
	request := new(model.ScheduleRequest)

	err := ctx.BodyParser(request)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "failed to parse request body")
	}

	// validate required fields
	if request.RouteID == 0 || request.DepartureDate.Time.IsZero() {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "All fields are required")
	}

	response, err := c.Usecase.CreateSchedule(ctx.UserContext(), request)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to create schedule")
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.BuildSuccessResponse(response, nil))
}

func (c *ScheduleController) UpdateSchedule(ctx *fiber.Ctx) error {
	id := ctx.QueryInt("id")
	request := new(model.ScheduleRequest)

	err := ctx.BodyParser(request)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "failed to parse request body")
	}

	// convert string to int64

	// validate required fields
	if id == 0 || request.DepartureDate.Time.IsZero() {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "All fields are required")
	}

	err = c.Usecase.UpdateSchedule(ctx.UserContext(), int64(id), *request)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to update schedule")
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse("Schedule updated successfully", nil))
}

func (c *ScheduleController) GetSchedule(ctx *fiber.Ctx) error {
	req := ctx.QueryInt("id")

	// validate required fields
	if req == 0 {
		return utils.HandleError(ctx, c.Log, nil, fiber.StatusBadRequest, "id is required")
	}

	// convert string to int64

	response, err := c.Usecase.GetSchedule(ctx.UserContext(), int64(req))
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to get schedule")
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse(response, nil))
}

func (c *ScheduleController) DeleteSchedule(ctx *fiber.Ctx) error {
	req := ctx.Query("id")

	// validate required fields
	if req == "" {
		return utils.HandleError(ctx, c.Log, nil, fiber.StatusBadRequest, "id is required")
	}

	// convert string to int64
	reqId, _ := strconv.ParseInt(req, 10, 64)

	err := c.Usecase.DeleteSchedule(ctx.UserContext(), reqId)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to delete schedule")
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse("Schedule deleted successfully", nil))
}

func (c *ScheduleController) SearchSchedules(ctx *fiber.Ctx) error {
	request := new(model.SearchScheduleRequest)

	err := ctx.BodyParser(request)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "failed to parse request body")
	}

	response, err := c.Usecase.SearchSchedules(ctx.UserContext(), request)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to search schedules")
	}

	// paging := &model.PageMetaData{
	// 	Page: 1,
	// 	Size: 10,
	// }

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse(response, nil))
}
