package model

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Phonenumber string    `json:"phonenumber"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Session struct {
	ID           string     `json:"id"`
	UserID       *uuid.UUID `json:"user_id,omitempty"`
	RefreshToken string     `json:"refresh_token"`
	Role         string     `json:"role"`
	UserAgent    string     `json:"user_agent"`
	ClientIP     string     `json:"client_ip"`
	IsBlocked    bool       `json:"is_blocked"`
	ExpiresAt    time.Time  `json:"expires_at"`
}

type ReservationResponse struct {
	Message     string    `json:"message"`
	ExpiresAt   time.Time `json:"expires_at"`
	ScheduleID  string    `json:"schedule_id"`
	SeatNumber  string    `json:"seat_number"`
	PassengerID string    `json:"passenger_id"`
}

type ReservationDetailResponse struct {
	ReservationID      uuid.UUID `json:"reservation_id"`
	PassengerName      string    `json:"passenger_name"`
	PassengerIDNumber  string    `json:"passenger_id_number"`
	UserEmail          string    `json:"user_email"`
	TrainName          string    `json:"train_name"`
	SeatNumber         string    `json:"seat_number"`
	BookingDate        time.Time `json:"booking_date"`
	DepartureTime      time.Time `json:"departure_time"`
	ArrivalTime        time.Time `json:"arrival_time"`
	ClassType          string    `json:"class_type"`
	Price              int64     `json:"price"`
	SourceStation      string    `json:"source_station"`
	DestinationStation string    `json:"destination_station"`
	PaymentStatus      string    `json:"payment_status"`
	DiscountCode       *string   `json:"discount_code"`
}
