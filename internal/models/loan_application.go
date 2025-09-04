package models

import (
	"gorm.io/gorm"
	"time"
)

type LoanApplicationStatus string

const (
	LoanApplicationPending           LoanApplicationStatus = "pending"
	LoanApplicationKYBRequired       LoanApplicationStatus = "kyb_required"
	LoanApplicationFinancialRequired LoanApplicationStatus = "financial_required"
	LoanApplicationUnderReview       LoanApplicationStatus = "under_review"
	LoanApplicationProcessing        LoanApplicationStatus = "processing"
	LoanApplicationApproved          LoanApplicationStatus = "approved"
	LoanApplicationRejected          LoanApplicationStatus = "rejected"
	LoanApplicationDisbursed         LoanApplicationStatus = "disbursed"
	LoanApplicationCompleted         LoanApplicationStatus = "completed"
)

type LoanApplicationSource string

const (
	LoanApplicationSourceManual   LoanApplicationSource = "manual"
	LoanApplicationSourceInvoice  LoanApplicationSource = "invoice"
	LoanApplicationSourceContract LoanApplicationSource = "contract"
)

// LoanApplication represents a loan application with KYB and financial data
type LoanApplication struct {
	gorm.Model
	UserID uint `gorm:"not null;index"`
	User   User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// Application Details
	ApplicationReference string                `gorm:"type:varchar(50);uniqueIndex;not null"`
	Source               LoanApplicationSource `gorm:"type:varchar(20);not null"`
	RequestedAmount      float64               `gorm:"type:decimal(18,2);not null"`
	Currency             string                `gorm:"type:varchar(3);not null;default:'USD'"`
	Purpose              string                `gorm:"type:text;null"`
	
	// Status and Workflow
	Status           LoanApplicationStatus `gorm:"type:varchar(30);default:'pending';not null"`
	SubmittedAt      *time.Time
	LastUpdatedAt    time.Time `gorm:"autoUpdateTime"`
	
	// Admin Review
	ReviewedByID    *uint  `gorm:"null"`
	ReviewedBy      *Staff `gorm:"foreignKey:ReviewedByID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ReviewedAt      *time.Time
	ReviewNotes     *string `gorm:"type:text;null"`
	RejectionReason *string `gorm:"type:text;null"`

	// Approval and Disbursement
	ApprovedByID    *uint  `gorm:"null"`
	ApprovedBy      *Staff `gorm:"foreignKey:ApprovedByID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ApprovedAt      *time.Time
	ApprovedAmount  *float64 `gorm:"type:decimal(18,2);null"`
	
	DisbursedByID   *uint  `gorm:"null"`
	DisbursedBy     *Staff `gorm:"foreignKey:DisbursedByID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	DisbursedAt     *time.Time
	DisbursedAmount *float64 `gorm:"type:decimal(18,2);null"`

	// Financial Institution Matching
	MatchedInstitutions      string `gorm:"type:jsonb;null"` // Array of matched institutions with scores
	RecommendedInstitutionID *uint  `gorm:"null"`
	RecommendedInstitution   *FinancialInstitution `gorm:"foreignKey:RecommendedInstitutionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	// Supporting Documents
	SupportingDocumentPath *string `gorm:"type:varchar(500);null"` // Path to uploaded invoice/contract
	ProcessedDocumentData  string  `gorm:"type:jsonb;null"`        // N8N processed data
	
	// N8N Integration
	N8NWorkflowExecutionID *string `gorm:"type:varchar(100);null"`
	N8NProcessingStatus    *string `gorm:"type:varchar(50);null"`
	N8NProcessedAt         *time.Time

	// Relationships
	KYBInformation     *KYBInformation      `gorm:"foreignKey:LoanApplicationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	FinancialStatements []FinancialStatement `gorm:"foreignKey:LoanApplicationID"`
}

// KYBInformation represents Know Your Business information
type KYBInformation struct {
	gorm.Model
	LoanApplicationID uint            `gorm:"uniqueIndex;not null"`
	LoanApplication   LoanApplication `gorm:"foreignKey:LoanApplicationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// Business Information
	BusinessName            string `gorm:"type:varchar(255);not null"`
	BusinessRegistrationNo  string `gorm:"type:varchar(100);not null"`
	BusinessType            string `gorm:"type:varchar(100);not null"`
	IndustryType            string `gorm:"type:varchar(100);not null"`
	BusinessAddress         string `gorm:"type:text;not null"`
	TaxIdentificationNumber string `gorm:"type:varchar(100);not null"`
	YearsInOperation        int    `gorm:"not null"`
	
	// Key Personnel Information
	DirectorsInformation string `gorm:"type:jsonb;not null"` // Array of director details
	BeneficialOwners     string `gorm:"type:jsonb;not null"` // Array of beneficial owner details
	
	// Financial Overview
	AnnualRevenue       *float64 `gorm:"type:decimal(18,2);null"`
	MonthlyRevenue      *float64 `gorm:"type:decimal(18,2);null"`
	NumberOfEmployees   *int     `gorm:"null"`
	BusinessDescription string   `gorm:"type:text;not null"`
	
	// Documentation
	BusinessLicensePath       *string `gorm:"type:varchar(500);null"`
	TaxCertificatePath        *string `gorm:"type:varchar(500);null"`
	FinancialStatementsPath   *string `gorm:"type:varchar(500);null"`
	OwnershipStructurePath    *string `gorm:"type:varchar(500);null"`
	
	// Data Freshness Tracking
	SubmittedAt     time.Time  `gorm:"not null"`
	LastRefreshedAt *time.Time
	RefreshDueDate  time.Time  `gorm:"not null"` // 3 months from submission/last refresh
	IsExpired       bool       `gorm:"default:false"`
	
	// Verification Status
	VerificationStatus string     `gorm:"type:varchar(50);default:'pending';not null"`
	VerifiedAt         *time.Time
	VerifiedByID       *uint  `gorm:"null"`
	VerifiedBy         *Staff `gorm:"foreignKey:VerifiedByID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	VerificationNotes  *string `gorm:"type:text;null"`
}

type FinancialStatementType string

const (
	FinancialStatementTypeBank        FinancialStatementType = "bank"
	FinancialStatementTypeMobileMoney FinancialStatementType = "mobile_money"
)

// FinancialStatement represents financial statements (bank or mobile money)
type FinancialStatement struct {
	gorm.Model
	LoanApplicationID uint            `gorm:"not null;index"`
	LoanApplication   LoanApplication `gorm:"foreignKey:LoanApplicationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// Statement Details
	Type                FinancialStatementType `gorm:"type:varchar(20);not null"`
	ProviderName        string                 `gorm:"type:varchar(255);not null"` // Bank name or mobile money provider
	AccountNumber       string                 `gorm:"type:varchar(100);not null"`
	StatementPeriodFrom time.Time              `gorm:"not null"`
	StatementPeriodTo   time.Time              `gorm:"not null"`
	
	// Financial Data
	OpeningBalance    float64 `gorm:"type:decimal(18,2);not null"`
	ClosingBalance    float64 `gorm:"type:decimal(18,2);not null"`
	AverageBalance    float64 `gorm:"type:decimal(18,2);null"`
	TotalCredits      float64 `gorm:"type:decimal(18,2);not null"`
	TotalDebits       float64 `gorm:"type:decimal(18,2);not null"`
	TransactionCount  int     `gorm:"not null"`
	
	// Analysis from N8N
	CashFlowAnalysis    string `gorm:"type:jsonb;null"` // Processed cash flow data
	TransactionPatterns string `gorm:"type:jsonb;null"` // Transaction pattern analysis
	RiskIndicators      string `gorm:"type:jsonb;null"` // Risk assessment data
	
	// Document Storage
	OriginalFilePath      string  `gorm:"type:varchar(500);not null"`
	ProcessedDataPath     *string `gorm:"type:varchar(500);null"`
	
	// Data Freshness Tracking
	SubmittedAt     time.Time  `gorm:"not null"`
	LastRefreshedAt *time.Time
	RefreshDueDate  time.Time  `gorm:"not null"` // 3 months from submission/last refresh
	IsExpired       bool       `gorm:"default:false"`
	
	// N8N Processing
	N8NProcessingStatus *string `gorm:"type:varchar(50);null"`
	N8NProcessedAt      *time.Time
	N8NWorkflowID       *string `gorm:"type:varchar(100);null"`
}

// DataRefreshTracker manages the 3-month refresh cycle for KYB and financial data
type DataRefreshTracker struct {
	gorm.Model
	LoanApplicationID uint            `gorm:"not null;index"`
	LoanApplication   LoanApplication `gorm:"foreignKey:LoanApplicationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	
	// Tracking Information
	DataType        string    `gorm:"type:varchar(50);not null"` // "kyb" or "financial_statements"
	LastRefreshDate time.Time `gorm:"not null"`
	NextRefreshDate time.Time `gorm:"not null"`
	
	// Notification Management
	ReminderSentAt    *time.Time
	OverdueNotifiedAt *time.Time
	IsOverdue         bool `gorm:"default:false"`
	
	// Status
	Status string `gorm:"type:varchar(50);default:'active';not null"` // active, completed, suspended
}