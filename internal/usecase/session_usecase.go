package usecase

import (
	"railway-go/internal/repository"
)

type SessionUseCase interface {
}

type sessionUseCase struct {
	sessionRepo repository.SessionRepository
}
