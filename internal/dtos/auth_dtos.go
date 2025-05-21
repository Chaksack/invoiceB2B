package dtos

import "github.com/google/uuid"

// Auth DTOs (Data Transfer Objects)

type RegisterUserRequest struct {
	Email       string `json:"email" validate:"required,email"`
	FirstName   string `json:"firstName" validate:"required,min=2,max=50"`
	LastName    string `json:"lastName" validate:"required,min=2,max=50"`
	CompanyName string `json:"companyName" validate:"required,min=2,max=100"`
	Password    string `json:"password" validate:"required,min=8,max=72"` // Max 72 due to bcrypt limit
}

type UserResponse struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	CompanyName  string    `json:"companyName"`
	IsActive     bool      `json:"isActive"`
	TwoFAEnabled bool      `json:"twoFaEnabled"`
}

type RegisterUserResponse struct {
	User    UserResponse `json:"user"`
	Message string       `json:"message"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginUserResponse struct {
	User                 UserResponse `json:"user"`
	AccessToken          string       `json:"accessToken"`
	RefreshToken         string       `json:"refreshToken"`
	Message              string       `json:"message"` // e.g. "Login successful" or "OTP sent to your email"
	TwoFARequired        bool         `json:"twoFaRequired"`
	AccessTokenExpiresAt int64        `json:"accessTokenExpiresAt"` // Unix timestamp
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type RefreshTokenResponse struct {
	AccessToken          string `json:"accessToken"`
	RefreshToken         string `json:"refreshToken"`         // Optionally, refresh token can be rotated
	AccessTokenExpiresAt int64  `json:"accessTokenExpiresAt"` // Unix timestamp
	Message              string `json:"message"`
}

type RequestOTPRequest struct {
	Email string `json:"email" validate:"required,email"` // Could also be implicit if user is logged in
}

type RequestOTPResponse struct {
	Message string `json:"message"` // e.g., "OTP sent to your email if account exists and 2FA is enabled."
}

type VerifyOTPRequest struct {
	Email string `json:"email" validate:"required,email"` // Or UserID if already partially logged in
	OTP   string `json:"otp" validate:"required,len=6,numeric"`
}

// VerifyOTPResponse is similar to LoginUserResponse but after OTP verification
type VerifyOTPResponse struct {
	User                 UserResponse `json:"user"`
	AccessToken          string       `json:"accessToken"`
	RefreshToken         string       `json:"refreshToken"`
	Message              string       `json:"message"`
	AccessTokenExpiresAt int64        `json:"accessTokenExpiresAt"`
}

type Enable2FARequest struct {
	Enable bool `json:"enable"` // true to enable, false to disable
}

type Enable2FAResponse struct {
	Message      string `json:"message"`
	TwoFAEnabled bool   `json:"twoFaEnabled"`
}
