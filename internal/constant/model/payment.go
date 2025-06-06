package model

import "github.com/google/uuid"

type PaymentRequest struct {
	ReservationID uuid.UUID `json:"reservation_id" validate:"required"`
	PaymentMethod string    `json:"payment_method" validate:"required"`
	Amount        int64     `json:"amount"`
}

type PaymentResponse struct {
	Transaction uuid.UUID `json:"transaction_id"`
	Status      string    `json:"status"`
	Message     string    `json:"message"`
}
