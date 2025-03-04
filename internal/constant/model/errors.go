package model

import "errors"

var (
	ErrSessionNotFound = errors.New("Session not found")
	ErrSessionInvalid  = errors.New("Session token is invalid or has expired. please reauthenticate")
	ErrInvalidToken    = errors.New("token is invalid")
	ErrExpiredToken    = errors.New("token has expired")
)
