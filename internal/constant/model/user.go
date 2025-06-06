package model

import (
	"time"

	"github.com/google/uuid"
)

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type RegisterUserRequest struct {
	Name        string `json:"name" validate:"required,max=100"`
	Email       string `json:"email" validate:"required,email,max=50"`
	Password    string `json:"password" validate:"required,max=100"`
	PhoneNumber string `json:"phone_number" validate:"required,max=50"`
}

type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Phonenumber string    `json:"phonenumber"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RegisterAdminRequest struct {
	Name        string `json:"name" validate:"required,max=100"`
	Email       string `json:"email" validate:"required,max=50"`
	Password    string `json:"password" validate:"required,max=100"`
	PhoneNumber string `json:"phone_number" validate:"required,max=50"`
	Role        string `json:"role" validate:"required,oneof=admin general affair"`
}
