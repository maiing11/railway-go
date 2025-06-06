package model

import "github.com/jackc/pgx/v5/pgtype"

type TrainRequest struct {
	TrainName string `json:"name" validate:"required,max=100"`
	Capacity  int32  `json:"capacity"  validate:"required,min=1"`
}

type Train struct {
	ID        int64            `json:"id"`
	TrainName string           `json:"name"`
	Capacity  int32            `json:"capacity"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}
