package model

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type RegisterUserRequest struct {
	Name        string `json:"name" validate:"required,max=100"`
	Email       string `json:"email" validate:"required,max=50"`
	Password    string `json:"password" validate:"required,max=100"`
	PhoneNumber string `json:"phone_number" validate:"required,max=50"`
}

type ReserveSeatRequest struct {
	PassengerID uuid.UUID `json:"passenger_id" validate:"required"`
	ScheduleID  int64     `json:"schedule_id" validate:"required"`
	SeatNumber  int64     `json:"seat_number" validate:"required"`
	PaymentID   uuid.UUID `json:"payment_id" validate:"required"`
}

type ReservationRequest struct {
	PassengerID uuid.UUID   `json:"passenger_id"`
	ScheduleID  int64       `json:"schedule_id"`
	WagonID     int64       `json:"wagon_id"`
	Seat_id     int64       `json:"seat_id"`
	DiscountID  pgtype.UUID `json:"discount_id,omitempty"`
}
