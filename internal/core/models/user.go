package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName    string `gorm:"column:user_name;not null"`
	Email       string `gorm:"column:email;not null"`
	Password    string `gorm:"column:password;not null" json:"-"`
	PhoneNumber string `gorm:"column:phonenumber;not null"`
	Status      string `gorm:"check(status IN('Active', 'Inactive', 'Blocked'))"`
}

type TempUser struct {
	UserName    string `json:"username"`
	Address     string
	Email       string `json:"email"`
	Password    string
	PhoneNumber string
}
type SignupInput struct {
	UserName    string `json:"username" validate:"required,min=3,max=16,alphanum"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phonenumber" validate:"required,len=10,numeric"`
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
