package model

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type RegisterUserRequest struct {
	Name        string `json:"name" validate:"required,max=100"`
	Email       string `json:"email" validate:"required,max=50"`
	Password    string `json:"password" validate:"required,max=100"`
	PhoneNumber string `json:"phone_number" validate:"required,max=50"`
}
