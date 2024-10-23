package repository

import (
	"context"
)

type ReservationTxResult struct {
	Reservation Reservation `json:"reservation"`
	Schedule    *Schedule   `json:"schedule"`
}

func (store *SQLStore) ReservationSeatTX(ctx context.Context, arg CreateReservationParams) (ReservationTxResult, error) {
	var result ReservationTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Reservation, err = q.CreateReservation(ctx, CreateReservationParams{
			PassengerID: arg.PassengerID,
			ScheduleID:  arg.ScheduleID,
			SeatNumber:  arg.SeatNumber,
			BookingDate: arg.BookingDate,
			PaymentID:   arg.PaymentID,
		})
		if err != nil {
			return err
		}

		err = q.UpdateSchedule(ctx, UpdateScheduleParams{
			AvailableSeats: result.Schedule.AvailableSeats - 1,
		})
		return err
	})
	return result, err
}
