package model

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type ScheduleRequest struct {
	TrainID        int64            `json:"train_id" validate:"required"`
	DepartureDate  pgtype.Timestamp `json:"departure_time" validate:"required"`
	AvailableSeats int32            `json:"available_seats" validate:"required"`
	Price          int64            `json:"price" validate:"required"`
	RouteID        int64            `json:"route_id" validate:"required"`
}

type Schedule struct {
	ID             int64            `json:"id"`
	TrainID        int64            `json:"train_id"`
	RouteID        int64            `json:"route_id"`
	DepartureDate  pgtype.Timestamp `json:"departure_date"`
	ArrivalDate    pgtype.Timestamp `json:"arrival_date"`
	Price          int64            `json:"price"`
	AvailableSeats int32            `json:"available_seats"`
	CreatedAt      pgtype.Timestamp `json:"created_at"`
	UpdatedAt      pgtype.Timestamp `json:"updated_at"`
}

type SearchScheduleRequest struct {
	SourceStation      *string `json:"source_station" validate:"required"`
	DestinationStation *string `json:"destination_station" validate:"required"`
	DepartureDate      string  `json:"departure_date" validate:"required"`
}

type SearchScheduleResponse struct {
	ScheduleID         int64            `json:"schedule_id"`
	TrainName          string           `json:"train_name"`
	SourceStation      string           `json:"source_station"`
	DestinationStation string           `json:"destination_station"`
	DepartureDate      pgtype.Timestamp `json:"departure_date"`
	ArrivalDate        pgtype.Timestamp `json:"arrival_date"`
	AvailableSeats     int32            `json:"available_seats"`
	Price              int64            `json:"price"`
}
