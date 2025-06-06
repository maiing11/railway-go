package usecase

import (
	"context"
	"railway-go/internal/constant/model"
	"railway-go/internal/repository"
	"railway-go/internal/utils"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type RouteUC interface {
	GetRoute(ctx context.Context, id int64) (model.Route, error)
	GetRoutes(ctx context.Context) ([]model.Route, error)
	CreateRoute(ctx context.Context, request model.RouteRequest) (model.Route, error)
	UpdateRoute(ctx context.Context, id int64, request model.RouteRequest) error
	DeleteRoute(ctx context.Context, id int64) error
}

type RouteUsecase struct {
	*UseCase
}

func NewRouteUsecase(useCase *UseCase) RouteUC {
	return &RouteUsecase{
		UseCase: useCase,
	}
}

func (uc *RouteUsecase) GetRoute(ctx context.Context, id int64) (model.Route, error) {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return model.Route{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to begin transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	route, err := tx.GetRoute(ctx, id)
	if err != nil {
		return model.Route{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to get route")
	}

	if err := tx.Commit(ctx); err != nil {
		return model.Route{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to commit transaction")
	}

	response := model.Route{
		ID:                 route.ID,
		SourceStation:      route.SourceStation,
		DestinationStation: route.DestinationStation,
		TravelTime:         route.TravelTime,
		CreatedAt:          route.CreatedAt,
		UpdatedAt:          route.UpdatedAt,
	}
	return response, nil
}

func (uc *RouteUsecase) GetRoutes(ctx context.Context) ([]model.Route, error) {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return nil, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to begin transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	routes, err := tx.ListRoute(ctx)
	if err != nil {
		return nil, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to get routes")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to commit transaction")
	}

	response := make([]model.Route, len(routes))
	for i, r := range routes {
		response[i] = model.Route{
			ID:                 r.ID,
			SourceStation:      r.SourceStation,
			DestinationStation: r.DestinationStation,
			TravelTime:         r.TravelTime,
			CreatedAt:          r.CreatedAt,
			UpdatedAt:          r.UpdatedAt,
		}
	}
	return response, nil
}

func (uc *RouteUsecase) CreateRoute(ctx context.Context, request model.RouteRequest) (model.Route, error) {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return model.Route{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to begin transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	route := repository.CreateRouteParams{
		SourceStation:      request.SourceStation,
		DestinationStation: request.DestinationStation,
		TravelTime:         request.TravelTime,
	}

	newRoute, err := tx.CreateRoute(ctx, route)
	if err != nil {
		return model.Route{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to create route")
	}

	if err := tx.Commit(ctx); err != nil {
		return model.Route{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to commit transaction")
	}

	response := model.Route{
		ID:                 newRoute.ID,
		SourceStation:      newRoute.SourceStation,
		DestinationStation: newRoute.DestinationStation,
		TravelTime:         newRoute.TravelTime,
		CreatedAt:          newRoute.CreatedAt,
		UpdatedAt:          newRoute.UpdatedAt,
	}

	return response, nil
}

func (uc *RouteUsecase) UpdateRoute(ctx context.Context, id int64, request model.RouteRequest) error {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to begin transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	r, err := tx.GetRoute(ctx, id)
	if err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to get route")
	}

	route := repository.UpdateRouteParams{
		ID:                 r.ID,
		SourceStation:      request.SourceStation,
		DestinationStation: request.DestinationStation,
		TravelTime:         request.TravelTime,
	}

	if err := tx.UpdateRoute(ctx, route); err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to update route")
	}

	if err := tx.Commit(ctx); err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to commit transaction")
	}

	uc.Log.Info("route updated successfully", zap.Int64("id", r.ID))
	return nil
}

func (uc *RouteUsecase) DeleteRoute(ctx context.Context, id int64) error {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to begin transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	r, err := tx.GetRoute(ctx, id)
	if err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to get route")
	}

	err = tx.DeleteRoute(ctx, r.ID)
	if err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to delete route")
	}

	if err := tx.Commit(ctx); err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to commit transaction")
	}

	uc.Log.Info("route deleted successfully", zap.Int64("id", r.ID))
	return nil
}
