package repositories

import (
	"context"
	"errors"
	"invoiceB2B/internal/models"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type KYCRepository interface {
	CreateOrUpdate(ctx context.Context, kycDetail *models.KYCDetail) (*models.KYCDetail, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) (*models.KYCDetail, error)
	FindByID(ctx context.Context, kycID uuid.UUID) (*models.KYCDetail, error)
	// UpdateStatus(ctx context.Context, kycID uuid.UUID, status models.KYCStatus, reviewedByID *uuid.UUID, rejectionReason *string) error
}

type kycRepository struct {
	db *gorm.DB
}

func NewKYCRepository(db *gorm.DB) KYCRepository {
	return &kycRepository{db: db}
}

func (r *kycRepository) CreateOrUpdate(ctx context.Context, kycDetail *models.KYCDetail) (*models.KYCDetail, error) {
	err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}}, // Assuming user_id is a unique key for KYC details
		DoUpdates: clause.AssignmentColumns([]string{"status", "submitted_at", "reviewed_by_id", "reviewed_at", "rejection_reason", "documents_info", "updated_at"}),
	}).Create(kycDetail).Error

	if err != nil {
		log.Printf("Error creating/updating KYC detail for user %s in DB: %v", kycDetail.UserID, err)
		return nil, err
	}
	return kycDetail, nil
}

func (r *kycRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*models.KYCDetail, error) {
	var kycDetail models.KYCDetail
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&kycDetail).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		log.Printf("Error finding KYC detail by user ID %s in DB: %v", userID, err)
		return nil, err
	}
	return &kycDetail, nil
}

func (r *kycRepository) FindByID(ctx context.Context, kycID uuid.UUID) (*models.KYCDetail, error) {
	var kycDetail models.KYCDetail
	if err := r.db.WithContext(ctx).Where("id = ?", kycID).First(&kycDetail).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		log.Printf("Error finding KYC detail by ID %s in DB: %v", kycID, err)
		return nil, err
	}
	return &kycDetail, nil
}
