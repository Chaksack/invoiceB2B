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
	"strconv"
	"time"

	"gorm.io/gorm"
)

// --- Error Definitions (assuming these are or will be defined centrally) ---
var (
	ErrStaffNotFound       = errors.New("staff not found")
	ErrKYCNotFound         = errors.New("kyc record not found for user")
	ErrPDFGenerationFailed = errors.New("failed to generate PDF")
	ErrServiceNotAvailable = errors.New("a required service is not available")
)

// --- DTO Definitions (conceptual, should be in dtos package) ---

// dtos.AdminDashboardAnalytics
type AdminDashboardAnalytics struct {
	TotalUsers           int64            `json:"totalUsers"`
	KYCStats             KYCAnalytics     `json:"kycStats"`
	InvoiceStats         InvoiceAnalytics `json:"invoiceStats"`
	TotalFinancedAmount  float64          `json:"totalFinancedAmount"`
	TotalDisbursedAmount float64          `json:"totalDisbursedAmount"`
	// TotalRepaidAmount could be sum of transactions of type REPAYMENT
}

type KYCAnalytics struct {
	TotalSubmitted int64 `json:"totalSubmitted"`
	TotalApproved  int64 `json:"totalApproved"`
	TotalRejected  int64 `json:"totalRejected"`
	PendingReview  int64 `json:"pendingReview"`
}

type InvoiceAnalytics struct {
	TotalInvoices    int64 `json:"totalInvoices"`
	PendingApproval  int64 `json:"pendingApproval"`  // Status = PENDING_APPROVAL (or PENDING_REVIEW)
	Approved         int64 `json:"approved"`         // Status = APPROVED
	Rejected         int64 `json:"rejected"`         // Status = REJECTED
	Disbursed        int64 `json:"disbursed"`        // Status = DISBURSED
	RepaymentPending int64 `json:"repaymentPending"` // Status = REPAYMENT_PENDING
	Repaid           int64 `json:"repaid"`           // Status = REPAID
	Overdue          int64 `json:"overdue"`          // Status = DISBURSED and DueDate < Now
}

// dtos.InvoicePDFResponse
type InvoicePDFResponse struct {
	FilePath    string `json:"filePath"`
	FileName    string `json:"fileName"`
	DownloadURL string `json:"downloadUrl"`
}

// --- Service Interfaces (including new PDFService) ---

// PDFService interface (placeholder for actual PDF generation logic)
type PDFService interface {
	GenerateInvoicePDF(ctx context.Context, invoice *models.Invoice, user *models.User, kycDetails *models.KYCDetail, companyLogoPath string) (outputFilePath string, fileName string, err error)
}

// AdminService interface definition (updated)
type AdminService interface {
	// User & KYC Management
	GetAllUsers(ctx context.Context, page, pageSize int) ([]dtos.UserResponse, int64, error)
	GetUserByID(ctx context.Context, userID uint) (*dtos.UserResponse, error)
	GetUserKYCDetail(ctx context.Context, userID uint) (*dtos.AdminKYCDetailResponse, error)
	ReviewKYC(ctx context.Context, userID, reviewerStaffID uint, req dtos.AdminKYCReviewRequest) (*dtos.AdminKYCDetailResponse, error)

	// Invoice Management
	GetAllInvoices(ctx context.Context, page, pageSize int, statusFilter string) ([]dtos.InvoiceResponse, int64, error)
	GetInvoiceDetail(ctx context.Context, invoiceID uint) (*dtos.InvoiceResponse, error)
	UpdateInvoiceStatus(ctx context.Context, invoiceID, adminStaffID uint, req dtos.AdminInvoiceUpdateRequest) (*dtos.InvoiceResponse, error)
	UploadDisbursementReceipt(ctx context.Context, invoiceID, adminStaffID uint, req dtos.AdminUploadReceiptRequest) (*dtos.InvoiceResponse, error)
	DownloadInvoicePDF(ctx context.Context, invoiceID, adminStaffID uint) (*InvoicePDFResponse, error)

	// Staff Management
	CreateStaff(ctx context.Context, req dtos.CreateStaffRequest) (*dtos.StaffResponse, error)
	GetStaffByID(ctx context.Context, staffID uint) (*dtos.StaffResponse, error) // New method for GetAdminProfile
	GetAllStaff(ctx context.Context, page, pageSize int) ([]dtos.StaffResponse, int64, error)
	UpdateStaff(ctx context.Context, staffID uint, req dtos.UpdateStaffRequest) (*dtos.StaffResponse, error)
	DeleteStaff(ctx context.Context, staffID uint) error

	// Activity Logs & Analytics
	GetActivityLogs(ctx context.Context, page, pageSize int, filters map[string]string) ([]dtos.ActivityLogResponse, int64, error)
	GetUserActivityLogs(ctx context.Context, userID uint, page, pageSize int, filters map[string]string) ([]dtos.ActivityLogResponse, int64, error)
	GetAdminDashboardAnalytics(ctx context.Context) (*AdminDashboardAnalytics, error)
}

type adminService struct {
	userRepo        repositories.UserRepository
	kycRepo         repositories.KYCRepository
	staffRepo       repositories.StaffRepository
	invoiceRepo     repositories.InvoiceRepository
	transactionRepo repositories.TransactionRepository
	activityLogSvc  ActivityLogService
	emailService    EmailService
	notificationSvc NotificationService
	fileService     FileService
	pdfService      PDFService
	cfg             *config.Config
}

func NewAdminService(
	userRepo repositories.UserRepository,
	kycRepo repositories.KYCRepository,
	staffRepo repositories.StaffRepository,
	invoiceRepo repositories.InvoiceRepository,
	transactionRepo repositories.TransactionRepository,
	activityLogSvc ActivityLogService,
	emailService EmailService,
	notificationSvc NotificationService,
	fileService FileService,
	pdfService PDFService,
	cfg *config.Config,
) AdminService {
	return &adminService{
		userRepo:        userRepo,
		kycRepo:         kycRepo,
		staffRepo:       staffRepo,
		invoiceRepo:     invoiceRepo,
		transactionRepo: transactionRepo,
		activityLogSvc:  activityLogSvc,
		emailService:    emailService,
		notificationSvc: notificationSvc,
		fileService:     fileService,
		pdfService:      pdfService,
		cfg:             cfg,
	}
}

// --- Helper Functions (existing ones renamed with 'local' prefix) ---
func localMapModelUserToUserResponse(user *models.User) dtos.UserResponse {
	return dtos.UserResponse{
		ID:           user.ID,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		CompanyName:  user.CompanyName,
		IsActive:     user.IsActive,
		TwoFAEnabled: user.TwoFAEnabled,
	}
}

func localMapModelKYCToAdminResponse(kyc *models.KYCDetail, userEmail string, reviewerEmail *string) dtos.AdminKYCDetailResponse {
	resp := dtos.AdminKYCDetailResponse{
		ID:              kyc.ID,
		UserID:          kyc.UserID,
		UserEmail:       userEmail,
		Status:          kyc.Status,
		SubmittedAt:     kyc.SubmittedAt,
		ReviewedAt:      kyc.ReviewedAt,
		RejectionReason: kyc.RejectionReason,
		DocumentsInfo:   kyc.DocumentsInfo,
		CreatedAt:       kyc.CreatedAt,
		UpdatedAt:       kyc.UpdatedAt,
	}
	if reviewerEmail != nil {
		resp.ReviewedByEmail = *reviewerEmail
	}
	return resp
}

func localMapModelStaffToStaffResponse(staff *models.Staff) dtos.StaffResponse {
	return dtos.StaffResponse{
		ID:          staff.ID,
		Email:       staff.Email,
		FirstName:   staff.FirstName,
		LastName:    staff.LastName,
		Role:        staff.Role,
		IsActive:    staff.IsActive,
		LastLoginAt: staff.LastLoginAt,
		CreatedAt:   staff.CreatedAt,
		UpdatedAt:   staff.UpdatedAt,
	}
}

func localMapModelActivityLogToResponse(logEntry *models.ActivityLog, userEmail, staffEmail *string) dtos.ActivityLogResponse {
	return dtos.ActivityLogResponse{
		ID:         logEntry.ID,
		StaffEmail: staffEmail,
		UserEmail:  userEmail,
		Action:     logEntry.Action,
		Details:    logEntry.Details,
		IPAddress:  logEntry.IPAddress,
		Timestamp:  logEntry.CreatedAt,
	}
}

func localMapInvoiceToResponse(invoice *models.Invoice) dtos.InvoiceResponse {
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

// --- User & KYC Management (with GetUserByID) ---
func (s *adminService) GetAllUsers(ctx context.Context, page, pageSize int) ([]dtos.UserResponse, int64, error) {
	users, total, err := s.userRepo.FindAllWithPagination(ctx, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all users: %w", err)
	}
	var userResponses []dtos.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, localMapModelUserToUserResponse(&user))
	}
	return userResponses, total, nil
}

func (s *adminService) GetUserByID(ctx context.Context, userID uint) (*dtos.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // Assuming userRepo.FindByID returns gorm.ErrRecordNotFound
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user by ID %d: %w", userID, err)
	}
	resp := localMapModelUserToUserResponse(user)
	return &resp, nil
}

func (s *adminService) GetUserKYCDetail(ctx context.Context, userID uint) (*dtos.AdminKYCDetailResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user for KYC detail: %w", err)
	}
	kyc, err := s.kycRepo.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrKYCNotFound
		}
		return nil, fmt.Errorf("failed to find KYC by user ID %d: %w", userID, err)
	}

	var reviewerEmail *string
	if kyc.ReviewedByID != nil {
		reviewer, staffErr := s.staffRepo.FindByID(ctx, *kyc.ReviewedByID)
		if staffErr != nil {
			log.Printf("Warning: Reviewer staff ID %d for KYC %d not found: %v", *kyc.ReviewedByID, kyc.ID, staffErr)
		} else if reviewer != nil {
			reviewerEmail = &reviewer.Email
		}
	}

	resp := localMapModelKYCToAdminResponse(kyc, user.Email, reviewerEmail)
	return &resp, nil
}

func (s *adminService) ReviewKYC(ctx context.Context, userID, reviewerStaffID uint, req dtos.AdminKYCReviewRequest) (*dtos.AdminKYCDetailResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user for KYC review: %w", err)
	}
	kyc, err := s.kycRepo.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrKYCNotFound
		}
		return nil, fmt.Errorf("failed to find KYC for review by user ID %d: %w", userID, err)
	}
	reviewer, err := s.staffRepo.FindByID(ctx, reviewerStaffID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStaffNotFound // Use specific error for staff
		}
		return nil, fmt.Errorf("failed to find reviewer staff by ID %d: %w", reviewerStaffID, err)
	}

	oldStatus := kyc.Status
	kyc.Status = req.Status
	now := time.Now()
	kyc.ReviewedAt = &now
	kyc.ReviewedByID = &reviewerStaffID
	if req.Status == models.KYCRejected {
		if req.RejectionReason == nil || *req.RejectionReason == "" {
			return nil, errors.New("rejection reason is required for rejected KYC")
		}
		kyc.RejectionReason = req.RejectionReason
	} else {
		kyc.RejectionReason = nil
	}

	updatedKYC, err := s.kycRepo.CreateOrUpdate(ctx, kyc)
	if err != nil {
		return nil, fmt.Errorf("failed to update KYC record: %w", err)
	}

	_ = s.activityLogSvc.LogActivity(ctx, &reviewerStaffID, &userID, "ADMIN_KYC_REVIEWED",
		map[string]interface{}{"kyc_id": kyc.ID, "old_status": oldStatus, "new_status": kyc.Status, "reason": kyc.RejectionReason}, "")

	go func() {
		subject := fmt.Sprintf("Your KYC Application Status: %s", kyc.Status)
		body := fmt.Sprintf("Dear %s,\n\nYour KYC application status has been updated to: %s.", user.FirstName, kyc.Status)
		if kyc.Status == models.KYCRejected && kyc.RejectionReason != nil {
			body += fmt.Sprintf("\nReason: %s", *kyc.RejectionReason)
		}
		body += "\n\nPlease log in to your account for more details.\n\nRegards,\nThe Admin Team"
		if err := s.emailService.SendEmail(user.Email, subject, body); err != nil {
			log.Printf("Failed to send KYC status email to %s: %v", user.Email, err)
		}
	}()

	if s.notificationSvc != nil && s.cfg != nil {
		_ = s.notificationSvc.PublishEvent(s.cfg.RabbitMQEventExchangeName, s.cfg.RabbitMQKYCStatusUpdatedRoutingKey,
			map[string]interface{}{"user_id": userID, "kyc_id": kyc.ID, "status": kyc.Status, "rejection_reason": kyc.RejectionReason})
	}

	resp := localMapModelKYCToAdminResponse(updatedKYC, user.Email, &reviewer.Email)
	return &resp, nil
}

// --- Invoice Management (with DownloadInvoicePDF) ---
func (s *adminService) GetAllInvoices(ctx context.Context, page, pageSize int, statusFilter string) ([]dtos.InvoiceResponse, int64, error) {
	filters := make(map[string]string)
	if statusFilter != "" {
		filters["status"] = statusFilter
	}

	invoices, total, err := s.invoiceRepo.FindAll(ctx, page, pageSize, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all invoices: %w", err)
	}
	var responses []dtos.InvoiceResponse
	for _, inv := range invoices {
		responses = append(responses, localMapInvoiceToResponse(&inv))
	}
	return responses, total, nil
}

func (s *adminService) GetInvoiceDetail(ctx context.Context, invoiceID uint) (*dtos.InvoiceResponse, error) {
	invoice, err := s.invoiceRepo.FindByIDWithRelations(ctx, invoiceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvoiceNotFound
		}
		return nil, fmt.Errorf("failed to find invoice by ID %d: %w", invoiceID, err)
	}
	resp := localMapInvoiceToResponse(invoice)
	return &resp, nil
}

func (s *adminService) UpdateInvoiceStatus(ctx context.Context, invoiceID, adminStaffID uint, req dtos.AdminInvoiceUpdateRequest) (*dtos.InvoiceResponse, error) {
	invoice, err := s.invoiceRepo.FindByID(ctx, invoiceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvoiceNotFound
		}
		return nil, fmt.Errorf("failed to find invoice for status update: %w", err)
	}
	user, userErr := s.userRepo.FindByID(ctx, invoice.UserID)
	if userErr != nil {
		log.Printf("Warning: User for invoice %d not found during status update: %v", invoiceID, userErr)
	}

	oldStatus := invoice.Status
	invoice.Status = req.Status
	now := time.Now()

	logDetails := map[string]interface{}{
		"invoice_id":     invoice.ID,
		"old_status":     oldStatus,
		"new_status":     invoice.Status,
		"admin_staff_id": adminStaffID,
	}

	switch req.Status {
	case models.InvoiceApproved:
		invoice.ApprovedByID = &adminStaffID
		invoice.ApprovedAt = &now
		logDetails["approved_by_staff_id"] = adminStaffID
	case models.InvoiceRejected:
		invoice.ApprovedByID = &adminStaffID
		invoice.ApprovedAt = &now
		if req.RejectionReason == nil || *req.RejectionReason == "" {
			return nil, errors.New("rejection reason is required for rejected invoice")
		}
		invoice.ProcessingError = req.RejectionReason
		logDetails["processing_error"] = *req.RejectionReason
	case models.InvoiceDisbursed:
		if oldStatus != models.InvoiceApproved {
			return nil, ErrInvoiceNotApprovedForDisbursement
		}
		if req.FinancedAmount == nil || *req.FinancedAmount <= 0 {
			return nil, errors.New("valid financed amount is required for disbursement")
		}
		if req.FinancingFeePercentage == nil || *req.FinancingFeePercentage < 0 {
			return nil, errors.New("valid financing fee percentage is required for disbursement")
		}
		invoice.DisbursedByID = &adminStaffID
		invoice.DisbursedAt = &now
		invoice.FinancingFeePercentage = req.FinancingFeePercentage
		invoice.FinancedAmount = req.FinancedAmount
		logDetails["financed_amount"] = *req.FinancedAmount
		logDetails["financing_fee_percentage"] = *req.FinancingFeePercentage

		disbursementTx := &models.Transaction{
			InvoiceID:       invoiceID,
			Type:            models.TransactionDisbursement,
			Amount:          *req.FinancedAmount,
			TransactionDate: now,
		}
		if err := s.transactionRepo.Create(ctx, disbursementTx); err != nil {
			log.Printf("Error creating disbursement transaction for invoice %d: %v", invoiceID, err)
			return nil, fmt.Errorf("failed to record disbursement transaction: %w", err)
		}
		logDetails["disbursement_transaction_id"] = disbursementTx.ID
	case models.InvoiceRepaid:
		if oldStatus != models.InvoiceDisbursed && oldStatus != models.InvoiceRepaymentPending {
			return nil, ErrInvalidInvoiceStatusForOperation
		}
		logDetails["repayment_confirmed_by_staff_id"] = adminStaffID
	default:
		return nil, fmt.Errorf("invalid status '%s' for admin update", req.Status)
	}

	if err := s.invoiceRepo.Update(ctx, invoice); err != nil {
		return nil, fmt.Errorf("failed to update invoice: %w", err)
	}

	_ = s.activityLogSvc.LogActivity(ctx, &adminStaffID, &invoice.UserID, "ADMIN_INVOICE_STATUS_UPDATE", logDetails, "")

	if user != nil {
		go func() {
			subject := fmt.Sprintf("Invoice #%s Status Update: %s", invoice.InvoiceNumber, invoice.Status)
			body := fmt.Sprintf("Dear %s,\n\nYour invoice #%s (Amount: %.2f %s) has been updated to: %s.",
				user.FirstName, invoice.InvoiceNumber, invoice.Amount, invoice.Currency, invoice.Status)

			if invoice.Status == models.InvoiceRejected && req.RejectionReason != nil {
				body += fmt.Sprintf("\nReason: %s", *req.RejectionReason)
			}
			if invoice.Status == models.InvoiceDisbursed {
				body += fmt.Sprintf("\nFinanced Amount: %.2f %s.", *invoice.FinancedAmount, invoice.Currency)
				if invoice.DisbursementReceiptPath != nil && *invoice.DisbursementReceiptPath != "" {
					body += "\nA disbursement receipt is available for viewing/download from your dashboard."
				}
			}
			body += "\n\nPlease log in to your account for more details.\n\nRegards,\nThe Admin Team"
			if err := s.emailService.SendEmail(user.Email, subject, body); err != nil {
				log.Printf("Failed to send invoice status email to %s: %v", user.Email, err)
			}
		}()
	}

	if s.notificationSvc != nil && s.cfg != nil {
		eventPayload := map[string]interface{}{
			"user_id":    invoice.UserID,
			"invoice_id": invoice.ID,
			"status":     invoice.Status,
		}
		if req.Status == models.InvoiceRejected && req.RejectionReason != nil {
			eventPayload["rejection_reason"] = *req.RejectionReason
		}
		_ = s.notificationSvc.PublishEvent(s.cfg.RabbitMQEventExchangeName, s.cfg.RabbitMQInvoiceStatusUpdatedRoutingKey, eventPayload)
	}

	updatedInvoice, _ := s.invoiceRepo.FindByIDWithRelations(ctx, invoiceID)
	resp := localMapInvoiceToResponse(updatedInvoice)
	return &resp, nil
}

func (s *adminService) UploadDisbursementReceipt(ctx context.Context, invoiceID, adminStaffID uint, req dtos.AdminUploadReceiptRequest) (*dtos.InvoiceResponse, error) {
	invoice, err := s.invoiceRepo.FindByID(ctx, invoiceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvoiceNotFound
		}
		return nil, fmt.Errorf("failed to find invoice for receipt upload: %w", err)
	}

	if s.fileService == nil {
		log.Println("AdminService: FileService is not initialized.")
		return nil, ErrServiceNotAvailable
	}

	relativePath, uniqueFileName, err := s.fileService.SaveFile(req.File, filepath.Join("receipts", "disbursements"))
	if err != nil {
		return nil, fmt.Errorf("failed to save receipt file: %w", err)
	}

	invoice.DisbursementReceiptPath = &relativePath
	if err := s.invoiceRepo.Update(ctx, invoice); err != nil {
		return nil, fmt.Errorf("failed to update invoice with receipt path: %w", err)
	}

	_ = s.activityLogSvc.LogActivity(ctx, &adminStaffID, &invoice.UserID, "ADMIN_UPLOAD_DISBURSEMENT_RECEIPT",
		map[string]interface{}{"invoice_id": invoice.ID, "receipt_path": relativePath, "file_name": uniqueFileName}, "")

	user, userErr := s.userRepo.FindByID(ctx, invoice.UserID)
	if userErr == nil && user != nil {
		go func() {
			subject := fmt.Sprintf("Disbursement Receipt Uploaded for Invoice #%s", invoice.InvoiceNumber)
			body := fmt.Sprintf("Dear %s,\n\nA disbursement receipt has been uploaded for your invoice #%s (Amount: %.2f %s). You can view or download it from your dashboard.\n\nRegards,\nThe Admin Team", user.FirstName, invoice.InvoiceNumber, invoice.Amount, invoice.Currency)

			absAttachmentPath, pathErr := s.fileService.GetAbsPath(relativePath)
			if pathErr != nil {
				log.Printf("Failed to get absolute path for receipt attachment %s: %v. Sending email without attachment.", relativePath, pathErr)
				if emailErr := s.emailService.SendEmail(user.Email, subject, body); emailErr != nil {
					log.Printf("Failed to send receipt upload notification email (no attachment) to %s: %v", user.Email, emailErr)
				}
				return
			}

			if emailErr := s.emailService.SendEmailWithAttachment(user.Email, subject, body, absAttachmentPath, uniqueFileName); emailErr != nil {
				log.Printf("Failed to send receipt upload notification email with attachment to %s: %v", user.Email, emailErr)
			}
		}()
	} else {
		log.Printf("Warning: User for invoice %d not found, cannot send receipt upload email: %v", invoiceID, userErr)
	}
	updatedInvoice, _ := s.invoiceRepo.FindByIDWithRelations(ctx, invoiceID)
	resp := localMapInvoiceToResponse(updatedInvoice)
	return &resp, nil
}

func (s *adminService) DownloadInvoicePDF(ctx context.Context, invoiceID, adminStaffID uint) (*InvoicePDFResponse, error) {
	if s.pdfService == nil {
		log.Println("AdminService: PDFService is not initialized.")
		return nil, ErrServiceNotAvailable
	}
	if s.fileService == nil {
		log.Println("AdminService: FileService is not initialized for PDF storage.")
	}

	invoice, err := s.invoiceRepo.FindByIDWithRelations(ctx, invoiceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvoiceNotFound
		}
		return nil, fmt.Errorf("failed to find invoice %d for PDF generation: %w", invoiceID, err)
	}

	user, err := s.userRepo.FindByID(ctx, invoice.UserID)
	if err != nil {
		log.Printf("Warning: User %d for invoice %d not found during PDF generation: %v", invoice.UserID, invoiceID, err)
	}

	kycDetails, err := s.kycRepo.FindByUserID(ctx, invoice.UserID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Warning: KYC details for user %d not found during PDF generation for invoice %d: %v", invoice.UserID, invoiceID, err)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		kycDetails = nil
	}

	companyLogoPath := "" // e.g., s.cfg.CompanyLogoPathForPDF

	generatedFilePath, generatedFileName, err := s.pdfService.GenerateInvoicePDF(ctx, invoice, user, kycDetails, companyLogoPath)
	if err != nil {
		log.Printf("Error generating PDF for invoice %d by staff %d: %v", invoiceID, adminStaffID, err)
		return nil, fmt.Errorf("%w: %v", ErrPDFGenerationFailed, err)
	}

	_ = s.activityLogSvc.LogActivity(ctx, &adminStaffID, &invoice.UserID, "ADMIN_DOWNLOAD_INVOICE_PDF",
		map[string]interface{}{"invoice_id": invoice.ID, "generated_file_path": generatedFilePath, "file_name": generatedFileName}, "")

	pdfResponse := &InvoicePDFResponse{
		FilePath: generatedFilePath,
		FileName: generatedFileName,
	}
	if s.fileService != nil {
		// Example: publicURL, urlErr := s.fileService.GetPublicURL(generatedFilePath)
	}

	return pdfResponse, nil
}

// --- Staff Management ---
func (s *adminService) CreateStaff(ctx context.Context, req dtos.CreateStaffRequest) (*dtos.StaffResponse, error) {
	existing, _ := s.staffRepo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("staff with this email already exists")
	}
	staff := &models.Staff{
		Email:        req.Email,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		PasswordHash: req.Password,
		Role:         req.Role,
		IsActive:     true,
	}
	if err := s.staffRepo.Create(ctx, staff); err != nil {
		return nil, fmt.Errorf("failed to create staff: %w", err)
	}

	resp := localMapModelStaffToStaffResponse(staff)
	return &resp, nil
}

// GetStaffByID retrieves a specific staff member by their ID.
func (s *adminService) GetStaffByID(ctx context.Context, staffID uint) (*dtos.StaffResponse, error) {
	staff, err := s.staffRepo.FindByID(ctx, staffID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // Assuming staffRepo.FindByID returns gorm.ErrRecordNotFound
			return nil, ErrStaffNotFound // Use specific error
		}
		return nil, fmt.Errorf("failed to find staff by ID %d: %w", staffID, err)
	}
	resp := localMapModelStaffToStaffResponse(staff)
	return &resp, nil
}

func (s *adminService) GetAllStaff(ctx context.Context, page, pageSize int) ([]dtos.StaffResponse, int64, error) {
	staffList, total, err := s.staffRepo.FindAll(ctx, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all staff: %w", err)
	}
	var responses []dtos.StaffResponse
	for _, staff := range staffList {
		responses = append(responses, localMapModelStaffToStaffResponse(&staff))
	}
	return responses, total, nil
}

func (s *adminService) UpdateStaff(ctx context.Context, staffID uint, req dtos.UpdateStaffRequest) (*dtos.StaffResponse, error) {
	staff, err := s.staffRepo.FindByID(ctx, staffID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStaffNotFound
		}
		return nil, fmt.Errorf("failed to find staff for update: %w", err)
	}
	if req.FirstName != nil {
		staff.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		staff.LastName = *req.LastName
	}
	if req.Role != nil {
		staff.Role = *req.Role
	}
	if req.IsActive != nil {
		staff.IsActive = *req.IsActive
	}

	if err := s.staffRepo.Update(ctx, staff); err != nil {
		return nil, fmt.Errorf("failed to update staff: %w", err)
	}
	resp := localMapModelStaffToStaffResponse(staff)
	return &resp, nil
}

func (s *adminService) DeleteStaff(ctx context.Context, staffID uint) error {
	// First, check if staff exists to return a more specific error
	_, err := s.staffRepo.FindByID(ctx, staffID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrStaffNotFound
		}
		return fmt.Errorf("failed to find staff for deletion check: %w", err)
	}

	if err := s.staffRepo.Delete(ctx, staffID); err != nil {
		return fmt.Errorf("failed to delete staff %d: %w", staffID, err)
	}
	log.Printf("Staff member %d deleted.", staffID)
	return nil
}

// --- Activity Logs & Analytics ---

func (s *adminService) GetActivityLogs(ctx context.Context, page, pageSize int, filters map[string]string) ([]dtos.ActivityLogResponse, int64, error) {
	logs, total, err := s.activityLogSvc.GetActivityLogs(ctx, page, pageSize, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get activity logs: %w", err)
	}

	var responses []dtos.ActivityLogResponse
	for _, logEntry := range logs {
		var userEmail, staffEmail *string
		if logEntry.UserID != nil {
			user, _ := s.userRepo.FindByID(ctx, *logEntry.UserID)
			if user != nil {
				userEmail = &user.Email
			}
		}
		if logEntry.StaffID != nil {
			staff, _ := s.staffRepo.FindByID(ctx, *logEntry.StaffID)
			if staff != nil {
				staffEmail = &staff.Email
			}
		}
		responses = append(responses, localMapModelActivityLogToResponse(&logEntry, userEmail, staffEmail))
	}
	return responses, total, nil
}

func (s *adminService) GetUserActivityLogs(ctx context.Context, userID uint, page, pageSize int, inputFilters map[string]string) ([]dtos.ActivityLogResponse, int64, error) {
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, ErrUserNotFound
		}
		return nil, 0, fmt.Errorf("failed to verify user for activity logs: %w", err)
	}

	currentFilters := make(map[string]string)
	if inputFilters != nil {
		for k, v := range inputFilters {
			currentFilters[k] = v
		}
	}
	currentFilters["user_id"] = strconv.FormatUint(uint64(userID), 10)

	logs, total, err := s.activityLogSvc.GetActivityLogs(ctx, page, pageSize, currentFilters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get activity logs for user %d: %w", userID, err)
	}

	var responses []dtos.ActivityLogResponse
	for _, logEntry := range logs {
		var userEmail, staffEmail *string
		if logEntry.UserID != nil && *logEntry.UserID == userID {
			user, userFetchErr := s.userRepo.FindByID(ctx, *logEntry.UserID)
			if userFetchErr == nil && user != nil {
				userEmail = &user.Email
			} else {
				tempEmail := fmt.Sprintf("user_id_%d", *logEntry.UserID)
				userEmail = &tempEmail
			}
		}

		if logEntry.StaffID != nil {
			staff, _ := s.staffRepo.FindByID(ctx, *logEntry.StaffID)
			if staff != nil {
				staffEmail = &staff.Email
			}
		}
		responses = append(responses, localMapModelActivityLogToResponse(&logEntry, userEmail, staffEmail))
	}
	return responses, total, nil
}

func (s *adminService) GetAdminDashboardAnalytics(ctx context.Context) (*AdminDashboardAnalytics, error) {
	analytics := AdminDashboardAnalytics{}
	var err error

	analytics.TotalUsers, err = s.userRepo.CountAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get total users count: %w", err)
	}

	analytics.KYCStats.TotalSubmitted, err = s.kycRepo.CountByStatus(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get total submitted KYC count: %w", err)
	}
	analytics.KYCStats.TotalApproved, err = s.kycRepo.CountByStatus(ctx, models.KYCApproved)
	if err != nil {
		return nil, fmt.Errorf("failed to get approved KYC count: %w", err)
	}
	analytics.KYCStats.TotalRejected, err = s.kycRepo.CountByStatus(ctx, models.KYCRejected)
	if err != nil {
		return nil, fmt.Errorf("failed to get rejected KYC count: %w", err)
	}
	analytics.KYCStats.PendingReview, err = s.kycRepo.CountByStatus(ctx, models.KYCPending)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending review KYC count: %w", err)
	}

	analytics.InvoiceStats.TotalInvoices, err = s.invoiceRepo.CountByStatus(ctx, "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get total invoices count: %w", err)
	}
	analytics.InvoiceStats.PendingApproval, err = s.invoiceRepo.CountByStatus(ctx, models.InvoicePendingApproval, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending approval invoice count: %w", err)
	}
	analytics.InvoiceStats.Approved, err = s.invoiceRepo.CountByStatus(ctx, models.InvoiceApproved, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get approved invoice count: %w", err)
	}
	analytics.InvoiceStats.Rejected, err = s.invoiceRepo.CountByStatus(ctx, models.InvoiceRejected, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get rejected invoice count: %w", err)
	}
	analytics.InvoiceStats.Disbursed, err = s.invoiceRepo.CountByStatus(ctx, models.InvoiceDisbursed, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get disbursed invoice count: %w", err)
	}
	analytics.InvoiceStats.RepaymentPending, err = s.invoiceRepo.CountByStatus(ctx, models.InvoiceRepaymentPending, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get repayment pending invoice count: %w", err)
	}
	analytics.InvoiceStats.Repaid, err = s.invoiceRepo.CountByStatus(ctx, models.InvoiceRepaid, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get repaid invoice count: %w", err)
	}
	analytics.InvoiceStats.Overdue, err = s.invoiceRepo.CountOverdue(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue invoice count: %w", err)
	}

	analytics.TotalFinancedAmount, err = s.invoiceRepo.SumAmountByStatus(ctx,
		[]models.InvoiceStatus{models.InvoiceApproved, models.InvoiceDisbursed, models.InvoiceRepaymentPending, models.InvoiceRepaid},
		"amount")
	if err != nil {
		return nil, fmt.Errorf("failed to get total financed amount: %w", err)
	}
	analytics.TotalDisbursedAmount, err = s.invoiceRepo.SumAmountByStatus(ctx,
		[]models.InvoiceStatus{models.InvoiceDisbursed, models.InvoiceRepaymentPending, models.InvoiceRepaid},
		"financed_amount")
	if err != nil {
		return nil, fmt.Errorf("failed to get total disbursed amount: %w", err)
	}

	return &analytics, nil
}
