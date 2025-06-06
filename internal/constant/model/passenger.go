package model

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Passenger struct {
	ID        uuid.UUID        `json:"id"`
	Name      string           `json:"name"`
	IDNumber  string           `json:"id_number"`
	UserID    pgtype.UUID      `json:"user_id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}

type PassengerRequest struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name" validate:"required"`
	IDNumber string    `json:"id_number" validate:"required,max=36"`
	UserID   uuid.UUID `json:"user_id" validate:"omitempty"`
}
