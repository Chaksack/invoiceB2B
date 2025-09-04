package services

import (
	"encoding/json"
	"fmt"
	"time"
	"strings"
	"math/rand"

	"invoiceB2B/internal/models"
	"invoiceB2B/internal/dtos"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

// N8NIntegrationService interface for N8N integration
type N8NIntegrationService interface {
	ProcessLoanApplication(request *dtos.N8NProcessingRequest) (*dtos.N8NProcessingResponse, error)
}

// LoanApplicationService interface defines loan application operations
type LoanApplicationService interface {
	CreateLoanApplication(userID uint, request *dtos.CreateLoanApplicationRequest) (*models.LoanApplication, error)
	CreateManualLoanApplication(userID uint, request *dtos.ManualLoanInputRequest) (*models.LoanApplication, error)
	GetLoanApplication(id uint, userID uint) (*models.LoanApplication, error)
	GetLoanApplicationsByUser(userID uint, page, pageSize int) ([]models.LoanApplication, int64, error)
	UpdateLoanApplicationStatus(id uint, status models.LoanApplicationStatus, notes string) error
	SubmitKYBInformation(loanApplicationID uint, request *dtos.SubmitKYBInformationRequest) (*models.KYBInformation, error)
	SubmitFinancialStatement(loanApplicationID uint, request *dtos.SubmitFinancialStatementRequest, filePath string) (*models.FinancialStatement, error)
	ProcessDocumentUpload(loanApplicationID uint, filePath string, source models.LoanApplicationSource) error
	SendToN8NProcessing(loanApplicationID uint, processingType string) (*dtos.N8NProcessingResponse, error)
	UpdateFromN8NProcessing(loanApplicationID uint, response *dtos.N8NProcessingResponse) error
	AdminReview(loanApplicationID uint, adminID uint, request *dtos.AdminReviewRequest) error
	SendToFinancialInstitution(loanApplicationID uint, adminID uint, request *dtos.SendToFinancialInstitutionRequest) error
	CheckDataExpiration() ([]models.LoanApplication, error)
	RefreshExpiredData(loanApplicationID uint, dataType string) error
	GetLoanApplicationStats() (*dtos.LoanApplicationStatsResponse, error)
	GetApplicationsRequiringRefresh() ([]models.DataRefreshTracker, error)
	SendDataRefreshReminder(trackerID uint, reminderType string) error
}

type loanApplicationService struct {
	db                  *gorm.DB
	n8nService          N8NIntegrationService
	fileService         FileService
	notificationService NotificationService
	emailService        EmailService
}

// NewLoanApplicationService creates a new loan application service
func NewLoanApplicationService(
	db *gorm.DB,
	n8nService N8NIntegrationService,
	fileService FileService,
	notificationService NotificationService,
	emailService EmailService,
) LoanApplicationService {
	return &loanApplicationService{
		db:                  db,
		n8nService:          n8nService,
		fileService:         fileService,
		notificationService: notificationService,
		emailService:        emailService,
	}
}

// CreateLoanApplication creates a new loan application
func (s *loanApplicationService) CreateLoanApplication(userID uint, request *dtos.CreateLoanApplicationRequest) (*models.LoanApplication, error) {
	// Generate unique application reference
	reference, err := s.generateApplicationReference()
	if err != nil {
		return nil, fmt.Errorf("failed to generate application reference: %w", err)
	}

	loanApp := &models.LoanApplication{
		UserID:               userID,
		ApplicationReference: reference,
		Source:               models.LoanApplicationSource(request.Source),
		RequestedAmount:      request.RequestedAmount,
		Currency:             request.Currency,
		Purpose:              request.Purpose,
		Status:               models.LoanApplicationPending,
		LastUpdatedAt:        time.Now(),
	}

	// Determine initial status based on source
	switch request.Source {
	case "manual":
		loanApp.Status = models.LoanApplicationKYBRequired
	case "invoice", "contract":
		loanApp.Status = models.LoanApplicationPending
	}

	if err := s.db.Create(loanApp).Error; err != nil {
		return nil, fmt.Errorf("failed to create loan application: %w", err)
	}

	// Create data refresh trackers
	if err := s.createDataRefreshTrackers(loanApp.ID); err != nil {
		return nil, fmt.Errorf("failed to create data refresh trackers: %w", err)
	}

	return loanApp, nil
}

// CreateManualLoanApplication creates a loan application with all data provided manually
func (s *loanApplicationService) CreateManualLoanApplication(userID uint, request *dtos.ManualLoanInputRequest) (*models.LoanApplication, error) {
	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create loan application
	createReq := &dtos.CreateLoanApplicationRequest{
		Source:          request.Source,
		RequestedAmount: request.RequestedAmount,
		Currency:        request.Currency,
		Purpose:         request.Purpose,
	}

	reference, err := s.generateApplicationReference()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to generate application reference: %w", err)
	}

	loanApp := &models.LoanApplication{
		UserID:               userID,
		ApplicationReference: reference,
		Source:               models.LoanApplicationSource(createReq.Source),
		RequestedAmount:      createReq.RequestedAmount,
		Currency:             createReq.Currency,
		Purpose:              createReq.Purpose,
		Status:               models.LoanApplicationUnderReview,
		SubmittedAt:          &time.Time{},
		LastUpdatedAt:        time.Now(),
	}

	now := time.Now()
	loanApp.SubmittedAt = &now

	if err := tx.Create(loanApp).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create loan application: %w", err)
	}

	// Create KYB Information
	kybInfo, err := s.createKYBInformationFromRequest(tx, loanApp.ID, &request.KYBInformation)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create KYB information: %w", err)
	}

	// Create Financial Statements
	for _, stmtReq := range request.FinancialStatements {
		if _, err := s.createFinancialStatementFromRequest(tx, loanApp.ID, &stmtReq, ""); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create financial statement: %w", err)
		}
	}

	// Create data refresh trackers
	if err := s.createDataRefreshTrackersWithTx(tx, loanApp.ID); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create data refresh trackers: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Load complete application with associations
	if err := s.db.Preload("KYBInformation").Preload("FinancialStatements").First(loanApp, loanApp.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load complete application: %w", err)
	}

	// Send to N8N for processing
	go func() {
		if _, err := s.SendToN8NProcessing(loanApp.ID, "complete_application"); err != nil {
			fmt.Printf("Failed to send application to N8N: %v\n", err)
		}
	}()

	// Update with KYB reference for response
	loanApp.KYBInformation = kybInfo

	return loanApp, nil
}

// GetLoanApplication retrieves a loan application by ID
func (s *loanApplicationService) GetLoanApplication(id uint, userID uint) (*models.LoanApplication, error) {
	var loanApp models.LoanApplication
	query := s.db.Preload("User").Preload("KYBInformation").Preload("FinancialStatements").
		Preload("RecommendedInstitution").Preload("ReviewedBy").Preload("ApprovedBy").Preload("DisbursedBy")

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&loanApp, id).Error; err != nil {
		return nil, fmt.Errorf("loan application not found: %w", err)
	}

	return &loanApp, nil
}

// GetLoanApplicationsByUser retrieves loan applications for a specific user
func (s *loanApplicationService) GetLoanApplicationsByUser(userID uint, page, pageSize int) ([]models.LoanApplication, int64, error) {
	var applications []models.LoanApplication
	var totalCount int64

	// Count total records
	if err := s.db.Model(&models.LoanApplication{}).Where("user_id = ?", userID).Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count loan applications: %w", err)
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Fetch paginated results
	query := s.db.Where("user_id = ?", userID).
		Preload("KYBInformation").
		Preload("FinancialStatements").
		Preload("RecommendedInstitution").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset)

	if err := query.Find(&applications).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch loan applications: %w", err)
	}

	return applications, totalCount, nil
}

// UpdateLoanApplicationStatus updates the status of a loan application
func (s *loanApplicationService) UpdateLoanApplicationStatus(id uint, status models.LoanApplicationStatus, notes string) error {
	updates := map[string]interface{}{
		"status":            status,
		"last_updated_at":   time.Now(),
	}

	if notes != "" {
		updates["review_notes"] = notes
	}

	if err := s.db.Model(&models.LoanApplication{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update loan application status: %w", err)
	}

	return nil
}

// SubmitKYBInformation submits KYB information for a loan application
func (s *loanApplicationService) SubmitKYBInformation(loanApplicationID uint, request *dtos.SubmitKYBInformationRequest) (*models.KYBInformation, error) {
	// Check if loan application exists
	var loanApp models.LoanApplication
	if err := s.db.First(&loanApp, loanApplicationID).Error; err != nil {
		return nil, fmt.Errorf("loan application not found: %w", err)
	}

	// Create or update KYB information
	return s.createKYBInformationFromRequest(s.db, loanApplicationID, request)
}

// SubmitFinancialStatement submits a financial statement for a loan application
func (s *loanApplicationService) SubmitFinancialStatement(loanApplicationID uint, request *dtos.SubmitFinancialStatementRequest, filePath string) (*models.FinancialStatement, error) {
	// Check if loan application exists
	var loanApp models.LoanApplication
	if err := s.db.First(&loanApp, loanApplicationID).Error; err != nil {
		return nil, fmt.Errorf("loan application not found: %w", err)
	}

	// Create financial statement
	return s.createFinancialStatementFromRequest(s.db, loanApplicationID, request, filePath)
}

// ProcessDocumentUpload processes an uploaded invoice or contract document
func (s *loanApplicationService) ProcessDocumentUpload(loanApplicationID uint, filePath string, source models.LoanApplicationSource) error {
	// Update loan application with document path
	updates := map[string]interface{}{
		"supporting_document_path": filePath,
		"last_updated_at":         time.Now(),
	}

	if err := s.db.Model(&models.LoanApplication{}).Where("id = ?", loanApplicationID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update loan application with document: %w", err)
	}

	// Send to N8N for document processing
	go func() {
		if _, err := s.SendToN8NProcessing(loanApplicationID, "document"); err != nil {
			fmt.Printf("Failed to send document to N8N: %v\n", err)
		}
	}()

	return nil
}

// SendToN8NProcessing sends loan application data to N8N for processing
func (s *loanApplicationService) SendToN8NProcessing(loanApplicationID uint, processingType string) (*dtos.N8NProcessingResponse, error) {
	// Load complete loan application data
	var loanApp models.LoanApplication
	if err := s.db.Preload("KYBInformation").Preload("FinancialStatements").First(&loanApp, loanApplicationID).Error; err != nil {
		return nil, fmt.Errorf("failed to load loan application: %w", err)
	}

	// Prepare N8N request
	request := &dtos.N8NProcessingRequest{
		LoanApplicationID: loanApplicationID,
		ApplicationData:   dtos.ToLoanApplicationResponse(&loanApp),
		ProcessingType:    processingType,
	}

	if loanApp.KYBInformation != nil {
		kybResponse := dtos.ToKYBInformationResponse(loanApp.KYBInformation)
		request.KYBInformation = &kybResponse
	}

	for _, stmt := range loanApp.FinancialStatements {
		request.FinancialStatements = append(request.FinancialStatements, 
			dtos.ToFinancialStatementResponse(&stmt))
	}

	if loanApp.SupportingDocumentPath != nil {
		request.DocumentPath = loanApp.SupportingDocumentPath
	}

	// Send to N8N service
	response, err := s.n8nService.ProcessLoanApplication(request)
	if err != nil {
		return nil, fmt.Errorf("failed to process with N8N: %w", err)
	}

	// Update loan application with N8N execution details
	updates := map[string]interface{}{
		"n8n_workflow_execution_id": response.WorkflowExecutionID,
		"n8n_processing_status":     response.Status,
		"last_updated_at":          time.Now(),
	}

	if response.Status == "completed" {
		now := time.Now()
		updates["n8n_processed_at"] = &now
	}

	if err := s.db.Model(&models.LoanApplication{}).Where("id = ?", loanApplicationID).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update N8N processing status: %w", err)
	}

	return response, nil
}

// UpdateFromN8NProcessing updates loan application with N8N processing results
func (s *loanApplicationService) UpdateFromN8NProcessing(loanApplicationID uint, response *dtos.N8NProcessingResponse) error {
	updates := map[string]interface{}{
		"n8n_processing_status": response.Status,
		"last_updated_at":      time.Now(),
	}

	if response.Status == "completed" {
		now := time.Now()
		updates["n8n_processed_at"] = &now

		// Update matched institutions
		if len(response.MatchedInstitutions) > 0 {
			matchedData, _ := json.Marshal(response.MatchedInstitutions)
			updates["matched_institutions"] = string(matchedData)
		}

		// Set recommended institution
		if response.RecommendedInstitution != nil {
			updates["recommended_institution_id"] = response.RecommendedInstitution.InstitutionID
		}

		// Update status to under review if processing completed successfully
		updates["status"] = models.LoanApplicationUnderReview
	}

	if err := s.db.Model(&models.LoanApplication{}).Where("id = ?", loanApplicationID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update loan application from N8N: %w", err)
	}

	return nil
}

// AdminReview processes admin review of a loan application
func (s *loanApplicationService) AdminReview(loanApplicationID uint, adminID uint, request *dtos.AdminReviewRequest) error {
	now := time.Now()
	
	updates := map[string]interface{}{
		"reviewed_by_id":   adminID,
		"reviewed_at":      &now,
		"review_notes":     request.ReviewNotes,
		"last_updated_at":  now,
	}

	switch request.Action {
	case "approve":
		updates["status"] = models.LoanApplicationApproved
		updates["approved_by_id"] = adminID
		updates["approved_at"] = &now
		if request.ApprovedAmount != nil {
			updates["approved_amount"] = *request.ApprovedAmount
		}
	case "reject":
		updates["status"] = models.LoanApplicationRejected
		if request.RejectionReason != nil {
			updates["rejection_reason"] = *request.RejectionReason
		}
	case "request_info":
		updates["status"] = models.LoanApplicationKYBRequired
	}

	if err := s.db.Model(&models.LoanApplication{}).Where("id = ?", loanApplicationID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update loan application review: %w", err)
	}

	// Send notification to user
	go func() {
		if s.notificationService != nil {
			eventPayload := map[string]interface{}{
				"loan_application_id": loanApplicationID,
				"status": string(models.LoanApplicationStatus(updates["status"].(models.LoanApplicationStatus))),
				"admin_id": adminID,
				"timestamp": time.Now(),
			}
			if err := s.notificationService.PublishEvent("loan.events", "loan.status.updated", eventPayload); err != nil {
				fmt.Printf("Failed to send status notification: %v\n", err)
			}
		}
	}()

	return nil
}

// Helper methods

// generateApplicationReference generates a unique application reference
func (s *loanApplicationService) generateApplicationReference() (string, error) {
	for i := 0; i < 10; i++ { // Try up to 10 times
		reference := fmt.Sprintf("LA-%d-%04d", time.Now().Year(), rand.Intn(10000))
		
		var count int64
		if err := s.db.Model(&models.LoanApplication{}).Where("application_reference = ?", reference).Count(&count).Error; err != nil {
			return "", err
		}
		
		if count == 0 {
			return reference, nil
		}
	}
	
	// Fallback to UUID-based reference
	return fmt.Sprintf("LA-%s", strings.Replace(uuid.New().String(), "-", "", -1)[:8]), nil
}

// createKYBInformationFromRequest creates KYB information from DTO request
func (s *loanApplicationService) createKYBInformationFromRequest(db *gorm.DB, loanApplicationID uint, request *dtos.SubmitKYBInformationRequest) (*models.KYBInformation, error) {
	directorsJSON, _ := json.Marshal(request.DirectorsInformation)
	ownersJSON, _ := json.Marshal(request.BeneficialOwners)
	
	now := time.Now()
	refreshDate := now.AddDate(0, 3, 0) // 3 months from now

	kybInfo := &models.KYBInformation{
		LoanApplicationID:       loanApplicationID,
		BusinessName:            request.BusinessName,
		BusinessRegistrationNo:  request.BusinessRegistrationNo,
		BusinessType:            request.BusinessType,
		IndustryType:            request.IndustryType,
		BusinessAddress:         request.BusinessAddress,
		TaxIdentificationNumber: request.TaxIdentificationNumber,
		YearsInOperation:        request.YearsInOperation,
		DirectorsInformation:    string(directorsJSON),
		BeneficialOwners:        string(ownersJSON),
		AnnualRevenue:           request.AnnualRevenue,
		MonthlyRevenue:          request.MonthlyRevenue,
		NumberOfEmployees:       request.NumberOfEmployees,
		BusinessDescription:     request.BusinessDescription,
		SubmittedAt:             now,
		RefreshDueDate:          refreshDate,
		VerificationStatus:      "pending",
	}

	// Use upsert: create if not exists, update if exists
	if err := db.Where("loan_application_id = ?", loanApplicationID).Save(kybInfo).Error; err != nil {
		return nil, fmt.Errorf("failed to save KYB information: %w", err)
	}

	return kybInfo, nil
}

// createFinancialStatementFromRequest creates financial statement from DTO request
func (s *loanApplicationService) createFinancialStatementFromRequest(db *gorm.DB, loanApplicationID uint, request *dtos.SubmitFinancialStatementRequest, filePath string) (*models.FinancialStatement, error) {
	now := time.Now()
	refreshDate := now.AddDate(0, 3, 0) // 3 months from now

	avgBalance := request.AverageBalance
	if avgBalance == nil {
		avg := (request.OpeningBalance + request.ClosingBalance) / 2
		avgBalance = &avg
	}

	statement := &models.FinancialStatement{
		LoanApplicationID:   loanApplicationID,
		Type:                models.FinancialStatementType(request.Type),
		ProviderName:        request.ProviderName,
		AccountNumber:       request.AccountNumber,
		StatementPeriodFrom: request.StatementPeriodFrom,
		StatementPeriodTo:   request.StatementPeriodTo,
		OpeningBalance:      request.OpeningBalance,
		ClosingBalance:      request.ClosingBalance,
		AverageBalance:      *avgBalance,
		TotalCredits:        request.TotalCredits,
		TotalDebits:         request.TotalDebits,
		TransactionCount:    request.TransactionCount,
		OriginalFilePath:    filePath,
		SubmittedAt:         now,
		RefreshDueDate:      refreshDate,
	}

	if err := db.Create(statement).Error; err != nil {
		return nil, fmt.Errorf("failed to create financial statement: %w", err)
	}

	return statement, nil
}

// createDataRefreshTrackers creates data refresh trackers for a loan application
func (s *loanApplicationService) createDataRefreshTrackers(loanApplicationID uint) error {
	return s.createDataRefreshTrackersWithTx(s.db, loanApplicationID)
}

// createDataRefreshTrackersWithTx creates data refresh trackers with transaction
func (s *loanApplicationService) createDataRefreshTrackersWithTx(db *gorm.DB, loanApplicationID uint) error {
	now := time.Now()
	nextRefresh := now.AddDate(0, 3, 0) // 3 months from now

	trackers := []models.DataRefreshTracker{
		{
			LoanApplicationID: loanApplicationID,
			DataType:         "kyb",
			LastRefreshDate:  now,
			NextRefreshDate:  nextRefresh,
			Status:           "active",
		},
		{
			LoanApplicationID: loanApplicationID,
			DataType:         "financial_statements",
			LastRefreshDate:  now,
			NextRefreshDate:  nextRefresh,
			Status:           "active",
		},
	}

	for _, tracker := range trackers {
		if err := db.Create(&tracker).Error; err != nil {
			return fmt.Errorf("failed to create data refresh tracker: %w", err)
		}
	}

	return nil
}

// CheckDataExpiration checks for loan applications with expired data
func (s *loanApplicationService) CheckDataExpiration() ([]models.LoanApplication, error) {
	var applications []models.LoanApplication
	now := time.Now()

	// Find applications with expired KYB data
	var kybApplicationIDs []uint
	if err := s.db.Model(&models.KYBInformation{}).
		Where("refresh_due_date < ? AND is_expired = false", now).
		Pluck("loan_application_id", &kybApplicationIDs).Error; err != nil {
		return nil, fmt.Errorf("failed to find applications with expired KYB data: %w", err)
	}

	// Find applications with expired financial statement data
	var financialApplicationIDs []uint
	if err := s.db.Model(&models.FinancialStatement{}).
		Where("refresh_due_date < ? AND is_expired = false", now).
		Pluck("loan_application_id", &financialApplicationIDs).Error; err != nil {
		return nil, fmt.Errorf("failed to find applications with expired financial data: %w", err)
	}

	// Combine IDs and remove duplicates
	allIDs := make(map[uint]bool)
	for _, id := range kybApplicationIDs {
		allIDs[id] = true
	}
	for _, id := range financialApplicationIDs {
		allIDs[id] = true
	}

	var uniqueIDs []uint
	for id := range allIDs {
		uniqueIDs = append(uniqueIDs, id)
	}

	if len(uniqueIDs) == 0 {
		return applications, nil
	}

	// Find applications with expired data
	query := s.db.Preload("KYBInformation").Preload("FinancialStatements").
		Where("id IN (?)", uniqueIDs)

	if err := query.Find(&applications).Error; err != nil {
		return nil, fmt.Errorf("failed to find applications with expired data: %w", err)
	}

	// Mark data as expired
	if err := s.db.Model(&models.KYBInformation{}).Where("refresh_due_date < ? AND is_expired = false", now).Update("is_expired", true).Error; err != nil {
		return nil, fmt.Errorf("failed to mark KYB data as expired: %w", err)
	}

	if err := s.db.Model(&models.FinancialStatement{}).Where("refresh_due_date < ? AND is_expired = false", now).Update("is_expired", true).Error; err != nil {
		return nil, fmt.Errorf("failed to mark financial statements as expired: %w", err)
	}

	return applications, nil
}

// RefreshExpiredData refreshes expired data for a loan application
func (s *loanApplicationService) RefreshExpiredData(loanApplicationID uint, dataType string) error {
	now := time.Now()
	nextRefresh := now.AddDate(0, 3, 0)

	// Update the refresh tracker
	updates := map[string]interface{}{
		"last_refresh_date": now,
		"next_refresh_date": nextRefresh,
		"is_overdue":       false,
		"reminder_sent_at": nil,
		"overdue_notified_at": nil,
	}

	if err := s.db.Model(&models.DataRefreshTracker{}).
		Where("loan_application_id = ? AND data_type = ?", loanApplicationID, dataType).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update refresh tracker: %w", err)
	}

	// Mark the specific data as not expired
	switch dataType {
	case "kyb":
		if err := s.db.Model(&models.KYBInformation{}).
			Where("loan_application_id = ?", loanApplicationID).
			Updates(map[string]interface{}{
				"last_refreshed_at": &now,
				"refresh_due_date": nextRefresh,
				"is_expired": false,
			}).Error; err != nil {
			return fmt.Errorf("failed to refresh KYB data: %w", err)
		}
	case "financial_statements":
		if err := s.db.Model(&models.FinancialStatement{}).
			Where("loan_application_id = ?", loanApplicationID).
			Updates(map[string]interface{}{
				"last_refreshed_at": &now,
				"refresh_due_date": nextRefresh,
				"is_expired": false,
			}).Error; err != nil {
			return fmt.Errorf("failed to refresh financial statements: %w", err)
		}
	}

	return nil
}

// GetLoanApplicationStats retrieves loan application statistics
func (s *loanApplicationService) GetLoanApplicationStats() (*dtos.LoanApplicationStatsResponse, error) {
	var stats dtos.LoanApplicationStatsResponse

	// Total applications
	if err := s.db.Model(&models.LoanApplication{}).Count(&stats.TotalApplications).Error; err != nil {
		return nil, fmt.Errorf("failed to count total applications: %w", err)
	}

	// Applications by status
	statusCounts := []struct {
		Status string
		Count  int64
	}{}

	if err := s.db.Model(&models.LoanApplication{}).
		Select("status, count(*) as count").
		Group("status").
		Find(&statusCounts).Error; err != nil {
		return nil, fmt.Errorf("failed to count applications by status: %w", err)
	}

	for _, sc := range statusCounts {
		switch sc.Status {
		case string(models.LoanApplicationPending):
			stats.PendingApplications = sc.Count
		case string(models.LoanApplicationApproved):
			stats.ApprovedApplications = sc.Count
		case string(models.LoanApplicationRejected):
			stats.RejectedApplications = sc.Count
		case string(models.LoanApplicationDisbursed):
			stats.DisbursedApplications = sc.Count
		}
	}

	// Amount statistics
	amountStats := struct {
		TotalRequested float64 `gorm:"column:total_requested"`
		TotalApproved  float64 `gorm:"column:total_approved"`
		TotalDisbursed float64 `gorm:"column:total_disbursed"`
	}{}

	if err := s.db.Model(&models.LoanApplication{}).
		Select(`
			COALESCE(SUM(requested_amount), 0) as total_requested,
			COALESCE(SUM(approved_amount), 0) as total_approved,
			COALESCE(SUM(disbursed_amount), 0) as total_disbursed
		`).
		First(&amountStats).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate amount statistics: %w", err)
	}

	stats.TotalAmountRequested = amountStats.TotalRequested
	stats.TotalAmountApproved = amountStats.TotalApproved
	stats.TotalAmountDisbursed = amountStats.TotalDisbursed

	// Expired data count
	now := time.Now()
	if err := s.db.Model(&models.KYBInformation{}).
		Where("refresh_due_date < ? OR is_expired = true", now).
		Count(&stats.ExpiredDataCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count expired data: %w", err)
	}

	return &stats, nil
}

// GetApplicationsRequiringRefresh retrieves applications requiring data refresh
func (s *loanApplicationService) GetApplicationsRequiringRefresh() ([]models.DataRefreshTracker, error) {
	var trackers []models.DataRefreshTracker
	now := time.Now()
	reminderThreshold := now.AddDate(0, 0, 7) // 7 days before due date

	query := s.db.Preload("LoanApplication").Preload("LoanApplication.User").
		Where("next_refresh_date <= ? AND status = 'active'", reminderThreshold).
		Where("reminder_sent_at IS NULL OR (is_overdue = true AND overdue_notified_at IS NULL)")

	if err := query.Find(&trackers).Error; err != nil {
		return nil, fmt.Errorf("failed to find applications requiring refresh: %w", err)
	}

	return trackers, nil
}

// SendDataRefreshReminder sends a data refresh reminder
func (s *loanApplicationService) SendDataRefreshReminder(trackerID uint, reminderType string) error {
	now := time.Now()
	
	updates := map[string]interface{}{}
	
	switch reminderType {
	case "due_soon":
		updates["reminder_sent_at"] = &now
	case "overdue":
		updates["is_overdue"] = true
		updates["overdue_notified_at"] = &now
	}

	if err := s.db.Model(&models.DataRefreshTracker{}).Where("id = ?", trackerID).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update refresh reminder status: %w", err)
	}

	return nil
}

// SendToFinancialInstitution sends loan application data to a financial institution
func (s *loanApplicationService) SendToFinancialInstitution(loanApplicationID uint, adminID uint, request *dtos.SendToFinancialInstitutionRequest) error {
	// Get the loan application with all related data
	var loanApp models.LoanApplication
	if err := s.db.Preload("User").
		Preload("KYBInformation").
		Preload("FinancialStatements").
		Preload("RecommendedInstitution").
		First(&loanApp, loanApplicationID).Error; err != nil {
		return fmt.Errorf("loan application not found: %w", err)
	}

	// Prepare raw data (original application information)
	rawData := fmt.Sprintf(`Application Reference: %s
Source: %s
Requested Amount: %s %.2f
Purpose: %s
Status: %s
Submitted At: %s

User Information:
- Company: %s
- Email: %s`, 
		loanApp.ApplicationReference,
		string(loanApp.Source),
		loanApp.Currency,
		loanApp.RequestedAmount,
		loanApp.Purpose,
		string(loanApp.Status),
		func() string {
			if loanApp.SubmittedAt != nil {
				return loanApp.SubmittedAt.Format("2006-01-02 15:04:05")
			}
			return "N/A"
		}(),
		loanApp.User.CompanyName,
		loanApp.User.Email)

	// Add KYB information to raw data
	if loanApp.KYBInformation != nil {
		rawData += fmt.Sprintf(`

KYB Information:
- Business Name: %s
- Registration No: %s
- Business Type: %s
- Industry: %s
- Tax ID: %s
- Years in Operation: %d`,
			loanApp.KYBInformation.BusinessName,
			loanApp.KYBInformation.BusinessRegistrationNo,
			loanApp.KYBInformation.BusinessType,
			loanApp.KYBInformation.IndustryType,
			loanApp.KYBInformation.TaxIdentificationNumber,
			loanApp.KYBInformation.YearsInOperation)
	}

	// Prepare processed data (N8N workflow results and analysis)
	processedData := "N8N Processing Results:\n"
	if loanApp.ProcessedDocumentData != "" {
		processedData += fmt.Sprintf("Document Processing: %s\n", loanApp.ProcessedDocumentData)
	}
	if loanApp.MatchedInstitutions != "" {
		processedData += fmt.Sprintf("Matched Institutions: %s\n", loanApp.MatchedInstitutions)
	}
	if loanApp.RecommendedInstitution != nil {
		processedData += fmt.Sprintf(`Recommended Institution:
- Name: %s
- Code: %s
- Interest Rate: %.4f%% - %.4f%%`,
			loanApp.RecommendedInstitution.Name,
			loanApp.RecommendedInstitution.Code,
			loanApp.RecommendedInstitution.InterestRateMin,
			loanApp.RecommendedInstitution.InterestRateMax)
	}

	// Prepare email subject and body
	subject := fmt.Sprintf("Loan Application Submission - %s (%s)", loanApp.User.CompanyName, loanApp.ApplicationReference)
	body := fmt.Sprintf(`<html>
<body>
	<h2>Loan Application Submission</h2>
	<p>Dear %s,</p>
	<p>Please find below a loan application for your review and consideration.</p>
	
	<h3>Application Summary:</h3>
	<ul>
		<li><strong>Reference:</strong> %s</li>
		<li><strong>Company:</strong> %s</li>
		<li><strong>Requested Amount:</strong> %s %.2f</li>
		<li><strong>Purpose:</strong> %s</li>
	</ul>
	
	<h3>Raw Application Data:</h3>
	<pre>%s</pre>
	
	<h3>Processed Data & Analysis:</h3>
	<pre>%s</pre>
	
	<p><strong>Admin Notes:</strong> %s</p>
	
	<p>Please review the application and provide your decision at your earliest convenience.</p>
	
	<p>Best regards,<br>Invoice B2B Platform Administration</p>
</body>
</html>`,
		request.FinancialInstitutionName,
		loanApp.ApplicationReference,
		loanApp.User.CompanyName,
		loanApp.Currency,
		loanApp.RequestedAmount,
		loanApp.Purpose,
		rawData,
		processedData,
		request.Notes)

	// Send email to financial institution
	if err := s.emailService.SendEmail(request.FinancialInstitutionEmail, subject, body); err != nil {
		return fmt.Errorf("failed to send email to financial institution: %w", err)
	}

	return nil
}