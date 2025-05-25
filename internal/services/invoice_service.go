package services

import (
	"context"
	"errors"
	"fmt"
	"invoiceB2B/internal/config"
	"invoiceB2B/internal/dtos"
	"invoiceB2B/internal/models"
	"invoiceB2B/internal/repositories"
	"log"
	"path/filepath"
	"time"
)

// Errors defined in auth_service.go
var (
	ErrInvoiceNotFound                   = errors.New("invoice not found")
	ErrInvoiceAccessDenied               = errors.New("access to invoice denied")
	ErrReceiptNotFound                   = errors.New("receipt not found for this invoice")
	ErrKYCNotApprovedForInvoiceUpload    = errors.New("kyc not approved, cannot upload invoice")
	ErrInvoiceNotApprovedForDisbursement = errors.New("invoice not approved for disbursement")
	ErrInvalidInvoiceStatusForOperation  = errors.New("invalid invoice status for this operation")
	ErrInvoiceNotDisbursedForRepayment   = errors.New("invoice has not been disbursed, cannot process repayment")
	ErrRepaymentAmountMismatch           = errors.New("repayment amount does not match financed amount")
)

type InvoiceService interface {
	CreateInvoice(ctx context.Context, userID uint, req dtos.InvoiceUploadRequest) (*dtos.InvoiceResponse, error)
	GetUserInvoices(ctx context.Context, userID uint, page, pageSize int) ([]dtos.InvoiceResponse, int64, error)
	GetInvoiceByIDForUser(ctx context.Context, invoiceID, userID uint) (*dtos.InvoiceResponse, error)
	GetReceiptPathForUser(ctx context.Context, invoiceID, userID uint) (string, string, error)
}

type invoiceService struct {
	invoiceRepo     repositories.InvoiceRepository
	userRepo        repositories.UserRepository
	transactionRepo repositories.TransactionRepository
	fileService     FileService
	notificationSvc NotificationService
	activityLogSvc  ActivityLogService
	emailService    EmailService
	cfg             *config.Config
}

func NewInvoiceService(
	invoiceRepo repositories.InvoiceRepository,
	userRepo repositories.UserRepository,
	transactionRepo repositories.TransactionRepository,
	fileService FileService,
	notificationSvc NotificationService,
	activityLogSvc ActivityLogService,
	emailService EmailService,
	cfg *config.Config,
) InvoiceService {
	return &invoiceService{
		invoiceRepo:     invoiceRepo,
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
		fileService:     fileService,
		notificationSvc: notificationSvc,
		activityLogSvc:  activityLogSvc,
		emailService:    emailService,
		cfg:             cfg,
	}
}

// mapInvoiceToResponse helper function
func mapInvoiceToResponse(invoice *models.Invoice) dtos.InvoiceResponse {
	var receiptPath string
	if invoice.DisbursementReceiptPath != nil {
		receiptPath = *invoice.DisbursementReceiptPath
	}
	return dtos.InvoiceResponse{
		ID:                      invoice.ID,
		UserID:                  invoice.UserID,
		InvoiceNumber:           invoice.InvoiceNumber,
		IssuerName:              invoice.IssuerName,
		IssuerBankAccount:       invoice.IssuerBankAccount,
		IssuerBankName:          invoice.IssuerBankName,
		DebtorName:              invoice.DebtorName,
		Amount:                  invoice.Amount,
		Currency:                invoice.Currency,
		DueDate:                 invoice.DueDate,
		Status:                  invoice.Status,
		OriginalFilePath:        invoice.OriginalFilePath,
		JSONData:                invoice.JSONData,
		UploadedAt:              invoice.UploadedAt,
		ApprovedAt:              invoice.ApprovedAt,
		DisbursedAt:             invoice.DisbursedAt,
		FinancingFeePercentage:  invoice.FinancingFeePercentage,
		FinancedAmount:          invoice.FinancedAmount,
		DisbursementReceiptPath: receiptPath,
		CreatedAt:               invoice.CreatedAt,
		UpdatedAt:               invoice.UpdatedAt,
	}
}

func (s *invoiceService) CreateInvoice(ctx context.Context, userID uint, req dtos.InvoiceUploadRequest) (*dtos.InvoiceResponse, error) {
	user, err := s.userRepo.FindByIDWithKYC(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// --- START DEBUG LOGS --- (Can be removed after confirming fix)
	log.Printf("CreateInvoice DEBUG: For UserID %d:", userID)
	if user.KYCDetail == nil {
		log.Printf("CreateInvoice DEBUG: user.KYCDetail IS NIL.")
	} else {
		log.Printf("CreateInvoice DEBUG: user.KYCDetail.ID is '%d', user.KYCDetail.Status is '%s'. models.KYCApproved is '%s'.", user.KYCDetail.ID, user.KYCDetail.Status, models.KYCApproved)
		log.Printf("CreateInvoice DEBUG: Is user.KYCDetail.Status == models.KYCApproved? %t", user.KYCDetail.Status == models.KYCApproved)
	}
	// --- END DEBUG LOGS ---

	if user.KYCDetail == nil || user.KYCDetail.Status != models.KYCApproved {
		log.Printf("CreateInvoice INFO: KYC check failed for UserID %d. KYCDetail present: %t. Status from DB: '%s'. Expected status: '%s'.",
			userID,
			user.KYCDetail != nil,
			func() string {
				if user.KYCDetail != nil {
					return string(user.KYCDetail.Status)
				}
				return "N/A"
			}(),
			models.KYCApproved)
		return nil, ErrKYCNotApprovedForInvoiceUpload
	}

	// Assuming req.File is of type *multipart.FileHeader or compatible with fileService.SaveFile
	relativePath, originalFileName, err := s.fileService.SaveFile(req.File, "invoices")
	if err != nil {
		log.Printf("Error saving invoice file for user %d: %v", userID, err)
		return nil, fmt.Errorf("failed to save invoice file: %w", err)
	}
	log.Printf("Invoice file saved: relativePath=%s, originalFileName=%s", relativePath, originalFileName)

	now := time.Now()
	invoice := &models.Invoice{
		UserID:           userID,
		Status:           models.InvoicePendingReview, // Ensure models.InvoicePendingReview is defined
		OriginalFilePath: relativePath,
		UploadedAt:       now,
		InvoiceNumber:    originalFileName, // Or a generated one
		JSONData:         "{}",             // Corrected: Initialize with an empty JSON object string
		// Initialize other nullable fields to their zero values or specific defaults if needed
		// GORM handles zero values for basic types (0 for float64, "" for string, nil for *time.Time)
		// which will translate to NULL in the DB if the column is nullable.
	}

	if err := s.invoiceRepo.Create(ctx, invoice); err != nil {
		// The log from the repository might be more specific, but this adds service-level context.
		log.Printf("Error creating invoice record in service for user %d, file %s: %v", userID, originalFileName, err)
		return nil, fmt.Errorf("failed to create invoice record: %w", err)
	}

	// Send email confirmation for invoice submission
	go func() {
		subject := "Invoice Submission Confirmation"
		body := fmt.Sprintf("Hi %s,\n\nYour invoice '%s' (ID: %d) has been successfully submitted and is pending review.\n\nThanks,\nThe Team",
			user.FirstName, originalFileName, invoice.ID)
		if s.emailService != nil { // Check if emailService is initialized
			if emailErr := s.emailService.SendEmail(user.Email, subject, body); emailErr != nil {
				log.Printf("Failed to send invoice submission confirmation email to %s for invoice %d: %v", user.Email, invoice.ID, emailErr)
			}
		} else {
			log.Println("EmailService is not initialized, skipping invoice submission email.")
		}
	}()

	eventPayload := map[string]interface{}{
		"invoice_id":        invoice.ID,
		"user_id":           userID,
		"user_email":        user.Email,
		"company_name":      user.CompanyName,
		"file_path":         relativePath,
		"original_filename": originalFileName, // Use the filename obtained from SaveFile
		"uploaded_at":       invoice.UploadedAt.Format(time.RFC3339),
	}
	if s.notificationSvc != nil { // Check if notificationSvc is initialized
		err = s.notificationSvc.PublishEvent(
			s.cfg.RabbitMQEventExchangeName, // Ensure cfg fields are correct
			s.cfg.RabbitMQInvoiceUploadedRoutingKey,
			eventPayload,
		)
		if err != nil {
			log.Printf("Failed to publish invoice.uploaded event for invoice %d: %v", invoice.ID, err)
		} else {
			log.Printf("Published invoice.uploaded event for admin notification: Invoice ID %d", invoice.ID)
		}
	} else {
		log.Println("NotificationService is not initialized, skipping invoice.uploaded event.")
	}

	if s.activityLogSvc != nil { // Check if activityLogSvc is initialized
		_ = s.activityLogSvc.LogActivity(ctx, nil, &userID, "INVOICE_UPLOADED",
			map[string]interface{}{"invoice_id": invoice.ID, "filename": originalFileName}, "")
	} else {
		log.Println("ActivityLogService is not initialized, skipping INVOICE_UPLOADED log.")
	}

	resp := mapInvoiceToResponse(invoice)
	return &resp, nil
}

func (s *invoiceService) GetUserInvoices(ctx context.Context, userID uint, page, pageSize int) ([]dtos.InvoiceResponse, int64, error) {
	invoices, total, err := s.invoiceRepo.FindByUserID(ctx, userID, page, pageSize)
	if err != nil {
		log.Printf("Error fetching invoices for user %d: %v", userID, err)
		return nil, 0, fmt.Errorf("could not retrieve invoices: %w", err)
	}

	var responses []dtos.InvoiceResponse
	for _, inv := range invoices {
		responses = append(responses, mapInvoiceToResponse(&inv))
	}
	return responses, total, nil
}

func (s *invoiceService) GetInvoiceByIDForUser(ctx context.Context, invoiceID, userID uint) (*dtos.InvoiceResponse, error) {
	invoice, err := s.invoiceRepo.FindByID(ctx, invoiceID)
	if err != nil {
		return nil, ErrInvoiceNotFound
	}
	if invoice.UserID != userID {
		return nil, ErrInvoiceAccessDenied
	}
	resp := mapInvoiceToResponse(invoice)
	return &resp, nil
}

func (s *invoiceService) GetReceiptPathForUser(ctx context.Context, invoiceID, userID uint) (string, string, error) {
	invoice, err := s.invoiceRepo.FindByID(ctx, invoiceID)
	if err != nil {
		return "", "", ErrInvoiceNotFound
	}
	if invoice.UserID != userID {
		return "", "", ErrInvoiceAccessDenied
	}
	if invoice.DisbursementReceiptPath == nil || *invoice.DisbursementReceiptPath == "" {
		return "", "", ErrReceiptNotFound
	}

	// Ensure fileService is not nil before calling
	if s.fileService == nil {
		log.Println("FileService is not initialized in GetReceiptPathForUser.")
		return "", "", errors.New("file service unavailable")
	}
	absPath, err := s.fileService.GetAbsPath(*invoice.DisbursementReceiptPath)
	if err != nil {
		return "", "", fmt.Errorf("receipt file path error: %w", err)
	}
	fileName := filepath.Base(*invoice.DisbursementReceiptPath)
	return absPath, fileName, nil
}
