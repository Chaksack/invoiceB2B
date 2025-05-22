package services

import (
	"context"
	"errors"
	"fmt"
	"invoiceB2B/internal/config" // Added for cfg access
	"invoiceB2B/internal/dtos"
	"invoiceB2B/internal/models"
	"invoiceB2B/internal/repositories"
	"log"
	"time"

	"gorm.io/gorm"
)

// AdminService interface definition
type AdminService interface {
	GetAllUsers(ctx context.Context, page, pageSize int) ([]dtos.UserResponse, int64, error)
	GetUserKYCDetail(ctx context.Context, userID uint) (*dtos.AdminKYCDetailResponse, error)
	ReviewKYC(ctx context.Context, userID, reviewerStaffID uint, req dtos.AdminKYCReviewRequest) (*dtos.AdminKYCDetailResponse, error)
	GetAllInvoices(ctx context.Context, page, pageSize int, statusFilter string) ([]dtos.InvoiceResponse, int64, error)
	GetInvoiceDetail(ctx context.Context, invoiceID uint) (*dtos.InvoiceResponse, error)
	UpdateInvoiceStatus(ctx context.Context, invoiceID, adminStaffID uint, req dtos.AdminInvoiceUpdateRequest) (*dtos.InvoiceResponse, error)
	UploadDisbursementReceipt(ctx context.Context, invoiceID, adminStaffID uint, req dtos.AdminUploadReceiptRequest) (*dtos.InvoiceResponse, error)
	CreateStaff(ctx context.Context, req dtos.CreateStaffRequest) (*dtos.StaffResponse, error)
	GetAllStaff(ctx context.Context, page, pageSize int) ([]dtos.StaffResponse, int64, error)
	UpdateStaff(ctx context.Context, staffID uint, req dtos.UpdateStaffRequest) (*dtos.StaffResponse, error)
	DeleteStaff(ctx context.Context, staffID uint) error
	GetActivityLogs(ctx context.Context, page, pageSize int) ([]dtos.ActivityLogResponse, int64, error)
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
		cfg:             cfg,
	}
}

// Helper functions mapModelUserToUserResponse, mapModelKYCToAdminResponse, etc.
func localMapModelUserToUserResponse(user *models.User) dtos.UserResponse { // Renamed to avoid conflict if defined elsewhere
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

func localMapModelKYCToAdminResponse(kyc *models.KYCDetail, userEmail string, reviewerEmail *string) dtos.AdminKYCDetailResponse { // Renamed
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

func localMapModelStaffToStaffResponse(staff *models.Staff) dtos.StaffResponse { // Renamed
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

func localMapModelActivityLogToResponse(logEntry *models.ActivityLog, userEmail, staffEmail *string) dtos.ActivityLogResponse { // Renamed
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

func localMapInvoiceToResponse(invoice *models.Invoice) dtos.InvoiceResponse { // Renamed
	var receiptPath string
	if invoice.DisbursementReceiptPath != nil {
		receiptPath = *invoice.DisbursementReceiptPath // Dereference pointer
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
		DisbursementReceiptPath: receiptPath, // Use dereferenced value
		CreatedAt:               invoice.CreatedAt,
		UpdatedAt:               invoice.UpdatedAt,
	}
}

// --- User & KYC Management ---
func (s *adminService) GetAllUsers(ctx context.Context, page, pageSize int) ([]dtos.UserResponse, int64, error) {
	users, total, err := s.userRepo.FindAllWithPagination(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	var userResponses []dtos.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, localMapModelUserToUserResponse(&user))
	}
	return userResponses, total, nil
}

func (s *adminService) GetUserKYCDetail(ctx context.Context, userID uint) (*dtos.AdminKYCDetailResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	kyc, err := s.kycRepo.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dtos.AdminKYCDetailResponse{UserID: userID, UserEmail: user.Email, Status: "Not Submitted"}, nil
		}
		return nil, err
	}

	var reviewerEmail *string
	if kyc.ReviewedByID != nil {
		reviewer, _ := s.staffRepo.FindByID(ctx, *kyc.ReviewedByID)
		if reviewer != nil {
			reviewerEmail = &reviewer.Email
		}
	}

	resp := localMapModelKYCToAdminResponse(kyc, user.Email, reviewerEmail)
	return &resp, nil
}

func (s *adminService) ReviewKYC(ctx context.Context, userID, reviewerStaffID uint, req dtos.AdminKYCReviewRequest) (*dtos.AdminKYCDetailResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	kyc, err := s.kycRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("kyc record not found for user")
	}
	reviewer, err := s.staffRepo.FindByID(ctx, reviewerStaffID)
	if err != nil {
		return nil, errors.New("reviewer staff account not found")
	}

	oldStatus := kyc.Status
	kyc.Status = req.Status
	now := time.Now()
	kyc.ReviewedAt = &now
	kyc.ReviewedByID = &reviewerStaffID
	if req.Status == models.KYCRejected {
		kyc.RejectionReason = req.RejectionReason
	} else {
		kyc.RejectionReason = nil
	}

	updatedKYC, err := s.kycRepo.CreateOrUpdate(ctx, kyc)
	if err != nil {
		return nil, err
	}

	_ = s.activityLogSvc.LogActivity(ctx, &reviewerStaffID, &userID, "ADMIN_KYC_REVIEWED",
		map[string]interface{}{"kyc_id": kyc.ID, "old_status": oldStatus, "new_status": kyc.Status, "reason": kyc.RejectionReason}, "")

	go func() {
		subject := fmt.Sprintf("KYC Status Update: %s", kyc.Status)
		body := fmt.Sprintf("Hi %s,\n\nYour KYC status has been updated to: %s.", user.FirstName, kyc.Status)
		if kyc.Status == models.KYCRejected && kyc.RejectionReason != nil {
			body += fmt.Sprintf("\nReason: %s", *kyc.RejectionReason)
		}
		body += "\n\nThanks,\nThe Admin Team"
		if err := s.emailService.SendEmail(user.Email, subject, body); err != nil {
			log.Printf("Failed to send KYC status email to %s: %v", user.Email, err)
		}
	}()
	if s.notificationSvc != nil {
		_ = s.notificationSvc.PublishEvent(s.cfg.RabbitMQEventExchangeName, s.cfg.RabbitMQKYCStatusUpdatedRoutingKey,
			map[string]interface{}{"user_id": userID, "kyc_id": kyc.ID, "status": kyc.Status})
	}

	resp := localMapModelKYCToAdminResponse(updatedKYC, user.Email, &reviewer.Email)
	return &resp, nil
}

// --- Invoice Management ---
func (s *adminService) GetAllInvoices(ctx context.Context, page, pageSize int, statusFilter string) ([]dtos.InvoiceResponse, int64, error) {
	filters := make(map[string]string)
	if statusFilter != "" {
		filters["status"] = statusFilter
	}
	invoices, total, err := s.invoiceRepo.FindAll(ctx, page, pageSize, filters)
	if err != nil {
		return nil, 0, err
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
		return nil, ErrInvoiceNotFound
	}
	resp := localMapInvoiceToResponse(invoice)
	return &resp, nil
}

func (s *adminService) UpdateInvoiceStatus(ctx context.Context, invoiceID, adminStaffID uint, req dtos.AdminInvoiceUpdateRequest) (*dtos.InvoiceResponse, error) {
	invoice, err := s.invoiceRepo.FindByID(ctx, invoiceID)
	if err != nil {
		return nil, ErrInvoiceNotFound
	}
	user, err := s.userRepo.FindByID(ctx, invoice.UserID)
	if err != nil {
		log.Printf("Warning: User for invoice %d not found during status update: %v", invoiceID, err)
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
		logDetails["approved_by"] = adminStaffID
	case models.InvoiceRejected:
		invoice.ApprovedByID = &adminStaffID
		invoice.ApprovedAt = &now
		if req.RejectionReason != nil {
			logDetails["rejection_reason"] = *req.RejectionReason
		}
	case models.InvoiceDisbursed:
		if oldStatus != models.InvoiceApproved {
			return nil, ErrInvoiceNotApprovedForDisbursement
		}
		invoice.DisbursedByID = &adminStaffID
		invoice.DisbursedAt = &now
		if req.FinancingFeePercentage != nil {
			invoice.FinancingFeePercentage = req.FinancingFeePercentage
			logDetails["financing_fee_percentage"] = *req.FinancingFeePercentage
		}
		if req.FinancedAmount != nil {
			invoice.FinancedAmount = req.FinancedAmount
			logDetails["financed_amount"] = *req.FinancedAmount

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
		} else {
			return nil, errors.New("financed amount is required for disbursement")
		}
	case models.InvoiceRepaid:
		if oldStatus != models.InvoiceDisbursed && oldStatus != models.InvoiceRepaymentPending {
			return nil, ErrInvalidInvoiceStatusForOperation
		}
		logDetails["confirmed_by_admin"] = adminStaffID
	default:
		return nil, errors.New("invalid status for admin update")
	}

	if err := s.invoiceRepo.Update(ctx, invoice); err != nil {
		return nil, err
	}

	_ = s.activityLogSvc.LogActivity(ctx, &adminStaffID, &invoice.UserID, "ADMIN_INVOICE_STATUS_UPDATE", logDetails, "")

	if user != nil {
		go func() {
			subject := fmt.Sprintf("Invoice #%s Status Update: %s", invoice.InvoiceNumber, invoice.Status)
			body := fmt.Sprintf("Hi %s,\n\nYour invoice #%s status has been updated to: %s.", user.FirstName, invoice.InvoiceNumber, invoice.Status)
			if invoice.Status == models.InvoiceDisbursed && invoice.DisbursementReceiptPath != nil {
				body += "\nA disbursement receipt is available for viewing/download."
			}
			if invoice.Status == models.InvoiceRejected && req.RejectionReason != nil {
				body += fmt.Sprintf("\nReason: %s", *req.RejectionReason)
			}
			body += "\n\nThanks,\nThe Admin Team"
			if err := s.emailService.SendEmail(user.Email, subject, body); err != nil {
				log.Printf("Failed to send invoice status email to %s: %v", user.Email, err)
			}
		}()
	}
	if s.notificationSvc != nil {
		_ = s.notificationSvc.PublishEvent(s.cfg.RabbitMQEventExchangeName, s.cfg.RabbitMQInvoiceStatusUpdatedRoutingKey,
			map[string]interface{}{"user_id": invoice.UserID, "invoice_id": invoice.ID, "status": invoice.Status})
	}

	resp := localMapInvoiceToResponse(invoice)
	return &resp, nil
}
func (s *adminService) UploadDisbursementReceipt(ctx context.Context, invoiceID, adminStaffID uint, req dtos.AdminUploadReceiptRequest) (*dtos.InvoiceResponse, error) {
	invoice, err := s.invoiceRepo.FindByID(ctx, invoiceID)
	if err != nil {
		return nil, ErrInvoiceNotFound
	}
	if invoice.Status != models.InvoiceDisbursed {
		log.Printf("Warning: Uploading receipt for invoice %d not in 'disbursed' state (current: %s)", invoiceID, invoice.Status)
	}

	if s.fileService == nil {
		log.Println("AdminService: FileService is not initialized.")
		return nil, errors.New("file service not available")
	}

	// SaveFile now returns relativePath and the generated unique filename
	relativePath, uniqueFileName, err := s.fileService.SaveFile(req.File, "receipts")
	if err != nil {
		return nil, fmt.Errorf("failed to save receipt file: %w", err)
	}

	invoice.DisbursementReceiptPath = &relativePath
	if err := s.invoiceRepo.Update(ctx, invoice); err != nil {
		return nil, err
	}

	_ = s.activityLogSvc.LogActivity(ctx, &adminStaffID, &invoice.UserID, "ADMIN_UPLOAD_DISBURSEMENT_RECEIPT",
		map[string]interface{}{"invoice_id": invoice.ID, "receipt_path": relativePath}, "")

	user, _ := s.userRepo.FindByID(ctx, invoice.UserID)
	if user != nil {
		go func() {
			subject := fmt.Sprintf("Disbursement Receipt Uploaded for Invoice #%s", invoice.InvoiceNumber)
			body := fmt.Sprintf("Hi %s,\n\nA disbursement receipt has been uploaded for your invoice #%s. You can view or download it from your dashboard.\n\nThanks,\nThe Admin Team", user.FirstName, invoice.InvoiceNumber)

			// Get absolute path for attachment
			absAttachmentPath, err := s.fileService.GetAbsPath(relativePath)
			if err != nil {
				log.Printf("Failed to get absolute path for receipt attachment %s: %v. Sending email without attachment.", relativePath, err)
				if emailErr := s.emailService.SendEmail(user.Email, subject, body); emailErr != nil {
					log.Printf("Failed to send receipt upload notification email (no attachment) to %s: %v", user.Email, emailErr)
				}
				return
			}

			if emailErr := s.emailService.SendEmailWithAttachment(user.Email, subject, body, absAttachmentPath, uniqueFileName); emailErr != nil {
				log.Printf("Failed to send receipt upload notification email with attachment to %s: %v", user.Email, emailErr)
			}
		}()
	}

	resp := localMapInvoiceToResponse(invoice)
	return &resp, nil
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
		return nil, err
	}

	resp := localMapModelStaffToStaffResponse(staff)
	return &resp, nil
}

func (s *adminService) GetAllStaff(ctx context.Context, page, pageSize int) ([]dtos.StaffResponse, int64, error) {
	staffList, total, err := s.staffRepo.FindAll(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
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
		return nil, errors.New("staff not found")
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
		return nil, err
	}
	resp := localMapModelStaffToStaffResponse(staff)
	return &resp, nil
}

func (s *adminService) DeleteStaff(ctx context.Context, staffID uint) error {
	return s.staffRepo.Delete(ctx, staffID)
}

// --- Activity Logs ---
func (s *adminService) GetActivityLogs(ctx context.Context, page, pageSize int) ([]dtos.ActivityLogResponse, int64, error) {
	filters := make(map[string]string)
	logs, total, err := s.activityLogSvc.GetActivityLogs(ctx, page, pageSize, filters)
	if err != nil {
		return nil, 0, err
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
