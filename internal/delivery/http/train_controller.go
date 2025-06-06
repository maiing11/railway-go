package http

import (
	"railway-go/internal/constant/model"
	"railway-go/internal/usecase"
	"railway-go/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type TrainControllers interface {
	CreateTrain(ctx *fiber.Ctx) error
	GetTrain(ctx *fiber.Ctx) error
	GetTrains(ctx *fiber.Ctx) error
	UpdateTrain(ctx *fiber.Ctx) error
	DeleteTrain(ctx *fiber.Ctx) error
}

type TrainController struct {
	Log     *zap.Logger
	Usecase usecase.TrainUC
}

func NewTrainController(usecase usecase.TrainUC, log *zap.Logger) TrainControllers {
	return &TrainController{
		Log:     log,
		Usecase: usecase,
	}
}

func (c *TrainController) CreateTrain(ctx *fiber.Ctx) error {
	request := new(model.TrainRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "failed to parse request body")
	}

	// validate required fields
	if request.TrainName == "" || request.Capacity == 0 {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "All fields are required")
	}

	response, err := c.Usecase.CreateTrain(ctx.UserContext(), *request)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to create train")
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.BuildSuccessResponse(response, nil))
}

func (c *TrainController) GetTrain(ctx *fiber.Ctx) error {
	request := ctx.Query("id")
	if request == "" {
		return utils.HandleError(ctx, c.Log, nil, fiber.StatusBadRequest, "train id is required")
	}

	// convert string to int64
	trainID, err := strconv.ParseInt(request, 10, 64)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "invalid train id")
	}

	response, err := c.Usecase.GetTrain(ctx.UserContext(), trainID)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to get train")
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse(response, nil))
}

func (c *TrainController) GetTrains(ctx *fiber.Ctx) error {
	response, err := c.Usecase.GetAllTrains(ctx.UserContext())
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to get trains")
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse(response, nil))
}

func (c *TrainController) UpdateTrain(ctx *fiber.Ctx) error {
	request := ctx.Query("id")
	trainRequest := new(model.TrainRequest)

	err := ctx.BodyParser(trainRequest)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "failed to parse request body")
	}

	// convert string to int64
	trainID, _ := strconv.ParseInt(request, 10, 64)

	// validate required fields
	if trainID == 0 || trainRequest.TrainName == "" || trainRequest.Capacity == 0 {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "All fields are required")
	}

	err = c.Usecase.UpdateTrain(ctx.UserContext(), trainID, *trainRequest)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to update train")
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse("success updating train", nil))
}

func (c *TrainController) DeleteTrain(ctx *fiber.Ctx) error {
	request := ctx.Query("id")

	// validate required fields
	if request == "" {
		return utils.HandleError(ctx, c.Log, nil, fiber.StatusBadRequest, "id is required")
	}

	// convert string to int64
	trainID, _ := strconv.ParseInt(request, 10, 64)

	err := c.Usecase.DeleteTrain(ctx.UserContext(), trainID)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to delete train")
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse("success deleting train", nil))
}
