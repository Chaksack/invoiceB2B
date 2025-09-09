package dtos

import (
	"time"
	"invoiceB2B/internal/models"
)

// CreateLoanApplicationRequest represents the request to create a new loan application
type CreateLoanApplicationRequest struct {
	Source          string  `json:"source" validate:"required,oneof=manual invoice contract"`
	RequestedAmount float64 `json:"requested_amount" validate:"required,gt=0"`
	Currency        string  `json:"currency" validate:"required,len=3"`
	Purpose         string  `json:"purpose"`
}

// LoanApplicationResponse represents a loan application response
type LoanApplicationResponse struct {
	ID                   uint                              `json:"id"`
	ApplicationReference string                            `json:"application_reference"`
	Source               string                            `json:"source"`
	RequestedAmount      float64                           `json:"requested_amount"`
	Currency             string                            `json:"currency"`
	Purpose              string                            `json:"purpose"`
	Status               string                            `json:"status"`
	SubmittedAt          *time.Time                        `json:"submitted_at"`
	LastUpdatedAt        time.Time                         `json:"last_updated_at"`
	ReviewNotes          *string                           `json:"review_notes,omitempty"`
	ApprovedAmount       *float64                          `json:"approved_amount,omitempty"`
	ApprovedAt           *time.Time                        `json:"approved_at,omitempty"`
	DisbursedAmount      *float64                          `json:"disbursed_amount,omitempty"`
	DisbursedAt          *time.Time                        `json:"disbursed_at,omitempty"`
	RecommendedInstitution *FinancialInstitutionSummary    `json:"recommended_institution,omitempty"`
	KYBInformation       *KYBInformationResponse           `json:"kyb_information,omitempty"`
	FinancialStatements  []FinancialStatementResponse      `json:"financial_statements,omitempty"`
	NextRefreshDate      *time.Time                        `json:"next_refresh_date,omitempty"`
	IsDataExpired        bool                              `json:"is_data_expired"`
	CreatedAt            time.Time                         `json:"created_at"`
	UpdatedAt            time.Time                         `json:"updated_at"`
}

// FinancialInstitutionSummary represents a simplified financial institution response
type FinancialInstitutionSummary struct {
	ID              uint    `json:"id"`
	Name            string  `json:"name"`
	Code            string  `json:"code"`
	InterestRateMin float64 `json:"interest_rate_min"`
	InterestRateMax float64 `json:"interest_rate_max"`
	ProcessingFee   float64 `json:"processing_fee"`
}

// SubmitKYBInformationRequest represents KYB information submission
type SubmitKYBInformationRequest struct {
	BusinessName            string                    `json:"business_name" validate:"required,min=2,max=255"`
	BusinessRegistrationNo  string                    `json:"business_registration_no" validate:"required,max=100"`
	BusinessType            string                    `json:"business_type" validate:"required,max=100"`
	IndustryType            string                    `json:"industry_type" validate:"required,max=100"`
	BusinessAddress         string                    `json:"business_address" validate:"required"`
	PostalAddress           *string                   `json:"postal_address" validate:"omitempty,max=255"`
	DigitalAddress          *string                   `json:"digital_address" validate:"omitempty,max=255"`
	TaxIdentificationNumber string                    `json:"tax_identification_number" validate:"required,max=100"`
	YearsInOperation        int                       `json:"years_in_operation" validate:"required,min=0"`
	DirectorsInformation    []DirectorInformation     `json:"directors_information" validate:"required,min=1"`
	BeneficialOwners        []BeneficialOwnerInfo     `json:"beneficial_owners" validate:"required,min=1"`
	AnnualRevenue           *float64                  `json:"annual_revenue,omitempty" validate:"omitempty,gt=0"`
	MonthlyRevenue          *float64                  `json:"monthly_revenue,omitempty" validate:"omitempty,gt=0"`
	NumberOfEmployees       *int                      `json:"number_of_employees,omitempty" validate:"omitempty,min=1"`
	BusinessDescription     string                    `json:"business_description" validate:"required,min=10"`
}

// DirectorInformation represents director details
type DirectorInformation struct {
	FullName        string   `json:"full_name" validate:"required,min=2,max=255"`
	Position        string   `json:"position" validate:"required,max=100"`
	IDNumber        string   `json:"id_number" validate:"required,max=50"`
	ContactNumber   string   `json:"contact_number" validate:"required,max=20"`
	Address         string   `json:"address" validate:"required"`
	SharePercentage *float64 `json:"share_percentage,omitempty" validate:"omitempty,min=0,max=100"`
	IDPicturePath   string   `json:"id_picture_path" validate:"required"`
}

// BeneficialOwnerInfo represents beneficial owner details
type BeneficialOwnerInfo struct {
	FullName        string  `json:"full_name" validate:"required,min=2,max=255"`
	IDNumber        string  `json:"id_number" validate:"required,max=50"`
	ContactNumber   string  `json:"contact_number" validate:"required,max=20"`
	Address         string  `json:"address" validate:"required"`
	SharePercentage float64 `json:"share_percentage" validate:"required,min=0,max=100"`
	Nationality     string  `json:"nationality" validate:"required,max=100"`
	IDPicturePath   string  `json:"id_picture_path" validate:"required"`
}

// KYBInformationResponse represents KYB information response
type KYBInformationResponse struct {
	ID                      uint                    `json:"id"`
	BusinessName            string                  `json:"business_name"`
	BusinessRegistrationNo  string                  `json:"business_registration_no"`
	BusinessType            string                  `json:"business_type"`
	IndustryType            string                  `json:"industry_type"`
	BusinessAddress         string                  `json:"business_address"`
	PostalAddress           *string                 `json:"postal_address,omitempty"`
	DigitalAddress          *string                 `json:"digital_address,omitempty"`
	TaxIdentificationNumber string                  `json:"tax_identification_number"`
	YearsInOperation        int                     `json:"years_in_operation"`
	DirectorsInformation    []DirectorInformation   `json:"directors_information"`
	BeneficialOwners        []BeneficialOwnerInfo   `json:"beneficial_owners"`
	AnnualRevenue           *float64                `json:"annual_revenue,omitempty"`
	MonthlyRevenue          *float64                `json:"monthly_revenue,omitempty"`
	NumberOfEmployees       *int                    `json:"number_of_employees,omitempty"`
	BusinessDescription     string                  `json:"business_description"`
	SubmittedAt             time.Time               `json:"submitted_at"`
	LastRefreshedAt         *time.Time              `json:"last_refreshed_at,omitempty"`
	RefreshDueDate          time.Time               `json:"refresh_due_date"`
	IsExpired               bool                    `json:"is_expired"`
	VerificationStatus      string                  `json:"verification_status"`
	VerifiedAt              *time.Time              `json:"verified_at,omitempty"`
	CreatedAt               time.Time               `json:"created_at"`
	UpdatedAt               time.Time               `json:"updated_at"`
}

// SubmitFinancialStatementRequest represents financial statement submission
type SubmitFinancialStatementRequest struct {
	Type                string    `json:"type" validate:"required,oneof=bank mobile_money"`
	ProviderName        string    `json:"provider_name" validate:"required,min=2,max=255"`
	AccountNumber       string    `json:"account_number" validate:"required,max=100"`
	StatementPeriodFrom time.Time `json:"statement_period_from" validate:"required"`
	StatementPeriodTo   time.Time `json:"statement_period_to" validate:"required"`
	OpeningBalance      float64   `json:"opening_balance" validate:"required"`
	ClosingBalance      float64   `json:"closing_balance" validate:"required"`
	AverageBalance      *float64  `json:"average_balance,omitempty"`
	TotalCredits        float64   `json:"total_credits" validate:"required,gte=0"`
	TotalDebits         float64   `json:"total_debits" validate:"required,gte=0"`
	TransactionCount    int       `json:"transaction_count" validate:"required,min=0"`
}

// FinancialStatementResponse represents financial statement response
type FinancialStatementResponse struct {
	ID                  uint       `json:"id"`
	Type                string     `json:"type"`
	ProviderName        string     `json:"provider_name"`
	AccountNumber       string     `json:"account_number"`
	StatementPeriodFrom time.Time  `json:"statement_period_from"`
	StatementPeriodTo   time.Time  `json:"statement_period_to"`
	OpeningBalance      float64    `json:"opening_balance"`
	ClosingBalance      float64    `json:"closing_balance"`
	AverageBalance      *float64   `json:"average_balance,omitempty"`
	TotalCredits        float64    `json:"total_credits"`
	TotalDebits         float64    `json:"total_debits"`
	TransactionCount    int        `json:"transaction_count"`
	SubmittedAt         time.Time  `json:"submitted_at"`
	LastRefreshedAt     *time.Time `json:"last_refreshed_at,omitempty"`
	RefreshDueDate      time.Time  `json:"refresh_due_date"`
	IsExpired           bool       `json:"is_expired"`
	N8NProcessingStatus *string    `json:"n8n_processing_status,omitempty"`
	N8NProcessedAt      *time.Time `json:"n8n_processed_at,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

// ManualLoanInputRequest represents manual loan application input
type ManualLoanInputRequest struct {
	CreateLoanApplicationRequest
	KYBInformation      SubmitKYBInformationRequest       `json:"kyb_information"`
	FinancialStatements []SubmitFinancialStatementRequest `json:"financial_statements" validate:"required,min=1"`
}

// DocumentUploadResponse represents document upload response
type DocumentUploadResponse struct {
	FileName   string `json:"file_name"`
	FilePath   string `json:"file_path"`
	FileSize   int64  `json:"file_size"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// N8NProcessingRequest represents data sent to N8N for processing
type N8NProcessingRequest struct {
	LoanApplicationID   uint                        `json:"loan_application_id"`
	ApplicationData     LoanApplicationResponse     `json:"application_data"`
	KYBInformation      *KYBInformationResponse     `json:"kyb_information,omitempty"`
	FinancialStatements []FinancialStatementResponse `json:"financial_statements,omitempty"`
	DocumentPath        *string                     `json:"document_path,omitempty"`
	ProcessingType      string                      `json:"processing_type"` // "kyb", "financial", "document"
}

// N8NProcessingResponse represents N8N processing response
type N8NProcessingResponse struct {
	WorkflowExecutionID string                          `json:"workflow_execution_id"`
	Status              string                          `json:"status"`
	ProcessedData       interface{}                     `json:"processed_data,omitempty"`
	MatchedInstitutions []InstitutionMatchResult        `json:"matched_institutions,omitempty"`
	RecommendedInstitution *InstitutionMatchResult      `json:"recommended_institution,omitempty"`
	RiskAssessment      *RiskAssessmentResult           `json:"risk_assessment,omitempty"`
	Errors              []string                        `json:"errors,omitempty"`
	ProcessingTime      float64                         `json:"processing_time_seconds"`
}

// InstitutionMatchResult represents financial institution matching result
type InstitutionMatchResult struct {
	InstitutionID   uint    `json:"institution_id"`
	InstitutionName string  `json:"institution_name"`
	InstitutionCode string  `json:"institution_code"`
	MatchScore      float64 `json:"match_score"`
	MatchReasons    []string `json:"match_reasons"`
	EstimatedRate   *float64 `json:"estimated_rate,omitempty"`
	EstimatedFees   *float64 `json:"estimated_fees,omitempty"`
	ProductSuggestions []string `json:"product_suggestions,omitempty"`
}

// RiskAssessmentResult represents risk assessment from N8N
type RiskAssessmentResult struct {
	OverallRiskScore    float64           `json:"overall_risk_score"`
	RiskLevel           string            `json:"risk_level"` // "low", "medium", "high"
	CreditworthinessScore float64         `json:"creditworthiness_score"`
	CashFlowStability   string            `json:"cash_flow_stability"`
	BusinessRiskFactors []string          `json:"business_risk_factors"`
	FinancialIndicators map[string]float64 `json:"financial_indicators"`
	Recommendations     []string          `json:"recommendations"`
}

// DataRefreshReminderRequest represents request to refresh expired data
type DataRefreshReminderRequest struct {
	LoanApplicationID uint   `json:"loan_application_id"`
	DataType          string `json:"data_type"` // "kyb" or "financial_statements"
	ReminderType      string `json:"reminder_type"` // "due_soon", "overdue"
}

// AdminReviewRequest represents admin review request
type AdminReviewRequest struct {
	Action          string  `json:"action" validate:"required,oneof=approve reject request_info"`
	ReviewNotes     string  `json:"review_notes"`
	ApprovedAmount  *float64 `json:"approved_amount,omitempty" validate:"omitempty,gt=0"`
	RejectionReason *string `json:"rejection_reason,omitempty"`
}

// LoanApplicationListResponse represents paginated loan application list
type LoanApplicationListResponse struct {
	Applications []LoanApplicationResponse `json:"applications"`
	TotalCount   int64                     `json:"total_count"`
	Page         int                       `json:"page"`
	PageSize     int                       `json:"page_size"`
	TotalPages   int                       `json:"total_pages"`
}

// SendToFinancialInstitutionRequest represents admin request to send loan application to financial institution
type SendToFinancialInstitutionRequest struct {
	FinancialInstitutionID    uint   `json:"financial_institution_id" validate:"required"`
	FinancialInstitutionName  string `json:"financial_institution_name" validate:"required"`
	FinancialInstitutionEmail string `json:"financial_institution_email" validate:"required,email"`
	Notes                     string `json:"notes"`
}

// LoanApplicationStatsResponse represents loan application statistics
type LoanApplicationStatsResponse struct {
	TotalApplications     int64   `json:"total_applications"`
	PendingApplications   int64   `json:"pending_applications"`
	ApprovedApplications  int64   `json:"approved_applications"`
	RejectedApplications  int64   `json:"rejected_applications"`
	DisbursedApplications int64   `json:"disbursed_applications"`
	TotalAmountRequested  float64 `json:"total_amount_requested"`
	TotalAmountApproved   float64 `json:"total_amount_approved"`
	TotalAmountDisbursed  float64 `json:"total_amount_disbursed"`
	ExpiredDataCount      int64   `json:"expired_data_count"`
}

// Helper methods to convert between models and DTOs

// ToLoanApplicationResponse converts a loan application model to response DTO
func ToLoanApplicationResponse(app *models.LoanApplication) LoanApplicationResponse {
	response := LoanApplicationResponse{
		ID:                   app.ID,
		ApplicationReference: app.ApplicationReference,
		Source:               string(app.Source),
		RequestedAmount:      app.RequestedAmount,
		Currency:             app.Currency,
		Purpose:              app.Purpose,
		Status:               string(app.Status),
		SubmittedAt:          app.SubmittedAt,
		LastUpdatedAt:        app.LastUpdatedAt,
		ReviewNotes:          app.ReviewNotes,
		ApprovedAmount:       app.ApprovedAmount,
		ApprovedAt:           app.ApprovedAt,
		DisbursedAmount:      app.DisbursedAmount,
		DisbursedAt:          app.DisbursedAt,
		CreatedAt:            app.CreatedAt,
		UpdatedAt:            app.UpdatedAt,
	}

	// Check for expired data
	response.IsDataExpired = false
	if app.KYBInformation != nil && app.KYBInformation.IsExpired {
		response.IsDataExpired = true
		response.NextRefreshDate = &app.KYBInformation.RefreshDueDate
	}

	// Add recommended institution if available
	if app.RecommendedInstitution != nil {
		response.RecommendedInstitution = &FinancialInstitutionSummary{
			ID:              app.RecommendedInstitution.ID,
			Name:            app.RecommendedInstitution.Name,
			Code:            app.RecommendedInstitution.Code,
			InterestRateMin: app.RecommendedInstitution.InterestRateMin,
			InterestRateMax: app.RecommendedInstitution.InterestRateMax,
			ProcessingFee:   app.RecommendedInstitution.ProcessingFee,
		}
	}

	// Add KYB information if available
	if app.KYBInformation != nil {
		kybResponse := ToKYBInformationResponse(app.KYBInformation)
		response.KYBInformation = &kybResponse
	}

	// Add financial statements if available
	for _, stmt := range app.FinancialStatements {
		response.FinancialStatements = append(response.FinancialStatements, 
			ToFinancialStatementResponse(&stmt))
	}

	return response
}

// ToKYBInformationResponse converts KYB model to response DTO
func ToKYBInformationResponse(kyb *models.KYBInformation) KYBInformationResponse {
	response := KYBInformationResponse{
		ID:                      kyb.ID,
		BusinessName:            kyb.BusinessName,
		BusinessRegistrationNo:  kyb.BusinessRegistrationNo,
		BusinessType:            kyb.BusinessType,
		IndustryType:            kyb.IndustryType,
		BusinessAddress:         kyb.BusinessAddress,
		TaxIdentificationNumber: kyb.TaxIdentificationNumber,
		YearsInOperation:        kyb.YearsInOperation,
		AnnualRevenue:           kyb.AnnualRevenue,
		MonthlyRevenue:          kyb.MonthlyRevenue,
		NumberOfEmployees:       kyb.NumberOfEmployees,
		BusinessDescription:     kyb.BusinessDescription,
		SubmittedAt:             kyb.SubmittedAt,
		LastRefreshedAt:         kyb.LastRefreshedAt,
		RefreshDueDate:          kyb.RefreshDueDate,
		IsExpired:               kyb.IsExpired,
		VerificationStatus:      kyb.VerificationStatus,
		VerifiedAt:              kyb.VerifiedAt,
		CreatedAt:               kyb.CreatedAt,
		UpdatedAt:               kyb.UpdatedAt,
	}

	// Parse JSON fields (directors and beneficial owners would need proper JSON unmarshaling)
	// This is a simplified version - in production, you'd want proper JSON handling
	
	return response
}

// ToFinancialStatementResponse converts financial statement model to response DTO
func ToFinancialStatementResponse(stmt *models.FinancialStatement) FinancialStatementResponse {
	return FinancialStatementResponse{
		ID:                  stmt.ID,
		Type:                string(stmt.Type),
		ProviderName:        stmt.ProviderName,
		AccountNumber:       stmt.AccountNumber,
		StatementPeriodFrom: stmt.StatementPeriodFrom,
		StatementPeriodTo:   stmt.StatementPeriodTo,
		OpeningBalance:      stmt.OpeningBalance,
		ClosingBalance:      stmt.ClosingBalance,
		AverageBalance:      &stmt.AverageBalance,
		TotalCredits:        stmt.TotalCredits,
		TotalDebits:         stmt.TotalDebits,
		TransactionCount:    stmt.TransactionCount,
		SubmittedAt:         stmt.SubmittedAt,
		LastRefreshedAt:     stmt.LastRefreshedAt,
		RefreshDueDate:      stmt.RefreshDueDate,
		IsExpired:           stmt.IsExpired,
		N8NProcessingStatus: stmt.N8NProcessingStatus,
		N8NProcessedAt:      stmt.N8NProcessedAt,
		CreatedAt:           stmt.CreatedAt,
		UpdatedAt:           stmt.UpdatedAt,
	}
}

// KYBReviewRequest represents an admin request to review KYB information
type KYBReviewRequest struct {
	Status string  `json:"status" validate:"required,oneof=approved rejected pending"`
	Notes  *string `json:"notes,omitempty"`
}