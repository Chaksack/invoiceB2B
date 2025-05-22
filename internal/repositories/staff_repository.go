package repositories

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"invoiceB2B/internal/models"
	"log"
)

type StaffRepository interface {
	Create(ctx context.Context, staff *models.Staff) error
	Update(ctx context.Context, staff *models.Staff) error
	FindByID(ctx context.Context, id uint) (*models.Staff, error)
	FindByEmail(ctx context.Context, email string) (*models.Staff, error)
	FindAll(ctx context.Context, page, pageSize int) ([]models.Staff, int64, error)
	Delete(ctx context.Context, id uint) error
}

type staffRepository struct {
	db *gorm.DB
}

func NewStaffRepository(db *gorm.DB) StaffRepository {
	return &staffRepository{db: db}
}

func (r *staffRepository) Create(ctx context.Context, staff *models.Staff) error {
	return r.db.WithContext(ctx).Create(staff).Error
}

func (r *staffRepository) Update(ctx context.Context, staff *models.Staff) error {
	return r.db.WithContext(ctx).Save(staff).Error
}

func (r *staffRepository) FindByID(ctx context.Context, id uint) (*models.Staff, error) {
	var staff models.Staff
	if err := r.db.WithContext(ctx).First(&staff, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("staff not found")
		}
		return nil, err
	}
	return &staff, nil
}

func (r *staffRepository) FindByEmail(ctx context.Context, email string) (*models.Staff, error) {
	var staff models.Staff
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&staff).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &staff, nil
}

func (r *staffRepository) FindAll(ctx context.Context, page, pageSize int) ([]models.Staff, int64, error) {
	var staffList []models.Staff
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Staff{})

	if err := query.Count(&total).Error; err != nil {
		log.Printf("Error counting staff: %v", err)
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&staffList).Error; err != nil {
		log.Printf("Error fetching staff: %v", err)
		return nil, 0, err
	}
	return staffList, total, nil
}

func (r *staffRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Staff{}, id).Error
}
