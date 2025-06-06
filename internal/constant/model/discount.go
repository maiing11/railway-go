package model

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type DiscountRequest struct {
	Code            string           `json:"code" validate:"required"`
	DiscountPercent int32            `json:"discount" validate:"required"`
	ExpiresAt       pgtype.Timestamp `json:"expires_at" validate:"required"`
	MaxUses         int32            `json:"max_uses" validate:"required"`
}

type DiscountResponseRow struct {
	ID              uuid.UUID        `json:"id"`
	Code            string           ` json:"code"`
	DiscountPercent int32            ` json:"discount_percent"`
	ExpiresAt       pgtype.Timestamp ` json:"expires_at"`
	MaxUses         int32            `json:"max_uses"`
}
