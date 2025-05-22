package repositories

import (
	"context"
	"invoiceB2B/internal/models"
	"log"

	"gorm.io/gorm"
)

type ActivityLogRepository interface {
	Create(ctx context.Context, logEntry *models.ActivityLog) error
	FindAll(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.ActivityLog, int64, error)
}

type activityLogRepository struct {
	db *gorm.DB
}

func NewActivityLogRepository(db *gorm.DB) ActivityLogRepository {
	return &activityLogRepository{db: db}
}

func (r *activityLogRepository) Create(ctx context.Context, logEntry *models.ActivityLog) error {
	if err := r.db.WithContext(ctx).Create(logEntry).Error; err != nil {
		log.Printf("Error creating activity log in DB: %v", err)
		return err
	}
	return nil
}

func (r *activityLogRepository) FindAll(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.ActivityLog, int64, error) {
	var logs []models.ActivityLog
	var total int64

	query := r.db.WithContext(ctx).Model(&models.ActivityLog{})

	// Apply filters (example)
	if userID, ok := filters["user_id"]; ok && userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if staffID, ok := filters["staff_id"]; ok && staffID != "" {
		query = query.Where("staff_id = ?", staffID)
	}
	if action, ok := filters["action"]; ok && action != "" {
		query = query.Where("action = ?", action)
	}
	// Add date range filters if needed

	// Count total records for pagination
	if err := query.Count(&total).Error; err != nil {
		log.Printf("Error counting activity logs: %v", err)
		return nil, 0, err
	}

	// Apply pagination and order
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		log.Printf("Error fetching activity logs: %v", err)
		return nil, 0, err
	}

	return logs, total, nil
}
