package usecase

import (
	"context"
	"railway-go/internal/constant/model"
	"railway-go/internal/repository"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ReservationUsecase struct {
	Repo   repository.Store
	Logger *zap.Logger
}

func (uc *ReservationUsecase) ReserveSeat(ctx context.Context, req *model.ReservationRequest) (uuid.UUID, error) {
	lockttl := 5 * time.Minute

	booked, err := uc.Repo.CheckSeatAvailability(ctx, repository.CheckSeatAvailabilityParams{
		ScheduleID: req.ScheduleID,
		WagonID:    req.WagonID,
		SeatID:     req.Seat_id,
	})

	if err != nil {
		uc.Logger.Error("Failed to check seat availability", zap.Error(err))
		return uuid.Nil, fiber.ErrInternalServerError
	}

	if booked > 0 {
		uc.Logger.Info("Seat is already booked")
		return uuid.Nil, fiber.NewError(fiber.StatusConflict, "Seat is already booked")
	}

	locked, err := uc.Repo.LockSeat(ctx, req.ScheduleID, req.WagonID, req.Seat_id, lockttl)

	if err != nil {
		uc.Logger.Error("Failed to acquire seat lock", zap.Error(err))
		return uuid.Nil, fiber.ErrInternalServerError
	}

	if !locked {
		uc.Logger.Info("Seat is already locked")
		return uuid.Nil, fiber.NewError(fiber.StatusConflict, "Seat is already locked")
	}

	seat, err := uc.Repo.CreateReservation(ctx, repository.CreateReservationParams{
		PassengerID: req.PassengerID,
		ScheduleID:  req.ScheduleID,
		WagonID:     req.WagonID,
		SeatID:      req.Seat_id,
		DiscountID:  req.DiscountID,
	})

	if err != nil {
		uc.Logger.Error("Failed to create reservation", zap.Error(err))
		// ensure seat lock is removed if reservation fails
		uc.Repo.UnlockSeat(ctx, req.ScheduleID, req.WagonID, req.Seat_id)
		return uuid.Nil, fiber.ErrInternalServerError
	}

	return seat.ID, nil

}

// func (uc *ReservationUsecase) ConfitmPayment(ctx context.Context, reservationID uuid.UUID, paymentData model.PaymentRequest)

// response, err := uc.Repo.GetFullReservation(ctx, seat.ID)
// seatNumber := fmt.Sprintf("%s/%d-%s", response.WagonNumber, response.SeatNumber, response.SeatRow)

// uc.Logger.Info("Reservation created", zap.Any("req", req))
// return &model.ReservationDetailResponse{
// 	ReservationID:      response.ReservationID,
// 	PassengerName:      *response.PassengerName,
// 	PassengerIDNumber:  *response.PassengerIDNumber,
// 	UserEmail:          *response.UserEmail,
// 	TrainName:          *response.TrainName,
// 	SeatNumber:         seatNumber,
// 	BookingDate:        time.Now(),
// 	DepartureTime:      response.DepartureTime,
// 	ArrivalTime:        response.ArrivalTime,
// 	ClassType:          response.ClassType,
// 	Price:              response.Price,
// 	SourceStation:      response.SourceStation,
// 	DestinationStation: response.DestinationStation,
// 	PaymentStatus:      "pending",
// 	DiscountCode:       response.DiscountCode,
// }, nil
