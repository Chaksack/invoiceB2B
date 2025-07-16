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
	"math"
	"path/filepath"
	"sort"
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
	SelectFinancialInstitution(ctx context.Context, invoiceID, userID uint, req dtos.SelectFinancialInstitutionRequest) (*dtos.InvoiceResponse, error)
	GetSuggestedFinancialInstitutions(ctx context.Context, invoiceID, userID uint) (*dtos.SuggestedFinancialInstitutionsResponse, error)
}

type invoiceService struct {
	invoiceRepo              repositories.InvoiceRepository
	userRepo                 repositories.UserRepository
	transactionRepo          repositories.TransactionRepository
	financialInstitutionRepo repositories.FinancialInstitutionRepository
	fileService              FileService
	notificationSvc          NotificationService
	activityLogSvc           ActivityLogService
	emailService             EmailService
	cfg                      *config.Config
}

func NewInvoiceService(
	invoiceRepo repositories.InvoiceRepository,
	userRepo repositories.UserRepository,
	transactionRepo repositories.TransactionRepository,
	financialInstitutionRepo repositories.FinancialInstitutionRepository,
	fileService FileService,
	notificationSvc NotificationService,
	activityLogSvc ActivityLogService,
	emailService EmailService,
	cfg *config.Config,
) InvoiceService {
	return &invoiceService{
		invoiceRepo:              invoiceRepo,
		userRepo:                 userRepo,
		transactionRepo:          transactionRepo,
		financialInstitutionRepo: financialInstitutionRepo,
		fileService:              fileService,
		notificationSvc:          notificationSvc,
		activityLogSvc:           activityLogSvc,
		emailService:             emailService,
		cfg:                      cfg,
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

func (s *invoiceService) SelectFinancialInstitution(ctx context.Context, invoiceID, userID uint, req dtos.SelectFinancialInstitutionRequest) (*dtos.InvoiceResponse, error) {
	// Get the invoice
	invoice, err := s.invoiceRepo.FindByID(ctx, invoiceID)
	if err != nil {
		log.Printf("Error finding invoice %d: %v", invoiceID, err)
		return nil, ErrInvoiceNotFound
	}

	// Verify that the invoice belongs to the user
	if invoice.UserID != userID {
		log.Printf("User %d attempted to access invoice %d which belongs to user %d", userID, invoiceID, invoice.UserID)
		return nil, ErrInvoiceAccessDenied
	}

	// Get the user for email notification
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		log.Printf("Error finding user %d: %v", userID, err)
		return nil, ErrUserNotFound
	}

	// Get the absolute path of the original invoice file
	originalFilePath, err := s.fileService.GetAbsPath(invoice.OriginalFilePath)
	if err != nil {
		log.Printf("Error getting absolute path for invoice file %s: %v", invoice.OriginalFilePath, err)
		return nil, fmt.Errorf("invoice file path error: %w", err)
	}

	// Prepare email subject and body
	subject := fmt.Sprintf("New Invoice Financing Request from %s", user.CompanyName)
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>New Invoice Financing Request</h2>
			<p>A new invoice financing request has been submitted by %s.</p>
			<h3>Invoice Details:</h3>
			<ul>
				<li><strong>Invoice Number:</strong> %s</li>
				<li><strong>Amount:</strong> %s %.2f</li>
				<li><strong>Issuer:</strong> %s</li>
				<li><strong>Issuer Bank Account:</strong> %s</li>
				<li><strong>Issuer Bank Name:</strong> %s</li>
				<li><strong>Debtor:</strong> %s</li>
				<li><strong>Due Date:</strong> %s</li>
			</ul>
			<h3>Additional Extracted Data:</h3>
			<pre>%s</pre>
			<p>The original invoice file and extracted data are attached to this email.</p>
			<p>Please review the request and respond accordingly.</p>
			<p>Thank you,<br>Invoice B2B Platform</p>
		</body>
		</html>
	`, 
	user.CompanyName, 
	invoice.InvoiceNumber, 
	invoice.Currency, 
	invoice.Amount, 
	invoice.IssuerName, 
	invoice.IssuerBankAccount,
	invoice.IssuerBankName,
	invoice.DebtorName, 
	func() string {
		if invoice.DueDate != nil {
			return invoice.DueDate.Format("2006-01-02")
		}
		return "N/A"
	}(),
	invoice.JSONData)

	// Send email with invoice data and raw invoice to the financial institution
	fileName := filepath.Base(invoice.OriginalFilePath)
	err = s.emailService.SendEmailWithAttachment(
		req.FinancialInstitutionEmail,
		subject,
		body,
		originalFilePath,
		fileName,
	)
	if err != nil {
		log.Printf("Error sending email to financial institution %s: %v", req.FinancialInstitutionName, err)
		return nil, fmt.Errorf("failed to send email to financial institution: %w", err)
	}

	// Log the activity
	logDetails := map[string]interface{}{
		"invoice_id":                 invoiceID,
		"financial_institution_id":   req.FinancialInstitutionID,
		"financial_institution_name": req.FinancialInstitutionName,
	}
	_ = s.activityLogSvc.LogActivity(ctx, nil, &userID, "INVOICE_SENT_TO_FINANCIAL_INSTITUTION", logDetails, "")

	// Return the invoice response
	resp := mapInvoiceToResponse(invoice)
	return &resp, nil
}

// GetSuggestedFinancialInstitutions returns a list of suggested financial institutions for an invoice
// based on the invoice details and financial institution products.
func (s *invoiceService) GetSuggestedFinancialInstitutions(ctx context.Context, invoiceID, userID uint) (*dtos.SuggestedFinancialInstitutionsResponse, error) {
	// Get the invoice
	invoice, err := s.invoiceRepo.FindByID(ctx, invoiceID)
	if err != nil {
		log.Printf("Error finding invoice %d: %v", invoiceID, err)
		return nil, ErrInvoiceNotFound
	}

	// Verify that the invoice belongs to the user
	if invoice.UserID != userID {
		log.Printf("User %d attempted to access invoice %d which belongs to user %d", userID, invoiceID, invoice.UserID)
		return nil, ErrInvoiceAccessDenied
	}

	// Get all active financial institutions
	filters := map[string]string{"is_active": "true"}
	financialInstitutions, _, err := s.financialInstitutionRepo.FindAllFinancialInstitutions(ctx, 1, 100, filters)
	if err != nil {
		log.Printf("Error finding financial institutions: %v", err)
		return nil, fmt.Errorf("failed to retrieve financial institutions: %w", err)
	}

	// Calculate compatibility scores for each financial institution
	var suggestions []dtos.SuggestedFinancialInstitution
	for _, fi := range financialInstitutions {
		// Skip if the financial institution doesn't have min/max invoice amount set
		if fi.MinInvoiceAmount > 0 && invoice.Amount < fi.MinInvoiceAmount {
			continue
		}
		if fi.MaxInvoiceAmount > 0 && invoice.Amount > fi.MaxInvoiceAmount {
			continue
		}

		// Calculate base compatibility score
		score := 0.0

		// Higher score if the invoice amount is within the financial institution's preferred range
		if fi.MinInvoiceAmount > 0 && fi.MaxInvoiceAmount > 0 {
			// If the amount is in the middle of the range, give a higher score
			rangeSize := fi.MaxInvoiceAmount - fi.MinInvoiceAmount
			if rangeSize > 0 {
				position := (invoice.Amount - fi.MinInvoiceAmount) / rangeSize
				// Score is highest (1.0) when position is 0.5 (middle of range)
				score += 1.0 - math.Abs(position - 0.5) * 2.0
			}
		}

		// Check if the financial institution has products that match the invoice
		products, _, err := s.financialInstitutionRepo.FindFinancialInstitutionProductsByFinancialInstitutionID(ctx, fi.ID, 1, 100)
		if err != nil {
			log.Printf("Error finding products for financial institution %d: %v", fi.ID, err)
			continue
		}

		// Find the best matching product
		bestProductScore := 0.0
		for _, product := range products {
			if !product.IsActive {
				continue
			}

			// Skip if the product doesn't have min/max invoice amount set
			if product.MinInvoiceAmount > 0 && invoice.Amount < product.MinInvoiceAmount {
				continue
			}
			if product.MaxInvoiceAmount > 0 && invoice.Amount > product.MaxInvoiceAmount {
				continue
			}

			// Calculate product compatibility score
			productScore := 0.0

			// Higher score if the invoice amount is within the product's preferred range
			if product.MinInvoiceAmount > 0 && product.MaxInvoiceAmount > 0 {
				rangeSize := product.MaxInvoiceAmount - product.MinInvoiceAmount
				if rangeSize > 0 {
					position := (invoice.Amount - product.MinInvoiceAmount) / rangeSize
					productScore += 1.0 - math.Abs(position - 0.5) * 2.0
				}
			}

			// Check if the invoice has a due date and the product has terms days
			if invoice.DueDate != nil && product.TermsDays > 0 {
				daysUntilDue := int(invoice.DueDate.Sub(time.Now()).Hours() / 24)
				// Higher score if the days until due is close to the product's terms days
				if daysUntilDue > 0 {
					termRatio := float64(daysUntilDue) / float64(product.TermsDays)
					if termRatio <= 1.0 {
						// Score is highest (1.0) when termRatio is 1.0 (exact match)
						productScore += termRatio
					} else {
						// Score decreases as termRatio increases beyond 1.0
						productScore += 1.0 / termRatio
					}
				}
			}

			// Update best product score
			if productScore > bestProductScore {
				bestProductScore = productScore
			}
		}

		// Add product score to overall score
		score += bestProductScore

		// Normalize score to be between 0 and 1
		score = math.Min(1.0, score / 2.0)

		// Add to suggestions if score is above threshold
		if score > 0.3 {
			suggestions = append(suggestions, dtos.SuggestedFinancialInstitution{
				ID:                fi.ID,
				Name:              fi.Name,
				Code:              fi.Code,
				Description:       fi.Description,
				InterestRateMin:   fi.InterestRateMin,
				InterestRateMax:   fi.InterestRateMax,
				ProcessingFee:     fi.ProcessingFee,
				TermsDays:         fi.TermsDays,
				CompatibilityScore: score,
				// Email field would need to be populated from a separate source
				// For now, we'll leave it empty and let the frontend handle it
			})
		}
	}

	// Sort suggestions by compatibility score (highest first)
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].CompatibilityScore > suggestions[j].CompatibilityScore
	})

	// Limit to top 5 suggestions
	if len(suggestions) > 5 {
		suggestions = suggestions[:5]
	}

	return &dtos.SuggestedFinancialInstitutionsResponse{
		Suggestions: suggestions,
	}, nil
}
