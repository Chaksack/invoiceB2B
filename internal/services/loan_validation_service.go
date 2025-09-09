package services

import (
	"fmt"
	"strings"

	"invoiceB2B/internal/dtos"
	"invoiceB2B/internal/models"
)

// LoanValidationService handles business rule validation for loan applications
type LoanValidationService interface {
	ValidateCreateLoanRequest(request *dtos.CreateLoanApplicationRequest) error
	ValidateManualLoanRequest(request *dtos.ManualLoanInputRequest) error
	ValidateStatusTransition(currentStatus, newStatus models.LoanApplicationStatus) error
	ValidateKYBInformation(request *dtos.SubmitKYBInformationRequest) error
	ValidateFinancialStatement(request *dtos.SubmitFinancialStatementRequest) error
}

type loanValidationService struct {
	// Add dependencies like configuration service if needed
}

// NewLoanValidationService creates a new loan validation service
func NewLoanValidationService() LoanValidationService {
	return &loanValidationService{}
}

// Business rules constants
const (
	MinLoanAmount = 1000.0    // Minimum loan amount in USD
	MaxLoanAmount = 1000000.0 // Maximum loan amount in USD
	
	MinYearsInOperation = 1   // Minimum years business must be in operation
	MaxBusinessNameLength = 255
	MaxPurposeLength = 1000
)

// ValidateCreateLoanRequest validates loan application creation request
func (v *loanValidationService) ValidateCreateLoanRequest(request *dtos.CreateLoanApplicationRequest) error {
	if request == nil {
		return fmt.Errorf("loan application request cannot be nil")
	}

	// Validate loan amount
	if request.RequestedAmount <= 0 {
		return fmt.Errorf("requested amount must be greater than 0")
	}
	
	if request.RequestedAmount < MinLoanAmount {
		return fmt.Errorf("requested amount must be at least $%.2f", MinLoanAmount)
	}
	
	if request.RequestedAmount > MaxLoanAmount {
		return fmt.Errorf("requested amount cannot exceed $%.2f", MaxLoanAmount)
	}

	// Validate currency
	if err := v.validateCurrency(request.Currency); err != nil {
		return err
	}

	// Validate source
	validSources := []string{"manual", "invoice", "contract"}
	if !v.isValidSource(request.Source, validSources) {
		return fmt.Errorf("invalid source: %s. Valid sources are: %s", 
			request.Source, strings.Join(validSources, ", "))
	}

	// Validate purpose length
	if len(request.Purpose) > MaxPurposeLength {
		return fmt.Errorf("purpose cannot exceed %d characters", MaxPurposeLength)
	}

	return nil
}

// ValidateManualLoanRequest validates manual loan application request
func (v *loanValidationService) ValidateManualLoanRequest(request *dtos.ManualLoanInputRequest) error {
	if request == nil {
		return fmt.Errorf("manual loan application request cannot be nil")
	}

	// Validate basic loan application fields
	createReq := &dtos.CreateLoanApplicationRequest{
		Source:          request.Source,
		RequestedAmount: request.RequestedAmount,
		Currency:        request.Currency,
		Purpose:         request.Purpose,
	}
	
	if err := v.ValidateCreateLoanRequest(createReq); err != nil {
		return err
	}

	// Validate KYB information
	if err := v.ValidateKYBInformation(&request.KYBInformation); err != nil {
		return fmt.Errorf("invalid KYB information: %w", err)
	}

	// Validate financial statements
	if len(request.FinancialStatements) == 0 {
		return fmt.Errorf("at least one financial statement is required")
	}

	for i, stmt := range request.FinancialStatements {
		if err := v.ValidateFinancialStatement(&stmt); err != nil {
			return fmt.Errorf("invalid financial statement at index %d: %w", i, err)
		}
	}

	return nil
}

// ValidateStatusTransition validates if a status transition is allowed
func (v *loanValidationService) ValidateStatusTransition(currentStatus, newStatus models.LoanApplicationStatus) error {
	validTransitions := map[models.LoanApplicationStatus][]models.LoanApplicationStatus{
		models.LoanApplicationPending: {
			models.LoanApplicationKYBRequired,
			models.LoanApplicationFinancialRequired,
			models.LoanApplicationUnderReview,
			models.LoanApplicationRejected,
		},
		models.LoanApplicationKYBRequired: {
			models.LoanApplicationFinancialRequired,
			models.LoanApplicationUnderReview,
			models.LoanApplicationRejected,
		},
		models.LoanApplicationFinancialRequired: {
			models.LoanApplicationUnderReview,
			models.LoanApplicationRejected,
		},
		models.LoanApplicationUnderReview: {
			models.LoanApplicationProcessing,
			models.LoanApplicationApproved,
			models.LoanApplicationRejected,
		},
		models.LoanApplicationProcessing: {
			models.LoanApplicationApproved,
			models.LoanApplicationRejected,
		},
		models.LoanApplicationApproved: {
			models.LoanApplicationDisbursed,
			models.LoanApplicationRejected, // Can be rejected even after approval
		},
		models.LoanApplicationDisbursed: {
			models.LoanApplicationCompleted,
		},
		// Terminal states - no transitions allowed
		models.LoanApplicationRejected:  {},
		models.LoanApplicationCompleted: {},
	}

	validNextStates, exists := validTransitions[currentStatus]
	if !exists {
		return fmt.Errorf("unknown current status: %s", currentStatus)
	}

	for _, validStatus := range validNextStates {
		if newStatus == validStatus {
			return nil
		}
	}

	return fmt.Errorf("invalid status transition from %s to %s", currentStatus, newStatus)
}

// ValidateKYBInformation validates KYB information submission
func (v *loanValidationService) ValidateKYBInformation(request *dtos.SubmitKYBInformationRequest) error {
	if request == nil {
		return fmt.Errorf("KYB information request cannot be nil")
	}

	// Validate business name
	if strings.TrimSpace(request.BusinessName) == "" {
		return fmt.Errorf("business name is required")
	}
	
	if len(request.BusinessName) > MaxBusinessNameLength {
		return fmt.Errorf("business name cannot exceed %d characters", MaxBusinessNameLength)
	}

	// Validate business registration number
	if strings.TrimSpace(request.BusinessRegistrationNo) == "" {
		return fmt.Errorf("business registration number is required")
	}

	// Validate business type
	if strings.TrimSpace(request.BusinessType) == "" {
		return fmt.Errorf("business type is required")
	}

	// Validate industry type
	if strings.TrimSpace(request.IndustryType) == "" {
		return fmt.Errorf("industry type is required")
	}

	// Validate address
	if strings.TrimSpace(request.BusinessAddress) == "" {
		return fmt.Errorf("business address is required")
	}

	// Validate tax identification number
	if strings.TrimSpace(request.TaxIdentificationNumber) == "" {
		return fmt.Errorf("tax identification number is required")
	}

	// Validate years in operation
	if request.YearsInOperation < MinYearsInOperation {
		return fmt.Errorf("business must be in operation for at least %d year(s)", MinYearsInOperation)
	}

	// Validate business description
	if strings.TrimSpace(request.BusinessDescription) == "" {
		return fmt.Errorf("business description is required")
	}

	// Validate financial information if provided
	if request.AnnualRevenue != nil && *request.AnnualRevenue < 0 {
		return fmt.Errorf("annual revenue cannot be negative")
	}

	if request.MonthlyRevenue != nil && *request.MonthlyRevenue < 0 {
		return fmt.Errorf("monthly revenue cannot be negative")
	}

	if request.NumberOfEmployees != nil && *request.NumberOfEmployees < 0 {
		return fmt.Errorf("number of employees cannot be negative")
	}

	return nil
}

// ValidateFinancialStatement validates financial statement submission
func (v *loanValidationService) ValidateFinancialStatement(request *dtos.SubmitFinancialStatementRequest) error {
	if request == nil {
		return fmt.Errorf("financial statement request cannot be nil")
	}

	// Validate statement type
	validTypes := []string{"bank", "mobile_money"}
	if !v.contains(validTypes, request.Type) {
		return fmt.Errorf("invalid statement type: %s. Valid types are: %s", 
			request.Type, strings.Join(validTypes, ", "))
	}

	// Validate provider name
	if strings.TrimSpace(request.ProviderName) == "" {
		return fmt.Errorf("provider name is required")
	}

	// Validate account number
	if strings.TrimSpace(request.AccountNumber) == "" {
		return fmt.Errorf("account number is required")
	}

	// Validate statement periods
	if request.StatementPeriodFrom.IsZero() {
		return fmt.Errorf("statement period from date is required")
	}

	if request.StatementPeriodTo.IsZero() {
		return fmt.Errorf("statement period to date is required")
	}

	if request.StatementPeriodFrom.After(request.StatementPeriodTo) {
		return fmt.Errorf("statement period from date cannot be after to date")
	}

	// Validate balances
	if request.OpeningBalance < 0 {
		return fmt.Errorf("opening balance cannot be negative")
	}

	if request.ClosingBalance < 0 {
		return fmt.Errorf("closing balance cannot be negative")
	}

	if request.TotalCredits < 0 {
		return fmt.Errorf("total credits cannot be negative")
	}

	if request.TotalDebits < 0 {
		return fmt.Errorf("total debits cannot be negative")
	}

	if request.TransactionCount < 0 {
		return fmt.Errorf("transaction count cannot be negative")
	}

	return nil
}

// validateCurrency validates currency code
func (v *loanValidationService) validateCurrency(currency string) error {
	if currency == "" {
		return fmt.Errorf("currency is required")
	}

	// List of supported currencies
	supportedCurrencies := []string{"USD", "EUR", "GBP", "KES", "UGX", "TZS", "GHS", "NGN", "ZAR"}
	
	for _, supported := range supportedCurrencies {
		if strings.ToUpper(currency) == supported {
			return nil
		}
	}

	return fmt.Errorf("unsupported currency: %s. Supported currencies are: %s", 
		currency, strings.Join(supportedCurrencies, ", "))
}

// isValidSource checks if the source is valid
func (v *loanValidationService) isValidSource(source string, validSources []string) bool {
	return v.contains(validSources, source)
}

// contains checks if a slice contains a string
func (v *loanValidationService) contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}