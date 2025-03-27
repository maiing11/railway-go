package model

import "errors"

var (
	ErrSessionNotFound            = errors.New("Session not found")
	ErrSessionInvalid             = errors.New("Session token is invalid or has expired. please reauthenticate")
	ErrInvalidToken               = errors.New("token is invalid")
	ErrExpiredToken               = errors.New("token has expired")
	ErrSeatAlreadyLocked          = errors.New("seat is already locked")
	ErrSeatNotLocked              = errors.New("seat is not locked")
	ErrSeatNotFound               = errors.New("seat not found")
	ErrReservationNotFound        = errors.New("reservation not found")
	ErrReservationAlreadyExist    = errors.New("reservation already exist")
	ErrReservationNotAvailable    = errors.New("reservation not available")
	ErrReservationNotPaid         = errors.New("reservation not paid")
	ErrReservationAlreadyPaid     = errors.New("reservation already paid")
	ErrReservationNotCanceled     = errors.New("reservation not canceled")
	ErrReservationAlreadyCanceled = errors.New("reservation already canceled")
	ErrReservationNotRefunded     = errors.New("reservation not refunded")
	ErrReservationAlreadyRefunded = errors.New("reservation already refunded")
	ErrReservationNotExpired      = errors.New("reservation not expired")
	ErrReservationAlreadyExpired  = errors.New("reservation already expired")
	ErrReservationNotUpdated      = errors.New("reservation not updated")
	ErrReservationNotDeleted      = errors.New("reservation not deleted")
	ErrReservationNotCreated      = errors.New("reservation not created")
)
