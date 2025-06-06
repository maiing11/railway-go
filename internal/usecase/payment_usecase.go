package usecase

import (
	"context"
	"math/rand"
	"railway-go/internal/constant/model"
	"railway-go/internal/repository"
	"railway-go/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type PaymentUC interface {
	ProcessMockPayment(ctx context.Context, req model.PaymentRequest) (model.PaymentResponse, error)
	AutoCancelExpiredPayments(ctx context.Context) error
}

type PaymentUsecase struct {
	*UseCase
	ReservationUC
}

func NewPaymentUsecase(useCase *UseCase) PaymentUC {
	return &PaymentUsecase{UseCase: useCase}
}

// simulate payment processing
func (uc *PaymentUsecase) ProcessMockPayment(ctx context.Context, req model.PaymentRequest) (model.PaymentResponse, error) {

	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return model.PaymentResponse{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to begin transaction")
	}

	// validate reservation exists
	reservation, err := tx.GetReservation(ctx, req.ReservationID)
	if err != nil {
		return model.PaymentResponse{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to get reservation")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	if reservation.ReservationStatus != "pending" {
		return model.PaymentResponse{}, utils.WrapError(fiber.StatusBadRequest, uc.Log, utils.Warn, err, "reservation already paid or canceled")
	}

	// simulating processing delay
	time.Sleep(1500 * time.Millisecond)

	success := rand.Intn(100) < 80

	var status string
	var message string
	transactionID := uuid.New()

	if success {
		status = "success"
		message = "Payment successful!"
		err = tx.ConfirmReservation(ctx, req.ReservationID)
		if err != nil {
			return model.PaymentResponse{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to confirm reservation")
		}
	} else {
		status = "failed"
		message = "Payment failed!"
	}

	uc.Log.Info("payment status", zap.Any("status", status))
	if err = tx.CreatePayment(ctx, repository.CreatePaymentParams{
		ReservationID:   req.ReservationID,
		PaymentMethod:   req.PaymentMethod,
		PaymentStatus:   status,
		Amount:          req.Amount,
		GatewayResponse: &status,
		PaymentDate:     pgtype.Timestamp{Time: time.Now(), Valid: true},
		TransactionID:   transactionID.String(),
	}); err != nil {
		return model.PaymentResponse{Message: "failed"}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to create payment")
	}

	if err := uc.Repo.UnlockSeat(ctx, reservation.ScheduleID, reservation.WagonID, reservation.SeatID); err != nil {
		return model.PaymentResponse{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to unlock reservation")
	}

	if err := tx.DecreaseWagonSeat(ctx, reservation.WagonID); err != nil {
		return model.PaymentResponse{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to decrease wagon seat")
	}

	if err := tx.Commit(ctx); err != nil {
		return model.PaymentResponse{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to commit transaction")
	}

	return model.PaymentResponse{
		Transaction: transactionID,
		Status:      status,
		Message:     message,
	}, nil
}

func (uc *PaymentUsecase) AutoCancelExpiredPayments(ctx context.Context) error {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to begin transaction")
	}
	expiredReservation, err := tx.GetExpiredPayments(ctx)
	if err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to get expired payments")
	}

	for _, res := range expiredReservation {
		err := uc.Repo.CancelReservation(ctx, res)
		if err != nil {
			return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to cancel reservation")
		}

		uc.Log.Info("Auto-canceling expired reservation", zap.String("reservation_id", res.String()))
	}

	if err := tx.Commit(ctx); err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to commit transaction")
	}
	uc.Log.Info("Auto-cancel expired payments completed successfully")
	return nil
}
