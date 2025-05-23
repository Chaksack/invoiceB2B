package dtos

import (
	"invoiceB2B/internal/models"
	"mime/multipart"
	"time"
)

type InvoiceUploadRequest struct {
	File *multipart.FileHeader `form:"invoiceFile" validate:"required"`
}

type InvoiceResponse struct {
	ID                      uint                 `json:"id"`
	UserID                  uint                 `json:"userId"`
	InvoiceNumber           string               `json:"invoiceNumber,omitempty"` // omitempty if null
	IssuerName              string               `json:"issuerName,omitempty"`
	IssuerBankAccount       string               `json:"issuerBankAccount,omitempty"`
	IssuerBankName          string               `json:"issuerBankName,omitempty"`
	DebtorName              string               `json:"debtorName,omitempty"`
	Amount                  float64              `json:"amount,omitempty"`   // omitempty if null
	Currency                string               `json:"currency,omitempty"` // omitempty if null
	DueDate                 *time.Time           `json:"dueDate,omitempty"`
	Status                  models.InvoiceStatus `json:"status"`
	OriginalFilePath        string               `json:"originalFilePath,omitempty"`
	JSONData                string               `json:"jsonData,omitempty"`
	UploadedAt              time.Time            `json:"uploadedAt"`
	ApprovedAt              *time.Time           `json:"approvedAt,omitempty"`
	DisbursedAt             *time.Time           `json:"disbursedAt,omitempty"`
	FinancingFeePercentage  *float64             `json:"financingFeePercentage,omitempty"`
	FinancedAmount          *float64             `json:"financedAmount,omitempty"`
	DisbursementReceiptPath string               `json:"disbursementReceiptPath,omitempty"`
	ProcessingError         *string              `json:"processingError,omitempty"`
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
	Status                 models.InvoiceStatus `json:"status" validate:"required,oneof=approved rejected disbursed repaid"`
	FinancingFeePercentage *float64             `json:"financingFeePercentage,omitempty" validate:"omitempty,gt=0,lte=100"`
	FinancedAmount         *float64             `json:"financedAmount,omitempty" validate:"omitempty,gt=0"`
	RejectionReason        *string              `json:"rejectionReason,omitempty"`
}

type AdminUploadReceiptRequest struct {
	File *multipart.FileHeader `form:"receiptFile" validate:"required"`
}

// New DTO for n8n to update invoice data
type UpdateInvoiceProcessedDataRequest struct {
	JSONData                   string                `json:"jsonData" validate:"required,json"`
	ExtractedInvoiceNumber     *string               `json:"extractedInvoiceNumber,omitempty"`
	ExtractedAmount            *float64              `json:"extractedAmount,omitempty" validate:"omitempty,gt=0"`
	ExtractedCurrency          *string               `json:"extractedCurrency,omitempty" validate:"omitempty,len=3"`
	ExtractedDueDate           *string               `json:"extractedDueDate,omitempty" validate:"omitempty,datetime=2006-01-02"` // Expect YYYY-MM-DD
	ExtractedDebtorName        *string               `json:"extractedDebtorName,omitempty"`
	ExtractedIssuerName        *string               `json:"extractedIssuerName,omitempty"`
	ExtractedIssuerBankAccount *string               `json:"extractedIssuerBankAccount,omitempty"`
	ExtractedIssuerBankName    *string               `json:"extractedIssuerBankName,omitempty"`
	ProcessingError            *string               `json:"processingError,omitempty"` // For n8n to report errors
	NewStatus                  *models.InvoiceStatus `json:"newStatus,omitempty"`       // Optional: n8n can suggest a new status
}
