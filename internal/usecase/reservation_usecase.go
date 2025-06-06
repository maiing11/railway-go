package usecase

import (
	"context"
	"fmt"
	"railway-go/internal/constant/model"
	"railway-go/internal/repository"
	"railway-go/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type ReservationUC interface {
	CreateReservation(ctx context.Context, req model.ReservationRequest) (model.Reservation, error)
	GetDetailReservation(ctx context.Context, id uuid.UUID) (model.ListReservationsResponse, error)
	CancelReservation(ctx context.Context, id uuid.UUID) error
	ConfirmReservation(ctx context.Context, id uuid.UUID) error
	GetReservationById(ctx context.Context, id uuid.UUID) (repository.Reservation, error)
	DeleteReservation(ctx context.Context, id uuid.UUID) error
	GetAllReservations(ctx context.Context) ([]model.ListReservationsResponse, int64, error)
	AutoDeleteReservations(ctx context.Context) error
}
type ReservationUsecase struct {
	*UseCase
}

func NewReservationUsecase(useCase *UseCase) ReservationUC {
	return &ReservationUsecase{UseCase: useCase}
}

// func (uc *ReservationUsecase) StartReservationCleanup(ctx context.Context) {
// go func() {
// 		ticker := time.NewTicker(5 * time.Minute)
// 		defer ticker.Stop()

// 		for {
// 			select {
// 			case <-ticker.C:
// 				ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
// 				defer cancel()

// 				err := uc.Repo.cancel(ctx)
// 				if err != nil {
// 					log.Warn("auto-cancel cleanup failed", zap.Error(err))
// 				} else {
// 					log.Info("auto-cancel expired reservations executed")
// 				}
// 			}
// 		}
// 	}()
// }

// CreateReservation handles the process of creating a new reservation for a train seat.
// It performs the following steps:
//  1. Begins a database transaction.
//  2. Retrieves the schedule, wagon, and seat based on the request.
//  3. Validates that the seat belongs to the selected wagon.
//  4. Validates the request body.
//  5. Checks if the seat is already booked or locked.
//  6. Locks the seat for a specified TTL to prevent race conditions.
//  7. Calculates the price, applying a discount if provided.
//  8. Retrieves or creates a passenger associated with the user.
//  9. Creates the reservation record in the database.
//  10. Applies the discount to the reservation if applicable.
//  11. Commits the transaction and returns the reservation details.
//
// If any step fails, the transaction is rolled back and an appropriate error is returned.
//
// Parameters:
//   - ctx: context.Context for request-scoped values and cancellation.
//   - req: model.ReservationRequest containing reservation details.
//
// Returns:
//   - model.Reservation: The created reservation object.
//   - error: An error if the reservation could not be created.
func (uc *ReservationUsecase) CreateReservation(ctx context.Context, req model.ReservationRequest) (model.Reservation, error) {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return model.Reservation{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to begin transaction")
	}

	defer func() {
		if p := recover(); p != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				uc.Log.Error("rollback failed during panic", zap.Error(rollbackErr))
			}
			panic(p) // Re-throw panic
		} else if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				uc.Log.Error("rollback failed", zap.Error(rollbackErr))
			}
		}
	}()

	// validate request
	if err := uc.Validate.Struct(req); err != nil {
		return model.Reservation{}, utils.WrapError(fiber.StatusBadRequest, uc.Log, utils.Error, err, "validation failed")
	}

	var passenger repository.Passenger
	passengerID, _ := utils.ToUUID(req.PassengerID)
	// this condition allows all role except guest to auto get passenger with user_id
	if passengerID == uuid.Nil {
		userID := utils.ToPgUUID(req.UserId)
		passenger, err = tx.GetPassengerByUser(ctx, userID)
		if err != nil {
			err := uc.Repo.UnlockSeat(ctx, req.ScheduleID, req.WagonID, req.Seat_id)
			if err != nil {
				return model.Reservation{}, err
			}
			return model.Reservation{}, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("No passenger found for the provided user ID: %s. Please ensure the passenger exists or register a new passenger.", req.UserId))
		}
	} else {
		passenger, err = tx.GetPassenger(ctx, passengerID)
		if err != nil {
			err := uc.Repo.UnlockSeat(ctx, req.ScheduleID, req.WagonID, req.Seat_id)
			if err != nil {
				return model.Reservation{}, err
			}
			return model.Reservation{}, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("No passenger found for the provided passenger ID: %s. Please ensure the passenger_id is correct or register a new passenger", passengerID))
		}
	}

	// uc.Log.Info("got passenger", zap.Any("passenger id: ", passenger.ID))

	// get the Schedule
	schedule, err := tx.GetSchedule(ctx, req.ScheduleID)
	if err != nil {
		return model.Reservation{}, fiber.NewError(fiber.StatusBadRequest, "failed fetch schedule")
	}

	wagon, err := tx.GetWagon(ctx, req.WagonID)
	if err != nil {
		return model.Reservation{}, fiber.NewError(fiber.StatusBadRequest, "failed to fetch wagon")
	}

	seat, err := tx.GetSeat(ctx, req.Seat_id)
	if err != nil {
		return model.Reservation{}, fiber.NewError(fiber.StatusBadRequest, "failed to fetch seat")
	}

	bookedParams := repository.CheckSeatAvailabilityParams{
		ScheduleID: schedule.ID,
		WagonID:    wagon.ID,
		SeatID:     seat.ID,
	}
	// counting from reservations where schedule, wagon , seat
	booked, err := tx.CheckSeatAvailability(ctx, bookedParams)
	if err != nil {
		return model.Reservation{},
			fiber.NewError(fiber.StatusInternalServerError, "failed to check seat availability")
	}

	if booked > 0 {
		return model.Reservation{}, fiber.NewError(fiber.StatusConflict, "seat already booked")
	}

	// default price
	price := schedule.Price

	var discountID pgtype.UUID
	if req.DiscountID == uuid.Nil {
		discountID = pgtype.UUID{Valid: false}
	} else {
		discountID = utils.ToPgUUID(req.DiscountID)
		discount, err := tx.GetDiscountByID(ctx, req.DiscountID)
		if err != nil {
			return model.Reservation{}, utils.WrapError(fiber.StatusNotFound, uc.Log, utils.Warn, err, "failed to get discount")
		}
		if discount.DiscountPercent > 0 {
			price -= price * int64(discount.DiscountPercent) / 100
		}
		expired := time.Now()
		if discount.ExpiresAt.Valid && discount.ExpiresAt.Time.Before(expired) {
			return model.Reservation{}, utils.WrapError(fiber.StatusRequestTimeout, uc.Log, utils.Warn, nil, "discount expired")
		}

	}

	// uc.Log.Info("reservation params", zap.Any("discount", discount))

	// check if a discount is provided

	bookingTime := pgtype.Timestamp{
		Time:  time.Now(),
		Valid: true,
	}
	expiresAt := pgtype.Timestamp{
		Time:  time.Now().Add(15 * time.Minute),
		Valid: true,
	}

	params := repository.CreateReservationParams{
		PassengerID:       passenger.ID,
		ScheduleID:        schedule.ID,
		WagonID:           wagon.ID,
		SeatID:            seat.ID,
		BookingDate:       bookingTime,
		ReservationStatus: "pending",
		ExpiresAt:         expiresAt,
		DiscountID:        discountID,
		Price:             &price,
	}

	lockttl := 5 * time.Minute
	if err := uc.Repo.LockSeat(ctx, schedule.ID, wagon.ID, seat.ID, lockttl); err != nil {
		return model.Reservation{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to lock seat")
	}

	reserve, err := tx.CreateReservation(ctx, params)
	if err != nil {
		// var pgErr *pgconn.PgError
		// if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		// 	_ = uc.Repo.UnlockSeat(ctx, req.ScheduleID, req.WagonID, req.Seat_id)
		// 	return model.Reservation{}, utils.WrapError(fiber.StatusConflict, uc.Log, utils.Warn, err, fmt.Sprintf("seat already booked:%v", booked))
		// }
		// ensure seat lock is removed if reservation fails
		return model.Reservation{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, fmt.Sprintf("failed to create reservation for passengerID %v, scheduleID %v, wagonID %v, seatID %v",
			passenger.ID,
			req.ScheduleID,
			req.WagonID,
			req.Seat_id,
		),
		)
	}

	if reserve.DiscountID.Valid {
		err := tx.ApplyDiscountToReservation(ctx, repository.ApplyDiscountToReservationParams{
			ReservationID: reserve.ID,
			DiscountID:    req.DiscountID,
		})
		if err != nil {
			err := uc.Repo.UnlockSeat(ctx, req.ScheduleID, req.WagonID, req.Seat_id)
			if err != nil {
				return model.Reservation{}, err
			}
			return model.Reservation{}, fiber.NewError(fiber.StatusBadRequest, "failed to apply discount")
		}
	}

	response := model.Reservation{
		ID:                reserve.ID,
		PassengerID:       reserve.ID,
		ScheduleID:        reserve.ScheduleID,
		WagonID:           reserve.WagonID,
		SeatID:            reserve.SeatID,
		BookingDate:       reserve.BookingDate,
		DiscountID:        reserve.DiscountID,
		Price:             reserve.Price,
		ReservationStatus: string(reserve.ReservationStatus),
		ExpiresAt:         reserve.ExpiresAt,
		CreatedAt:         reserve.CreatedAt,
		UpdatedAt:         reserve.UpdatedAt,
	}

	// commit transaction
	if err := tx.Commit(ctx); err != nil {
		return model.Reservation{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to commit transaction")
	}

	return response, nil
}

func (uc *ReservationUsecase) GetDetailReservation(ctx context.Context, id uuid.UUID) (model.ListReservationsResponse, error) {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return model.ListReservationsResponse{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to begin transaction")
	}
	defer func() {
		if err != nil {
			err := tx.Rollback(ctx)
			if err != nil {
				return
			}
		}
	}()

	reservation, err := tx.GetFullReservation(ctx, id)
	if err != nil {
		return model.ListReservationsResponse{},
			utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to get reservation")
	}

	seatNumber := fmt.Sprintf("Gerbong %d/%s-%d", *reservation.WagonNumber, reservation.SeatRow.SeatRow, *reservation.SeatNumber)

	response := model.ListReservationsResponse{
		ReservationID:      reservation.ReservationID,
		PassengerName:      reservation.PassengerName,
		PassengerIDNumber:  reservation.PassengerIDNumber,
		UserEmail:          reservation.UserEmail,
		DepartureDate:      reservation.DepartureDate,
		ArrivalDate:        reservation.ArrivalDate,
		TicketPrice:        reservation.TicketPrice,
		TrainName:          reservation.TrainName,
		ClassType:          string(reservation.ClassType.TipeClass), // Extract the string value from NullTipeClass
		SeatNumber:         seatNumber,
		BookingDate:        reservation.BookingDate, // Extract the string value from NullTipeClass
		SourceStation:      reservation.SourceStation,
		DestinationStation: reservation.DestinationStation,
		DiscountCode:       reservation.DiscountCode,
		PaymentAmount:      reservation.PaymentAmount,
		PaymentMethod:      reservation.PaymentMethod,
		PaymentStatus:      reservation.PaymentStatus,
	}

	return response, nil
}

func (uc *ReservationUsecase) GetReservationById(ctx context.Context, id uuid.UUID) (repository.Reservation, error) {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return repository.Reservation{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to begin transaction")
	}

	defer func() {
		if err != nil {
			err := tx.Rollback(ctx)
			if err != nil {
				return
			}
		}
	}()

	reservation, err := tx.GetReservation(ctx, id)
	if err != nil {
		return repository.Reservation{}, utils.WrapError(fiber.StatusNotFound, uc.Log, utils.Error, err, "failed to get reservation")
	}

	if err := tx.Commit(ctx); err != nil {
		return repository.Reservation{}, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to commit get reservation")
	}

	return reservation, nil
}

// func (uc *ReservationUsecase) DeleteReservation(ctx context.Context) error {
// 	tx, err := uc.Repo.BeginTransaction(ctx)
// 	if err != nil {
// 		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to begin transaction")

// 	}

// 	defer func() {
// 		if err != nil {
// 			err := tx.Rollback(ctx)
// 			if err != nil {
// 				return
// 			}
// 		}
// 	}()
// }

func (uc *ReservationUsecase) CancelReservation(ctx context.Context, id uuid.UUID) error {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to begin transaction")
	}
	defer func() {
		if err != nil {
			err := tx.Rollback(ctx)
			if err != nil {
				return
			}
		}
	}()

	reservation, err := uc.GetReservationById(ctx, id)
	if err != nil {
		return utils.WrapError(fiber.StatusNotFound, uc.Log, utils.Warn, err, "failed to get reservation / unkown id")
	}

	if reservation.ReservationStatus == "canceled" {
		return utils.WrapError(fiber.StatusConflict, uc.Log, utils.Info, nil, "Reservation already canceled")
	}

	err = tx.CancelReservation(ctx, id)
	if err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to canceled reservation")
	}

	if err := tx.Commit(ctx); err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to commit trancsaction")
	}

	uc.Log.Info("successfully canceled status reservation")
	return nil
}

func (uc *ReservationUsecase) DeleteReservation(ctx context.Context, id uuid.UUID) error {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to begin transaction")
	}
	defer func() {
		if err != nil {
			err := tx.Rollback(ctx)
			if err != nil {
				return
			}
		}
	}()

	reservation, err := tx.GetReservation(ctx, id)
	if err != nil {
		return utils.WrapError(fiber.StatusNotFound, uc.Log, utils.Warn, err, "failed unkown id or invalid reservation_id")
	}

	if reservation.ReservationStatus != "canceled" {
		return fiber.NewError(fiber.StatusBadRequest, "failed reservation status isn't canceled yet ")
	}

	err = tx.DeleteReservation(ctx, reservation.ID)
	if err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, err.Error())
	}

	if err := tx.Commit(ctx); err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to commit trancsaction")
	}
	uc.Log.Info("successfully deleted reservation", zap.Any("id", reservation.ID))
	return nil
}

func (uc *ReservationUsecase) AutoDeleteReservations(ctx context.Context) error {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to begin transaction")
	}
	defer func() {
		if err != nil {
			err := tx.Rollback(ctx)
			if err != nil {
				return
			}
		}
	}()

	if err := tx.ExpireUndpaidReservations(ctx); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to delete expired reservation")
	}

	if err := tx.Commit(ctx); err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to commit trancsaction")
	}

	return nil
}

func (uc *ReservationUsecase) ConfirmReservation(ctx context.Context, id uuid.UUID) error {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to begin transaction")
	}
	defer func() {
		if err != nil {
			err := tx.Rollback(ctx)
			if err != nil {
				return
			}
		}
	}()

	reservation, err := tx.GetReservation(ctx, id)
	if err != nil {
		uc.Log.Warn("failed to get reservation", zap.Error(err))
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to get reservation")
	}

	if reservation.ReservationStatus == "confirmed" {
		return utils.WrapError(fiber.StatusBadRequest, uc.Log, utils.Warn, err, "Reservation already confirmed")
	}

	err = tx.ConfirmReservation(ctx, id)
	if err != nil {
		uc.Log.Warn("failed to confirm reservation", zap.Error(err))
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to confirm reservation")
	}

	if err := tx.Commit(ctx); err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to commit transaction")
	}

	uc.Log.Info("reservation confirmed successfully", zap.String("reservation_id", id.String()))
	return nil
}

func (uc *ReservationUsecase) GetAllReservations(ctx context.Context) ([]model.ListReservationsResponse, int64, error) {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		uc.Log.Warn("failed to begin transaction", zap.Error(err))
		return nil, 0, fiber.NewError(fiber.StatusInternalServerError, "failed to begin transaction")
	}
	defer func() {
		if err != nil {
			err := tx.Rollback(ctx)
			if err != nil {
				return
			}
		}
	}()

	arg := repository.ListReservationsParams{
		Limit:  100,
		Offset: 0,
	}

	reservations, err := tx.ListReservations(ctx, arg)
	if err != nil {
		return nil, 0, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to get all reservations")
	}
	// get total count
	totalItem, err := tx.CountReservations(ctx)
	if err != nil {
		return nil, 0, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Warn, err, "failed to count reservations")
	}

	var response []model.ListReservationsResponse
	for _, reservation := range reservations {
		seatNumber := fmt.Sprintf("Gerbong %d/%s-%d", *reservation.WagonNumber, reservation.SeatRow.SeatRow, *reservation.SeatNumber)

		response = append(response, model.ListReservationsResponse{
			ReservationID:      reservation.ReservationID,
			PassengerName:      reservation.PassengerName,
			PassengerIDNumber:  reservation.PassengerIDNumber,
			UserEmail:          reservation.UserEmail,
			DepartureDate:      reservation.DepartureDate,
			ArrivalDate:        reservation.ArrivalDate,
			TicketPrice:        reservation.TicketPrice,
			TrainName:          reservation.TrainName,
			ClassType:          string(reservation.ClassType.TipeClass), // Extract the string value from NullTipeClass
			SeatNumber:         seatNumber,
			BookingDate:        reservation.BookingDate, // Extract the string value from NullTipeClass
			SourceStation:      reservation.SourceStation,
			DestinationStation: reservation.DestinationStation,
			DiscountCode:       reservation.DiscountCode,
			PaymentAmount:      reservation.PaymentAmount,
			PaymentMethod:      reservation.PaymentMethod,
			PaymentStatus:      reservation.PaymentStatus,
		})
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, 0, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to commit transaction")
	}
	return response, totalItem, nil
}
