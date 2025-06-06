package http

import (
	"railway-go/internal/constant/model"
	"railway-go/internal/usecase"
	"railway-go/internal/utils"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type PaymentControllers interface {
	MockPaymentWebhook(ctx *fiber.Ctx) error
}

type PaymentController struct {
	Log     *zap.Logger
	Usecase usecase.PaymentUC
}

func NewPaymentController(usecase usecase.PaymentUC, log *zap.Logger) PaymentControllers {
	return &PaymentController{
		Usecase: usecase,
		Log:     log,
	}
}

func (c *PaymentController) MockPaymentWebhook(ctx *fiber.Ctx) error {
	req := new(model.PaymentRequest)

	err := ctx.BodyParser(req)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "failed to parse request body")
	}

	// validate required fields
	if req.ReservationID.String() == "" || req.PaymentMethod == "" || req.Amount == 0 {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "All fields are required")
	}

	response, err := c.Usecase.ProcessMockPayment(ctx.UserContext(), *req)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to process payment")
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse(response, nil))
}
