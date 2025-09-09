package repositories

import (
	"fmt"
	"time"

	"invoiceB2B/internal/models"
	"gorm.io/gorm"
)

// LoanApplicationRepository interface defines loan application data access operations
type LoanApplicationRepository interface {
	Create(loanApp *models.LoanApplication) error
	GetByID(id uint) (*models.LoanApplication, error)
	GetByIDWithUserID(id uint, userID uint) (*models.LoanApplication, error)
	GetByUserID(userID uint, offset, limit int) ([]models.LoanApplication, int64, error)
	Update(loanApp *models.LoanApplication) error
	UpdateStatus(id uint, status models.LoanApplicationStatus, notes string) error
	GetExpiredApplications() ([]models.LoanApplication, error)
	GetStats() (*LoanApplicationStats, error)
	GetApplicationsRequiringRefresh() ([]models.DataRefreshTracker, error)
	
	// KYB Information operations
	CreateKYBInformation(kyb *models.KYBInformation) error
	UpdateKYBInformation(kyb *models.KYBInformation) error
	
	// Financial Statement operations
	CreateFinancialStatement(stmt *models.FinancialStatement) error
	UpdateFinancialStatement(stmt *models.FinancialStatement) error
	
	// Data Refresh Tracker operations
	CreateDataRefreshTracker(tracker *models.DataRefreshTracker) error
	UpdateDataRefreshTracker(tracker *models.DataRefreshTracker) error
	
	// Transaction support
	WithTransaction(tx *gorm.DB) LoanApplicationRepository
}

// LoanApplicationStats represents aggregated loan application statistics
type LoanApplicationStats struct {
	TotalApplications    int64
	PendingApplications  int64
	ApprovedApplications int64
	RejectedApplications int64
	TotalApprovedAmount  float64
	AverageProcessingTime float64
}

type loanApplicationRepository struct {
	db *gorm.DB
}

// NewLoanApplicationRepository creates a new loan application repository
func NewLoanApplicationRepository(db *gorm.DB) LoanApplicationRepository {
	return &loanApplicationRepository{db: db}
}

// WithTransaction returns a new repository instance using the provided transaction
func (r *loanApplicationRepository) WithTransaction(tx *gorm.DB) LoanApplicationRepository {
	return &loanApplicationRepository{db: tx}
}

// Create creates a new loan application
func (r *loanApplicationRepository) Create(loanApp *models.LoanApplication) error {
	if err := r.db.Create(loanApp).Error; err != nil {
		return fmt.Errorf("failed to create loan application: %w", err)
	}
	return nil
}

// GetByID retrieves a loan application by ID with all associations
func (r *loanApplicationRepository) GetByID(id uint) (*models.LoanApplication, error) {
	var loanApp models.LoanApplication
	err := r.db.Preload("User").
		Preload("KYBInformation").
		Preload("FinancialStatements").
		Preload("RecommendedInstitution").
		Preload("ReviewedBy").
		Preload("ApprovedBy").
		Preload("DisbursedBy").
		First(&loanApp, id).Error
	
	if err != nil {
		return nil, fmt.Errorf("loan application not found: %w", err)
	}
	return &loanApp, nil
}

// GetByIDWithUserID retrieves a loan application by ID and user ID
func (r *loanApplicationRepository) GetByIDWithUserID(id uint, userID uint) (*models.LoanApplication, error) {
	var loanApp models.LoanApplication
	err := r.db.Preload("User").
		Preload("KYBInformation").
		Preload("FinancialStatements").
		Preload("RecommendedInstitution").
		Preload("ReviewedBy").
		Preload("ApprovedBy").
		Preload("DisbursedBy").
		Where("user_id = ?", userID).
		First(&loanApp, id).Error
	
	if err != nil {
		return nil, fmt.Errorf("loan application not found: %w", err)
	}
	return &loanApp, nil
}

// GetByUserID retrieves loan applications for a specific user with pagination
func (r *loanApplicationRepository) GetByUserID(userID uint, offset, limit int) ([]models.LoanApplication, int64, error) {
	var loanApps []models.LoanApplication
	var total int64
	
	// Count total records
	if err := r.db.Model(&models.LoanApplication{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count loan applications: %w", err)
	}
	
	// Fetch paginated records
	err := r.db.Preload("KYBInformation").
		Preload("FinancialStatements").
		Preload("RecommendedInstitution").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&loanApps).Error
	
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch loan applications: %w", err)
	}
	
	return loanApps, total, nil
}

// Update updates a loan application
func (r *loanApplicationRepository) Update(loanApp *models.LoanApplication) error {
	loanApp.LastUpdatedAt = time.Now()
	if err := r.db.Save(loanApp).Error; err != nil {
		return fmt.Errorf("failed to update loan application: %w", err)
	}
	return nil
}

// UpdateStatus updates the status of a loan application
func (r *loanApplicationRepository) UpdateStatus(id uint, status models.LoanApplicationStatus, notes string) error {
	updates := map[string]interface{}{
		"status":           status,
		"last_updated_at": time.Now(),
	}
	
	if notes != "" {
		updates["review_notes"] = notes
	}
	
	if err := r.db.Model(&models.LoanApplication{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update loan application status: %w", err)
	}
	return nil
}

// GetExpiredApplications retrieves loan applications with expired data
func (r *loanApplicationRepository) GetExpiredApplications() ([]models.LoanApplication, error) {
	var loanApps []models.LoanApplication
	
	err := r.db.Preload("KYBInformation").
		Preload("FinancialStatements").
		Joins("LEFT JOIN kyb_informations ON kyb_informations.loan_application_id = loan_applications.id").
		Joins("LEFT JOIN financial_statements ON financial_statements.loan_application_id = loan_applications.id").
		Where("kyb_informations.is_expired = ? OR financial_statements.is_expired = ?", true, true).
		Distinct().
		Find(&loanApps).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to fetch expired applications: %w", err)
	}
	
	return loanApps, nil
}

// GetStats retrieves loan application statistics
func (r *loanApplicationRepository) GetStats() (*LoanApplicationStats, error) {
	var stats LoanApplicationStats
	
	// Total applications
	if err := r.db.Model(&models.LoanApplication{}).Count(&stats.TotalApplications).Error; err != nil {
		return nil, fmt.Errorf("failed to count total applications: %w", err)
	}
	
	// Pending applications
	if err := r.db.Model(&models.LoanApplication{}).
		Where("status IN ?", []string{"pending", "kyb_required", "financial_required", "under_review", "processing"}).
		Count(&stats.PendingApplications).Error; err != nil {
		return nil, fmt.Errorf("failed to count pending applications: %w", err)
	}
	
	// Approved applications
	if err := r.db.Model(&models.LoanApplication{}).
		Where("status IN ?", []string{"approved", "disbursed", "completed"}).
		Count(&stats.ApprovedApplications).Error; err != nil {
		return nil, fmt.Errorf("failed to count approved applications: %w", err)
	}
	
	// Rejected applications
	if err := r.db.Model(&models.LoanApplication{}).
		Where("status = ?", "rejected").
		Count(&stats.RejectedApplications).Error; err != nil {
		return nil, fmt.Errorf("failed to count rejected applications: %w", err)
	}
	
	// Total approved amount
	var totalAmount *float64
	if err := r.db.Model(&models.LoanApplication{}).
		Where("status IN ? AND approved_amount IS NOT NULL", []string{"approved", "disbursed", "completed"}).
		Select("SUM(approved_amount)").
		Scan(&totalAmount).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate total approved amount: %w", err)
	}
	
	if totalAmount != nil {
		stats.TotalApprovedAmount = *totalAmount
	}
	
	// Average processing time (in days)
	var avgProcessingTime *float64
	if err := r.db.Model(&models.LoanApplication{}).
		Where("status IN ? AND submitted_at IS NOT NULL AND approved_at IS NOT NULL", 
			[]string{"approved", "disbursed", "completed"}).
		Select("AVG(EXTRACT(epoch FROM (approved_at - submitted_at))/86400)").
		Scan(&avgProcessingTime).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate average processing time: %w", err)
	}
	
	if avgProcessingTime != nil {
		stats.AverageProcessingTime = *avgProcessingTime
	}
	
	return &stats, nil
}

// GetApplicationsRequiringRefresh retrieves applications that require data refresh
func (r *loanApplicationRepository) GetApplicationsRequiringRefresh() ([]models.DataRefreshTracker, error) {
	var trackers []models.DataRefreshTracker
	
	err := r.db.Preload("LoanApplication").
		Where("next_refresh_date <= ? AND status = ?", time.Now(), "active").
		Find(&trackers).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to fetch applications requiring refresh: %w", err)
	}
	
	return trackers, nil
}

// CreateKYBInformation creates KYB information
func (r *loanApplicationRepository) CreateKYBInformation(kyb *models.KYBInformation) error {
	if err := r.db.Create(kyb).Error; err != nil {
		return fmt.Errorf("failed to create KYB information: %w", err)
	}
	return nil
}

// UpdateKYBInformation updates KYB information
func (r *loanApplicationRepository) UpdateKYBInformation(kyb *models.KYBInformation) error {
	if err := r.db.Save(kyb).Error; err != nil {
		return fmt.Errorf("failed to update KYB information: %w", err)
	}
	return nil
}

// CreateFinancialStatement creates a financial statement
func (r *loanApplicationRepository) CreateFinancialStatement(stmt *models.FinancialStatement) error {
	if err := r.db.Create(stmt).Error; err != nil {
		return fmt.Errorf("failed to create financial statement: %w", err)
	}
	return nil
}

// UpdateFinancialStatement updates a financial statement
func (r *loanApplicationRepository) UpdateFinancialStatement(stmt *models.FinancialStatement) error {
	if err := r.db.Save(stmt).Error; err != nil {
		return fmt.Errorf("failed to update financial statement: %w", err)
	}
	return nil
}

// CreateDataRefreshTracker creates a data refresh tracker
func (r *loanApplicationRepository) CreateDataRefreshTracker(tracker *models.DataRefreshTracker) error {
	if err := r.db.Create(tracker).Error; err != nil {
		return fmt.Errorf("failed to create data refresh tracker: %w", err)
	}
	return nil
}

// UpdateDataRefreshTracker updates a data refresh tracker
func (r *loanApplicationRepository) UpdateDataRefreshTracker(tracker *models.DataRefreshTracker) error {
	if err := r.db.Save(tracker).Error; err != nil {
		return fmt.Errorf("failed to update data refresh tracker: %w", err)
	}
	return nil
}