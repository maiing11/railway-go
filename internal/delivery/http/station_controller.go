package http

import (
	"railway-go/internal/constant/model"
	"railway-go/internal/usecase"
	"railway-go/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type StationControllers interface {
	CreateStation(ctx *fiber.Ctx) error
	GetStation(ctx *fiber.Ctx) error
	GetStationByCode(ctx *fiber.Ctx) error
	GetStationByName(ctx *fiber.Ctx) error
	GetStations(ctx *fiber.Ctx) error
	UpdateStation(ctx *fiber.Ctx) error
	DeleteStation(ctx *fiber.Ctx) error
}

type StationController struct {
	Log     *zap.Logger
	Usecase usecase.StationUC
}

func NewStationController(usecase usecase.StationUC, log *zap.Logger) StationControllers {
	return &StationController{
		Usecase: usecase,
		Log:     log,
	}
}

func (c *StationController) CreateStation(ctx *fiber.Ctx) error {
	request := new(model.StationRequest)

	if err := ctx.BodyParser(request); err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "failed to parse body requrest")
	}

	if request.Code == "" || request.StationName == "" {
		return utils.HandleError(ctx, c.Log, nil, fiber.StatusBadRequest, "all field are required")
	}

	response, err := c.Usecase.CreateStation(ctx.UserContext(), *request)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to create station")
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.BuildSuccessResponse(response, nil))
}

func (c *StationController) GetStation(ctx *fiber.Ctx) error {
	reqId := ctx.Query("id")
	if reqId == "" {
		return utils.HandleError(ctx, c.Log, nil, fiber.StatusBadRequest, "station id is required")
	}

	id, _ := strconv.ParseInt(reqId, 10, 64)

	station, err := c.Usecase.GetStation(ctx.UserContext(), id)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to get station")
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse(station, nil))
}

func (c *StationController) GetStationByCode(ctx *fiber.Ctx) error {
	req := ctx.Params("code")
	if req == "" {
		return utils.HandleError(ctx, c.Log, nil, fiber.StatusBadRequest, "station code is required")
	}

	station, err := c.Usecase.GetStationByCode(ctx.UserContext(), req)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusNotFound, "station not found")
	}

	return ctx.Status(fiber.StatusFound).JSON(model.BuildSuccessResponse(station, nil))
}

func (c *StationController) GetStationByName(ctx *fiber.Ctx) error {
	req := ctx.Query("station_name")
	if req == "" {
		return utils.HandleError(ctx, c.Log, nil, fiber.StatusBadRequest, "station name is required")
	}

	station, err := c.Usecase.GetStationByName(ctx.UserContext(), req)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusNotFound, "station not found")
	}

	return ctx.Status(fiber.StatusFound).JSON(model.BuildSuccessResponse(station, nil))
}

func (c *StationController) GetStations(ctx *fiber.Ctx) error {

	station, err := c.Usecase.GetAllStations(ctx.UserContext())
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusNotFound, "station not found")
	}

	paging := &model.PageMetaData{
		Page: 1,
		Size: 10,
	}
	return ctx.Status(fiber.StatusFound).JSON(model.BuildSuccessResponse(station, paging))
}

func (c *StationController) UpdateStation(ctx *fiber.Ctx) error {
	reqId := ctx.Query("id")
	stationReq := new(model.StationRequest)

	if err := ctx.BodyParser(stationReq); err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "failed to parse body request")
	}

	if reqId == "" {
		return utils.HandleError(ctx, c.Log, nil, fiber.StatusBadRequest, "station id are required")
	}

	stationId, _ := strconv.ParseInt(reqId, 10, 64)

	if err := c.Usecase.UpdateStation(ctx.UserContext(), stationId, *stationReq); err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "error cannot update station")
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse("station updated successfully", nil))
}

func (c *StationController) DeleteStation(ctx *fiber.Ctx) error {
	reqId := ctx.Query("id")
	if reqId == "" {
		return utils.HandleError(ctx, c.Log, nil, fiber.StatusBadRequest, "station id is required")
	}

	id, _ := strconv.ParseInt(reqId, 10, 64)

	if err := c.Usecase.DeleteStation(ctx.UserContext(), id); err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to get station")
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse("station deleted successfully", nil))
}
