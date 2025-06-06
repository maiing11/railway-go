package usecase

import (
	"context"
	"railway-go/internal/constant/model"
	"railway-go/internal/repository"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type SeatUC interface {
	CreateSeat(ctx context.Context, request *model.SeatRequest) (repository.Seat, error)
	GetSeat(ctx context.Context, id int64) (repository.Seat, error)
	DeleteSeat(ctx context.Context, id int64) error
	UpdateSeat(ctx context.Context, id int64, request *model.SeatRequest) error
	GetAllSeat(ctx context.Context, wagonID int64) ([]model.Seat, error)
}

type SeatUsecase struct {
	*UseCase
	WagonUC
}

func NewSeatUsecase(useCase *UseCase) SeatUC {
	return &SeatUsecase{UseCase: useCase}
}

func (uc *SeatUsecase) CreateSeat(ctx context.Context, request *model.SeatRequest) (repository.Seat, error) {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		uc.Log.Warn("failed to begin transaction", zap.Error(err))
		return repository.Seat{}, fiber.NewError(fiber.StatusInternalServerError, "failed to begin transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// validate request
	if err := uc.Validate.Struct(request); err != nil {
		uc.Log.Warn("error validating request body", zap.Error(err))
		return repository.Seat{}, fiber.ErrBadRequest
	}

	// check if wagon exists
	wagon, err := tx.GetWagon(ctx, request.WagonID)
	if err != nil {
		uc.Log.Warn("invalid wagon id or wagon not found", zap.Error(err))
		return repository.Seat{}, fiber.ErrInternalServerError
	}

	seat := repository.CreateSeatParams{
		WagonID:     &wagon.ID,
		SeatNumber:  request.SeatNumber,
		SeatRow:     repository.SeatRow(request.SeatRow),
		IsAvailable: &request.IsAvailable,
	}

	response, err := tx.CreateSeat(ctx, seat)
	if err != nil {
		uc.Log.Warn("failed to create seat", zap.Error(err))
		return repository.Seat{}, fiber.ErrInternalServerError
	}

	if err := tx.Commit(ctx); err != nil {
		uc.Log.Warn("failed to commit transaction", zap.Error(err))
		return repository.Seat{}, fiber.NewError(fiber.StatusInternalServerError, "failed to commit transaction")
	}

	return response, nil
}

func (uc *SeatUsecase) GetSeat(ctx context.Context, id int64) (repository.Seat, error) {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		uc.Log.Warn("failed to begin transaction", zap.Error(err))
		return repository.Seat{}, fiber.NewError(fiber.StatusInternalServerError, "failed to begin transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	seat, err := tx.GetSeat(ctx, id)
	if err != nil {
		uc.Log.Warn("failed to get seat", zap.Error(err))
		return repository.Seat{}, fiber.ErrInternalServerError
	}

	if err := tx.Commit(ctx); err != nil {
		uc.Log.Warn("failed to commit transaction", zap.Error(err))
		return repository.Seat{}, fiber.NewError(fiber.StatusInternalServerError, "failed to commit transaction")
	}

	return seat, nil
}

func (uc *SeatUsecase) DeleteSeat(ctx context.Context, id int64) error {
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

	seat, err := uc.GetSeat(ctx, id)
	if err != nil {
		uc.Log.Warn("error seat not found", zap.Error(err))
		return fiber.NewError(fiber.StatusNotFound, "error seat not found")
	}

	if err := tx.DeleteSeat(ctx, seat.ID); err != nil {
		uc.Log.Error("failed to delete seat", zap.Error(err))
		return nil
	}

	if err := tx.Commit(ctx); err != nil {
		uc.Log.Error("failed to commit transaction", zap.Error(err))
		return fiber.ErrInternalServerError
	}

	return nil
}

func (uc *SeatUsecase) UpdateSeat(ctx context.Context, id int64, request *model.SeatRequest) error {
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
		uc.Log.Warn("error validating request body", zap.Error(err))
		return fiber.ErrBadRequest
	}

	// check if wagon exists
	wagon, err := uc.GetWagon(ctx, request.WagonID)
	if err != nil {
		uc.Log.Warn("invalid wagon id or wagon not found", zap.Error(err))
		return fiber.ErrInternalServerError
	}

	seat := repository.UpdateSeatParams{
		ID:          id,
		WagonID:     &wagon.ID,
		IsAvailable: &request.IsAvailable,
	}

	if err := tx.UpdateSeat(ctx, seat); err != nil {
		uc.Log.Warn("failed to update seat", zap.Error(err))
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit(ctx); err != nil {
		uc.Log.Warn("failed to commit transaction", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "failed to commit transaction")
	}

	return nil
}

func (uc *SeatUsecase) GetAllSeat(ctx context.Context, wagonID int64) ([]model.Seat, error) {
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

	seats, err := tx.ListSeats(ctx, &wagonID)
	if err != nil {
		uc.Log.Warn("failed to get all seats", zap.Error(err))
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit(ctx); err != nil {
		uc.Log.Warn("failed to commit transaction", zap.Error(err))
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to commit transaction")
	}

	response := make([]model.Seat, len(seats))
	for i, seat := range seats {
		response[i] = model.Seat{
			ID:          seat.ID,
			WagonID:     seat.WagonID,
			SeatNumber:  seat.SeatNumber,
			SeatRow:     string(seat.SeatRow),
			IsAvailable: seat.IsAvailable,
			CreatedAt:   seat.CreatedAt,
			UpdatedAt:   seat.UpdatedAt,
		}
	}

	return response, nil
}
