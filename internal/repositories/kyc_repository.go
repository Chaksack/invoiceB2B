package repositories

import (
	"context"
	"errors"
	"invoiceB2B/internal/models"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// KYCRepository interface defines operations for KYC data management.
type KYCRepository interface {
	CreateOrUpdate(ctx context.Context, kycDetail *models.KYCDetail) (*models.KYCDetail, error)
	FindByUserID(ctx context.Context, userID uint) (*models.KYCDetail, error)
	FindByID(ctx context.Context, kycID uint) (*models.KYCDetail, error)
	CountByStatus(ctx context.Context, status models.KYCStatus) (int64, error) // Added for analytics
}

type kycRepository struct {
	db *gorm.DB
}

// NewKYCRepository creates a new instance of KYCRepository.
func NewKYCRepository(db *gorm.DB) KYCRepository {
	return &kycRepository{db: db}
}

// CreateOrUpdate creates a new KYC record or updates an existing one based on UserID.
func (r *kycRepository) CreateOrUpdate(ctx context.Context, kycDetail *models.KYCDetail) (*models.KYCDetail, error) {
	// Ensure all relevant fields are included in DoUpdates for a proper update on conflict.
	err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"status", "submitted_at", "reviewed_by_id", "reviewed_at", "rejection_reason", "documents_info", "updated_at"}),
	}).Create(kycDetail).Error

	if err != nil {
		log.Printf("Error creating/updating KYC detail for user %d in DB: %v", kycDetail.UserID, err)
		return nil, err
	}
	// Re-fetch to ensure all fields (like ID, CreatedAt if new) are populated.
	// This is especially useful if the Create operation happened.
	// If performance is critical and only updated fields are needed, this can be omitted.
	// For now, let's re-fetch for consistency.
	var fetchedKYC models.KYCDetail
	if findErr := r.db.WithContext(ctx).Where("id = ?", kycDetail.ID).First(&fetchedKYC).Error; findErr != nil {
		log.Printf("Error re-fetching KYC detail ID %d after create/update: %v", kycDetail.ID, findErr)
		return nil, findErr // Or return kycDetail if re-fetch fails but create/update succeeded
	}
	return &fetchedKYC, nil
}

// FindByUserID retrieves a KYC record by UserID.
func (r *kycRepository) FindByUserID(ctx context.Context, userID uint) (*models.KYCDetail, error) {
	var kycDetail models.KYCDetail
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&kycDetail).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound // Propagate gorm.ErrRecordNotFound
		}
		log.Printf("Error finding KYC detail by user ID %d in DB: %v", userID, err)
		return nil, err
	}
	return &kycDetail, nil
}

// FindByID retrieves a KYC record by its primary ID.
func (r *kycRepository) FindByID(ctx context.Context, kycID uint) (*models.KYCDetail, error) {
	var kycDetail models.KYCDetail
	if err := r.db.WithContext(ctx).Where("id = ?", kycID).First(&kycDetail).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound // Propagate gorm.ErrRecordNotFound
		}
		log.Printf("Error finding KYC detail by ID %d in DB: %v", kycID, err)
		return nil, err
	}
	return &kycDetail, nil
}

// CountByStatus counts KYC records based on their status.
// If status is an empty string, it counts all KYC records.
func (r *kycRepository) CountByStatus(ctx context.Context, status models.KYCStatus) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&models.KYCDetail{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Count(&count).Error; err != nil {
		log.Printf("Error counting KYC records by status '%s' in DB: %v", status, err)
		return 0, err
	}
	return count, nil
}
