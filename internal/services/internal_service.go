package services

import (
	"context"
	"fmt"
	"invoiceB2B/internal/dtos"
	"invoiceB2B/internal/models"
	"invoiceB2B/internal/repositories"
	"log"
	"time"
)

type InternalService interface {
	UpdateInvoiceWithProcessedData(ctx context.Context, invoiceID uint, req dtos.UpdateInvoiceProcessedDataRequest, ipAddress string) (*dtos.InvoiceResponse, error)
}

type internalService struct {
	invoiceRepo    repositories.InvoiceRepository
	activityLogSvc ActivityLogService
}

func NewInternalService(invoiceRepo repositories.InvoiceRepository, activityLogSvc ActivityLogService) InternalService {
	return &internalService{
		invoiceRepo:    invoiceRepo,
		activityLogSvc: activityLogSvc,
	}
}

func (s *internalService) UpdateInvoiceWithProcessedData(ctx context.Context, invoiceID uint, req dtos.UpdateInvoiceProcessedDataRequest, ipAddress string) (*dtos.InvoiceResponse, error) {
	invoice, err := s.invoiceRepo.FindByID(ctx, invoiceID)
	if err != nil {
		log.Printf("Error finding invoice %d for processed data update: %v", invoiceID, err)
		return nil, ErrInvoiceNotFound // Assuming ErrInvoiceNotFound is defined
	}

	// Update fields from request
	invoice.JSONData = req.JSONData
	if req.ExtractedInvoiceNumber != nil {
		invoice.InvoiceNumber = *req.ExtractedInvoiceNumber
	}
	if req.ExtractedAmount != nil {
		invoice.Amount = *req.ExtractedAmount
	}
	if req.ExtractedCurrency != nil {
		invoice.Currency = *req.ExtractedCurrency
	}
	if req.ExtractedDebtorName != nil {
		invoice.DebtorName = *req.ExtractedDebtorName
	}
	if req.ExtractedIssuerName != nil {
		invoice.IssuerName = *req.ExtractedIssuerName
	}
	if req.ExtractedIssuerBankAccount != nil {
		invoice.IssuerBankAccount = *req.ExtractedIssuerBankAccount
	}
	if req.ExtractedIssuerBankName != nil {
		invoice.IssuerBankName = *req.ExtractedIssuerBankName
	}

	if req.ExtractedDueDate != nil && *req.ExtractedDueDate != "" {
		dueDate, err := time.Parse("2006-01-02", *req.ExtractedDueDate)
		if err == nil {
			invoice.DueDate = &dueDate
		} else {
			log.Printf("Warning: Could not parse ExtractedDueDate '%s' for invoice %d: %v", *req.ExtractedDueDate, invoiceID, err)
		}
	}

	if req.ProcessingError != nil {
		invoice.ProcessingError = req.ProcessingError
		invoice.Status = models.InvoiceProcessingFailed // Or a similar status
	} else if req.NewStatus != nil {
		invoice.Status = *req.NewStatus // n8n can suggest a status like "PendingApproval"
	} else {
		// Default status after successful processing if not specified by n8n
		if invoice.Status == models.InvoicePendingReview { // Only if it was pending review
			invoice.Status = models.InvoicePendingApproval
		}
	}

	if err := s.invoiceRepo.Update(ctx, invoice); err != nil {
		log.Printf("Error updating invoice %d with processed data: %v", invoiceID, err)
		return nil, fmt.Errorf("failed to update invoice: %w", err)
	}

	// Log activity
	logDetails := map[string]interface{}{
		"invoice_id": invoiceID,
		"source":     "n8n_processor",
		"new_status": invoice.Status,
	}
	if req.ProcessingError != nil {
		logDetails["processing_error"] = *req.ProcessingError
	}

	_ = s.activityLogSvc.LogActivity(ctx, nil, &invoice.UserID, "INVOICE_DATA_PROCESSED", logDetails, ipAddress)

	// We need mapInvoiceToResponse here or use a similar function
	// For now, just returning the updated model, but a DTO is better.
	// This should ideally return a dtos.InvoiceResponse
	// Let's use the one from admin_service.go for consistency (renamed to localMapInvoiceToResponse there)
	// For this service, let's define a similar mapper or make it shared.
	// For simplicity, returning the model directly for now but this is not ideal for API response.

	// To return dtos.InvoiceResponse, we need the mapper.
	// Let's assume mapInvoiceToResponse from invoice_service.go is available or re-defined.
	// For this self-contained update, I will redefine it locally.
	mappedResp := mapInvoiceToResponse(invoice)
	return &mappedResp, nil
}
