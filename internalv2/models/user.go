package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID           string    `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Role         string    `json:"role" db:"role"`
	IsApproved   bool      `json:"is_approved" db:"is_approved"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// UserLoginRequest represents the login request
type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserRegisterRequest represents the registration request
type UserRegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"oneof=business admin"`
}

// UserResponse represents the user response
type UserResponse struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	Role       string    `json:"role"`
	IsApproved bool      `json:"is_approved"`
	CreatedAt  time.Time `json:"created_at"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
} 