package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type ReservationRequest struct {
	PassengerID pgtype.UUID `json:"passenger_id"`
	UserId      uuid.UUID   `json:"user_id"`
	ScheduleID  int64       `json:"schedule_id" validate:"required,max=50"`
	WagonID     int64       `json:"wagon_id" validate:"required,max=50"`
	Seat_id     int64       `json:"seat_id" validate:"required,max=50"`
	DiscountID  uuid.UUID   `json:"discount_id" validate:"max=50"`
}

type ReservationResponse struct {
	Message     string    `json:"message"`
	ExpiresAt   time.Time `json:"expires_at"`
	ScheduleID  int64     `json:"schedule_id"`
	SeatNumber  string    `json:"seat_number"`
	PassengerID uuid.UUID `json:"passenger_id"`
}

type Reservation struct {
	ID                uuid.UUID        `json:"id"`
	PassengerID       uuid.UUID        `json:"passenger_id"`
	ScheduleID        int64            `json:"schedule_id"`
	WagonID           int64            `json:"wagon_id"`
	SeatID            int64            `json:"seat_id"`
	BookingDate       pgtype.Timestamp `json:"booking_date"`
	DiscountID        pgtype.UUID      `json:"discount_id"`
	Price             *int64           `json:"price"`
	ReservationStatus string           `json:"reservation_status"`
	ExpiresAt         pgtype.Timestamp `json:"expires_at"`
	CreatedAt         pgtype.Timestamp `json:"created_at"`
	UpdatedAt         pgtype.Timestamp `json:"updated_at"`
}

type ListReservationsResponse struct {
	ReservationID      uuid.UUID        `json:"reservation_id"`
	PassengerName      *string          `json:"passenger_name"`
	PassengerIDNumber  *string          `json:"passenger_id_number"`
	UserName           *string          `json:"user_name"`
	UserEmail          *string          `json:"user_email"`
	DepartureDate      pgtype.Timestamp `json:"departure_date"`
	ArrivalDate        pgtype.Timestamp `json:"arrival_date"`
	TicketPrice        *int64           `json:"ticket_price"`
	TrainName          *string          `json:"train_name"`
	ClassType          string           `json:"class_type"`
	SeatNumber         string           `json:"seat_number"`
	BookingDate        pgtype.Timestamp `json:"booking_date"`
	ReservationStatus  string           `json:"reservation_status"`
	SourceStation      *string          `json:"source_station"`
	DestinationStation *string          `json:"destination_station"`
	DiscountCode       *string          `json:"discount_code"`
	DiscountPercent    *int32           `json:"discount_percent"`
	PaymentAmount      *int64           `json:"payment_amount"`
	PaymentMethod      *string          `json:"payment_method"`
	PaymentStatus      *string          `json:"payment_status"`
}
