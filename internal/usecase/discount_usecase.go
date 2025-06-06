package usecase

import (
	"context"
	"railway-go/internal/constant/model"
	"railway-go/internal/repository"
	"railway-go/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type DiscountUC interface {
	CreateDiscount(ctx context.Context, request model.DiscountRequest) (repository.DiscountCode, error)
	GetDiscountByID(ctx context.Context, id uuid.UUID) (repository.DiscountCode, error)
	GetDiscountByCode(ctx context.Context, code string) (repository.DiscountCode, error)
	ReduceDiscountUsage(ctx context.Context, id uuid.UUID) error
}

type DiscountUsecase struct {
	*UseCase
	WagonUC
}

func NewDiscountUsecase(useCase *UseCase) DiscountUC {
	return &DiscountUsecase{UseCase: useCase}
}

func (uc *DiscountUsecase) CreateDiscount(ctx context.Context, request model.DiscountRequest) (repository.DiscountCode, error) {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return repository.DiscountCode{},
			utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to begin transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// validate request
	if err := uc.Validate.Struct(request); err != nil {
		return repository.DiscountCode{},
			utils.WrapError(fiber.StatusBadRequest, uc.Log, utils.Warn, err, "error validating request body")
	}

	discount := repository.CreateDiscountCodeParams{
		Code:            request.Code,
		DiscountPercent: request.DiscountPercent,
		ExpiresAt:       request.ExpiresAt,
		MaxUses:         request.MaxUses,
	}

	discountCode, err := tx.CreateDiscountCode(ctx, discount)
	if err != nil {
		return repository.DiscountCode{},
			utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to create discount code")
	}

	if err := tx.Commit(ctx); err != nil {
		return repository.DiscountCode{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to commit")
	}

	uc.Log.Info("discount code created successfully", zap.String("code", discountCode.Code))
	return discountCode, nil
}

func (uc *DiscountUsecase) GetDiscountByID(ctx context.Context, id uuid.UUID) (repository.DiscountCode, error) {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return repository.DiscountCode{},
			utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to begin transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	discountCode, err := tx.GetDiscountByID(ctx, id)
	if err != nil {
		return repository.DiscountCode{},
			utils.WrapError(fiber.StatusNotFound, uc.Log, utils.Warn, err, "failed to get discount code")
	}

	return discountCode, nil
}

func (uc *DiscountUsecase) GetDiscountByCode(ctx context.Context, code string) (repository.DiscountCode, error) {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return repository.DiscountCode{},
			utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to begin transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	discountCode, err := tx.GetDiscountByCode(ctx, code)
	if err != nil {
		return repository.DiscountCode{},
			utils.WrapError(fiber.StatusNotFound, uc.Log, utils.Warn, err, "failed to get discount code")
	}

	return discountCode, nil
}

func (uc *DiscountUsecase) ReduceDiscountUsage(ctx context.Context, id uuid.UUID) error {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to begin transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	err = tx.ReduceDiscountUsage(ctx, id)
	if err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to decrease discount usage")
	}

	if err := tx.Commit(ctx); err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to commit transaction")
	}

	return nil
}

func (uc *DiscountUsecase) ApplyDiscountToReservation(ctx context.Context, reservationID, discountID uuid.UUID) error {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to begin transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	reservation, err := tx.GetReservation(ctx, reservationID)
	if err != nil {
		return utils.WrapError(fiber.StatusNotFound, uc.Log, utils.Warn, err, "failed to get reservation")
	}

	if reservation.DiscountID.Valid {
		return utils.WrapError(fiber.StatusBadRequest, uc.Log, utils.Warn, err, "discount code already applied")
	}

	discountData, err := uc.GetDiscountByID(ctx, discountID)
	if err != nil {
		return utils.WrapError(fiber.StatusNotFound, uc.Log, utils.Warn, err, "failed to get discount code")
	}

	if discountData.MaxUses <= 0 {
		return utils.WrapError(fiber.StatusBadRequest, uc.Log, utils.Warn, err,
			"discount code has reached maximum usage limit")
	}

	err = tx.ApplyDiscountToReservation(ctx, repository.ApplyDiscountToReservationParams{
		ReservationID: reservation.ID,
		DiscountID:    discountData.ID,
	})
	if err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to apply discount code to reservation")
	}

	err = uc.ReduceDiscountUsage(ctx, discountData.ID)
	if err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to reduce discount usage")
	}

	if err := tx.Commit(ctx); err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log,
			utils.Warn,
			err,
			"failed to commit transaction",
		)
	}

	return nil
}

func (uc *DiscountUsecase) GetDiscountFromReservation(ctx context.Context, discountID uuid.UUID) ([]model.DiscountResponseRow, error) {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return nil,
			utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to begin transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	discountCode, err := tx.GetDiscountsForReservation(ctx, discountID)
	if err != nil {
		return nil,
			utils.WrapError(fiber.StatusNotFound, uc.Log, utils.Warn, err, "failed to get discount code")
	}

	var discountCodes []model.DiscountResponseRow
	for _, code := range discountCode {
		discountCodes = append(discountCodes, model.DiscountResponseRow{
			ID:              code.ID,
			Code:            code.Code,
			DiscountPercent: code.DiscountPercent,
			ExpiresAt:       code.ExpiresAt,
			MaxUses:         code.MaxUses,
		})
	}

	return discountCodes, nil
}
