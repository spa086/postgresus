package users

import "github.com/google/uuid"

type SignUpRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type SignInRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SignInResponse struct {
	UserID uuid.UUID `json:"userId"`
	Token  string    `json:"token"`
}
