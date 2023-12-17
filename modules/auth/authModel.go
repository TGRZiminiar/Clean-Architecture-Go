package auth

import (
	"time"

	"github.com/google/uuid"
)

type (
	UserJson struct {
		Id       uuid.UUID `json:"id"`
		Username string    `json:"username"`
		Email    string    `json:"email"`
		Image    string    `json:"image"`
		// Password  string    `json:"password"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	CreateUser struct {
		Username string `json:"username" validate:"required,max=32"`
		Email    string `json:"email" validate:"required,email,max=255"`
		// Image    string `json:"image"`
		Password string `json:"password" validate:"required,max=32"`
	}
)
