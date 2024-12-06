package entity

import (
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Session struct {
	ID           string       `json:"id"`
	UserID       *pgtype.UUID `json:"user_id,omitempty"`
	RefreshToken string       `json:"refresh_token"`
	Role         string       `json:"role"`
	UserAgent    string       `json:"user_agent"`
	ClientIP     string       `json:"client_ip"`
	IsBlocked    bool         `json:"is_blocked"`
	ExpiresAt    time.Time    `json:"expires_at"`
}

var (
	ErrSessionNotFound = errors.New("Session not found")
	ErrSessionInvalid  = errors.New("Session token is invalid or has expired. Please reauthenticate.")
)
