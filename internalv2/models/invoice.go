package models

import (
	"time"
)

// Invoice represents an invoice in the system
type Invoice struct {
	ID            string    `json:"id" db:"id"`
	BusinessID    string    `json:"business_id" db:"business_id"`
	InvoiceNumber string    `json:"invoice_number" db:"invoice_number"`
	CustomerName  string    `json:"customer_name" db:"customer_name"`
	Amount        float64   `json:"amount" db:"amount"`
	Currency      string    `json:"currency" db:"currency"`
	Status        string    `json:"status" db:"status"`
	IssueDate     time.Time `json:"issue_date" db:"issue_date"`
	DueDate       time.Time `json:"due_date" db:"due_date"`
	Description   string    `json:"description" db:"description"`
	FileURL       string    `json:"file_url" db:"file_url"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at" db:"deleted_at"`
}

// InvoiceCreateRequest represents the invoice creation request
type InvoiceCreateRequest struct {
	InvoiceNumber string    `json:"invoice_number" validate:"required"`
	CustomerName  string    `json:"customer_name" validate:"required"`
	Amount        float64   `json:"amount" validate:"required,gt=0"`
	Currency      string    `json:"currency" validate:"required"`
	IssueDate     time.Time `json:"issue_date" validate:"required"`
	DueDate       time.Time `json:"due_date" validate:"required"`
	Description   string    `json:"description"`
}

// InvoiceUpdateRequest represents the invoice update request
type InvoiceUpdateRequest struct {
	InvoiceNumber string    `json:"invoice_number" validate:"required"`
	CustomerName  string    `json:"customer_name" validate:"required"`
	Amount        float64   `json:"amount" validate:"required,gt=0"`
	Currency      string    `json:"currency" validate:"required"`
	IssueDate     time.Time `json:"issue_date" validate:"required"`
	DueDate       time.Time `json:"due_date" validate:"required"`
	Description   string    `json:"description"`
}

// InvoiceListRequest represents the invoice list request
type InvoiceListRequest struct {
	Page   int    `json:"page" validate:"min=1"`
	Limit  int    `json:"limit" validate:"min=1,max=100"`
	Status string `json:"status" validate:"omitempty,oneof=pending submitted approved funded rejected paid"`
	Search string `json:"search"`
}

// InvoiceListResponse represents the invoice list response
type InvoiceListResponse struct {
	Invoices   []Invoice `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// Pagination represents pagination information
type Pagination struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Total     int `json:"total"`
	TotalPages int `json:"totalPages"`
	HasNext   bool `json:"hasNext"`
	HasPrev   bool `json:"hasPrev"`
} 