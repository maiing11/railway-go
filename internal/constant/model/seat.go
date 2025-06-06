package model

import "github.com/jackc/pgx/v5/pgtype"

type SeatRequest struct {
	WagonID     int64  `json:"wagon_id" validate:"required"`
	SeatNumber  int32  `json:"seat_number" validate:"required"`
	SeatRow     string `json:"seat_row" validate:"required,max=1"`
	IsAvailable bool   `json:"is_available" default:"true"`
}

type Seat struct {
	ID          int64            `json:"id"`
	WagonID     *int64           `json:"wagon_id"`
	SeatNumber  int32            `json:"seat_number"`
	SeatRow     string           `json:"seat_row"`
	IsAvailable *bool            `json:"is_available"`
	CreatedAt   pgtype.Timestamp `json:"created_at"`
	UpdatedAt   pgtype.Timestamp `json:"updated_at"`
}
