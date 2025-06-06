package http

import (
	"railway-go/internal/constant/model"
	"railway-go/internal/usecase"
	"railway-go/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type RouteControllers interface {
	CreateRoute(ctx *fiber.Ctx) error
	GetRoute(ctx *fiber.Ctx) error
	GetRoutes(ctx *fiber.Ctx) error
	UpdateRoute(ctx *fiber.Ctx) error
	DeleteRoute(ctx *fiber.Ctx) error
}

type RouteController struct {
	Log     *zap.Logger
	Usecase usecase.RouteUC
}

func NewRouteController(usecase usecase.RouteUC, log *zap.Logger) RouteControllers {
	return &RouteController{
		Log:     log,
		Usecase: usecase,
	}
}

func (c *RouteController) CreateRoute(ctx *fiber.Ctx) error {
	request := new(model.RouteRequest)

	err := ctx.BodyParser(request)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "failed to parse request body")
	}

	// validate required fields
	if request.SourceStation == "" || request.DestinationStation == "" || request.TravelTime == 0 {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "All fields are required")
	}

	response, err := c.Usecase.CreateRoute(ctx.UserContext(), *request)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to create route")
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.BuildSuccessResponse(response, nil))
}

func (c *RouteController) GetRoute(ctx *fiber.Ctx) error {
	request := ctx.Query("id")
	if request == "" {
		return utils.HandleError(ctx, c.Log, nil, fiber.StatusBadRequest, "route id is required")
	}

	// convert string to int64
	routeID, _ := strconv.ParseInt(request, 10, 64)

	response, err := c.Usecase.GetRoute(ctx.UserContext(), routeID)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to get route")
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse(response, nil))
}

func (c *RouteController) GetRoutes(ctx *fiber.Ctx) error {

	response, err := c.Usecase.GetRoutes(ctx.UserContext())
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to get routes")
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse(response, nil))
}

func (c *RouteController) UpdateRoute(ctx *fiber.Ctx) error {
	request := ctx.Query("id")
	routeRequest := new(model.RouteRequest)

	err := ctx.BodyParser(routeRequest)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "failed to parse request body")
	}

	// convert string to int64
	routeID, _ := strconv.ParseInt(request, 10, 64)

	// validate required fields
	if routeID == 0 || routeRequest.SourceStation == "" || routeRequest.DestinationStation == "" || routeRequest.TravelTime == 0 {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusBadRequest, "All fields are required")
	}

	err = c.Usecase.UpdateRoute(ctx.UserContext(), routeID, *routeRequest)
	if err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to update route")
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse("Route updated successfully", nil))
}

func (c *RouteController) DeleteRoute(ctx *fiber.Ctx) error {
	request := ctx.Query("id")
	if request == "" {
		return utils.HandleError(ctx, c.Log, nil, fiber.StatusBadRequest, "route id is required")
	}

	// convert string to int64
	routeID, _ := strconv.ParseInt(request, 10, 64)

	if err := c.Usecase.DeleteRoute(ctx.UserContext(), routeID); err != nil {
		return utils.HandleError(ctx, c.Log, err, fiber.StatusInternalServerError, "failed to delete route")
	}

	return ctx.Status(fiber.StatusOK).JSON(model.BuildSuccessResponse("Route deleted successfully", nil))
}
