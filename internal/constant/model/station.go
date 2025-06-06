package model

import "github.com/jackc/pgx/v5/pgtype"

type StationRequest struct {
	Code        string `json:"code" validate:"required,max=4"`
	StationName string `json:"station_name" validate:"required,max=50"`
}

type Station struct {
	ID          int64            `json:"id"`
	Code        string           `json:"code"`
	StationName string           `json:"station_name"`
	CreatedAt   pgtype.Timestamp `json:"created_at"`
	UpdatedAt   pgtype.Timestamp `json:"updated_ad"`
}
