package repositories

import (
	"context"
	"errors"
	"invoiceB2B/internal/models"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type KYCRepository interface {
	CreateOrUpdate(ctx context.Context, kycDetail *models.KYCDetail) (*models.KYCDetail, error)
	FindByUserID(ctx context.Context, userID uint) (*models.KYCDetail, error)
	FindByID(ctx context.Context, kycID uint) (*models.KYCDetail, error)
}

type kycRepository struct {
	db *gorm.DB
}

func NewKYCRepository(db *gorm.DB) KYCRepository {
	return &kycRepository{db: db}
}

func (r *kycRepository) CreateOrUpdate(ctx context.Context, kycDetail *models.KYCDetail) (*models.KYCDetail, error) {
	err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"status", "submitted_at", "reviewed_by_id", "reviewed_at", "rejection_reason", "documents_info", "updated_at"}),
	}).Create(kycDetail).Error

	if err != nil {
		log.Printf("Error creating/updating KYC detail for user %d in DB: %v", kycDetail.UserID, err)
		return nil, err
	}
	return kycDetail, nil
}

func (r *kycRepository) FindByUserID(ctx context.Context, userID uint) (*models.KYCDetail, error) {
	var kycDetail models.KYCDetail
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&kycDetail).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		log.Printf("Error finding KYC detail by user ID %d in DB: %v", userID, err)
		return nil, err
	}
	return &kycDetail, nil
}

func (r *kycRepository) FindByID(ctx context.Context, kycID uint) (*models.KYCDetail, error) {
	var kycDetail models.KYCDetail
	if err := r.db.WithContext(ctx).Where("id = ?", kycID).First(&kycDetail).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		log.Printf("Error finding KYC detail by ID %d in DB: %v", kycID, err)
		return nil, err
	}
	return &kycDetail, nil
}
