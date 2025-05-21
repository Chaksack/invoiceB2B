package dtos

import (
	"github.com/google/uuid"
	"time"
)

// User DTOs

type UserProfileResponse struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	CompanyName  string    `json:"companyName"`
	IsActive     bool      `json:"isActive"`
	TwoFAEnabled bool      `json:"twoFaEnabled"`
	KYCStatus    string    `json:"kycStatus"` // e.g., "pending", "approved", "rejected", "Not Submitted"
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type UpdateUserProfileRequest struct {
	FirstName   string `json:"firstName" validate:"omitempty,min=2,max=50"`
	LastName    string `json:"lastName" validate:"omitempty,min=2,max=50"`
	CompanyName string `json:"companyName" validate:"omitempty,min=2,max=100"`
	// Password updates should be a separate endpoint for security
}

type SubmitKYCRequest struct {
	// Fields depend on what KYC information you collect
	BusinessRegistrationNumber string `json:"businessRegistrationNumber" validate:"required"`
	AddressLine1               string `json:"addressLine1" validate:"required"`
	City                       string `json:"city" validate:"required"`
	Country                    string `json:"country" validate:"required"`
	DocumentsInfo              string `json:"documentsInfo" validate:"omitempty,json"`
}

type KYCStatusResponse struct {
	UserID          uuid.UUID  `json:"userId"`
	Status          string     `json:"status"`
	SubmittedAt     *time.Time `json:"submittedAt,omitempty"`
	ReviewedAt      *time.Time `json:"reviewedAt,omitempty"`
	RejectionReason *string    `json:"rejectionReason,omitempty"`
	Message         string     `json:"message"`
}
