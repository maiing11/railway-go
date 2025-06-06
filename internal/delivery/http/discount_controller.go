package http

import (
	"railway-go/internal/constant/model"
	"railway-go/internal/usecase"
	"railway-go/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type DiscountControllers interface {
	CreateDiscount(ctx *fiber.Ctx) error
	GetDiscount(ctx *fiber.Ctx) error
}

type DiscountController struct {
	Log     *zap.Logger
	Usecase usecase.DiscountUC
}

func NewDiscountController(log *zap.Logger, usecase usecase.DiscountUC) DiscountControllers {
	return &DiscountController{
		Log:     log,
		Usecase: usecase,
	}
}

func (c *DiscountController) CreateDiscount(ctx *fiber.Ctx) error {
	req := new(model.DiscountRequest)

	if err := ctx.BodyParser(req); err != nil {
		return utils.HandleError(ctx, c.Log, nil, fiber.StatusBadRequest, "failed to parse body requrest")
	}

	discount, err := c.Usecase.CreateDiscount(ctx.UserContext(), *req)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to create discount")
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.BuildSuccessResponse(discount, nil))
}

func (c *DiscountController) GetDiscount(ctx *fiber.Ctx) error {
	req := ctx.Query("id")

	if req == "" {
		return utils.HandleError(ctx, c.Log, nil, fiber.StatusBadRequest, "id is required")
	}
	reqID, err := uuid.Parse(req)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "invalid discount id")
	}

	discount, err := c.Usecase.GetDiscountByID(ctx.UserContext(), reqID)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to get Discount by id")
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse(discount, nil))
}

func (c *DiscountController) GetDiscountByCode(ctx *fiber.Ctx) error {
	req := ctx.Params("code")

	if req == "" {
		return utils.HandleError(ctx, c.Log, nil, fiber.StatusBadRequest, "code is required")
	}

	discount, err := c.Usecase.GetDiscountByCode(ctx.UserContext(), req)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to get Discount by code")
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse(discount, nil))
}
