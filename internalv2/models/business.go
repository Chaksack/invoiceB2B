package models

import (
	"time"
)

// Business represents a business in the system
type Business struct {
	ID            string    `json:"id" db:"id"`
	UserID        string    `json:"user_id" db:"user_id"`
	CompanyName   string    `json:"company_name" db:"company_name"`
	Industry      string    `json:"industry" db:"industry"`
	AnnualRevenue float64   `json:"annual_revenue" db:"annual_revenue"`
	EmployeeCount int       `json:"employee_count" db:"employee_count"`
	Status        string    `json:"status" db:"status"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// BusinessWithUser represents a business with user information
type BusinessWithUser struct {
	Business
	Email         string  `json:"email" db:"email"`
	TotalInvoices int     `json:"total_invoices" db:"total_invoices"`
	TotalFunded   float64 `json:"total_funded" db:"total_funded"`
}

// BusinessListRequest represents the business list request
type BusinessListRequest struct {
	Page   int    `json:"page" validate:"min=1"`
	Limit  int    `json:"limit" validate:"min=1,max=100"`
	Status string `json:"status" validate:"omitempty,oneof=pending approved rejected suspended"`
	Search string `json:"search"`
}

// BusinessListResponse represents the business list response
type BusinessListResponse struct {
	Businesses []BusinessWithUser `json:"data"`
	Pagination Pagination         `json:"pagination"`
}

// BusinessCreateRequest represents the business creation request
type BusinessCreateRequest struct {
	CompanyName   string  `json:"company_name" validate:"required"`
	Industry      string  `json:"industry" validate:"required"`
	AnnualRevenue float64 `json:"annual_revenue" validate:"required,gt=0"`
	EmployeeCount int     `json:"employee_count" validate:"required,gt=0"`
}

// BusinessUpdateRequest represents the business update request
type BusinessUpdateRequest struct {
	CompanyName   string  `json:"company_name" validate:"required"`
	Industry      string  `json:"industry" validate:"required"`
	AnnualRevenue float64 `json:"annual_revenue" validate:"required,gt=0"`
	EmployeeCount int     `json:"employee_count" validate:"required,gt=0"`
} 