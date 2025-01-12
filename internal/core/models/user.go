package models

import (
	"time"
)

type User struct {
	ID          int        `json:"id"`
	UserName    string     `json:"user_name"`
	Email       string     `json:"email"`
	Password    string     `json:"password"`
	PhoneNumber string     `json:"phone_number"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type TempUser struct {
	UserName    string `json:"username"`
	Address     string
	Email       string `json:"email"`
	Password    string
	PhoneNumber string
}
type SignupInput struct {
	UserName    string `json:"user_name" validate:"required,min=3,max=16,alphanum"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required,len=10,numeric"`
	Password    string `json:"password" validate:"required,min=8,max=32"`
}
type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type PasswordReset struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
	Reenter         string `json:"reenter"`
}

type UserProfileResponse struct {
	Name      string
	Email     string
	PhnNumber string
	Status    string
}
