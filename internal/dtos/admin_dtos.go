package dtos

import (
	"invoiceB2B/internal/models"
	"time"
)

// KYC DTOs for Admin
type AdminKYCReviewRequest struct {
	Status          models.KYCStatus `json:"status" validate:"required,oneof=approved rejected"`
	RejectionReason *string          `json:"rejectionReason,omitempty"` // Required if status is 'rejected'
}

type AdminKYCDetailResponse struct {
	ID              uint             `json:"id"`
	UserID          uint             `json:"userId"`
	UserEmail       string           `json:"userEmail"` // For context
	Status          models.KYCStatus `json:"status"`
	SubmittedAt     *time.Time       `json:"submittedAt,omitempty"`
	ReviewedAt      *time.Time       `json:"reviewedAt,omitempty"`
	ReviewedByEmail string           `json:"reviewedByEmail,omitempty"` // Staff email
	RejectionReason *string          `json:"rejectionReason,omitempty"`
	DocumentsInfo   string           `json:"documentsInfo,omitempty"` // JSON string or parsed
	CreatedAt       time.Time        `json:"createdAt"`
	UpdatedAt       time.Time        `json:"updatedAt"`
}

// Staff DTOs for Admin
type CreateStaffRequest struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"firstName" validate:"required,min=2"`
	LastName  string `json:"lastName" validate:"required,min=2"`
	Password  string `json:"password" validate:"required,min=8"`
	Role      string `json:"role" validate:"required,oneof=admin kyc_reviewer finance_manager"` // Example roles
}

type UpdateStaffRequest struct {
	FirstName *string `json:"firstName,omitempty" validate:"omitempty,min=2"`
	LastName  *string `json:"lastName,omitempty" validate:"omitempty,min=2"`
	Role      *string `json:"role,omitempty" validate:"omitempty,oneof=admin kyc_reviewer finance_manager"`
	IsActive  *bool   `json:"isActive,omitempty"`
	// Password change should be a separate, more secure endpoint
}

type StaffResponse struct {
	ID          uint       `json:"id"`
	Email       string     `json:"email"`
	FirstName   string     `json:"firstName"`
	LastName    string     `json:"lastName"`
	Role        string     `json:"role"`
	IsActive    bool       `json:"isActive"`
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// Activity Log DTOs
type ActivityLogResponse struct {
	ID         uint      `json:"id"`
	StaffEmail *string   `json:"staffEmail,omitempty"` // Email of staff who performed action
	UserEmail  *string   `json:"userEmail,omitempty"`  // Email of user related to action
	Action     string    `json:"action"`
	Details    string    `json:"details,omitempty"` // JSON string
	IPAddress  *string   `json:"ipAddress,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}

type ActivityLogListResponse struct {
	Logs     []ActivityLogResponse `json:"logs"`
	Total    int64                 `json:"total"`
	Page     int                   `json:"page"`
	PageSize int                   `json:"pageSize"`
}
