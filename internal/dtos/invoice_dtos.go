package dtos

import (
	"invoiceB2B/internal/models"
	"mime/multipart"
	"time"
)

type InvoiceUploadRequest struct {
	// Fields expected from multipart form, file handled separately
	// InvoiceNumber string `form:"invoiceNumber" validate:"required"`
	// Amount        float64 `form:"amount" validate:"required,gt=0"`
	// Currency      string  `form:"currency" validate:"required,len=3"`
	// DueDate       string  `form:"dueDate" validate:"required,datetime=2006-01-02"` // Expect YYYY-MM-DD
	// DebtorName    string  `form:"debtorName" validate:"required"`
	File *multipart.FileHeader `form:"invoiceFile" validate:"required"`
}

type InvoiceResponse struct {
	ID                      uint                 `json:"id"`
	UserID                  uint                 `json:"userId"`
	InvoiceNumber           string               `json:"invoiceNumber"`
	IssuerName              string               `json:"issuerName,omitempty"`
	IssuerBankAccount       string               `json:"issuerBankAccount,omitempty"`
	IssuerBankName          string               `json:"issuerBankName,omitempty"`
	DebtorName              string               `json:"debtorName,omitempty"`
	Amount                  float64              `json:"amount"`
	Currency                string               `json:"currency"`
	DueDate                 *time.Time           `json:"dueDate,omitempty"`
	Status                  models.InvoiceStatus `json:"status"`
	OriginalFilePath        string               `json:"originalFilePath,omitempty"` // Potentially hide in prod
	JSONData                string               `json:"jsonData,omitempty"`         // Or parsed object
	UploadedAt              time.Time            `json:"uploadedAt"`
	ApprovedAt              *time.Time           `json:"approvedAt,omitempty"`
	DisbursedAt             *time.Time           `json:"disbursedAt,omitempty"`
	FinancingFeePercentage  *float64             `json:"financingFeePercentage,omitempty"`
	FinancedAmount          *float64             `json:"financedAmount,omitempty"`
	DisbursementReceiptPath string               `json:"disbursementReceiptPath,omitempty"`
	CreatedAt               time.Time            `json:"createdAt"`
	UpdatedAt               time.Time            `json:"updatedAt"`
}

type InvoiceListResponse struct {
	Invoices []InvoiceResponse `json:"invoices"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"pageSize"`
}

type AdminInvoiceUpdateRequest struct {
	Status models.InvoiceStatus `json:"status" validate:"required,oneof=approved rejected disbursed repaid"` // Admin can set these
	// For 'disbursed', these might be set
	FinancingFeePercentage *float64 `json:"financingFeePercentage,omitempty" validate:"omitempty,gt=0,lte=100"`
	FinancedAmount         *float64 `json:"financedAmount,omitempty" validate:"omitempty,gt=0"`
	// For 'rejected', a reason might be needed
	RejectionReason *string `json:"rejectionReason,omitempty"`
}

type AdminUploadReceiptRequest struct {
	File *multipart.FileHeader `form:"receiptFile" validate:"required"`
}
