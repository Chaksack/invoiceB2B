package dtos

import (
	"time"
)

// ReportRequest represents a general report request with common filters
type ReportRequest struct {
	ReportType  string     `json:"report_type" validate:"required,oneof=loan_applications invoices users activity financial_institutions overview"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	GroupBy     string     `json:"group_by,omitempty"` // daily, weekly, monthly, yearly
	Status      string     `json:"status,omitempty"`
	UserID      *uint      `json:"user_id,omitempty"`
	ExportFormat string    `json:"export_format,omitempty" validate:"omitempty,oneof=json csv pdf"`
}

// LoanApplicationReportResponse represents loan application report data
type LoanApplicationReportResponse struct {
	Summary     LoanApplicationSummary        `json:"summary"`
	Trends      []LoanApplicationTrendData    `json:"trends"`
	StatusBreakdown []StatusBreakdownData     `json:"status_breakdown"`
	AmountRanges []AmountRangeData            `json:"amount_ranges"`
	SourceBreakdown []SourceBreakdownData     `json:"source_breakdown"`
	TopApplications []LoanApplicationSummaryItem `json:"top_applications"`
}

type LoanApplicationSummary struct {
	TotalApplications     int64   `json:"total_applications"`
	TotalAmountRequested  float64 `json:"total_amount_requested"`
	TotalAmountApproved   float64 `json:"total_amount_approved"`
	TotalAmountDisbursed  float64 `json:"total_amount_disbursed"`
	AverageProcessingDays float64 `json:"average_processing_days"`
	ApprovalRate          float64 `json:"approval_rate"`
	DisbursementRate      float64 `json:"disbursement_rate"`
}

type LoanApplicationTrendData struct {
	Date            time.Time `json:"date"`
	Applications    int64     `json:"applications"`
	AmountRequested float64   `json:"amount_requested"`
	AmountApproved  float64   `json:"amount_approved"`
	AmountDisbursed float64   `json:"amount_disbursed"`
}

type StatusBreakdownData struct {
	Status string  `json:"status"`
	Count  int64   `json:"count"`
	Amount float64 `json:"amount"`
	Percentage float64 `json:"percentage"`
}

type AmountRangeData struct {
	Range      string  `json:"range"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

type SourceBreakdownData struct {
	Source     string  `json:"source"`
	Count      int64   `json:"count"`
	Amount     float64 `json:"amount"`
	Percentage float64 `json:"percentage"`
}

type LoanApplicationSummaryItem struct {
	ID                   uint      `json:"id"`
	ApplicationReference string    `json:"application_reference"`
	CompanyName         string    `json:"company_name"`
	RequestedAmount     float64   `json:"requested_amount"`
	Status              string    `json:"status"`
	SubmittedAt         time.Time `json:"submitted_at"`
}

// InvoiceReportResponse represents invoice report data
type InvoiceReportResponse struct {
	Summary         InvoiceSummary        `json:"summary"`
	Trends          []InvoiceTrendData    `json:"trends"`
	StatusBreakdown []StatusBreakdownData `json:"status_breakdown"`
	AmountRanges    []AmountRangeData     `json:"amount_ranges"`
	TopInvoices     []InvoiceSummaryItem  `json:"top_invoices"`
}

type InvoiceSummary struct {
	TotalInvoices        int64   `json:"total_invoices"`
	TotalInvoiceAmount   float64 `json:"total_invoice_amount"`
	TotalProcessed       float64 `json:"total_processed"`
	AverageInvoiceAmount float64 `json:"average_invoice_amount"`
	ProcessingRate       float64 `json:"processing_rate"`
}

type InvoiceTrendData struct {
	Date           time.Time `json:"date"`
	Invoices       int64     `json:"invoices"`
	InvoiceAmount  float64   `json:"invoice_amount"`
	ProcessedAmount float64  `json:"processed_amount"`
}

type InvoiceSummaryItem struct {
	ID            uint      `json:"id"`
	InvoiceNumber string    `json:"invoice_number"`
	CompanyName   string    `json:"company_name"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

// UserReportResponse represents user report data
type UserReportResponse struct {
	Summary       UserSummary          `json:"summary"`
	Trends        []UserTrendData      `json:"trends"`
	KYCBreakdown  []KYCBreakdownData   `json:"kyc_breakdown"`
	TopUsers      []UserSummaryItem    `json:"top_users"`
}

type UserSummary struct {
	TotalUsers        int64   `json:"total_users"`
	ActiveUsers       int64   `json:"active_users"`
	VerifiedUsers     int64   `json:"verified_users"`
	KYCCompletedUsers int64   `json:"kyc_completed_users"`
	AverageApplications float64 `json:"average_applications_per_user"`
}

type UserTrendData struct {
	Date        time.Time `json:"date"`
	NewUsers    int64     `json:"new_users"`
	ActiveUsers int64     `json:"active_users"`
}

type KYCBreakdownData struct {
	Status     string  `json:"status"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

type UserSummaryItem struct {
	ID                    uint      `json:"id"`
	CompanyName          string    `json:"company_name"`
	Email                string    `json:"email"`
	TotalApplications    int64     `json:"total_applications"`
	TotalAmountRequested float64   `json:"total_amount_requested"`
	RegistrationDate     time.Time `json:"registration_date"`
}

// ActivityReportResponse represents activity report data
type ActivityReportResponse struct {
	Summary         ActivitySummary        `json:"summary"`
	Trends          []ActivityTrendData    `json:"trends"`
	ActionBreakdown []ActionBreakdownData  `json:"action_breakdown"`
	TopActivities   []ActivitySummaryItem  `json:"top_activities"`
}

type ActivitySummary struct {
	TotalActivities int64 `json:"total_activities"`
	UniqueUsers     int64 `json:"unique_users"`
	UniqueStaff     int64 `json:"unique_staff"`
	TopAction       string `json:"top_action"`
}

type ActivityTrendData struct {
	Date       time.Time `json:"date"`
	Activities int64     `json:"activities"`
}

type ActionBreakdownData struct {
	Action     string  `json:"action"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
}

type ActivitySummaryItem struct {
	ID        uint      `json:"id"`
	Action    string    `json:"action"`
	UserEmail string    `json:"user_email"`
	Timestamp time.Time `json:"timestamp"`
	Details   string    `json:"details"`
}

// FinancialInstitutionReportResponse represents financial institution report data
type FinancialInstitutionReportResponse struct {
	Summary                 FinancialInstitutionReportSummary  `json:"summary"`
	InstitutionBreakdown    []InstitutionBreakdownData         `json:"institution_breakdown"`
	ProductBreakdown        []ProductBreakdownData             `json:"product_breakdown"`
	TopInstitutions         []FinancialInstitutionSummaryItem  `json:"top_institutions"`
}

type FinancialInstitutionReportSummary struct {
	TotalInstitutions       int64   `json:"total_institutions"`
	ActiveInstitutions      int64   `json:"active_institutions"`
	TotalProducts           int64   `json:"total_products"`
	TotalApplicationsSent   int64   `json:"total_applications_sent"`
	AverageInterestRate     float64 `json:"average_interest_rate"`
}

type InstitutionBreakdownData struct {
	InstitutionName   string  `json:"institution_name"`
	ApplicationsSent  int64   `json:"applications_sent"`
	TotalAmount       float64 `json:"total_amount"`
	AverageAmount     float64 `json:"average_amount"`
}

type ProductBreakdownData struct {
	ProductName       string  `json:"product_name"`
	InstitutionName   string  `json:"institution_name"`
	ApplicationsCount int64   `json:"applications_count"`
	InterestRate      float64 `json:"interest_rate"`
}

type FinancialInstitutionSummaryItem struct {
	ID                uint    `json:"id"`
	Name              string  `json:"name"`
	Code              string  `json:"code"`
	ProductCount      int64   `json:"product_count"`
	ApplicationsCount int64   `json:"applications_count"`
	AverageRate       float64 `json:"average_rate"`
}

// OverviewReportResponse represents a comprehensive overview report
type OverviewReportResponse struct {
	LoanApplications LoanApplicationSummary         `json:"loan_applications"`
	Invoices         InvoiceSummary                 `json:"invoices"`
	Users            UserSummary                    `json:"users"`
	Activity         ActivitySummary                `json:"activity"`
	FinancialInstitutions FinancialInstitutionReportSummary `json:"financial_institutions"`
	SystemHealth     SystemHealthData               `json:"system_health"`
	RecentActivities []ActivitySummaryItem          `json:"recent_activities"`
}

type SystemHealthData struct {
	TotalTransactions    int64   `json:"total_transactions"`
	SuccessRate          float64 `json:"success_rate"`
	AverageResponseTime  float64 `json:"average_response_time"`
	ExpiredDataCount     int64   `json:"expired_data_count"`
	PendingRefreshCount  int64   `json:"pending_refresh_count"`
}

// ExportReportResponse represents export response metadata
type ExportReportResponse struct {
	FileName    string    `json:"file_name"`
	Format      string    `json:"format"`
	Size        int64     `json:"size"`
	GeneratedAt time.Time `json:"generated_at"`
	DownloadURL string    `json:"download_url"`
}