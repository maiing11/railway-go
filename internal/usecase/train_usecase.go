package usecase

import (
	"context"
	"railway-go/internal/constant/model"
	"railway-go/internal/repository"
	"railway-go/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type TrainUC interface {
	GetTrain(ctx context.Context, id int64) (model.Train, error)
	CreateTrain(ctx context.Context, request model.TrainRequest) (model.Train, error)
	UpdateTrain(ctx context.Context, id int64, request model.TrainRequest) error
	DeleteTrain(ctx context.Context, id int64) error
	GetAllTrains(ctx context.Context) ([]model.Train, error)
}

type TrainUsecase struct {
	*UseCase
}

func NewTrainUsecase(useCase *UseCase) TrainUC {
	return &TrainUsecase{UseCase: useCase}
}

func (uc *TrainUsecase) GetTrain(ctx context.Context, id int64) (model.Train, error) {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		uc.Log.Warn("failed to begin transaction", zap.Error(err))
		return model.Train{}, fiber.NewError(fiber.StatusInternalServerError, "failed to begin transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	train, err := tx.GetTrain(ctx, id)
	if err != nil {
		uc.Log.Warn("failed to get train", zap.Error(err))
		return model.Train{}, fiber.NewError(fiber.StatusInternalServerError, "failed to get train")
	}

	//commit transaction
	if err := tx.Commit(ctx); err != nil {
		uc.Log.Warn("failed to commit transaction", zap.Error(err))
		return model.Train{}, fiber.NewError(fiber.StatusInternalServerError, "failed to commit transaction")
	}

	response := model.Train{
		ID:        train.ID,
		TrainName: train.Name,
		Capacity:  train.Capacity,
	}
	return response, nil
}

func (uc *TrainUsecase) CreateTrain(ctx context.Context, request model.TrainRequest) (model.Train, error) {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		uc.Log.Warn("failed to begin transaction", zap.Error(err))
		return model.Train{}, fiber.NewError(fiber.StatusInternalServerError, "failed to begin transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// validate request
	if err := uc.Validate.Struct(request); err != nil {
		uc.Log.Warn("validation failed", zap.Error(err))
		return model.Train{}, fiber.NewError(fiber.StatusBadRequest, "validation failed")
	}

	trainParams := repository.CreateTrainParams{
		Name:     request.TrainName,
		Capacity: request.Capacity,
	}
	// create train
	train, err := tx.CreateTrain(ctx, trainParams)
	if err != nil {
		uc.Log.Warn("failed to create train", zap.Error(err))
		return model.Train{}, fiber.NewError(fiber.StatusInternalServerError, "failed to create train")
	}

	// commit transaction
	if err := tx.Commit(ctx); err != nil {
		uc.Log.Warn("failed to commit transaction", zap.Error(err))
		return model.Train{}, fiber.NewError(fiber.StatusInternalServerError, "failed to commit transaction")
	}

	response := model.Train{
		ID:        train.ID,
		TrainName: train.Name,
		Capacity:  train.Capacity,
		CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
	}
	return response, nil
}

func (uc *TrainUsecase) UpdateTrain(ctx context.Context, id int64, request model.TrainRequest) error {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		uc.Log.Warn("failed to begin transaction", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "failed to begin transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// validate request
	if err := uc.Validate.Struct(request); err != nil {
		uc.Log.Warn("validation failed", zap.Error(err))
		return fiber.NewError(fiber.StatusBadRequest, "validation failed")
	}

	trainParams := repository.UpdateTrainParams{
		ID:       id,
		Name:     request.TrainName,
		Capacity: request.Capacity,
	}

	err = tx.UpdateTrain(ctx, trainParams)
	if err != nil {
		uc.Log.Warn("failed to update train", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "failed to update train")
	}

	if err := tx.Commit(ctx); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to commit transaction")
	}

	uc.Log.Info("train updated successfully", zap.Int64("id", id))
	return nil
}

func (uc *TrainUsecase) DeleteTrain(ctx context.Context, id int64) error {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		uc.Log.Warn("failed to begin transaction", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "failed to begin transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	train, err := uc.GetTrain(ctx, id)
	if err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "unknown id for delete train")
	}

	if err = tx.DeleteTrain(ctx, train.ID); err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to delete train")
	}

	if err := tx.Commit(ctx); err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to commit transaction")
	}

	uc.Log.Info("train deleted successfully", zap.Any("id: ", id))
	return nil
}

func (uc *TrainUsecase) GetAllTrains(ctx context.Context) ([]model.Train, error) {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		uc.Log.Warn("failed to begin transaction", zap.Error(err))
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to begin transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	trains, err := tx.ListTrains(ctx)
	if err != nil {
		uc.Log.Warn("failed to get all trains", zap.Error(err))
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to get all trains")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to commit transaction")
	}

	response := make([]model.Train, len(trains))
	for i, train := range trains {
		response[i] = model.Train{
			ID:        train.ID,
			TrainName: train.Name,
			Capacity:  train.Capacity,
			CreatedAt: train.CreatedAt,
			UpdatedAt: train.UpdatedAt,
		}
	}

	return response, nil
}
