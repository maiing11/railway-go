package http

import (
	"railway-go/internal/constant/model"
	"railway-go/internal/usecase"
	"railway-go/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type WagonControllers interface {
	CreateWagon(ctx *fiber.Ctx) error
	GetWagon(ctx *fiber.Ctx) error
	GetWagons(ctx *fiber.Ctx) error
	UpdateWagon(ctx *fiber.Ctx) error
	DeleteWagon(ctx *fiber.Ctx) error
}

type WagonController struct {
	Log     *zap.Logger
	Usecase usecase.WagonUC
}

func NewWagonController(log *zap.Logger, uc usecase.WagonUC) WagonControllers {
	return &WagonController{
		Log:     log,
		Usecase: uc,
	}
}

func (c *WagonController) CreateWagon(ctx *fiber.Ctx) error {
	request := new(model.WagonRequest)

	if err := ctx.BodyParser(request); err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "failed to parse body request")
	}

	wagon, err := c.Usecase.CreateWagon(ctx.UserContext(), *request)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to create wagon")
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.BuildSuccessResponse(wagon, nil))
}

func (c *WagonController) GetWagon(ctx *fiber.Ctx) error {
	id := ctx.Query("id")

	if id == "" {
		return utils.HandleError(ctx, c.Log, nil, fiber.StatusBadRequest, "wagon id is required")
	}

	wagonId, _ := strconv.ParseInt(id, 10, 64)

	wagon, err := c.Usecase.GetWagon(ctx.UserContext(), wagonId)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to get wagon")
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse(wagon, nil))
}

func (c *WagonController) GetWagons(ctx *fiber.Ctx) error {
	req := ctx.Query("train_id")

	if req == "" {
		return utils.HandleError(ctx, c.Log, nil, fiber.StatusBadRequest, "train id is required")
	}

	trainId, _ := strconv.ParseInt(req, 10, 64)

	wagons, err := c.Usecase.GetWagonList(ctx.UserContext(), trainId)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to get wagons")
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse(wagons, nil))
}

func (c *WagonController) UpdateWagon(ctx *fiber.Ctx) error {
	req := new(model.WagonRequest)

	if err := ctx.BodyParser(req); err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "failed to parse body request")
	}

	wagonId := ctx.Query("id")
	if wagonId == "" {
		return utils.HandleError(ctx, c.Log, nil, fiber.StatusBadRequest, "wagon id is required")
	}

	wagonID, _ := strconv.ParseInt(wagonId, 10, 64)

	if err := c.Usecase.UpdateWagon(ctx.UserContext(), wagonID, *req); err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to update wagon")
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse("successfully updating wagon", nil))
}

func (c *WagonController) DeleteWagon(ctx *fiber.Ctx) error {
	req := ctx.Query("id")

	if req == "" {
		return utils.HandleError(ctx, c.Log, nil, fiber.StatusBadRequest, "wagon id is required")
	}

	wagonId, _ := strconv.ParseInt(req, 10, 64)

	if err := c.Usecase.DeleteWagon(ctx.UserContext(), wagonId); err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to Delete wagon")
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse("wagon successfully deleted", nil))
}
