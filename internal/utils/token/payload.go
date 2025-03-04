package token

import (
	"railway-go/internal/constant/model"
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	ID        uuid.UUID  `json:"id"`
	UserID    *uuid.UUID `json:"email,omitempty"` // can be nil for guest
	SessionID string     `json:"session_id,omitempty"`
	Role      string     `json:"role"`
	IssuedAt  time.Time  `json:"issued_at"`
	ExpiredAt time.Time  `json:"expired_at"`
}

func NewPayload(userID *uuid.UUID, sessionID, role string, duration time.Duration) (*Payload, error) {
	tokenID := uuid.New()

	now := time.Now()
	payload := &Payload{
		ID:        tokenID,
		UserID:    userID,
		SessionID: sessionID,
		Role:      role,
		IssuedAt:  now,
		ExpiredAt: now.Add(duration),
	}

	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return model.ErrExpiredToken
	}
	return nil
}
