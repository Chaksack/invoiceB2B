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
	if user.KYCDetail == nil || user.KYCDetail.Status != models.KYCApproved {
		return nil, ErrKYCNotApprovedForInvoiceUpload
	}

	relativePath, _, err := s.fileService.SaveFile(req.File, "invoices")
	if err != nil {
		log.Printf("Error saving invoice file for user %d: %v", userID, err)
		return nil, fmt.Errorf("failed to save invoice file: %w", err)
	}

	now := time.Now()
	invoice := &models.Invoice{
		UserID:           userID,
		Status:           models.InvoicePendingReview,
		OriginalFilePath: relativePath,
		UploadedAt:       now,
	}

	if err := s.invoiceRepo.Create(ctx, invoice); err != nil {
		log.Printf("Error creating invoice record for user %d: %v", userID, err)
		return nil, fmt.Errorf("failed to create invoice record: %w", err)
	}

	eventPayload := map[string]interface{}{
		"invoice_id":        invoice.ID,
		"user_id":           userID,
		"file_path":         relativePath,
		"original_filename": req.File.Filename,
	}
	if s.notificationSvc != nil {
		err = s.notificationSvc.PublishEvent(
			s.cfg.RabbitMQEventExchangeName,
			s.cfg.RabbitMQInvoiceUploadedRoutingKey,
			eventPayload,
		)
		if err != nil {
			log.Printf("Failed to publish invoice.uploaded event for invoice %d: %v", invoice.ID, err)
		}
	}

	_ = s.activityLogSvc.LogActivity(ctx, nil, &userID, "INVOICE_UPLOADED",
		map[string]interface{}{"invoice_id": invoice.ID, "filename": req.File.Filename}, "")

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

	absPath, err := s.fileService.GetAbsPath(*invoice.DisbursementReceiptPath)
	if err != nil {
		return "", "", fmt.Errorf("receipt file path error: %w", err)
	}
	fileName := filepath.Base(*invoice.DisbursementReceiptPath)
	return absPath, fileName, nil
}
