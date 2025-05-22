package services

import (
	"context"
	"encoding/json"
	"fmt"
	"invoiceB2B/internal/models"
	"invoiceB2B/internal/repositories"
	"log"
)

type ActivityLogService interface {
	LogActivity(ctx context.Context, staffID *uint, userID *uint, action string, details interface{}, ipAddress string) error
	GetActivityLogs(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.ActivityLog, int64, error)
}

type activityLogService struct {
	repo repositories.ActivityLogRepository
}

func NewActivityLogService(repo repositories.ActivityLogRepository) ActivityLogService {
	return &activityLogService{repo: repo}
}

func (s *activityLogService) LogActivity(ctx context.Context, staffID *uint, userID *uint, action string, details interface{}, ipAddress string) error {
	var detailsJSON *string
	if details != nil {
		detailBytes, err := json.Marshal(details)
		if err != nil {
			log.Printf("Error marshalling activity log details: %v. Details: %+v", err, details)
			// Decide if you want to log without details or return an error
			// For now, log with a placeholder if marshalling fails
			errMsg := fmt.Sprintf(`{"error": "failed to marshal details", "original_details_type": "%T"}`, details)
			detailsJSON = &errMsg
		} else {
			str := string(detailBytes)
			detailsJSON = &str
		}
	}

	var ip *string
	if ipAddress != "" {
		ip = &ipAddress
	}

	logEntry := models.ActivityLog{
		StaffID:   staffID,
		UserID:    userID,
		Action:    action,
		IPAddress: ip,
		// Timestamp is handled by gorm.Model's CreatedAt
	}
	if detailsJSON != nil {
		logEntry.Details = *detailsJSON
	}

	if err := s.repo.Create(ctx, &logEntry); err != nil {
		log.Printf("Failed to create activity log: Action=%s, UserID=%v, StaffID=%v, Error=%v", action, userID, staffID, err)
		return fmt.Errorf("could not create activity log: %w", err)
	}
	return nil
}

func (s *activityLogService) GetActivityLogs(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.ActivityLog, int64, error) {
	// Implementation for fetching logs with pagination and filters
	return s.repo.FindAll(ctx, page, pageSize, filters)
}
