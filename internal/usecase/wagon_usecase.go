package usecase

import (
	"context"
	"railway-go/internal/constant/model"
	"railway-go/internal/repository"
	"railway-go/internal/utils"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type WagonUC interface {
	CreateWagon(ctx context.Context, wagon model.WagonRequest) (repository.Wagon, error)
	GetWagon(ctx context.Context, wagonID int64) (repository.Wagon, error)
	UpdateWagon(ctx context.Context, wagonID int64, wagon model.WagonRequest) error
	DeleteWagon(ctx context.Context, wagonID int64) error
	GetWagonList(ctx context.Context, trainID int64) ([]repository.Wagon, error)
}

type WagonUsecase struct {
	*UseCase
	TrainUC
}

func NewWagonUsecase(useCase *UseCase) WagonUC {
	return &WagonUsecase{UseCase: useCase}
}

func (uc *WagonUsecase) CreateWagon(ctx context.Context, wagon model.WagonRequest) (repository.Wagon, error) {

	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return repository.Wagon{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to begin transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	train, err := tx.GetTrain(ctx, wagon.TrainID)
	if err != nil {
		return repository.Wagon{}, utils.WrapError(fiber.StatusNotFound, uc.Log, utils.Warn, err, "failed to get train / train id is unknown")
	}
	// validate request
	if err := uc.Validate.Struct(wagon); err != nil {
		uc.Log.Warn("error validating request body", zap.Error(err))
		return repository.Wagon{}, fiber.ErrBadRequest
	}

	wagonParams := repository.CreateWagonParams{
		TrainID:     train.ID,
		WagonNumber: wagon.WagonNumber,
		ClassType:   repository.TipeClass(wagon.ClassType),
		TotalSeats:  wagon.TotalSeats,
	}

	response, err := tx.CreateWagon(ctx, wagonParams)
	if err != nil {
		return repository.Wagon{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to create wagon")
	}

	if err := tx.Commit(ctx); err != nil {
		return repository.Wagon{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to commit transaction")
	}

	return response, nil
}

func (uc *WagonUsecase) GetWagon(ctx context.Context, wagonID int64) (repository.Wagon, error) {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		uc.Log.Warn("failed to begin transaction", zap.Error(err))
		return repository.Wagon{}, fiber.NewError(fiber.StatusInternalServerError, "failed to begin transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	response, err := tx.GetWagon(ctx, wagonID)
	if err != nil {
		uc.Log.Warn("failed to get wagon", zap.Error(err))
		return repository.Wagon{}, fiber.ErrInternalServerError
	}

	if err := tx.Commit(ctx); err != nil {
		uc.Log.Warn("failed to commit transaction", zap.Error(err))
		return repository.Wagon{}, fiber.NewError(fiber.StatusInternalServerError, "failed to commit transaction")
	}

	return response, nil
}

func (uc *WagonUsecase) UpdateWagon(ctx context.Context, wagonID int64, wagon model.WagonRequest) error {
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

	wagonParams := repository.UpdateWagonParams{
		ID:          wagonID,
		WagonNumber: wagon.WagonNumber,
		ClassType:   repository.TipeClass(wagon.ClassType),
		TotalSeats:  wagon.TotalSeats,
	}
	if err := tx.UpdateWagon(ctx, wagonParams); err != nil {
		uc.Log.Warn("failed to update wagon", zap.Error(err))
		return fiber.ErrInternalServerError
	}
	if err := tx.Commit(ctx); err != nil {
		uc.Log.Warn("failed to commit transaction", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "failed to commit transaction")
	}
	return nil
}

func (uc *WagonUsecase) DeleteWagon(ctx context.Context, wagonID int64) error {
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

	wagon, err := uc.GetWagon(ctx, wagonID)
	if err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "id wagon is unknown")
	}

	if err := tx.DeleteWagon(ctx, wagon.ID); err != nil {
		uc.Log.Warn("failed to delete wagon", zap.Error(err))
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit(ctx); err != nil {
		uc.Log.Warn("failed to commit transaction", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "failed to commit transaction")
	}
	return nil
}

func (uc *WagonUsecase) GetWagonList(ctx context.Context, trainID int64) ([]repository.Wagon, error) {
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

	response, err := tx.ListWagons(ctx, trainID)
	if err != nil {
		uc.Log.Warn("failed to get wagon list", zap.Error(err))
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit(ctx); err != nil {
		uc.Log.Warn("failed to commit transaction", zap.Error(err))
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to commit transaction")
	}

	return response, nil
}
