package model

import (
	"time"

	"github.com/google/uuid"
)

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
