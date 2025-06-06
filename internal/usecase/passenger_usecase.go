package usecase

import (
	"context"
	"railway-go/internal/constant/model"
	"railway-go/internal/repository"
	"railway-go/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type PassengerUC interface {
	GetPassenger(ctx context.Context, id uuid.UUID) (model.Passenger, error)
	CreatePassenger(ctx context.Context, request model.PassengerRequest) (model.Passenger, error)
	UpdatePassenger(ctx context.Context, request model.PassengerRequest) error
	DeletePassenger(ctx context.Context, id uuid.UUID) error
	GetPassengerByUserID(ctx context.Context, userID uuid.UUID) (model.Passenger, error)
}

type PassengerUsecase struct {
	*UseCase
}

func NewPassengerUsecase(useCase *UseCase) PassengerUC {
	return &PassengerUsecase{UseCase: useCase}
}

func (uc *PassengerUsecase) GetPassenger(ctx context.Context, id uuid.UUID) (model.Passenger, error) {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return model.Passenger{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to begin transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	passenger, err := tx.GetPassenger(ctx, id)
	if err != nil {
		return model.Passenger{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to get passenger")
	}

	if err := tx.Commit(ctx); err != nil {
		return model.Passenger{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to commit transaction")
	}

	response := model.Passenger{
		ID:       passenger.ID,
		Name:     passenger.Name,
		IDNumber: passenger.IDNumber,
		UserID:   passenger.UserID,
	}
	return response, nil
}

func (uc *PassengerUsecase) CreatePassenger(ctx context.Context, request model.PassengerRequest) (model.Passenger, error) {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		uc.Log.Warn("failed to begin transaction", zap.Error(err))
		return model.Passenger{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to begin transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// validate request
	if err := uc.Validate.Struct(request); err != nil {
		uc.Log.Warn("validation failed", zap.Error(err))
		return model.Passenger{}, utils.WrapError(fiber.StatusBadRequest, uc.Log, utils.Warn, err, "validation failed")
	}

	var userId pgtype.UUID
	var user repository.User
	if request.UserID == uuid.Nil {
		userId = pgtype.UUID{Valid: false}
	} else {
		user, err = tx.GetUser(ctx, request.UserID)
		if err != nil {
			return model.Passenger{}, utils.WrapError(fiber.StatusNotFound, uc.Log, utils.Warn, err, "failed to match passenger.userID")
		}
		userId = utils.ToPgUUID(user.ID)
	}

	passengerId := uuid.New()

	passenger, err := tx.CreatePassenger(ctx, repository.CreatePassengerParams{
		ID:       passengerId,
		Name:     request.Name,
		IDNumber: request.IDNumber,
		UserID:   userId,
	})

	if err != nil {
		uc.Log.Warn("failed to create passenger", zap.Error(err))
		return model.Passenger{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to create passenger")
	}

	if err := tx.Commit(ctx); err != nil {
		return model.Passenger{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to commit transaction")
	}

	response := model.Passenger{
		ID:        passenger.ID,
		Name:      passenger.Name,
		IDNumber:  passenger.IDNumber,
		UserID:    userId,
		CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
	}
	return response, nil

}

func (uc *PassengerUsecase) UpdatePassenger(ctx context.Context, request model.PassengerRequest) error {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to begin transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	r := repository.UpdatePassengerParams{
		ID:       request.ID,
		Name:     request.Name,
		IDNumber: request.IDNumber,
		UserID:   utils.ToPgUUID(request.UserID),
	}
	err = tx.UpdatePassenger(ctx, r)
	if err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to update passenger")
	}

	if err := tx.Commit(ctx); err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to commit transaction")
	}

	uc.Log.Info("passenger updated successfully")
	return nil
}

func (uc *PassengerUsecase) DeletePassenger(ctx context.Context, id uuid.UUID) error {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to begin transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	err = tx.DeletePassenger(ctx, id)
	if err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to delete passenger")
	}

	if err := tx.Commit(ctx); err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to commit transaction")
	}

	uc.Log.Info("passenger deleted successfully", zap.String("id", id.String()))
	return nil
}

func (uc *PassengerUsecase) GetPassengerByUserID(ctx context.Context, userID uuid.UUID) (model.Passenger, error) {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return model.Passenger{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to begin transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// get passenger by user id
	passenger, err := tx.GetPassengerByUser(ctx, utils.ToPgUUID(userID))
	if err != nil {
		return model.Passenger{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to get passenger")
	}

	if err := tx.Commit(ctx); err != nil {
		return model.Passenger{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to commit transaction")
	}

	response := model.Passenger{
		ID:       passenger.ID,
		Name:     passenger.Name,
		IDNumber: passenger.IDNumber,
		UserID:   passenger.UserID,
	}

	return response, nil
}
