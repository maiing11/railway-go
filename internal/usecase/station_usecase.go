package usecase

import (
	"context"
	"railway-go/internal/constant/model"
	"railway-go/internal/repository"
	"railway-go/internal/utils"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type StationUC interface {
	CreateStation(ctx context.Context, request model.StationRequest) (model.Station, error)
	GetStationByName(ctx context.Context, name string) (model.Station, error)
	GetStation(ctx context.Context, id int64) (model.Station, error)
	GetStationByCode(ctx context.Context, code string) (model.Station, error)
	GetAllStations(ctx context.Context) ([]model.Station, error)
	UpdateStation(ctx context.Context, id int64, request model.StationRequest) error
	DeleteStation(ctx context.Context, id int64) error
}

type StationUsecase struct {
	*UseCase
}

func NewStationUsecase(usecase *UseCase) StationUC {
	return &StationUsecase{
		UseCase: usecase,
	}
}

func (uc *StationUsecase) CreateStation(ctx context.Context, request model.StationRequest) (model.Station, error) {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return model.Station{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to begin transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// validate request
	if err := uc.Validate.Struct(request); err != nil {
		return model.Station{}, utils.WrapError(fiber.StatusBadRequest, uc.Log, utils.Warn, err, "error validating request body")
	}

	stationParams := repository.CreateStationParams{
		Code:        request.Code,
		StationName: request.StationName,
	}

	station, err := tx.CreateStation(ctx, stationParams)
	if err != nil {
		return model.Station{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to create station")
	}

	if err := tx.Commit(ctx); err != nil {
		return model.Station{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to commit transaction")
	}

	response := model.Station{
		ID:          station.ID,
		Code:        station.Code,
		StationName: station.StationName,
		CreatedAt:   station.CreatedAt,
		UpdatedAt:   station.UpdatedAt,
	}

	return response, nil
}

func (uc *StationUsecase) GetStation(ctx context.Context, id int64) (model.Station, error) {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return model.Station{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to begin transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	station, err := tx.GetStation(ctx, id)
	if err != nil {
		return model.Station{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to get station")
	}

	if err := tx.Commit(ctx); err != nil {
		return model.Station{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to commit transaction")
	}

	response := model.Station{
		ID:          station.ID,
		Code:        station.Code,
		StationName: station.StationName,
		CreatedAt:   station.CreatedAt,
		UpdatedAt:   station.UpdatedAt,
	}

	return response, nil
}

func (uc *StationUsecase) GetStationByCode(ctx context.Context, code string) (model.Station, error) {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return model.Station{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to begin transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	station, err := tx.GetStationByCode(ctx, code)
	if err != nil {
		return model.Station{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to get station by code")
	}

	if err := tx.Commit(ctx); err != nil {
		return model.Station{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to commit transaction")
	}

	response := model.Station{
		ID:          station.ID,
		Code:        station.Code,
		StationName: station.StationName,
		CreatedAt:   station.CreatedAt,
		UpdatedAt:   station.UpdatedAt,
	}

	return response, nil

}

func (uc *StationUsecase) GetStationByName(ctx context.Context, name string) (model.Station, error) {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return model.Station{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to begin transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	station, err := tx.GetStationByName(ctx, name)
	if err != nil {
		return model.Station{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to get station by name")
	}

	if err := tx.Commit(ctx); err != nil {
		return model.Station{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to commit transaction")
	}

	response := model.Station{
		ID:          station.ID,
		Code:        station.Code,
		StationName: station.StationName,
		CreatedAt:   station.CreatedAt,
		UpdatedAt:   station.UpdatedAt,
	}

	return response, nil

}

func (uc *StationUsecase) GetAllStations(ctx context.Context) ([]model.Station, error) {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return []model.Station{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to begin transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	stations, err := tx.ListStations(ctx)
	if err != nil {
		return []model.Station{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to get station by code")
	}

	if err := tx.Commit(ctx); err != nil {
		return []model.Station{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to commit transaction")
	}

	response := make([]model.Station, len(stations))

	for i, station := range stations {
		response[i] = model.Station{
			ID:          station.ID,
			Code:        station.Code,
			StationName: station.StationName,
			CreatedAt:   station.CreatedAt,
			UpdatedAt:   station.UpdatedAt,
		}
	}

	return response, nil

}

func (uc *StationUsecase) UpdateStation(ctx context.Context, id int64, request model.StationRequest) error {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to begin error")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	station, err := uc.GetStation(ctx, id)
	if err != nil {
		return utils.WrapError(fiber.StatusNotFound, uc.Log, utils.Warn, err, "unknown id or failed to get station")
	}

	stationParams := repository.UpdateStationParams{
		ID:          station.ID,
		Code:        request.Code,
		StationName: request.StationName,
	}

	if err := tx.UpdateStation(ctx, stationParams); err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to update station")
	}

	if err := tx.Commit(ctx); err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to commit transaction")
	}

	uc.Log.Info("station successfully updated", zap.String("station name: ", request.StationName))
	return nil
}

func (uc *StationUsecase) DeleteStation(ctx context.Context, id int64) error {

	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to begin transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	station, err := uc.GetStation(ctx, id)
	if err != nil {
		return utils.WrapError(fiber.StatusNotFound, uc.Log, utils.Warn, err, "unknown id for delete station")
	}

	if err := tx.DeleteStation(ctx, station.ID); err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to delete station")
	}

	if err := tx.Commit(ctx); err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to commit transaction")
	}

	uc.Log.Info("station successfully deleted", zap.Any("station", station))
	return nil
}
