package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type Payload struct {
	ID        uuid.UUID    `json:"id"`
	UserID    *pgtype.UUID `json:"email,omitempty"` // can be nil for guest
	Role      string       `json:"role"`
	IssuedAt  time.Time    `json:"issued_at"`
	ExpiredAt time.Time    `json:"expired_at"`
}

func NewPayload(userID *pgtype.UUID, role string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	payload := &Payload{
		ID:        tokenID,
		UserID:    userID,
		Role:      role,
		IssuedAt:  now,
		ExpiredAt: now.Add(duration),
	}

	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
